.PHONY: help build run dev clean test fmt lint docker-build docker-run deps info openapi generate swag-install

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
	@echo "  make openapi       - Generate OpenAPI spec (requires swag CLI)"
	@echo "  make generate      - Run all code generators (go generate)"
	@echo "  make swag-install  - Install swag CLI"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"

# Build the application (with OpenAPI spec generation)
build: openapi
	@echo "Building application..."
	@go build -o image-sys ./cmd/image-sys/
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

# Install swag CLI for OpenAPI spec generation
swag-install:
	@echo "Installing swag CLI..."
	@which swag > /dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	@echo "swag CLI installed"

# Generate OpenAPI spec from Swagger annotations
openapi: swag-install
	@echo "Generating OpenAPI spec..."
	@swag init -g cmd/image-sys/main.go -o docs --parseDependency --parseInternal
	@echo "OpenAPI spec generated: docs/swagger.yaml, docs/swagger.json"

# Run all code generators
generate:
	@echo "Running go generate..."
	@go generate ./...
	@echo "Code generation complete"
