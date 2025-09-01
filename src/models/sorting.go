package models

import "sort"

// SortTodosByOrder sorts todos by Order (ascending, zeros last) then by TodoTime descending,
// with a final tie-breaker on Name for stability.
func SortTodosByOrder(todos []*TodoItem) {
	sort.SliceStable(todos, func(i, j int) bool {
		a := todos[i]
		b := todos[j]

		if a.Order == 0 && b.Order == 0 {
			if a.TodoTime.Equal(b.TodoTime) {
				return a.Name < b.Name
			}
			return a.TodoTime.After(b.TodoTime)
		}

		if a.Order == 0 {
			return false
		}
		if b.Order == 0 {
			return true
		}
		if a.Order != b.Order {
			return a.Order < b.Order
		}

		if a.TodoTime.Equal(b.TodoTime) {
			return a.Name < b.Name
		}
		return a.TodoTime.After(b.TodoTime)
	})
}
