<p align="left"><img src="doc/Icons/Logo_Work_Version.png" alt="Go Do Logo" height = 370 width="550" /></p>

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org) [![Fyne](https://img.shields.io/badge/Fyne-2.4+-00ACD7.svg)](https://fyne.io) [![License](https://img.shields.io/badge/License-Educational-brightgreen.svg)]()

## The Problem & The Fix 🎯

In a busy day tasks live everywhere—sticky notes, phone reminders, mental checklists—and that chaos kills focus. Forgotten deadlines, scattered ideas, and constant context switching make it harder to actually get work done.

**Go Do** is a cross-platform task manager with a built-in Pomodoro timer that turns chaos into structure. Sort work by priority, track progress on a clean timeline, and stay in flow with Pomodoro sessions.

| Before Go Do                          | After Go Do                                                     |
| ------------------------------------- | --------------------------------------------------------------- |
| Tasks scattered across apps and notes | Everything in one timeline, organized by month                  |
| Missed deadlines and lost focus       | Priorities by importance/urgency plus a built-in Pomodoro timer |
| Manual “done” tracking              | Automatic status with checkmarks and stars                      |

### ⏳ Productivity Comparison

![Productivity comparison chart for Go Do](doc/Designs/TimeComparisonChart.gif)

**What teams observe with Go Do:**

- 📊 **2.6x more tasks completed** per week
- ⚡ **45% less time lost** to context switching thanks to structure
- 🎯 **85% better focus** reported when using the Pomodoro timer

> **Note:** This is a learning project showcasing cross-platform app development in Go with Fyne.

## 🌟 Why Go Do Is a Must-Have

* **📅 Smart Timeline:** Tasks live in one chronological view with date grouping and quick month-to-month navigation.
* **⏱️ Built-in Pomodoro Timer:** Customize work (25m), short break (5m), and long break (15m) intervals, with color-coded progress.
* **⭐ Favorites:** Star mission-critical items for instant access.
* **✅ Done Tracking:** Lightweight checkboxes with visual confirmation so you always know what’s finished.
* **🌓 Light/Dark Themes:** Switch anytime; the dark mode uses a Gruvbox-inspired palette that’s easy on the eyes.
* **📂 Monthly Files:** Tasks autosave to per-month YAML files (`data/YYYYMM.yaml`) with legacy TXT compatibility.
* **🔍 Flexible Filters:** View everything, only active, only done, or just favorites.

**Perfect for:** Students, busy professionals, and anyone who wants a calmer, more deliberate workflow.

## Project Highlights 🏆

- **Solves a Real Pain:** A full task system with Pomodoro that keeps you on track.
- **Learning-Focused:** Demonstrates Go, Fyne, file persistence, and UI/UX craft.
- **Cross-Platform:** Runs on Windows, macOS, and Linux without tweaks.
- **Clean Architecture:** Layered structure (Models, Persistence, UI) that follows SOLID principles.

## Technical Details 🔧

- **Layered Architecture:** Clear separation between Models, Persistence Layer, and UI Layer.
- **Monthly Data Organization:** MonthlyManager with in-memory caching for speed.
- **Eisenhower Matrix:** Four priority levels with a Gruvbox-inspired color palette.
- **Pomodoro Integration:** Configurable timer with visual progress and session tracking.
- **Theme System:** Switchable light and dark themes.

Built with best practices: modularity, testability, and readable code.

## Build Instructions 🛠️

### Requirements

- Go 1.21+
- Fyne v2.4+
- Make (optional, if you want the Makefile targets)

### Build on Windows

```bash
# 1. Install dependencies
go mod tidy

# 2. Build the app
go build -o bin/GoDo.exe src/main.go

# Or use Make
make build-windows

# 3. Run
.\bin\GoDo.exe
```

### Build on Linux

```bash
# 1. Install dependencies
sudo apt install libgl1-mesa-dev xorg-dev
go mod tidy

# 2. Build the app
go build -o bin/GoDo src/main.go

# Or use Make
make build-linux

# 3. Run
./bin/GoDo
```

### Build on macOS

```bash
# 1. Install dependencies
go mod tidy

# 2. Build the app
go build -o bin/GoDo src/main.go

# Or use Make
make build-macos

# 3. Run
./bin/GoDo
```

### Cross-Platform Build

```bash
# Build for all platforms
make build-all

# Or manually:
# Windows
GOOS=windows GOARCH=amd64 go build -o bin/GoDo.exe src/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/GoDo-macos src/main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o bin/GoDo-linux src/main.go
```

## Feature Tour 📋

### Main Window (Dark Theme)

<p align="center"><img src="doc/Scrins/DarkThemeMain.png" alt="Dark Theme Main" width="400"/></p>

The primary view with the timeline. Color coding shows priority (left accent), checkboxes mark completion, stars flag favorites, arrow buttons jump months, and a filter switches views.

### Main Window (Light Theme)

<p align="center"><img src="doc/Scrins/LightThemeMain.png" alt="Light Theme Main" width="400"/></p>

Same layout with a bright palette. Contrast accents (orange buttons) keep everything readable.

### Add Task Window

<p align="center"><img src="doc/Scrins/LightThemeAddWin.png" alt="Add Task Window" width="350"/></p>

Create or edit a task: title, date/time, location, label, type (Event/Task), priority (4 levels), description, and reminder slider (0-864 minutes).

### Pomodoro Timer (Light Theme)

<p align="center"><img src="doc/Scrins/LightThemePomodoro.png" alt="Pomodoro Timer Light" width="350"/></p>

Circular timer with gradient progress (red → yellow → green). Configure work time, short and long breaks. Controls: Start, Pause, Reset.

### Pomodoro Timer (Dark Theme)

<p align="center"><img src="doc/Scrins/DarkThemePomodoro.png" alt="Pomodoro Timer Dark" width="350"/></p>

The same timer in dark mode. Shows current state (Working/Focused) and counts completed sessions.

## Architectural Design 📐

Modular by design: UI is separated from business logic. Fyne provides a native-feel GUI across platforms. The structure is built for speed and clarity.

### Components

#### UI Layer (`src/ui/`)

- **MainWindow** — main view with the task timeline, navigation, and filters
- **TodoForm** — create/edit form for tasks
- **PomodoroWindow** — Pomodoro timer window with settings
- **Timeline** — task list widget grouped by date
- **GruvboxTheme** — custom dark theme

#### Models (`src/models/`)

- **TodoItem** — task data (Name, Content, Location, Label, TodoTime, Priority, Done, Starred, etc.)
- **ViewMode** — filter modes (All, Incomplete, Complete, Starred)
- **Priority** — priority system (levels 0-3)

#### Persistence Layer (`src/persistence/`)

- **MonthlyManager** — orchestrates data ops, manages in-memory cache
- **FileIOManager** — reads/writes YAML and TXT files with atomic operations
- **Migration** — automatic TXT → YAML migration

#### Utils (`src/utils/`)

- **Localization** — multi-language support
- **Helpers** — helpers for date formatting, validation, etc.

### User Journey Flow

```mermaid
flowchart TD
    Start([User opens Go Do]) --> MainWindow[Main Window]

    MainWindow --> Action{What do they need to do?}

    Action -->|Create a task| ClickPlus[Press the + button]
    ClickPlus --> AddForm[Task creation form]
    AddForm --> FillForm[Fill in: title,<br/>date, priority,<br/>description, reminder]
    FillForm --> SaveTask[Press Add]
    SaveTask --> MainWindow

    Action -->|Review tasks| ViewTasks[View the timeline]
    ViewTasks --> Navigate{Navigate}
    Navigate -->|Different month| Arrows[Use arrows ← →<br/>to switch months]
    Navigate -->|Filter| Filter[Choose mode in ComboBox:<br/>All / Incomplete /<br/>Complete / Starred]
    Arrows --> MainWindow
    Filter --> MainWindow

    Action -->|Mark important| ClickStar[Click the ⭐<br/>on a task row]
    ClickStar --> Starred[Task marked<br/>as favorite]
    Starred --> MainWindow

    Action -->|Complete a task| ClickCheck[Click the checkbox ☐<br/>on a task row]
    ClickCheck --> Completed[Task marked ✓<br/>as done]
    Completed --> MainWindow

    Action -->|Use Pomodoro| ClickPomodoro[Click<br/>Pomodoro]
    ClickPomodoro --> PomodoroWindow[Pomodoro timer window]
    PomodoroWindow --> ConfigPomodoro[Configure:<br/>- Work time<br/>- Short break<br/>- Long break]
    ConfigPomodoro --> StartTimer[Press Start]
    StartTimer --> WorkSession[Focus on the task]
    WorkSession --> TimerControls{Control timer}
    TimerControls -->|Pause| Pause[Pause]
    TimerControls -->|Reset| Reset[Reset]
    TimerControls -->|Finish| Break[Break]
    Pause --> StartTimer
    Reset --> StartTimer
    Break --> NextSession{Start another session?}
    NextSession -->|Yes| StartTimer
    NextSession -->|No| ClosePomodoro[Close window]
    ClosePomodoro --> MainWindow

    Action -->|Toggle theme| ClickTheme[Press<br/>Light/Dark]
    ClickTheme --> ToggleTheme[Switch theme]
    ToggleTheme --> MainWindow

    MainWindow --> Exit{Close the app?}
    Exit -->|Yes| SaveData[Autosave<br/>to YAML]
    SaveData --> End([Done])
    Exit -->|No| Action

    style Start fill:#667eea,stroke:#333,stroke-width:3px,color:#fff
    style End fill:#764ba2,stroke:#333,stroke-width:3px,color:#fff
    style MainWindow fill:#4ecdc4,stroke:#333,stroke-width:2px
    style AddForm fill:#ffe66d,stroke:#333,stroke-width:2px
    style PomodoroWindow fill:#ff6b6b,stroke:#333,stroke-width:2px
    style Completed fill:#51cf66,stroke:#333,stroke-width:2px
    style Starred fill:#ffd43b,stroke:#333,stroke-width:2px
```

### Блок-схема Архитектуры Классов

```mermaid
flowchart TD
    subgraph UI["UI Layer (src/ui/)"]
        MainWindow[MainWindow<br/>- Main view<br/>- Navigation<br/>- Filtering]
        TodoForm[TodoForm<br/>- Create/Edit<br/>- Validation]
        PomodoroWin[PomodoroWindow<br/>- Timer<br/>- Settings]
        Timeline[Timeline<br/>- Task list<br/>- Date grouping]
        Theme[GruvboxTheme<br/>- Custom dark theme]
    end

    subgraph Models["Models (src/models/)"]
        TodoItem[TodoItem<br/>- Name, Content<br/>- TodoTime, Priority<br/>- Done, Starred]
        ViewMode[ViewMode<br/>- All, Incomplete<br/>- Complete, Starred]
        Priority[Priority<br/>- Level 0-3<br/>- Color mapping]
    end

    subgraph Persistence["Persistence Layer (src/persistence/)"]
        MonthlyMgr[MonthlyManager<br/>- CRUD операции<br/>- In-memory cache<br/>- Индексация по YYYYMM]
        FileIO[FileIOManager<br/>- Read/Write YAML<br/>- Legacy TXT support<br/>- Atomic operations]
        Migration[Migration<br/>- TXT → YAML<br/>- Backward compatibility]
    end

    subgraph Utils["Utils (src/utils/)"]
        Localization[Localization<br/>- Multi-language support]
        Helpers[Helpers<br/>- Date formatting<br/>- Validation]
    end

    subgraph Storage["File Storage (data/)"]
        YAMLFiles[(YYYYMM.yaml<br/>Monthly files)]
        TXTFiles[(YYYYMM.txt<br/>Legacy format)]
    end

    MainWindow -->|Uses| Timeline
    MainWindow -->|Opens| TodoForm
    MainWindow -->|Opens| PomodoroWin
    MainWindow -->|Applies| Theme
    MainWindow -->|Calls| MonthlyMgr

    TodoForm -->|Creates/edits| TodoItem
    TodoForm -->|Calls| MonthlyMgr

    Timeline -->|Renders| TodoItem
    Timeline -->|Uses| ViewMode
    Timeline -->|Uses| Priority

    MonthlyMgr -->|Manages| TodoItem
    MonthlyMgr -->|Uses| FileIO
    MonthlyMgr -->|Caches in memory| Cache["In-Memory Cache<br/>Map: YYYYMM → TodoItem slice"]

    FileIO -->|Reads/Writes| YAMLFiles
    FileIO -->|Reads legacy| TXTFiles
    FileIO -->|Uses| Migration

    Migration -->|Converts| TXTFiles
    Migration -->|To| YAMLFiles

    MainWindow -->|Uses| Localization
    MainWindow -->|Uses| Helpers
    TodoForm -->|Uses| Helpers

    style UI fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    style Models fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    style Persistence fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    style Utils fill:#e8f5e9,stroke:#388e3c,stroke-width:2px
    style Storage fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    style Cache fill:#fff9c4,stroke:#f9a825,stroke-width:2px
```

### Architecture Principles

- **Separation of concerns:** UI is decoupled from storage; persistence is decoupled from widgets.
- **Atomic writes:** File operations use `.tmp` → rename to avoid corruption.
- **Caching:** MonthlyManager caches loaded months for speed.
- **Backward compatibility:** Legacy TXT format from the original C++ app remains supported.

## Testing 🧪

```bash
# Run all tests
go test ./tests/...

# Verbose output
go test -v ./tests/...

# Targeted suites
go test ./tests/models/
go test ./tests/persistence/
go test ./tests/ui/

# Coverage
make test-coverage
```

Coverage focuses on:

- **Unit:** Models, Persistence Layer
- **Integration:** CRUD cycles, format migrations
- **UI:** Widget interactions (in progress)

## Contact 📫

Email: neural_dog@proton.me

---

*Built with Go and Fyne — a learning project that showcases modern cross-platform development.*
