# CR_Refactoring: Go Do Codebase Refactoring Plan

**Document Type:** Change Requirements
**Created:** 2025-11-19
**Status:** Planning
**Goal:** Refactor Go Do codebase to follow Go best practices without changing functionality

---

## Overview

This document outlines a comprehensive refactoring plan for the Go Do application. The refactoring focuses on:
- Improving code organization and maintainability
- Adding interfaces for better testability
- Eliminating code duplication
- Simplifying complex functions
- Following Go naming conventions
- Extracting magic numbers to constants

**IMPORTANT:** No functionality or business logic changes - only code structure improvements.

---

## Phase 1: Foundation - Create Helper Packages and Interfaces

### 1.1 Create UI Helper Packages Structure
**Estimated time:** 30 minutes

- [x] Create directory structure:
  ```
  src/ui/helpers/
  src/ui/widgets/
  ```

### 1.2 Extract Color Utilities
**Files affected:** `src/ui/style_helpers.go`, `src/ui/forms/todoform.go`
**Estimated time:** 45 minutes

- [x] Create `src/ui/helpers/color.go`
- [x] Move `toNRGBA()` from style_helpers.go:162-167
- [x] Move `lighten()` from style_helpers.go:169-183
- [x] Move `darken()` from style_helpers.go:185-195
- [x] Move `hex()` from style_helpers.go:197-205
- [x] Remove duplicate `lighten()` from todoform.go:752-764
- [x] Update all imports in:
  - [x] src/ui/style_helpers.go
  - [x] src/ui/forms/todoform.go
  - [x] src/ui/timeline.go
  - [x] src/ui/mainwindow.go
  - [x] src/ui/pomodoro_window.go
  - [ ] src/ui/notes_window.go
- [x] Test: Run `go build` to ensure no compilation errors

### 1.3 Extract Theme Detection Utilities
**Files affected:** All UI files
**Estimated time:** 1 hour

- [x] Create `src/ui/helpers/theme.go`
- [x] Add function `IsLightTheme() bool`
- [x] Add function `GetBackgroundColor() color.Color`
- [x] Add function `GetForegroundColor() color.Color`
- [x] Add function `GetCardColor() color.Color`
- [x] Replace all instances of `if _, ok := currentTheme.(*LightSoftTheme)` in:
  - [x] src/ui/mainwindow.go (lines 178, 191, 224, 255, etc.)
  - [x] src/ui/timeline.go (lines 223, 237, 283, 522, 770, 811, 853)
  - [x] src/ui/forms/todoform.go (lines 519, 539, 554)
  - [x] src/ui/pomodoro_window.go (lines 84, 110, 184)
  - [ ] src/ui/notes_window.go (if applicable)
- [x] Test: Run application and toggle theme to verify functionality

### 1.4 Extract Constants
**Estimated time:** 30 minutes

- [x] Create `src/ui/constants.go`
- [x] Extract window dimensions:
  ```go
  const (
      MainWindowWidth       = 420
      MainWindowHeight      = 800
      PomodoroWindowWidth   = 400
      PomodoroWindowHeight  = 500
      NotesWindowWidth      = 400
      NotesWindowHeight     = 300
  )
  ```
- [x] Extract UI element sizes:
  ```go
  const (
      TimelineItemHeight    = 80
      ButtonHeight          = 44
      ButtonPadding         = 24
      BorderRadius          = 8
      ProgressRingSize      = 200
  )
  ```
- [x] Extract color constants (reference colors used in UI)
- [x] Replace magic numbers in:
  - [x] src/ui/mainwindow.go:150 (window size)
  - [x] src/ui/timeline.go:55 (item height)
  - [x] src/ui/pomodoro_window.go (window sizes)
  - [ ] src/ui/notes_window.go (window sizes)
  - [ ] Other hardcoded values
- [x] Test: Build and run to verify no visual changes

### 1.5 Create Persistence Interfaces
**Files affected:** `src/persistence/`
**Estimated time:** 1 hour

- [x] Create `src/persistence/interfaces.go`
- [x] Define `TodoRepository` interface:
  ```go
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
  ```
- [x] Define `ConfigRepository` interface:
  ```go
  type ConfigRepository interface {
      LoadConfig() (*models.Config, error)
      SaveConfig(config *models.Config) error
      GetConfigPath() string
  }
  ```
- [x] Verify `MonthlyManager` implements `TodoRepository`
- [x] Verify `ConfigManager` implements `ConfigRepository`
- [x] Test: Run `go build` to ensure interfaces are satisfied

### 1.6 Update UI to Use Interfaces
**Files affected:** All UI components
**Estimated time:** 1.5 hours

- [x] Update `src/ui/mainwindow.go`:
  - [x] Change `dataManager *persistence.MonthlyManager` to `dataManager persistence.TodoRepository` (line 28)
  - [x] Change `configManager *persistence.ConfigManager` to `configManager persistence.ConfigRepository` (line 29)
  - [x] Update `NewMainWindow()` signature
- [x] Update `src/ui/timeline.go`:
  - [x] Change `dataManager *persistence.MonthlyManager` to `dataManager persistence.TodoRepository` (line 26)
  - [x] Update `NewTimeline()` signature
- [x] Update `src/ui/forms/todoform.go`:
  - [x] Change `dataManager *persistence.MonthlyManager` to `dataManager persistence.TodoRepository` (line 27)
  - [x] Update `NewTodoForm()` signature
- [ ] Update `src/ui/pomodoro_window.go`:
  - [ ] Check for dataManager usage and update if needed
- [x] Update `src/main.go`:
  - [x] Verify dependency injection is explicit
  - [x] Pass interfaces instead of concrete types
- [x] Test: Full application test - create, edit, delete todos

---

## Phase 2: Extract Shared Logic

### 2.1 Extract Todo Sorting Logic
**Files affected:** `src/ui/mainwindow.go`, potentially others
**Estimated time:** 45 minutes

- [x] Create `src/models/sorting.go`
- [x] Add `SortTodosByOrder(todos []*TodoItem)` function
- [x] Extract sorting logic from mainwindow.go:373-398
- [x] Replace duplicate sorting code in:
  - [x] src/ui/mainwindow.go:373-398 (loadTodos)
  - [x] src/ui/mainwindow.go:516-538 (onTodoReorder)
  - [x] src/ui/mainwindow.go:594-616 (onTodoReorder UI update)
- [x] Test: Verify todo ordering works correctly after changes

### 2.2 Add ViewMode String Conversion Methods
**Files affected:** `src/models/viewmode.go`, `src/ui/mainwindow.go`
**Estimated time:** 30 minutes

- [x] Update `src/models/viewmode.go`:
  - [x] Add `String() string` method to ViewMode
  - [x] Add `ViewModeFromString(s string) ViewMode` function
- [x] Update `src/ui/mainwindow.go`:
  - [x] Replace manual conversion at lines 234-245 with `ViewModeFromString()`
  - [x] Replace manual conversion at lines 715-726 with `.String()`
  - [x] Replace manual conversion at lines 743-756 with `.String()`
- [x] Test: Change view modes and restart app to verify persistence

### 2.3 Extract Layout Helpers
**Files affected:** `src/ui/style_helpers.go`
**Estimated time:** 30 minutes

- [x] Create `src/ui/helpers/layout.go`
- [x] Move `CreateSpacer()` from style_helpers.go
- [x] Move `CreateCardStyle()` from style_helpers.go
- [x] Move other layout-related helper functions
- [x] Update imports across UI files
- [ ] Test: Build and verify UI renders correctly

---

## Phase 3: Split Large Files - Custom Widgets

### 3.1 Extract Widgets from style_helpers.go
**File affected:** `src/ui/style_helpers.go` (1036 lines)
**Estimated time:** 2 hours

- [ ] Create `src/ui/widgets/button.go`
  - [ ] Move `RoundIconButton` struct and methods (lines ~207-360)
  - [ ] Move `SimpleRectButton` struct and methods (lines ~362-495)
  - [ ] Update package to `package widgets`
- [ ] Create `src/ui/widgets/select.go`
  - [ ] Move `CustomSelect` struct and methods (lines ~497-680)
- [ ] Create `src/ui/widgets/spinner.go`
  - [ ] Move `NumberSpinner` struct and methods (lines ~682-850)
- [ ] Create `src/ui/widgets/progress.go`
  - [ ] Move `ProgressRing` struct and methods (lines ~852-1036)
- [ ] Create `src/ui/helpers/window.go`
  - [ ] Move `FlashWindow()` and window utilities
- [ ] Update `src/ui/style_helpers.go`:
  - [ ] Keep only generic helper functions
  - [ ] Add imports to new widget packages
- [ ] Update imports in all files using these widgets:
  - [ ] src/ui/mainwindow.go
  - [ ] src/ui/timeline.go
  - [ ] src/ui/forms/todoform.go
  - [ ] src/ui/pomodoro_window.go
  - [ ] src/ui/notes_window.go
- [ ] Test: Build and run full application, test all widgets

### 3.2 Extract Timeline Widgets
**File affected:** `src/ui/timeline.go` (907 lines)
**Estimated time:** 2.5 hours

- [ ] Create `src/ui/timeline/` directory
- [ ] Create `src/ui/timeline/timeline.go`
  - [ ] Move `Timeline` struct (lines 23-46)
  - [ ] Move main Timeline methods
  - [ ] Keep around 250-300 lines
- [ ] Create `src/ui/timeline/renderer.go`
  - [ ] Move `timelineRenderer` struct (lines 48-71)
  - [ ] Move rendering methods
  - [ ] Keep around 250-300 lines
- [ ] Create `src/ui/timeline/item.go`
  - [ ] Move `tappableTodo` widget (lines ~150-250)
  - [ ] Move item creation and interaction logic
- [ ] Create `src/ui/timeline/widgets/checkbox.go`
  - [ ] Move `squareCheckbox` widget
- [ ] Create `src/ui/timeline/widgets/status.go`
  - [ ] Move `statusIndicator` widget
- [ ] Create `src/ui/timeline/theme.go`
  - [ ] Move `foregroundOverrideTheme`
- [ ] Update `src/ui/mainwindow.go`:
  - [ ] Change import from `ui` to `ui/timeline`
  - [ ] Update `NewTimeline()` call
- [ ] Test: Full timeline functionality - scroll, tap, check, delete

---

## Phase 4: Split Large Files - TodoForm

### 4.1 Split todoform.go
**File affected:** `src/ui/forms/todoform.go` (917 lines)
**Estimated time:** 2.5 hours

- [ ] Create `src/ui/forms/todoform_dialogs.go`
  - [ ] Move `ShowCreateDialog()` method
  - [ ] Move `ShowEditDialog()` method
  - [ ] Move dialog-related helper methods
- [ ] Create `src/ui/forms/todoform_windows.go`
  - [ ] Move `ShowCreateWindow()` method
  - [ ] Move `ShowEditWindow()` method
  - [ ] Move window-related helper methods
- [ ] Create `src/ui/forms/widgets/reminder_slider.go`
  - [ ] Move `ReminderSlider` struct and methods
- [ ] Create `src/ui/forms/widgets/rect_button.go`
  - [ ] Move `rectButton` struct and methods (if not already moved)
- [ ] Keep in `src/ui/forms/todoform.go`:
  - [ ] TodoForm struct
  - [ ] NewTodoForm()
  - [ ] Core form building logic
  - [ ] Validation methods
  - [ ] Submit handlers
  - [ ] Should be ~300-400 lines
- [ ] Update imports in files using TodoForm
- [ ] Test: Create new todo, edit existing todo via both dialog and window

---

## Phase 5: Split Large Files - MainWindow

### 5.1 Split mainwindow.go
**File affected:** `src/ui/mainwindow.go` (763 lines)
**Estimated time:** 3 hours

- [ ] Create `src/ui/mainwindow_setup.go`
  - [ ] Move `setupUI()` method (lines 146-323)
  - [ ] Move `setupBottomButtons()` method
  - [ ] Move other UI construction methods
- [ ] Create `src/ui/mainwindow_handlers.go`
  - [ ] Move `onAddTodoClicked()` (lines 325-341)
  - [ ] Move `onPomodoroClicked()` (lines 343-358)
  - [ ] Move `onToggleTheme()` (lines 360-371)
  - [ ] Move `onViewModeChanged()` (lines 428-432)
  - [ ] Move `onTodoDeleted()` (lines 434-441)
  - [ ] Move `onDateChanged()` (lines 443-471)
  - [ ] Move `onTodoEdited()` (lines 473-491)
  - [ ] Move other event handlers
- [ ] Create `src/ui/mainwindow_reorder.go`
  - [ ] Move `onTodoReorder()` (lines 493-619)
  - [ ] Move `onReorderFinished()` (lines 621-690)
  - [ ] Move reorder helper methods
- [ ] Create `src/ui/mainwindow_config.go`
  - [ ] Move `loadConfig()` (lines 692-761)
  - [ ] Move `saveConfig()` (lines 763-end)
  - [ ] Move config-related methods
- [ ] Keep in `src/ui/mainwindow.go`:
  - [ ] MainWindow struct (lines 25-55)
  - [ ] NewMainWindow() (lines 57-144)
  - [ ] Show() method
  - [ ] loadTodos() (lines 373-426)
  - [ ] Core window logic
  - [ ] Should be ~250-300 lines
- [ ] Test: Full application test - all features

---

## Phase 6: Simplify Complex Functions

### 6.1 Refactor MainWindow.setupUI()
**File affected:** `src/ui/mainwindow_setup.go` (after Phase 5)
**Estimated time:** 1.5 hours

- [ ] Extract `buildHeader() fyne.CanvasObject`
  - [ ] Title label creation
  - [ ] Add button creation
  - [ ] Header container assembly
- [ ] Extract `buildControls() fyne.CanvasObject`
  - [ ] Date navigation buttons
  - [ ] View mode select
  - [ ] Theme toggle button
  - [ ] Controls container assembly
- [ ] Extract `buildTimeline() fyne.CanvasObject`
  - [ ] Timeline widget creation
  - [ ] Scroll container setup
- [ ] Extract `assembleLayout(...) fyne.CanvasObject`
  - [ ] Final layout composition
- [ ] Update `setupUI()` to call extracted methods
- [ ] Result: setupUI() should be ~30-40 lines
- [ ] Test: Build and verify UI layout

### 6.2 Refactor MainWindow.onTodoReorder()
**File affected:** `src/ui/mainwindow_reorder.go` (after Phase 5)
**Estimated time:** 1.5 hours

- [ ] Extract `getDayTodos() []*models.TodoItem`
  - [ ] Filter todos by current date
  - [ ] Return day's todos
- [ ] Extract `findTodoIndex(dayTodos []*models.TodoItem, todo *models.TodoItem) int`
  - [ ] Find todo index in slice
  - [ ] Return -1 if not found
- [ ] Extract `calculateNewIndex(currentIdx, delta, length int) int`
  - [ ] Calculate new position
  - [ ] Handle boundaries
- [ ] Extract `reorderTodos(todos []*models.TodoItem, oldIdx, newIdx int)`
  - [ ] Perform slice reordering
- [ ] Extract `updateOrderValues(todos []*models.TodoItem)`
  - [ ] Set Order field for each todo
- [ ] Extract `syncOrderToVisible(dayTodos []*models.TodoItem)`
  - [ ] Update visible todos with new order
- [ ] Update `onTodoReorder()` to call extracted methods
- [ ] Result: onTodoReorder() should be ~30-40 lines
- [ ] Test: Drag todos up and down, verify order persists

### 6.3 Refactor Timeline.createTodoItem()
**File affected:** `src/ui/timeline/renderer.go` (after Phase 3)
**Estimated time:** 1.5 hours

- [ ] Extract `buildColorIndicator(todo *models.TodoItem) fyne.CanvasObject`
- [ ] Extract `buildCheckbox(todo *models.TodoItem) fyne.CanvasObject`
- [ ] Extract `buildNameLabel(todo *models.TodoItem) fyne.CanvasObject`
- [ ] Extract `buildTimeDisplay(todo *models.TodoItem) fyne.CanvasObject`
- [ ] Extract `buildStatusIndicator(todo *models.TodoItem) fyne.CanvasObject`
- [ ] Extract `buildDeleteButton(todo *models.TodoItem) fyne.CanvasObject`
- [ ] Extract `assembleTodoRow(...) fyne.CanvasObject`
- [ ] Extract `wrapWithInteraction(content fyne.CanvasObject, todo *models.TodoItem) fyne.CanvasObject`
- [ ] Update `createTodoItem()` to call extracted methods
- [ ] Result: createTodoItem() should be ~20-30 lines
- [ ] Test: Verify all todo items render correctly

### 6.4 Refactor FileIOManager.readTodoItem()
**File affected:** `src/persistence/fileio.go`
**Estimated time:** 1 hour

- [ ] Extract `readName(scanner *bufio.Scanner, todo *models.TodoItem, lines *int) error`
- [ ] Extract `readLabel(scanner *bufio.Scanner, todo *models.TodoItem, lines *int) error`
- [ ] Extract `readLevel(scanner *bufio.Scanner, todo *models.TodoItem, lines *int) error`
- [ ] Extract `readTime(scanner *bufio.Scanner, todo *models.TodoItem, lines *int) error`
- [ ] Extract `readWarnTime(scanner *bufio.Scanner, todo *models.TodoItem, lines *int) error`
- [ ] Extract similar methods for other fields
- [ ] Update `readTodoItem()` to call extracted methods
- [ ] Result: readTodoItem() should be ~20-30 lines
- [ ] Test: Load legacy TXT files, verify data integrity

---

## Phase 7: Improve Naming Conventions

### 7.1 Update TodoItem Methods
**File affected:** `src/models/todo.go`
**Estimated time:** 1 hour

- [ ] Rename `HaveDone()` to `IsDone()` (line 116)
- [ ] Update all calls to `HaveDone()`:
  - [ ] Search codebase with `Grep`
  - [ ] Replace in all files
- [ ] Review getters/setters for removal:
  - [ ] Keep: `SetLevel()` (has validation)
  - [ ] Keep: `GetLevelString()`, `GetLevelColor()` (computed values)
  - [ ] Keep: `GetPriority()` (computed value)
  - [ ] Consider removing simple getters like:
    - `GetName()` → direct access to `todo.Name`
    - `GetContent()` → direct access to `todo.Content`
    - `GetLabel()` → direct access to `todo.Label`
    - `GetTodoTime()` → direct access to `todo.TodoTime`
  - [ ] **Decision point:** Discuss with team if breaking API change is acceptable
- [ ] If removing getters/setters:
  - [ ] Update all usage across codebase
  - [ ] Update timeline rendering
  - [ ] Update form handling
- [ ] Test: Full application test

### 7.2 Review CustomDate Usage
**File affected:** `src/utils/timeutils.go`
**Estimated time:** 30 minutes

- [ ] Search for `CustomDate` usage across codebase
- [ ] Verify if it's actually used or legacy code
- [ ] If unused:
  - [ ] Remove `CustomDate` struct
  - [ ] Remove related methods
  - [ ] Clean up imports
- [ ] If used:
  - [ ] Add documentation explaining why it exists
  - [ ] Consider migration path to `time.Time`
- [ ] Test: Build and run application

---

## Phase 8: Error Handling Improvements

### 8.1 Improve Error Handling in main.go
**File affected:** `src/main.go`
**Estimated time:** 15 minutes

- [ ] Add logging for migration errors (line 85):
  ```go
  if err := migrator.MigrateAllToYAML(); err != nil {
      log.Printf("Warning: Migration failed: %v", err)
  }
  ```
- [ ] Review other silent error handling
- [ ] Test: Trigger migration error and verify logging

### 8.2 Add Error Dialogs in Timeline
**File affected:** `src/ui/timeline/timeline.go` (after Phase 3)
**Estimated time:** 30 minutes

- [ ] Find all `_ = r.timeline.dataManager.UpdateTodo(...)` patterns
- [ ] Replace with proper error handling:
  ```go
  if err := r.timeline.dataManager.UpdateTodo(...); err != nil {
      dialog.ShowError(err, r.timeline.window)
  }
  ```
- [ ] Update similar patterns in:
  - [ ] Checkbox toggle
  - [ ] Todo editing
  - [ ] Todo deletion
- [ ] Test: Force errors and verify dialogs appear

### 8.3 Consistent Error Context
**Files affected:** `src/persistence/fileio.go` and others
**Estimated time:** 45 minutes

- [ ] Review all error returns in fileio.go
- [ ] Add context to errors that lack it
- [ ] Use `fmt.Errorf("operation: %w", err)` pattern consistently
- [ ] Review error handling in:
  - [ ] src/persistence/monthly.go
  - [ ] src/persistence/config.go
- [ ] Test: Trigger various errors and verify messages are clear

---

## Phase 9: Documentation and Final Touches

### 9.1 Add Package Documentation
**Estimated time:** 1 hour

- [ ] Create `src/models/doc.go`
  ```go
  /*
  Package models contains the core domain models for the todo application.

  The primary types are:
    - TodoItem: Represents a single todo with priority, time, and metadata
    - ViewMode: Filtering modes for displaying todos
    - Config: Application configuration and UI state
    - PomodoroTimer: Timer state for pomodoro sessions
  */
  package models
  ```
- [ ] Create `src/persistence/doc.go`
- [ ] Create `src/ui/doc.go`
- [ ] Create `src/utils/doc.go`
- [ ] Test: Run `go doc` to verify documentation

### 9.2 Update Project Documentation
**Files affected:** Documentation files
**Estimated time:** 30 minutes

- [ ] Update `CLAUDE.md` with new structure:
  - [ ] Document new package layout
  - [ ] Update architecture section
  - [ ] Add notes about interfaces
- [ ] Update `doc/project_structure.md`:
  - [ ] Reflect new file organization
  - [ ] Document new packages
- [ ] Create workflow log entry in `doc/WorkflowLogs/`

### 9.3 Code Review Checklist
**Estimated time:** 1 hour

- [ ] Run `go fmt ./...`
- [ ] Run `go vet ./...`
- [ ] Run `go build` - verify no errors
- [ ] Run `go test ./...` - verify all tests pass
- [ ] Manual testing checklist:
  - [ ] Launch application
  - [ ] Create new todo
  - [ ] Edit existing todo
  - [ ] Delete todo
  - [ ] Change view modes (all, incomplete, complete, starred)
  - [ ] Navigate between dates
  - [ ] Toggle theme
  - [ ] Reorder todos via drag-and-drop
  - [ ] Open Pomodoro window
  - [ ] Add notes to todo
  - [ ] Verify data persists after restart
  - [ ] Check YAML files are properly formatted

---

## Testing Strategy

### Unit Testing
- After each phase, run `go build` to catch compilation errors
- After each phase, run `go test ./...` to verify existing tests pass
- Consider adding new unit tests for extracted functions

### Integration Testing
- After Phase 5 (major splits complete), do full manual test
- After Phase 8 (error handling), test error scenarios
- After Phase 9, complete final testing checklist

### Regression Testing
- Keep a test todo list with various scenarios
- After each major phase, verify:
  - Data loads correctly
  - CRUD operations work
  - UI renders properly
  - Theme switching works
  - Config persists

---

## Risk Mitigation

### Backup Strategy
- [ ] Before starting, create git branch: `refactoring-2025-11`
- [ ] Commit after each completed phase
- [ ] Tag stable checkpoints: `refactor-phase-1`, `refactor-phase-2`, etc.

### Rollback Plan
- If issues arise, rollback to previous phase tag
- Each phase is independent enough to rollback individually

### Incremental Delivery
- Each phase produces a working application
- Can pause and deliver after any phase
- Priority: Complete Phases 1-3 first (foundation + big files)

---

## Success Criteria

### Code Quality Metrics
- [ ] No files over 500 lines
- [ ] No functions over 100 lines
- [ ] All magic numbers replaced with constants
- [ ] All persistence accessed via interfaces
- [ ] No code duplication for theme detection, sorting, or color utils

### Functionality Verification
- [ ] All existing features work identically
- [ ] No data loss or corruption
- [ ] Performance is same or better
- [ ] Application builds without warnings

### Maintainability Improvements
- [ ] New developers can navigate codebase easily
- [ ] Each file has single, clear responsibility
- [ ] Unit testing is straightforward
- [ ] Future features can be added without touching 10+ files

---

## Timeline Estimate

| Phase | Description | Estimated Time |
|-------|-------------|----------------|
| Phase 1 | Foundation - interfaces & helpers | 5-6 hours |
| Phase 2 | Extract shared logic | 2 hours |
| Phase 3 | Split style_helpers & timeline | 4-5 hours |
| Phase 4 | Split todoform | 2.5 hours |
| Phase 5 | Split mainwindow | 3 hours |
| Phase 6 | Simplify complex functions | 4-5 hours |
| Phase 7 | Improve naming | 1.5 hours |
| Phase 8 | Error handling | 1.5 hours |
| Phase 9 | Documentation & testing | 2.5 hours |
| **TOTAL** | | **26-31 hours** |

**Recommended schedule:** 3-4 days of focused work, or 1-2 weeks with other tasks.

---

## Notes

- This refactoring does NOT change any business logic or functionality
- All changes are purely structural and organizational
- The application should work identically before and after
- Each phase checkpoint should be tested and committed
- If any phase reveals issues, pause and address before continuing

---

## Completion Checklist

- [ ] All phases completed
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Git history is clean with clear commit messages
- [ ] Final manual testing completed
- [ ] Code review performed
- [ ] Ready to merge to main branch

---

**End of Document**
