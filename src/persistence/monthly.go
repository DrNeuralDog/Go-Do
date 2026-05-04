package persistence

import (
	"fmt"
	"os"
	"sort"
	"time"

	"godo/src/models"
	"godo/src/utils"
)

// MonthlyManager keeps todos split by month
type MonthlyManager struct {
	fileManager *FileIOManager
	cache       map[string][]*models.TodoItem
}

// NewMonthlyManager creates monthly todo storage
func NewMonthlyManager(dataDir string) *MonthlyManager {
	return &MonthlyManager{
		fileManager: NewFileIOManager(dataDir),
		cache:       make(map[string][]*models.TodoItem),
	}
}

// GetDataDir returns the data directory
func (m *MonthlyManager) GetDataDir() string {
	return m.fileManager.dataDir
}

// GetTodosForMonth returns todos for one month
func (m *MonthlyManager) GetTodosForMonth(year, month int) ([]*models.TodoItem, error) {
	dateKey := utils.FormatDateKey(year, month)

	if todos, exists := m.cache[dateKey]; exists {
		return todos, nil
	}

	todos, err := m.fileManager.LoadTodos(year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to load todos for %s: %w", dateKey, err)
	}

	sortTodosByTime(todos)
	m.cache[dateKey] = todos

	return todos, nil
}

// SaveTodosForMonth saves todos for one month
func (m *MonthlyManager) SaveTodosForMonth(year, month int, todos []*models.TodoItem) error {
	dateKey := utils.FormatDateKey(year, month)

	err := m.fileManager.SaveTodos(year, month, todos)
	if err != nil {
		return fmt.Errorf("failed to save todos for %s: %w", dateKey, err)
	}

	m.cache[dateKey] = todos

	return nil
}

// AddTodo adds todo to its month
func (m *MonthlyManager) AddTodo(todo *models.TodoItem) error {
	year, month := todo.TodoTime.Year(), int(todo.TodoTime.Month())

	todos, err := m.GetTodosForMonth(year, month)
	if err != nil {
		return err
	}

	todos = append(todos, todo)
	sortTodosByTime(todos)

	return m.SaveTodosForMonth(year, month, todos)
}

// UpdateTodo replaces todo found by its old time
func (m *MonthlyManager) UpdateTodo(todo *models.TodoItem, originalTime time.Time) error {
	originalYear, originalMonth := originalTime.Year(), int(originalTime.Month())
	newYear, newMonth := todo.TodoTime.Year(), int(todo.TodoTime.Month())

	if originalYear != newYear || originalMonth != newMonth {
		return m.moveTodoToMonth(todo, originalTime)
	}

	todos, err := m.GetTodosForMonth(originalYear, originalMonth)
	if err != nil {
		return err
	}

	found := false
	for i, existingTodo := range todos {
		if existingTodo.TodoTime.Equal(originalTime) {
			todos[i] = todo
			found = true

			break
		}
	}

	if !found {
		return fmt.Errorf("todo item not found for update")
	}

	sortTodosByTime(todos)

	return m.SaveTodosForMonth(originalYear, originalMonth, todos)
}

// RemoveTodo removes todo by its time
func (m *MonthlyManager) RemoveTodo(todoTime time.Time) error {
	year, month := todoTime.Year(), int(todoTime.Month())

	todos, err := m.GetTodosForMonth(year, month)
	if err != nil {
		return err
	}

	todos, found := removeTodoByTime(todos, todoTime)
	if !found {
		return fmt.Errorf("todo item not found for removal")
	}

	return m.SaveTodosForMonth(year, month, todos)
}

// RemoveTodos removes several todos by their times
func (m *MonthlyManager) RemoveTodos(todoTimes []time.Time) error {
	monthGroups := make(map[string]map[int64]struct{})

	for _, todoTime := range todoTimes {
		dateKey := utils.FormatDateKey(todoTime.Year(), int(todoTime.Month()))

		if monthGroups[dateKey] == nil {
			monthGroups[dateKey] = make(map[int64]struct{})
		}

		monthGroups[dateKey][todoTime.UnixNano()] = struct{}{}
	}

	for dateKey, removeTimes := range monthGroups {
		year, month := utils.ParseDateKey(dateKey)

		todos, err := m.GetTodosForMonth(year, month)
		if err != nil {
			return err
		}

		newTodos := make([]*models.TodoItem, 0, len(todos))
		removed := false

		for _, todo := range todos {
			if _, exists := removeTimes[todo.TodoTime.UnixNano()]; exists {
				removed = true
				continue
			}

			newTodos = append(newTodos, todo)
		}

		if !removed {
			continue
		}

		if err := m.SaveTodosForMonth(year, month, newTodos); err != nil {
			return err
		}
	}

	return nil
}

// GetTodoByTime finds todo by its time
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

// GetAllMonths returns months with data files
func (m *MonthlyManager) GetAllMonths() ([]string, error) {
	return m.fileManager.GetAllMonthlyFiles()
}

// ClearCache drops loaded month data
func (m *MonthlyManager) ClearCache() {
	m.cache = make(map[string][]*models.TodoItem)
}

// GetCacheSize returns cached month count
func (m *MonthlyManager) GetCacheSize() int {
	return len(m.cache)
}

// MigrateAllToYAML converts old TXT monthly files to YAML
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
			continue
		}

		todos, err := m.fileManager.loadTodosTxt(year, month)
		if err != nil {
			// Битый старый файл не стопает всю миграцию
			continue
		}

		if len(todos) == 0 {
			continue
		}

		if err := m.fileManager.SaveTodos(year, month, todos); err != nil {
			return fmt.Errorf("failed to migrate %s to YAML: %w", dateKey, err)
		}
	}

	m.ClearCache()

	return nil
}

// moveTodoToMonth moves todo between monthly files
func (m *MonthlyManager) moveTodoToMonth(todo *models.TodoItem, originalTime time.Time) error {
	originalTodo, err := m.GetTodoByTime(originalTime)
	if err != nil {
		return err
	}

	if err := m.RemoveTodo(originalTime); err != nil {
		return err
	}

	if err := m.AddTodo(todo); err != nil {
		if restoreErr := m.AddTodo(originalTodo); restoreErr != nil {
			return fmt.Errorf("failed to move todo: %w; restore failed: %v", err, restoreErr)
		}

		return err
	}

	return nil
}

// sortTodosByTime keeps newest todos first
func sortTodosByTime(todos []*models.TodoItem) {
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].TodoTime.After(todos[j].TodoTime)
	})
}

// removeTodoByTime cuts todo from slice by time
func removeTodoByTime(todos []*models.TodoItem, todoTime time.Time) ([]*models.TodoItem, bool) {
	for i, todo := range todos {
		if todo.TodoTime.Equal(todoTime) {
			return append(todos[:i], todos[i+1:]...), true
		}
	}

	return todos, false
}
