# Go Project Makefile for go-ctl-initializer
BINARY_NAME=go-ctl
MAIN_PATH=cmd/server
BUILD_DIR=bin
DOCKER_IMAGE=go-ctl:latest

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build flags
LDFLAGS=-ldflags="-w -s"
BUILD_FLAGS=-trimpath

.PHONY: all build clean test coverage deps fmt vet lint run dev docker-build docker-run help

# Default target
all: clean deps fmt vet test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)/main.go $(MAIN_PATH)/handlers.go $(MAIN_PATH)/templates.go

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)/main.go $(MAIN_PATH)/handlers.go $(MAIN_PATH)/templates.go
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)/main.go $(MAIN_PATH)/handlers.go $(MAIN_PATH)/templates.go
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)/main.go $(MAIN_PATH)/handlers.go $(MAIN_PATH)/templates.go
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)/main.go $(MAIN_PATH)/handlers.go $(MAIN_PATH)/templates.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f *.zip
	@rm -rf test-*

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Vet code
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1)
	golangci-lint run

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PATH)/main.go $(MAIN_PATH)/handlers.go $(MAIN_PATH)/templates.go

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development server with hot reload..."
	@which air > /dev/null || (echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; exit 1)
	air

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOCMD) install github.com/air-verse/air@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Initialize air configuration
init-air:
	@echo "Initializing air configuration..."
	@which air > /dev/null || (echo "Air not installed. Run 'make install-tools' first"; exit 1)
	air init

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --rm $(DOCKER_IMAGE)

# Docker compose up
docker-up:
	@echo "Starting with docker-compose..."
	docker-compose up --build

# Docker compose down
docker-down:
	@echo "Stopping docker-compose..."
	docker-compose down

# Generate project (for testing)
generate-test:
	@echo "Generating test project..."
	mkdir -p test-output
	curl -X POST "http://localhost:8080/generate" \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "projectName=test-project" \
		-d "goVersion=1.23" \
		-d "httpPackage=gin" \
		-d "databases=postgres" \
		-d "driver_postgres=gorm" \
		-d "features=config" \
		-d "features=gitignore" \
		--output "test-output/test-project.zip"
	@echo "Test project generated: test-output/test-project.zip"

# Start server in background for testing
start-bg:
	@echo "Starting server in background..."
	@$(GOCMD) run $(MAIN_PATH)/main.go $(MAIN_PATH)/handlers.go $(MAIN_PATH)/templates.go &
	@echo "Server started in background"

# Stop background server
stop-bg:
	@echo "Stopping background server..."
	@pkill -f "go run $(MAIN_PATH)/main.go" || true

# Full development setup
setup: install-tools init-air deps
	@echo "Development setup complete!"
	@echo "Run 'make dev' to start development server with hot reload"

# Check if server is running
health:
	@curl -s http://localhost:8080 > /dev/null && echo "✅ Server is running" || echo "❌ Server is not running"

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  coverage      - Run tests with coverage report"
	@echo "  bench         - Run benchmarks"
	@echo "  deps          - Download dependencies"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run golangci-lint"
	@echo "  run           - Run the application"
	@echo "  dev           - Start development server with hot reload"
	@echo "  install-tools - Install development tools (air, golangci-lint)"
	@echo "  init-air      - Initialize air configuration"
	@echo "  setup         - Full development setup"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-up     - Start with docker-compose"
	@echo "  docker-down   - Stop docker-compose"
	@echo "  generate-test - Generate test project (server must be running)"
	@echo "  start-bg      - Start server in background"
	@echo "  stop-bg       - Stop background server"
	@echo "  health        - Check if server is running"
	@echo "  help          - Show this help message"

# Quick development workflow
dev-full: clean deps fmt vet dev

# CI/CD workflow
ci: clean deps fmt vet lint test build
