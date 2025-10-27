# Project Structure for Todo List Application Migration

## Root Directory
```
todo-list-migration/
├── src/
│   ├── main.go                    # Application entry point
│   ├── app.go                     # Main application window and lifecycle
│   ├── models/
│   │   ├── todo.go                # TodoItem struct and related types
│   │   ├── priority.go            # Priority levels and color definitions
│   │   └── viewmode.go            # View mode filtering logic
│   ├── persistence/
│   │   ├── fileio.go              # File I/O operations for data persistence
│   │   ├── format.go              # Custom text format parsing/serialization
│   │   └── monthly.go             # Monthly data organization utilities
│   ├── ui/
│   │   ├── mainwindow.go          # Main application window
│   │   ├── timeline.go            # Timeline visualization widget
│   │   ├── todoitem.go            # Todo item display widget
│   │   ├── todolist.go            # Todo list container widget
│   │   ├── forms/
│   │   │   ├── todoform.go        # Todo creation/editing form
│   │   │   ├── priorityselector.go # Priority selection widget
│   │   │   └── datetimeselector.go # Date/time picker widget
│   │   └── dialogs/
│   │       ├── confirmdialog.go   # Confirmation dialogs
│   │       └── errordialog.go     # Error message dialogs
│   ├── localization/
│   │   └── english.go             # English language strings and translations
│   ├── utils/
│   │   ├── timeutils.go           # Date/time utility functions
│   │   ├── validation.go          # Input validation utilities
│   │   └── colors.go              # Color definitions and theme utilities
│   └── animations/
│       ├── scroll.go              # Smooth scrolling animations
│       └── transitions.go         # Page transition effects
├── data/                          # Data files (auto-generated)
│   └── YYYYMM.txt                 # Monthly todo data files
├── docs/                          # Documentation
│   ├── Implementation.md          # Implementation plan
│   ├── project_structure.md       # This file
│   ├── UI_UX_doc.md              # UI/UX specifications
│   ├── CR_ToDoList_Migration_.md  # Change Request document
│   ├── PRD.md                     # Product Requirements Document
│   ├── Archive/                   # Archived documents
│   └── WorkflowLogs/              # Development workflow logs
│       ├── DevelopmentLog.md       # Development progress log
│       ├── BugLog.md              # Bug reports and fixes
│       ├── GitLog.md              # Git operations log
│       └── UserInteractionLog.md # User interaction log
├── tests/                         # Test files
│   ├── models/
│   │   └── todo_test.go           # TodoItem model tests
│   ├── persistence/
│   │   └── fileio_test.go         # File I/O tests
│   ├── ui/
│   │   └── widgets_test.go        # UI widget tests
│   └── integration/
│       └── workflow_test.go       # End-to-end workflow tests
├── assets/                        # Static assets
│   ├── icons/                     # Application icons
│   └── styles/                    # UI styling files
├── build/                         # Build outputs (generated)
├── bin/                           # Binary outputs (generated)
│   ├── release/                   # Release builds
│   └── debug/                     # Debug builds
├── go.mod                         # Go module definition
├── go.sum                         # Go module checksums
├── Makefile                       # Build automation
├── README.md                      # Project overview and setup instructions
└── .gitignore                     # Git ignore patterns
```

## Detailed Structure

### Source Code Organization (`src/`)
**Purpose**: Contains all source code organized by functional layers

**models/**: Core data structures and business logic
- `todo.go`: TodoItem struct definition matching original C++ implementation
- `priority.go`: Priority level definitions and color mappings
- `viewmode.go`: View filtering logic (All, Incomplete, Reminders)

**persistence/**: Data storage and retrieval
- `fileio.go`: Core file I/O operations with error handling
- `format.go`: Custom text format parsing and serialization
- `monthly.go`: Monthly data organization and caching utilities

**ui/**: User interface components
- `mainwindow.go`: Main application window with menu and layout
- `timeline.go`: Timeline visualization with date grouping
- `todoitem.go`: Individual todo item display widget
- `todolist.go`: Container for todo list with scrolling support
- `forms/`: Form dialogs for data entry
- `dialogs/`: Modal dialogs for confirmations and errors

**localization/**: Language and internationalization
- `english.go`: English language strings and UI translations

**utils/**: Utility functions and helpers
- `timeutils.go`: Date/time manipulation and formatting
- `validation.go`: Input validation and sanitization
- `colors.go`: Color definitions and theme management

**animations/**: Animation and transition effects
- `scroll.go`: Smooth scrolling physics and animations
- `transitions.go`: Page transition and state change effects

### Data Storage (`data/`)
**Purpose**: File-based data persistence organized by month
- Monthly text files in `YYYYMM.txt` format matching original Qt implementation
- Automatic directory creation and file management
- Atomic write operations to prevent data corruption

### Documentation (`docs/`)
**Purpose**: Project documentation and development tracking

**Core Documentation**:
- `Implementation.md`: Detailed implementation plan and progress tracking
- `project_structure.md`: Project organization and file structure (this file)
- `UI_UX_doc.md`: UI/UX design specifications and guidelines
- `PRD.md`: Product Requirements Document
- `CR_ToDoList_Migration_.md`: Change Request document

**Archive** (`Archive/`):
- Completed CR documents and outdated development logs
- Historical project documentation

**Workflow Logs** (`WorkflowLogs/`):
- `DevelopmentLog.md`: Development progress and task completion tracking
- `BugLog.md`: Bug reports, debugging sessions, and resolution tracking
- `GitLog.md`: Git operations and version control activities
- `UserInteractionLog.md`: User interaction tracking and feedback

### Testing (`tests/`)
**Purpose**: Comprehensive test coverage for quality assurance

**Unit Tests**:
- Model tests for core data structures and business logic
- Persistence tests for file I/O operations
- UI component tests for widget behavior

**Integration Tests**:
- End-to-end workflow testing
- Data persistence and retrieval validation
- Cross-component interaction testing

### Build and Deployment
**Build Output** (`build/`): Intermediate build artifacts and temporary files
**Binary Output** (`bin/`): Compiled executables organized by build type and platform
**Assets** (`assets/`): Icons, styles, and static resources

### Project Configuration
- `go.mod`/`go.sum`: Go module dependencies and version management
- `Makefile`: Build automation and cross-compilation scripts
- `README.md`: Project overview, setup instructions, and development guide
- `.gitignore`: Git ignore patterns for build artifacts and temporary files

## Development Workflow Integration

### Standard Directory Purposes
This project follows established conventions for directory organization:

- **src/**: Contains the main project source files - core application code
- **docs/**: Project documentation including PRD, CR, and generated documents
- **docs/Archive/**: Historical documents and completed change requests
- **docs/WorkflowLogs/**: Development tracking and logging
- **tests/**: Test files for unit and integration testing
- **build/**: Intermediate build files and artifacts
- **bin/**: Executable binaries organized by build type and platform
- **assets/**: Static resources like icons and styling files

### Module Organization Patterns
- **Layered Architecture**: Clear separation between UI, business logic, and data layers
- **Feature-Based Grouping**: Related functionality grouped in subdirectories
- **Consistent Naming**: Go package naming conventions followed throughout
- **Import Organization**: Clean import paths with logical package hierarchy

### Build Structure
- **Cross-Platform Builds**: Support for Windows, macOS, and Linux compilation
- **Build Types**: Separate debug and release configurations
- **Dependency Management**: Go modules for reproducible builds
- **Automation**: Makefile for common build and test operations

This structure supports the migration requirements while maintaining clean separation of concerns, testability, and maintainability for the Go Fyne implementation.
