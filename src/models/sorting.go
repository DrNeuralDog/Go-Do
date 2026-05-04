package models

import "sort"

// SortTodosByOrder sorts todos by manual order and time
func SortTodosByOrder(todos []*TodoItem) {
	sort.SliceStable(todos, func(i, j int) bool {
		return todoLess(todos[i], todos[j])
	})
}

// todoLess compares two todos for visible sorting
func todoLess(a, b *TodoItem) bool {
	if a == nil {
		return false
	}

	if b == nil {
		return true
	}

	aHasOrder := hasTodoOrder(a)
	bHasOrder := hasTodoOrder(b)

	if !aHasOrder && !bHasOrder {
		return todoFallbackLess(a, b)
	}

	if !aHasOrder {
		return false
	}

	if !bHasOrder {
		return true
	}

	if a.Order != b.Order {
		return a.Order < b.Order
	}

	return todoFallbackLess(a, b)
}

// hasTodoOrder reports whether todo has manual order
func hasTodoOrder(todo *TodoItem) bool {
	return todo != nil && todo.Order > 0
}

// todoFallbackLess compares todos without manual order
func todoFallbackLess(a, b *TodoItem) bool {
	if a.TodoTime.Equal(b.TodoTime) {
		return a.Name < b.Name
	}

	return a.TodoTime.After(b.TodoTime)
}
