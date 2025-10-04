# NASA HDF Data Processor Makefile

# Variables
BINARY_NAME=nasa-hdf-processor
BUILD_DIR=build
INPUT_DIR=input
OUTPUT_DIR=output

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty) -X main.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"

.PHONY: all build clean test deps run docker-build docker-run help

# Default target
all: deps build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies installed"

# Run the application (continuous mode)
run: build
	@echo "Running $(BINARY_NAME) in continuous mode..."
	./$(BINARY_NAME) -input-dir $(INPUT_DIR) -output-dir $(OUTPUT_DIR)

# Run the application once
run-once: build
	@echo "Running $(BINARY_NAME) once..."
	./$(BINARY_NAME) -once -input-dir $(INPUT_DIR) -output-dir $(OUTPUT_DIR)

# Create directories
setup-dirs:
	@echo "Creating directories..."
	mkdir -p $(INPUT_DIR) $(OUTPUT_DIR)
	@echo "Directories created: $(INPUT_DIR), $(OUTPUT_DIR)"

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .
	@echo "Docker image built: $(BINARY_NAME)"

# Docker run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run --rm -v $(PWD)/$(INPUT_DIR):/app/input:ro -v $(PWD)/$(OUTPUT_DIR):/app/output $(BINARY_NAME)

# Docker run with custom interval
docker-run-interval: docker-build
	@echo "Running Docker container with custom interval..."
	docker run --rm -v $(PWD)/$(INPUT_DIR):/app/input:ro -v $(PWD)/$(OUTPUT_DIR):/app/output -e PROCESS_INTERVAL=10m $(BINARY_NAME)

# Docker compose
docker-compose-up:
	@echo "Starting with Docker Compose..."
	docker-compose up --build

# Docker compose down
docker-compose-down:
	@echo "Stopping Docker Compose..."
	docker-compose down

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Check for updates
check-updates:
	@echo "Checking for module updates..."
	$(GOCMD) list -u -m all

# Show help
help:
	@echo "NASA HDF Data Processor - Available commands:"
	@echo ""
	@echo "  build              - Build the application"
	@echo "  clean              - Clean build artifacts"
	@echo "  test               - Run tests"
	@echo "  deps               - Install dependencies"
	@echo "  run                - Run in continuous mode"
	@echo "  run-once           - Run once and exit"
	@echo "  setup-dirs         - Create input/output directories"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code (requires golangci-lint)"
	@echo "  check-updates      - Check for dependency updates"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-run-interval - Run with custom interval"
	@echo "  docker-compose-up  - Start with Docker Compose"
	@echo "  docker-compose-down - Stop Docker Compose"
	@echo ""
	@echo "Examples:"
	@echo "  make run                    # Run continuously"
	@echo "  make run-once               # Run once"
	@echo "  make docker-run             # Run in Docker"
	@echo "  make docker-run-interval    # Run with 10m interval"
	@echo ""
