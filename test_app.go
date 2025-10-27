package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"todo-list-migration/src/models"
	"todo-list-migration/src/persistence"
	"todo-list-migration/src/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	fmt.Println("Testing Todo List Application...")

	// Test 1: Create a test todo item
	fmt.Println("Creating test todo item...")
	todo := models.NewTodoItem()
	todo.SetName("Test Todo Item")
	todo.SetContent("This is a test todo item created programmatically")
	todo.SetPlace("Test Location")
	todo.SetLabel("Test Label")
	todo.SetKind(0)  // Event
	todo.SetLevel(2) // High priority
	todo.SetTime(time.Now().Add(time.Hour))
	todo.SetWarnTime(30) // 30 minutes reminder

	fmt.Printf("Created todo: %s\n", todo.GetName())
	fmt.Printf("Priority level: %d\n", todo.GetLevel())
	fmt.Printf("Priority color: %v\n", todo.GetLevelColor())

	// Test 2: Test file I/O
	fmt.Println("\nTesting file I/O...")
	dataDir := "./test_data"
	os.MkdirAll(dataDir, 0755)

	manager := persistence.NewMonthlyManager(dataDir)

	// Save the todo
	err := manager.AddTodo(todo)
	if err != nil {
		log.Printf("Error saving todo: %v", err)
	} else {
		fmt.Println("Todo saved successfully")
	}

	// Load todos back
	todos, err := manager.GetTodosForMonth(time.Now().Year(), int(time.Now().Month()))
	if err != nil {
		log.Printf("Error loading todos: %v", err)
	} else {
		fmt.Printf("Loaded %d todos\n", len(todos))
		for i, loadedTodo := range todos {
			fmt.Printf("  %d. %s\n", i+1, loadedTodo.GetName())
		}
	}

	// Test 3: Test GUI initialization (without showing window)
	fmt.Println("\nTesting GUI initialization...")
	myApp := app.New()
	myWindow := myApp.NewWindow("Test Window")

	testDataDir := filepath.Join(".", "test_data")
	ui.NewMainWindow(myWindow, testDataDir)

	fmt.Println("GUI initialized successfully")

	// Cleanup
	os.RemoveAll(dataDir)
	fmt.Println("\nTest completed successfully!")
}
