.PHONY: help build run clean test docker-build docker-run

# Default target
help:
	@echo "Available targets:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make dev           - Run in development mode"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make test          - Run tests"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"

# Build the application
build:
	@echo "Building application..."
	@go build -o image-sys main.go
	@echo "Build complete: image-sys"

# Run the application
run: build
	@echo "Starting application..."
	@./image-sys

# Run in development mode with hot reload (requires air)
dev:
	@echo "Starting application in development mode..."
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	@air

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f image-sys image-sys.exe
	@go clean
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run ./...

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t image-sys:latest .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 3128:3128 -v $(PWD)/files:/root/files -v $(PWD)/logs:/root/logs image-sys:latest

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Show help
info:
	@echo "Project: Go Image System"
	@echo "Language: Go $(shell go version | cut -d' ' -f3)"
	@echo "Module: github.com/gantoho/go-img-sys"
