# Todo List Application Build System

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names and paths
BINARY_NAME=todo-list
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_LINUX=$(BINARY_NAME)-linux
BINARY_MACOS=$(BINARY_NAME)-macos
BUILD_DIR=bin

# Build targets
.PHONY: all build clean test deps help

all: deps test build

# Build for current platform
build:
	@echo "Building for current platform..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) src/main.go

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_WINDOWS) src/main.go

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_LINUX) src/main.go

# Build for macOS
build-macos:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
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
	./$(BUILD_DIR)/$(BINARY_NAME)

# Development run (with rebuild)
dev: clean deps test build run

# Create data directory (if running for first time)
init:
	@echo "Initializing data directory..."
	@mkdir -p data

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
