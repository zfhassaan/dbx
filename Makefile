.PHONY: build build-windows build-linux build-all clean help

# Application name
APP_NAME=dbx

# Version (can be set via environment variable)
VERSION ?= 1.0.0

# Build directory
BUILD_DIR=dist

# Go build flags
LDFLAGS=-s -w -X main.version=$(VERSION)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build for current platform
	@echo "Building for current platform..."
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) main.go
	@echo "✅ Build complete: $(BUILD_DIR)/$(APP_NAME)"

build-windows: ## Build Windows executable (amd64)
	@echo "Building Windows executable..."
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe main.go
	@echo "✅ Windows build complete: $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe"

build-linux: ## Build Linux executable (amd64)
	@echo "Building Linux executable..."
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 main.go
	@echo "✅ Linux build complete: $(BUILD_DIR)/$(APP_NAME)-linux-amd64"

build-linux-arm64: ## Build Linux executable (ARM64)
	@echo "Building Linux ARM64 executable..."
	@GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 main.go
	@echo "✅ Linux ARM64 build complete: $(BUILD_DIR)/$(APP_NAME)-linux-arm64"

build-all: clean build-windows build-linux ## Build for all platforms
	@echo "✅ All builds complete!"

clean: ## Clean build artifacts
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)/$(APP_NAME)*
	@echo "✅ Clean complete"

install: build ## Build and install to system
	@echo "Installing $(APP_NAME)..."
	@go install -ldflags "$(LDFLAGS)"
	@echo "✅ Installation complete"

test: ## Run all tests with comprehensive reporting
	@echo "Running comprehensive test suite..."
	@go run scripts/runtests/main.go

test-simple: ## Run tests (simple)
	@echo "Running tests..."
	@go test -v ./...

coverage: ## Generate test coverage report
	@echo "Generating test coverage report..."
	@go run scripts/coverage/main.go

coverage-html: ## Generate HTML coverage report
	@echo "Generating HTML coverage report..."
	@go test -coverprofile=coverage.out ./tests/internal/cloud -coverpkg=./internal/cloud
	@go test -coverprofile=coverage_db.out ./tests/internal/db -coverpkg=./internal/db
	@go test -coverprofile=coverage_logs.out ./tests/internal/logs -coverpkg=./internal/logs
	@go test -coverprofile=coverage_notify.out ./tests/internal/notify -coverpkg=./internal/notify
	@go test -coverprofile=coverage_scheduler.out ./tests/internal/scheduler -coverpkg=./internal/scheduler
	@go test -coverprofile=coverage_utils.out ./tests/internal/utils -coverpkg=./internal/utils
	@echo "Coverage profiles generated. Opening HTML report..."
	@go tool cover -html=coverage.out

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies updated"

