# Todo List Application Build System

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names and paths
BINARY_NAME=GoDo
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_LINUX=$(BINARY_NAME)
BINARY_MACOS=$(BINARY_NAME)
BUILD_DIR=bin

# Build targets
.PHONY: all build clean test deps help package-windows package-all generate-ico embed-ico build-windows-ico

all: deps test build

# Build for current platform (Windows exe with embedded icon)
build:
	@echo "Building Windows executable with icon..."
	@$(MAKE) package-windows

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	@powershell -NoProfile -Command "if (-not (Test-Path '$(BUILD_DIR)')) { New-Item -ItemType Directory -Path '$(BUILD_DIR)' | Out-Null }"
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_WINDOWS) src/main.go

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@powershell -NoProfile -Command "if (-not (Test-Path '$(BUILD_DIR)')) { New-Item -ItemType Directory -Path '$(BUILD_DIR)' | Out-Null }"
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_LINUX) src/main.go

# Build for macOS
build-macos:
	@echo "Building for macOS..."
	@powershell -NoProfile -Command "if (-not (Test-Path '$(BUILD_DIR)')) { New-Item -ItemType Directory -Path '$(BUILD_DIR)' | Out-Null }"
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_MACOS) src/main.go

# Build for all platforms
build-all: build-windows build-linux build-macos

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./tests/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./tests/...
	$(GOCMD) tool cover -html=coverage.out

# Clean build files
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	$(GOMOD) tidy
	$(GOGET) -u ./...

# Run the application
run: build
	@echo "Running application..."
	.\bin\$(BINARY_WINDOWS)

# Development run (with rebuild)
dev: clean deps test build run

# Create data directory (if running for first time)
init:
	@echo "Initializing data directory..."
	@mkdir -p data

# Package Windows executable with embedded icon (taskbar/EXE icon)
package-windows:
	@echo "Packaging Windows EXE with icon..."
	@mkdir -p $(BUILD_DIR)
	@fyne package -os windows -icon doc/Icons/Icon_Work_Version.png -name GoDo -release
	@powershell -NoProfile -Command "Remove-Item -Force -ErrorAction SilentlyContinue .\\$(BUILD_DIR)\\GoDo.exe; Move-Item -Force .\\GoDo.exe .\\$(BUILD_DIR)\\GoDo.exe" || mv -f GoDo.exe $(BUILD_DIR)/GoDo.exe

# Package for all platforms (Windows packaging includes icon)
package-all: package-windows

# --- Alternative Windows build path with resource.syso (ICO) ---
# 1) Convert PNG icon to ICO via PowerShell (.NET) at 256x256
generate-ico:
	@echo "Generating .ico from PNG..."
	@mkdir -p build
	@powershell -NoProfile -Command "Add-Type -AssemblyName System.Drawing; $src='doc/Icons/Icon_Work_Version.png'; $dst='build/Icon_Work_Version.ico'; $bmp=New-Object System.Drawing.Bitmap($src); $bmp256=New-Object System.Drawing.Bitmap($bmp,256,256); $icon=[System.Drawing.Icon]::FromHandle(($bmp256.GetHicon())); $fs=New-Object System.IO.FileStream($dst,[System.IO.FileMode]::Create); $icon.Save($fs); $fs.Close(); $icon.Dispose(); $bmp256.Dispose(); $bmp.Dispose();"

# 2) Embed ICO into resource.syso using rsrc and then build EXE
embed-ico: generate-ico
	@echo "Embedding ICO into resource.syso..."
	@$(GOCMD) install github.com/akavel/rsrc@latest
	@rsrc -ico build/Icon_Work_Version.ico -o resource.syso

# 3) Build Windows exe that picks up resource.syso -> icon appears in Explorer
build-windows-ico: embed-ico
	@echo "Building Windows EXE with embedded .ico (resource.syso)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/GoDo.exe src/main.go

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Run deps, test, and build"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for Windows, Linux, and macOS"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  clean        - Clean build files"
	@echo "  deps         - Install/update dependencies"
	@echo "  run          - Build and run application"
	@echo "  dev          - Full development cycle (clean, deps, test, build, run)"
	@echo "  init         - Initialize data directory"
	@echo "  help         - Show this help message"
