package persistence

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"godo/src/models"
	"godo/src/utils"
)

// FileIOManager handles file operations for todo data persistence
type FileIOManager struct {
	dataDir string
}

// NewFileIOManager creates a new file I/O manager
func NewFileIOManager(dataDir string) *FileIOManager {
	return &FileIOManager{
		dataDir: dataDir,
	}
}

// EnsureDataDirectory creates the data directory if it doesn't exist
func (f *FileIOManager) EnsureDataDirectory() error {
	return os.MkdirAll(f.dataDir, 0755)
}

// getYamlFilePath returns YAML file path for a specific year/month
func (f *FileIOManager) getYamlFilePath(year, month int) string {
	dateKey := utils.FormatDateKey(year, month)
	return filepath.Join(f.dataDir, dateKey+".yaml")
}

// getTxtFilePath returns legacy TXT file path for a specific year/month
func (f *FileIOManager) getTxtFilePath(year, month int) string {
	dateKey := utils.FormatDateKey(year, month)
	return filepath.Join(f.dataDir, dateKey+".txt")
}

// GetFilePath returns the preferred file path (YAML) for a specific year/month
func (f *FileIOManager) GetFilePath(year, month int) string {
	return f.getYamlFilePath(year, month)
}

// SaveTodos saves todo items to a monthly file
func (f *FileIOManager) SaveTodos(year, month int, todos []*models.TodoItem) error {
	if err := f.EnsureDataDirectory(); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// YAML format with a simple wrapper for future extensions
	type monthlyYAML struct {
		Version int                `yaml:"version"`
		Todos   []*models.TodoItem `yaml:"todos"`
	}

	content := monthlyYAML{Version: 1, Todos: todos}

	data, err := yaml.Marshal(&content)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	filePath := f.getYamlFilePath(year, month)
	tempPath := filePath + ".tmp"

	// Write atomically
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp YAML: %w", err)
	}

	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to remove existing YAML: %w", err)
		}
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		return fmt.Errorf("failed to rename temp YAML: %w", err)
	}

	return nil
}

// writeTodoItem writes a single todo item to the file
func (f *FileIOManager) writeTodoItem(writer *bufio.Writer, todo *models.TodoItem) error {
	// Write name (with line count prefix)
	if err := f.writeMultiLineString(writer, todo.Name); err != nil {
		return err
	}

	// Write label (with line count prefix)
	if err := f.writeMultiLineString(writer, todo.Label); err != nil {
		return err
	}

	// Write level
	_, err := writer.WriteString(fmt.Sprintf("%d\n", todo.Level))
	if err != nil {
		return err
	}

	// Write date/time components
	date := todo.TodoTime
	_, err = writer.WriteString(fmt.Sprintf("%d %d %d %d %d\n",
		date.Year(), int(date.Month()), date.Day(), date.Hour(), date.Minute()))
	if err != nil {
		return err
	}

	// Write place (with line count prefix)
	if err := f.writeMultiLineString(writer, todo.Place); err != nil {
		return err
	}

	// Write content (with line count prefix)
	if err := f.writeMultiLineString(writer, todo.Content); err != nil {
		return err
	}

	// Write done status, kind, and warn time
	_, err = writer.WriteString(fmt.Sprintf("%t %d %d\n", todo.Done, todo.Kind, todo.WarnTime))
	if err != nil {
		return err
	}

	return nil
}

// writeMultiLineString writes a string that may contain newlines
func (f *FileIOManager) writeMultiLineString(writer *bufio.Writer, s string) error {
	// Count lines
	lineCount := 1
	for _, char := range s {
		if char == '\n' {
			lineCount++
		}
	}

	// Write line count
	_, err := writer.WriteString(fmt.Sprintf("%d\n", lineCount))
	if err != nil {
		return err
	}

	// Write the string
	_, err = writer.WriteString(s + "\n")
	if err != nil {
		return err
	}

	return nil
}

// LoadTodos loads todo items from a monthly file
func (f *FileIOManager) LoadTodos(year, month int) ([]*models.TodoItem, error) {
	// Prefer YAML
	yamlPath := f.getYamlFilePath(year, month)
	if file, err := os.ReadFile(yamlPath); err == nil {
		// Try wrapper format first
		type monthlyYAML struct {
			Version int                `yaml:"version"`
			Todos   []*models.TodoItem `yaml:"todos"`
		}
		var wrapper monthlyYAML
		if err := yaml.Unmarshal(file, &wrapper); err == nil && wrapper.Todos != nil {
			// Ensure deterministic order: newest first
			return wrapper.Todos, nil
		}
		// Fallback: direct list
		var list []*models.TodoItem
		if err := yaml.Unmarshal(file, &list); err == nil {
			return list, nil
		}
		return []*models.TodoItem{}, nil
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read YAML: %w", err)
	}

	// Fallback to legacy TXT
	return f.loadTodosTxt(year, month)
}

// loadTodosTxt loads legacy TXT format and is robust to trailing blank lines
func (f *FileIOManager) loadTodosTxt(year, month int) ([]*models.TodoItem, error) {
	filePath := f.getTxtFilePath(year, month)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.TodoItem{}, nil
		}
		return nil, fmt.Errorf("failed to open legacy TXT file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read first line: count
	if !scanner.Scan() {
		return []*models.TodoItem{}, nil
	}
	countStr := strings.TrimSpace(scanner.Text())
	count, err := strconv.Atoi(countStr)
	if err != nil {
		// corrupted file, ignore
		return []*models.TodoItem{}, nil
	}

	todos := make([]*models.TodoItem, 0, count)
	for i := 0; i < count; i++ {
		todo, _, err := f.readTodoItem(scanner)
		if err != nil {
			// stop on parse error and return what we have so far
			break
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// readTodoItem reads a single todo item from the scanner
func (f *FileIOManager) readTodoItem(scanner *bufio.Scanner) (*models.TodoItem, int, error) {
	todo := models.NewTodoItem()
	linesRead := 0

	// Read name
	name, lines, err := f.readMultiLineString(scanner)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read name: %w", err)
	}
	todo.Name = name
	linesRead += lines

	// Read label
	label, lines, err := f.readMultiLineString(scanner)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read label: %w", err)
	}
	todo.Label = label
	linesRead += lines

	// Read level
	if !scanner.Scan() {
		return nil, 0, fmt.Errorf("unexpected end of file reading level")
	}
	level, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, 0, fmt.Errorf("invalid level: %w", err)
	}
	todo.Level = level
	linesRead++

	// Read date/time
	if !scanner.Scan() {
		return nil, 0, fmt.Errorf("unexpected end of file reading date/time")
	}
	dateStr := scanner.Text()
	parts := strings.Fields(dateStr)
	if len(parts) != 5 {
		return nil, 0, fmt.Errorf("invalid date format: %s", dateStr)
	}

	year, err1 := strconv.Atoi(parts[0])
	month, err2 := strconv.Atoi(parts[1])
	day, err3 := strconv.Atoi(parts[2])
	hour, err4 := strconv.Atoi(parts[3])
	minute, err5 := strconv.Atoi(parts[4])

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		return nil, 0, fmt.Errorf("invalid date components in: %s", dateStr)
	}

	todo.TodoTime = time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	linesRead++

	// Read place
	place, lines, err := f.readMultiLineString(scanner)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read place: %w", err)
	}
	todo.Place = place
	linesRead += lines

	// Read content
	content, lines, err := f.readMultiLineString(scanner)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read content: %w", err)
	}
	todo.Content = content
	linesRead += lines

	// Read done status, kind, and warn time
	if !scanner.Scan() {
		return nil, 0, fmt.Errorf("unexpected end of file reading status")
	}
	statusStr := scanner.Text()
	parts = strings.Fields(statusStr)
	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("invalid status format: %s", statusStr)
	}

	done, err1 := strconv.ParseBool(parts[0])
	kind, err2 := strconv.Atoi(parts[1])
	warnTime, err3 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil || err3 != nil {
		return nil, 0, fmt.Errorf("invalid status components in: %s", statusStr)
	}

	todo.Done = done
	todo.Kind = kind
	todo.WarnTime = warnTime
	linesRead++

	return todo, linesRead, nil
}

// readMultiLineString reads a multi-line string from the scanner
func (f *FileIOManager) readMultiLineString(scanner *bufio.Scanner) (string, int, error) {
	if !scanner.Scan() {
		return "", 0, fmt.Errorf("unexpected end of file reading line count")
	}

	lineCount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return "", 0, fmt.Errorf("invalid line count: %w", err)
	}

	var result strings.Builder
	linesRead := 1

	for i := 0; i < lineCount && scanner.Scan(); i++ {
		line := scanner.Text()

		// Handle Windows/Unix line endings
		if i < lineCount-1 {
			// Not the last line, add back the newline
			result.WriteString(line)
			result.WriteString("\n")
		} else {
			// Last line, don't add trailing newline
			result.WriteString(line)
		}
		linesRead++
	}

	return result.String(), linesRead, nil
}

// DeleteFile removes a monthly data file
func (f *FileIOManager) DeleteFile(year, month int) error {
	// Try deleting both formats
	yamlErr := os.Remove(f.getYamlFilePath(year, month))
	txtErr := os.Remove(f.getTxtFilePath(year, month))
	if yamlErr == nil || txtErr == nil {
		return nil
	}
	// if both failed, return yamlErr (could be not-exist)
	return yamlErr
}

// FileExists checks if a monthly data file exists
func (f *FileIOManager) FileExists(year, month int) bool {
	if _, err := os.Stat(f.getYamlFilePath(year, month)); err == nil {
		return true
	}
	if _, err := os.Stat(f.getTxtFilePath(year, month)); err == nil {
		return true
	}
	return false
}

// GetAllMonthlyFiles returns a list of all monthly data files
func (f *FileIOManager) GetAllMonthlyFiles() ([]string, error) {
	files, err := os.ReadDir(f.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	months := make(map[string]struct{})
	// Prefer YAML when both exist
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasSuffix(name, ".yaml") {
			base := strings.TrimSuffix(name, ".yaml")
			if len(base) == 6 {
				if _, err := strconv.Atoi(base); err == nil {
					months[base] = struct{}{}
				}
			}
		}
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasSuffix(name, ".txt") {
			base := strings.TrimSuffix(name, ".txt")
			if len(base) == 6 {
				if _, err := strconv.Atoi(base); err == nil {
					if _, exists := months[base]; !exists {
						months[base] = struct{}{}
					}
				}
			}
		}
	}

	var monthlyFiles []string
	for k := range months {
		monthlyFiles = append(monthlyFiles, k)
	}
	return monthlyFiles, nil
}
