# Go Project Makefile for go-ctl-initializer
BINARY_NAME=go-ctl
SERVER_BINARY_NAME=go-ctl-server
CLI_BINARY_NAME=go-ctl
SERVER_MAIN_PATH=cmd/server
CLI_MAIN_PATH=cmd/cli
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

.PHONY: all build build-server build-cli clean test coverage deps fmt vet lint run run-server run-cli dev docker-build docker-run help

# Default target
all: clean deps fmt vet test build-server build-cli

# Build both applications
build: build-server build-cli

# Build the web server
build-server:
	@echo "Building $(SERVER_BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY_NAME) $(SERVER_MAIN_PATH)/main.go $(SERVER_MAIN_PATH)/handlers.go $(SERVER_MAIN_PATH)/templates.go

# Build the CLI application
build-cli:
	@echo "Building $(CLI_BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(CLI_BINARY_NAME) $(CLI_MAIN_PATH)/main.go

# Build for multiple platforms
build-all: clean
	@echo "Building server for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY_NAME)-linux-amd64 $(SERVER_MAIN_PATH)/main.go $(SERVER_MAIN_PATH)/handlers.go $(SERVER_MAIN_PATH)/templates.go
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY_NAME)-darwin-amd64 $(SERVER_MAIN_PATH)/main.go $(SERVER_MAIN_PATH)/handlers.go $(SERVER_MAIN_PATH)/templates.go
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY_NAME)-darwin-arm64 $(SERVER_MAIN_PATH)/main.go $(SERVER_MAIN_PATH)/handlers.go $(SERVER_MAIN_PATH)/templates.go
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY_NAME)-windows-amd64.exe $(SERVER_MAIN_PATH)/main.go $(SERVER_MAIN_PATH)/handlers.go $(SERVER_MAIN_PATH)/templates.go
	@echo "Building CLI for multiple platforms..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(CLI_BINARY_NAME)-linux-amd64 $(CLI_MAIN_PATH)/main.go
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(CLI_BINARY_NAME)-darwin-amd64 $(CLI_MAIN_PATH)/main.go
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(CLI_BINARY_NAME)-darwin-arm64 $(CLI_MAIN_PATH)/main.go
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(CLI_BINARY_NAME)-windows-amd64.exe $(CLI_MAIN_PATH)/main.go

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

# Run the web server application
run: run-server

# Run the web server
run-server:
	@echo "Running $(SERVER_BINARY_NAME)..."
	$(GOCMD) run $(SERVER_MAIN_PATH)/main.go $(SERVER_MAIN_PATH)/handlers.go $(SERVER_MAIN_PATH)/templates.go

# Run the CLI application
run-cli:
	@echo "Running $(CLI_BINARY_NAME)..."
	$(GOCMD) run $(CLI_MAIN_PATH)/main.go

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
	@$(GOCMD) run $(SERVER_MAIN_PATH)/main.go $(SERVER_MAIN_PATH)/handlers.go $(SERVER_MAIN_PATH)/templates.go &
	@echo "Server started in background"

# Stop background server
stop-bg:
	@echo "Stopping background server..."
	@pkill -f "go run $(SERVER_MAIN_PATH)/main.go" || true

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
	@echo "  build         - Build both server and CLI applications"
	@echo "  build-server  - Build the web server"
	@echo "  build-cli     - Build the CLI application"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  coverage      - Run tests with coverage report"
	@echo "  bench         - Run benchmarks"
	@echo "  deps          - Download dependencies"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run golangci-lint"
	@echo "  run           - Run the web server (alias for run-server)"
	@echo "  run-server    - Run the web server"
	@echo "  run-cli       - Run the CLI application"
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
ci: clean deps fmt vet lint test build-server build-cli
