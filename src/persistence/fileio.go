package persistence

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"godo/src/models"
	"godo/src/utils"
)

const (
	monthlyYAMLVersion = 1
	dataDirPermission  = 0755
	filePermission     = 0644
	yamlFileExt        = ".yaml"
	txtFileExt         = ".txt"
	tempFileExt        = ".tmp"
	backupFileExt      = ".bak"
)

// FileIOManager handles todo files
type FileIOManager struct {
	dataDir string
}

type monthlyYAML struct {
	Version int                `yaml:"version"`
	Todos   []*models.TodoItem `yaml:"todos"`
}

// NewFileIOManager creates file storage
func NewFileIOManager(dataDir string) *FileIOManager {
	return &FileIOManager{
		dataDir: dataDir,
	}
}

// EnsureDataDirectory creates data directory
func (f *FileIOManager) EnsureDataDirectory() error {
	return os.MkdirAll(f.dataDir, dataDirPermission)
}

// getYamlFilePath returns YAML path for a month
func (f *FileIOManager) getYamlFilePath(year, month int) string {
	dateKey := utils.FormatDateKey(year, month)

	return filepath.Join(f.dataDir, dateKey+yamlFileExt)
}

// getTxtFilePath returns old TXT path for a month
func (f *FileIOManager) getTxtFilePath(year, month int) string {
	dateKey := utils.FormatDateKey(year, month)

	return filepath.Join(f.dataDir, dateKey+txtFileExt)
}

// GetFilePath returns YAML path for a month
func (f *FileIOManager) GetFilePath(year, month int) string {
	return f.getYamlFilePath(year, month)
}

// SaveTodos saves todos to monthly YAML
func (f *FileIOManager) SaveTodos(year, month int, todos []*models.TodoItem) error {
	if err := f.EnsureDataDirectory(); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	content := monthlyYAML{Version: monthlyYAMLVersion, Todos: todos}

	data, err := yaml.Marshal(&content)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	filePath := f.getYamlFilePath(year, month)

	return writeFileAtomic(filePath, data)
}

// LoadTodos loads todos from YAML or old TXT
func (f *FileIOManager) LoadTodos(year, month int) ([]*models.TodoItem, error) {
	yamlPath := f.getYamlFilePath(year, month)
	file, err := os.ReadFile(yamlPath)

	if err == nil {
		todos, err := readTodosYAML(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse YAML %s: %w", yamlPath, err)
		}

		return todos, nil
	}

	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read YAML: %w", err)
	}

	return f.loadTodosTxt(year, month)
}

// loadTodosTxt reads old TXT month data
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

	if !scanner.Scan() {
		return []*models.TodoItem{}, nil
	}

	countStr := strings.TrimSpace(scanner.Text())
	count, err := strconv.Atoi(countStr)
	if err != nil {
		// Старый файл битый, считаем его пустым
		return []*models.TodoItem{}, nil
	}

	todos := make([]*models.TodoItem, 0, count)
	for i := 0; i < count; i++ {
		todo, err := f.readTodoItem(scanner)
		if err != nil {
			// Дальше уже не угадаешь, оставляем что прочитали
			break
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

// readTodoItem reads one old TXT todo
func (f *FileIOManager) readTodoItem(scanner *bufio.Scanner) (*models.TodoItem, error) {
	todo := models.NewTodoItem()

	name, err := f.readMultiLineString(scanner)
	if err != nil {
		return nil, fmt.Errorf("failed to read name: %w", err)
	}
	todo.Name = name

	label, err := f.readMultiLineString(scanner)
	if err != nil {
		return nil, fmt.Errorf("failed to read label: %w", err)
	}
	todo.Label = label

	if !scanner.Scan() {
		return nil, fmt.Errorf("unexpected end of file reading level")
	}

	level, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid level: %w", err)
	}
	todo.Level = level

	if !scanner.Scan() {
		return nil, fmt.Errorf("unexpected end of file reading date/time")
	}

	dateStr := scanner.Text()
	parts := strings.Fields(dateStr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid date format: %s", dateStr)
	}

	year, err1 := strconv.Atoi(parts[0])
	month, err2 := strconv.Atoi(parts[1])
	day, err3 := strconv.Atoi(parts[2])
	hour, err4 := strconv.Atoi(parts[3])
	minute, err5 := strconv.Atoi(parts[4])

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		return nil, fmt.Errorf("invalid date components in: %s", dateStr)
	}

	todo.TodoTime = time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)

	place, err := f.readMultiLineString(scanner)
	if err != nil {
		return nil, fmt.Errorf("failed to read place: %w", err)
	}
	todo.Place = place

	content, err := f.readMultiLineString(scanner)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}
	todo.Content = content

	if !scanner.Scan() {
		return nil, fmt.Errorf("unexpected end of file reading status")
	}

	statusStr := scanner.Text()
	parts = strings.Fields(statusStr)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid status format: %s", statusStr)
	}

	done, err1 := strconv.ParseBool(parts[0])
	kind, err2 := strconv.Atoi(parts[1])
	warnTime, err3 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil || err3 != nil {
		return nil, fmt.Errorf("invalid status components in: %s", statusStr)
	}

	todo.Done = done
	todo.Kind = kind
	todo.WarnTime = warnTime

	return todo, nil
}

// readMultiLineString reads old counted string
func (f *FileIOManager) readMultiLineString(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		return "", fmt.Errorf("unexpected end of file reading line count")
	}

	lineCount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return "", fmt.Errorf("invalid line count: %w", err)
	}

	if lineCount < 0 {
		return "", fmt.Errorf("invalid line count: %d", lineCount)
	}

	var result strings.Builder

	for i := 0; i < lineCount; i++ {
		if !scanner.Scan() {
			return "", fmt.Errorf("unexpected end of file reading string")
		}

		if i > 0 {
			result.WriteString("\n")
		}

		result.WriteString(scanner.Text())
	}

	return result.String(), nil
}

// DeleteFile removes monthly data files
func (f *FileIOManager) DeleteFile(year, month int) error {
	yamlPath := f.getYamlFilePath(year, month)
	paths := []string{
		yamlPath,
		yamlPath + tempFileExt,
		yamlPath + backupFileExt,
		f.getTxtFilePath(year, month),
	}

	removed := false
	var firstErr error

	for _, path := range paths {
		err := os.Remove(path)
		if err == nil {
			removed = true
			continue
		}

		if !os.IsNotExist(err) && firstErr == nil {
			firstErr = err
		}
	}

	if removed {
		return nil
	}

	if firstErr != nil {
		return firstErr
	}

	return &os.PathError{Op: "remove", Path: yamlPath, Err: os.ErrNotExist}
}

// FileExists checks monthly data file
func (f *FileIOManager) FileExists(year, month int) bool {
	if _, err := os.Stat(f.getYamlFilePath(year, month)); err == nil {
		return true
	}

	if _, err := os.Stat(f.getTxtFilePath(year, month)); err == nil {
		return true
	}

	return false
}

// GetAllMonthlyFiles returns sorted month keys
func (f *FileIOManager) GetAllMonthlyFiles() ([]string, error) {
	files, err := os.ReadDir(f.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}

		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	months := make(map[string]struct{})
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		dateKey, ok := monthlyFileDateKey(file.Name())
		if !ok {
			continue
		}

		months[dateKey] = struct{}{}
	}

	monthlyFiles := make([]string, 0, len(months))

	for dateKey := range months {
		monthlyFiles = append(monthlyFiles, dateKey)
	}

	sort.Strings(monthlyFiles)

	return monthlyFiles, nil
}

// readTodosYAML reads current and old YAML shapes
func readTodosYAML(data []byte) ([]*models.TodoItem, error) {
	if strings.TrimSpace(string(data)) == "" {
		return []*models.TodoItem{}, nil
	}

	var wrapper monthlyYAML
	wrapperErr := yaml.Unmarshal(data, &wrapper)
	if wrapperErr == nil && wrapper.Todos != nil {
		return wrapper.Todos, nil
	}

	var list []*models.TodoItem
	listErr := yaml.Unmarshal(data, &list)

	if listErr == nil {
		if list == nil {
			return []*models.TodoItem{}, nil
		}

		return list, nil
	}

	return nil, fmt.Errorf("wrapper=%v list=%v", wrapperErr, listErr)
}

// writeFileAtomic replaces file and restores old one on failure
func writeFileAtomic(filePath string, data []byte) error {
	tempPath := filePath + tempFileExt
	backupPath := filePath + backupFileExt

	if err := os.WriteFile(tempPath, data, filePermission); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	tempExists := true
	defer func() {
		if tempExists {
			_ = os.Remove(tempPath)
		}
	}()

	if err := removeFileIfExists(backupPath); err != nil {
		return fmt.Errorf("failed to remove old backup file: %w", err)
	}

	hasOldFile, err := backupExistingFile(filePath, backupPath)

	if err != nil {
		return err
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		if hasOldFile {
			if restoreErr := os.Rename(backupPath, filePath); restoreErr != nil {
				return fmt.Errorf("failed to replace file: %w; restore failed: %v", err, restoreErr)
			}
		}

		return fmt.Errorf("failed to replace file: %w", err)
	}

	tempExists = false
	if hasOldFile {
		_ = os.Remove(backupPath)
	}

	return nil
}

// backupExistingFile moves current file aside
func backupExistingFile(filePath, backupPath string) (bool, error) {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("failed to stat existing file: %w", err)
	}

	if err := os.Rename(filePath, backupPath); err != nil {
		return false, fmt.Errorf("failed to backup existing file: %w", err)
	}

	return true, nil
}

// removeFileIfExists removes file when it is present
func removeFileIfExists(filePath string) error {
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// monthlyFileDateKey extracts YYYYMM from monthly file name
func monthlyFileDateKey(name string) (string, bool) {
	var dateKey string

	switch {
	case strings.HasSuffix(name, yamlFileExt):
		dateKey = strings.TrimSuffix(name, yamlFileExt)
	case strings.HasSuffix(name, txtFileExt):
		dateKey = strings.TrimSuffix(name, txtFileExt)
	default:
		return "", false
	}

	year, month := utils.ParseDateKey(dateKey)

	return dateKey, year != 0 && month != 0
}
