# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a cross-platform todo list application built with Go and Fyne, migrated from a C++ Qt implementation. It provides complete todo management with priority levels, reminders, and timeline visualization. Data is persisted in monthly YAML files (with backward compatibility for legacy TXT format).

## Essential Commands

### Development
```bash
# Install/update dependencies
make deps
# or
go mod tidy

# Build for current platform
make build
# or
go build -o bin/GoDo.exe src/main.go

# Run tests
make test
# or
go test ./tests/...

# Run tests with coverage
make test-coverage

# Clean build artifacts
make clean

# Build and run
make run

# Full development cycle
make dev
```

### Cross-Platform Builds
```bash
# Windows
make build-windows
# or
GOOS=windows GOARCH=amd64 go build -o bin/GoDo.exe src/main.go

# macOS
make build-macos

# Linux
make build-linux

# All platforms
make build-all
```

### Testing
```bash
# Run all tests
go test ./tests/...

# Run specific test package
go test ./tests/models/
go test ./tests/persistence/
go test ./tests/ui/
go test ./tests/integration/

# Run tests with verbose output
go test -v ./tests/...

# Run a single test
go test -v -run TestSpecificFunction ./tests/models/
```

## Architecture Overview

### Core Design Principles

**Layered Architecture**: The application uses a clear separation of concerns:
- **Models** (`src/models/`): Core data structures with no external dependencies
- **Persistence** (`src/persistence/`): File I/O and monthly data organization with caching
- **UI** (`src/ui/`): Fyne-based interface components
- **Utils/Localization**: Supporting functionality

**Data Flow**: User interactions flow through MainWindow -> MonthlyManager -> FileIOManager -> Disk. The MonthlyManager maintains an in-memory cache indexed by "YYYYMM" keys for performance.

### Key Architectural Components

#### 1. Monthly Data Organization
The core persistence pattern organizes todos by month:
- Each month's todos stored in `data/YYYYMM.yaml` (e.g., `data/202501.yaml`)
- **MonthlyManager** (`src/persistence/monthly.go`) orchestrates all data operations
- In-memory cache prevents redundant file reads (cache key: "YYYYMM")
- Automatic migration from legacy TXT format to YAML on startup

#### 2. Todo Item Identity
TodoItems don't have explicit IDs. Identity is determined by:
- **Primary key**: Combination of `TodoTime` (time.Time) and `Name` (string)
- When editing, the `originalTime` parameter tracks the item's previous time
- Moving a todo between months is handled by remove-from-old + add-to-new

#### 3. File Format Evolution
- **Current format**: YAML with wrapper structure (`{version: 1, todos: [...]}`)
- **Legacy format**: Custom multi-line text format (TXT)
- FileIOManager provides transparent reading of both formats
- Migration happens automatically on first run via `MigrateAllToYAML()`

#### 4. UI State Management
MainWindow (`src/ui/mainwindow.go`) manages:
- **Current date**: Which month is being viewed (CustomDate)
- **View mode**: Filter (All, Incomplete, Complete, Starred)
- **Theme toggle**: Light vs. Gruvbox dark theme
- Timeline widget displays filtered todos

The Timeline widget (`src/ui/timeline.go`) renders todos grouped by date with smooth scrolling.

### Priority System
Four-level priority (Eisenhower Matrix):
- Level 0: Not Important, Not Urgent (Green - #b8bb26)
- Level 1: Not Important, Urgent (Blue - #83a598)
- Level 2: Important, Not Urgent (Orange - #fe8019)
- Level 3: Important, Urgent (Red - #fb4934)

Colors use Gruvbox palette for visual consistency.

### View Modes
Filtering logic in `src/models/viewmode.go`:
- **ViewAll**: Show all todos
- **ViewIncomplete**: Only todos with `Done == false`
- **ViewComplete**: Only todos with `Done == true`
- **ViewStarred**: Only todos with `Starred == true`

Filters applied at load time in MainWindow.loadTodos().

## Important Implementation Details

### Data Persistence Patterns

**Adding a Todo**:
1. Create TodoItem with all fields
2. Call `MonthlyManager.AddTodo(todo)`
3. Manager determines month from `todo.TodoTime`
4. Loads existing todos for that month
5. Appends new todo, sorts by time (reverse chronological)
6. Saves to YAML and updates cache

**Updating a Todo**:
1. Call `MonthlyManager.UpdateTodo(todo, originalTime)`
2. If month changed: remove from old month, add to new month
3. If same month: find by originalTime+name, replace in slice
4. Resort and save

**Critical**: Always provide `originalTime` when editing to handle month changes correctly.

### File I/O Atomicity
FileIOManager uses atomic writes:
1. Write to `{filename}.tmp`
2. Delete old file if exists
3. Rename `.tmp` to actual filename

This prevents data corruption from crashes/interrupts.

### Caching Strategy
MonthlyManager cache behavior:
- Cache key: `utils.FormatDateKey(year, month)` returns "YYYYMM"
- Cache populated on first load of a month
- Cache updated on every save operation
- `ClearCache()` available for testing or forcing reload

### Time Handling
All times use Go's `time.Time` in UTC:
- Parse date components: `time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)`
- Format for display: Use `utils` package helpers or `time.Format()`
- Reminder calculation: `remindTime = todoTime.Add(-time.Duration(warnTime) * time.Minute)`

### Theme System
Two themes available:
- Light theme: Fyne's default `theme.LightTheme()`
- Gruvbox Black: Custom theme in `src/ui/gruvbox_theme.go`

Toggle via `fyne.CurrentApp().Settings().SetTheme(...)` in MainWindow.

## Common Development Patterns

### Adding a New TodoItem Field
1. Add field to `TodoItem` struct in `src/models/todo.go` with JSON tag
2. Add getter/setter methods following existing pattern
3. Update YAML serialization (automatic via struct tags)
4. Add form input in `src/ui/forms/todoform.go`
5. Update display logic in Timeline/UI widgets
6. Add tests in `tests/models/todo_test.go`
7. Consider migration path for existing data files

### Creating a New View Mode
1. Add constant to `src/models/viewmode.go`
2. Implement `FilterItems()` logic
3. Add label via `GetLabel()`
4. Update `GetNextMode()` cycle
5. Add option to Select widget in MainWindow.setupUI()

### Modifying File Format
**Important**: Changes require careful migration strategy:
1. Increment version number in YAML wrapper
2. Maintain backward compatibility with older versions
3. Add migration logic to handle old → new format
4. Test with existing data files
5. Consider adding a `MigrateToV2()` function similar to `MigrateAllToYAML()`

## Testing Guidelines

### Unit Test Structure
Follow existing patterns in `tests/`:
- Test file names: `{package}_test.go`
- Use table-driven tests for multiple cases
- Mock file I/O by creating temporary directories
- Clean up test artifacts in teardown

### Integration Testing
Full workflow tests in `tests/integration/`:
- Test complete CRUD cycles
- Verify file persistence across app restarts
- Test month transitions and data integrity
- Validate UI interactions end-to-end

## Project Documentation

The `doc/` directory contains comprehensive project documentation:
- **PRD.md**: Product requirements (original project goals)
- **CR_ToDoList_Migration_.md**: Change request describing migration from C++ to Go
- **Implementation.md**: Detailed implementation plan with staged tasks
- **project_structure.md**: Full directory structure explanation
- **UI_UX_doc.md**: UI/UX design specifications

**Workflow Logs** (`doc/WorkflowLogs/`):
- Development progress tracking
- Bug reports and resolutions
- Git operations log

When making significant changes, update relevant documentation and workflow logs as described in `.cursor/rules/workflow.mdc`.

## Code Style Notes

The project uses `.cursor/rules/` for AI-assisted development guidelines. Key points:
- Follow Go conventions: `golint`, `go fmt`
- Use descriptive names: `GetTodosForMonth()` not `GetTodos()`
- Error handling: Always propagate errors with context via `fmt.Errorf("...: %w", err)`
- Comments: Document exported functions and complex logic

## Compatibility Note

This Go implementation maintains **data file compatibility** with the original C++ Qt application. The legacy TXT format is still readable but new writes use YAML. Do not break the TXT parsing logic in `fileio.go:loadTodosTxt()` as users may have existing data.
