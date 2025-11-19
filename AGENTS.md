# Repository Guidelines

## Project Structure & Module Organization
Source code lives under `src/`, split into focused packages such as `ui` (Fyne widgets), `models` (domain structs), `persistence` (file I/O), `data` (seed & parsing helpers), `localization`, `animations`, and `utils`. `src/main.go` wires these modules into the desktop app. Tests mirror those modules inside `tests/`, so add new suites under the same relative path (for example, `tests/ui/...`). Assets, icons, and product docs belong in `doc/`, while compiled artifacts stay in `bin/`. Generated monthly todo files go in `data/` and should remain untracked; run `make init` before first launch if the folder is missing.

## Build, Test, and Development Commands
Use the Makefile to stay consistent:

```bash
make deps            # tidy and download modules
make build           # build current platform binary into bin/todo-list
make run             # build then launch the UI
make test            # run go test -v ./tests/...
make test-coverage   # produce coverage.out and open HTML report
make package-windows # create bin/GoDo.exe with icon
```

`go build ./src` or `go test ./tests/...` are acceptable for quick edits, but check in binaries built via the scripted targets only.

## Coding Style & Naming Conventions
Format Go code with `gofmt ./...` (the CI expects standard Go style: tabs for indentation, imports grouped stdlib/external/internal). Use Go’s mixedCaps for exported identifiers and lowerCamelCase for package-private helpers. File names should describe the feature (`task_editor.go`, `pomodoro_controller.go`). UI resource names follow the existing `Icon_*` or `Theme*` patterns inside `doc/Icons`. Avoid committing generated data; if sample payloads are needed, use anonymized fixtures inside `tests/data/`.

## Testing Guidelines
Add table-driven tests in `*_test.go` files that sit under `tests/` but import the real packages from `src/`. Name functions `Test<Module><Behavior>` to keep `go test -run` filters useful. Maintain or raise coverage when touching persistence or scheduling logic; use `make test-coverage` locally and attach the resulting summary to review notes whenever new subsystems land. Integration tests that spin up real data files should use the `tests/tmpdata` helper paths and clean up via `t.Cleanup`.

## Commit & Pull Request Guidelines
The history uses release-style messages (`Version 0.0.2.3 - <summary>`). Continue that scheme for user-facing work; internal refactors can use `Chore:` or `Fix:` prefixes but still provide a precise scope. Each PR should describe motivation, the key modules touched, and any UI changes (include screenshots if the Fyne layout changed). Link issues or task IDs in the description, mention new Make targets, and call out manual steps (e.g., re-running `make package-windows`). Ensure the checklist includes `make test` passing and note any coverage regressions up front.

## Data & Configuration Tips
Todo content is serialized into `data/YYYYMM.txt`. Keep samples sanitized and never commit personal tasks. When debugging persistence, use a disposable workspace by pointing the `GO_DO_DATA_DIR` env var (if introduced) or temporarily overriding paths inside `data/` helpers—never edit production files in-place without backups.

## Animation System

### Window Opening Animations

Added smooth rotation and scale animations when opening Todo form and Pomodoro windows. Windows now appear to unfold from a central point with an elegant scaling effect.

#### Implementation Details

**Animation Function**: `showWindowWithRotationAnimation()`
- **Duration**: 400ms
- **Frame Rate**: 60 FPS
- **Easing**: Cubic ease-out for smooth deceleration

**Animation Phases**:
1. Initial State: Window starts at 1x1 pixel (invisible)
2. Expansion: Window smoothly scales from 0% to 100% of target size
3. Centering: Window remains centered during entire animation
4. Final State: Window reaches exact target dimensions

**Easing Curve**: Uses cubic ease-out `f(t) = (t-1)³ + 1` for natural, organic feel

#### Files Modified
- `src/ui/forms/todoform.go` - Added animation to Create/Edit windows
- `src/ui/pomodoro_window.go` - Added animation to Pomodoro window

#### Usage
Windows automatically animate when opened - no user configuration needed.

#### Performance
- Minimal CPU usage (60 FPS)
- No blocking operations
- Animation runs asynchronously
- Does not interfere with main UI thread
