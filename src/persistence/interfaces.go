package persistence

import (
	"time"

	"godo/src/models"
)

type TodoRepository interface {
	GetTodosForMonth(year, month int) ([]*models.TodoItem, error)
	SaveTodosForMonth(year, month int, todos []*models.TodoItem) error
	AddTodo(todo *models.TodoItem) error
	UpdateTodo(todo *models.TodoItem, originalTime time.Time) error
	RemoveTodo(todoTime time.Time) error
	RemoveTodos(todoTimes []time.Time) error
	GetTodoByTime(todoTime time.Time) (*models.TodoItem, error)
	GetAllMonths() ([]string, error)
	ClearCache()
	MigrateAllToYAML() error
}

type ConfigRepository interface {
	LoadConfig() (*models.Config, error)
	SaveConfig(config *models.Config) error
	GetConfigPath() string
}

var (
	_ TodoRepository   = (*MonthlyManager)(nil)
	_ ConfigRepository = (*ConfigManager)(nil)
)
