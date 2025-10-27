package persistence

import (
	"fmt"
	"os"
	"sort"
	"time"

	"todo-list-migration/src/models"
	"todo-list-migration/src/utils"
)

// MonthlyManager handles monthly organization of todo data
type MonthlyManager struct {
	fileManager *FileIOManager
	cache       map[string][]*models.TodoItem // Cache for loaded monthly data
}

// NewMonthlyManager creates a new monthly manager
func NewMonthlyManager(dataDir string) *MonthlyManager {
	return &MonthlyManager{
		fileManager: NewFileIOManager(dataDir),
		cache:       make(map[string][]*models.TodoItem),
	}
}

// GetDataDir returns the data directory path
func (m *MonthlyManager) GetDataDir() string {
	return m.fileManager.dataDir
}

// GetTodosForMonth retrieves todos for a specific month, loading from file if necessary
func (m *MonthlyManager) GetTodosForMonth(year, month int) ([]*models.TodoItem, error) {
	dateKey := utils.FormatDateKey(year, month)

	// Check cache first
	if todos, exists := m.cache[dateKey]; exists {
		return todos, nil
	}

	// Load from file
	todos, err := m.fileManager.LoadTodos(year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to load todos for %s: %w", dateKey, err)
	}

	// Sort todos by time (reverse chronological order like original)
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].TodoTime.After(todos[j].TodoTime)
	})

	// Cache the results
	m.cache[dateKey] = todos

	return todos, nil
}

// SaveTodosForMonth saves todos for a specific month
func (m *MonthlyManager) SaveTodosForMonth(year, month int, todos []*models.TodoItem) error {
	dateKey := utils.FormatDateKey(year, month)

	err := m.fileManager.SaveTodos(year, month, todos)
	if err != nil {
		return fmt.Errorf("failed to save todos for %s: %w", dateKey, err)
	}

	// Update cache
	m.cache[dateKey] = todos

	return nil
}

// AddTodo adds a new todo item to the appropriate month
func (m *MonthlyManager) AddTodo(todo *models.TodoItem) error {
	year, month := todo.TodoTime.Year(), int(todo.TodoTime.Month())

	// Get existing todos for the month
	todos, err := m.GetTodosForMonth(year, month)
	if err != nil {
		return err
	}

	// Add new todo
	todos = append(todos, todo)

	// Sort and save
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].TodoTime.After(todos[j].TodoTime)
	})

	return m.SaveTodosForMonth(year, month, todos)
}

// UpdateTodo updates an existing todo item
func (m *MonthlyManager) UpdateTodo(todo *models.TodoItem, originalTime time.Time) error {
	originalYear, originalMonth := originalTime.Year(), int(originalTime.Month())
	newYear, newMonth := todo.TodoTime.Year(), int(todo.TodoTime.Month())

	// If the month changed, we need to move the todo
	if originalYear != newYear || originalMonth != newMonth {
		// Remove from original month
		if err := m.RemoveTodo(originalTime); err != nil {
			return err
		}
		// Add to new month
		return m.AddTodo(todo)
	}

	// Update within the same month
	todos, err := m.GetTodosForMonth(originalYear, originalMonth)
	if err != nil {
		return err
	}

	// Find and update the todo
	found := false
	for i, existingTodo := range todos {
		if existingTodo.TodoTime.Equal(originalTime) &&
			existingTodo.Name == todo.Name {
			todos[i] = todo
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("todo item not found for update")
	}

	// Sort and save
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].TodoTime.After(todos[j].TodoTime)
	})

	return m.SaveTodosForMonth(originalYear, originalMonth, todos)
}

// RemoveTodo removes a todo item by its time
func (m *MonthlyManager) RemoveTodo(todoTime time.Time) error {
	year, month := todoTime.Year(), int(todoTime.Month())

	todos, err := m.GetTodosForMonth(year, month)
	if err != nil {
		return err
	}

	// Find and remove the todo
	for i, todo := range todos {
		if todo.TodoTime.Equal(todoTime) {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}

	return m.SaveTodosForMonth(year, month, todos)
}

// RemoveTodos removes multiple todos by their times
func (m *MonthlyManager) RemoveTodos(todoTimes []time.Time) error {
	// Group by month for efficient processing
	monthGroups := make(map[string][]time.Time)

	for _, todoTime := range todoTimes {
		dateKey := utils.FormatDateKey(todoTime.Year(), int(todoTime.Month()))
		monthGroups[dateKey] = append(monthGroups[dateKey], todoTime)
	}

	// Remove from each month
	for dateKey, times := range monthGroups {
		year, month := utils.ParseDateKey(dateKey)

		todos, err := m.GetTodosForMonth(year, month)
		if err != nil {
			return err
		}

		// Remove todos
		newTodos := make([]*models.TodoItem, 0, len(todos))
		for _, todo := range todos {
			shouldRemove := false
			for _, removeTime := range times {
				if todo.TodoTime.Equal(removeTime) {
					shouldRemove = true
					break
				}
			}
			if !shouldRemove {
				newTodos = append(newTodos, todo)
			}
		}

		if err := m.SaveTodosForMonth(year, month, newTodos); err != nil {
			return err
		}
	}

	return nil
}

// GetTodoByTime finds a todo item by its time (for editing)
func (m *MonthlyManager) GetTodoByTime(todoTime time.Time) (*models.TodoItem, error) {
	year, month := todoTime.Year(), int(todoTime.Month())

	todos, err := m.GetTodosForMonth(year, month)
	if err != nil {
		return nil, err
	}

	for _, todo := range todos {
		if todo.TodoTime.Equal(todoTime) {
			return todo, nil
		}
	}

	return nil, fmt.Errorf("todo item not found")
}

// GetAllMonths returns all months that have data files
func (m *MonthlyManager) GetAllMonths() ([]string, error) {
	return m.fileManager.GetAllMonthlyFiles()
}

// ClearCache clears the internal cache
func (m *MonthlyManager) ClearCache() {
	m.cache = make(map[string][]*models.TodoItem)
}

// GetCacheSize returns the number of cached months
func (m *MonthlyManager) GetCacheSize() int {
	return len(m.cache)
}

// MigrateAllToYAML converts existing legacy TXT monthly files to YAML format.
// If a YAML file already exists for a month, it will be left untouched.
func (m *MonthlyManager) MigrateAllToYAML() error {
	months, err := m.GetAllMonths()
	if err != nil {
		return err
	}

	for _, dateKey := range months {
		year, month := utils.ParseDateKey(dateKey)
		if year == 0 {
			continue
		}

		yamlPath := m.fileManager.getYamlFilePath(year, month)
		if _, err := os.Stat(yamlPath); err == nil {
			// Already migrated
			continue
		}

		// Try loading from legacy TXT directly
		todos, err := m.fileManager.loadTodosTxt(year, month)
		if err != nil {
			// Skip problematic month but continue others
			continue
		}

		if len(todos) == 0 {
			// Nothing to migrate
			continue
		}

		// Save in YAML using current saver
		if err := m.fileManager.SaveTodos(year, month, todos); err != nil {
			return fmt.Errorf("failed to migrate %s to YAML: %w", dateKey, err)
		}
	}

	// Clear cache to ensure fresh loads from YAML
	m.ClearCache()
	return nil
}
