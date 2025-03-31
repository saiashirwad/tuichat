.PHONY: build run clean test deps

# Binary name and path
BINARY_NAME=gochat
BINARY_PATH=bin/$(BINARY_NAME)

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Main package path
MAIN_PACKAGE=cmd/gochat/main.go

# Build variables
BUILD_TIME=$(shell date +%FT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION?=dev

# Environment variables
export GOCHAT_LLM_API_KEY?=$(GROQ_API_KEY)

# Default target
all: clean build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(GOBIN)
	@go build -o $(BINARY_PATH) $(MAIN_PACKAGE)
	@echo "Build complete"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BINARY_PATH)

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@if ! command -v air > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@air -c .air.toml

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(GOBIN)
	@go clean
	@echo "Clean complete"

# Download and tidy dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  make          - Clean and build the application"
	@echo "  make build    - Build the application"
	@echo "  make run      - Build and run the application"
	@echo "  make dev      - Run with hot reload (requires air)"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make deps     - Download and tidy dependencies"
	@echo "  make test     - Run tests"
	@echo "  make help     - Show this help message" 