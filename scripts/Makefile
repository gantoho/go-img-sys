.PHONY: help build run dev clean test fmt lint docker-build docker-run deps info openapi generate swag-install build-linux build-windows build-darwin build-all

ROOT := $(realpath $(dir $(lastword $(MAKEFILE_LIST)))/..)
BUILD_DIR := $(ROOT)/build

# Default target
help:
	@echo "Available targets:"
	@echo "  make build         - Build the application (current OS)"
	@echo "  make run           - Run the application"
	@echo "  make dev           - Run in development mode"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make test          - Run tests"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make openapi       - Generate OpenAPI spec (requires swag CLI)"
	@echo "  make generate      - Run all code generators"
	@echo "  make build-linux   - Cross-compile for Linux (amd64)"
	@echo "  make build-windows - Cross-compile for Windows (amd64)"
	@echo "  make build-darwin  - Cross-compile for macOS (amd64)"
	@echo "  make build-all     - Cross-compile for all platforms"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"

# Build the application (with OpenAPI spec generation)
build: openapi
	@echo "Building application..."
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@cd /d "$(ROOT)" && go build -o "$(BUILD_DIR)/image-sys" ./cmd/image-sys/
	@echo "Build complete: $(BUILD_DIR)/image-sys"

# Run the application
run: build
	@echo "Starting application..."
	@"$(BUILD_DIR)/image-sys"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@if exist "$(BUILD_DIR)" rmdir /s /q "$(BUILD_DIR)"
	@cd /d "$(ROOT)" && go clean
	@echo "Clean complete: removed $(BUILD_DIR)"

# Run tests
test:
	@echo "Running tests..."
	@cd /d "$(ROOT)" && go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	@cd /d "$(ROOT)" && go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@cd /d "$(ROOT)" && golangci-lint run ./...

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@cd /d "$(ROOT)" && docker build -t image-sys:latest -f deployments/Dockerfile .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@cd /d "$(ROOT)" && docker run -p 3128:3128 -v "$(ROOT)/files:/root/files" -v "$(ROOT)/logs:/root/logs" image-sys:latest

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@cd /d "$(ROOT)" && go mod download && go mod tidy

# Show help
info:
	@echo "Project: Go Image System"
	@echo "Language: Go $(shell go version | cut -d' ' -f3)"
	@echo "Module: github.com/gantoho/go-img-sys"

# Install swag CLI for OpenAPI spec generation
swag-install:
	@echo "Installing swag CLI..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "swag CLI installed"

# Generate OpenAPI spec from Swagger annotations
openapi: swag-install
	@echo "Generating OpenAPI spec..."
	@cd /d "$(ROOT)" && swag init -g cmd/image-sys/main.go -o docs --parseDependency --parseInternal
	@echo "OpenAPI spec generated: $(ROOT)/docs/swagger.yaml, $(ROOT)/docs/swagger.json"

# Run all code generators
generate:
	@echo "Running go generate..."
	@cd /d "$(ROOT)" && go generate ./...
	@echo "Code generation complete"

# Cross-compile for Linux (amd64)
build-linux: openapi
	@echo "Cross-compiling for Linux (amd64)..."
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@cd /d "$(ROOT)" && set GOOS=linux&& set GOARCH=amd64&& go build -o "$(BUILD_DIR)/image-sys-linux" ./cmd/image-sys/
	@echo "Build complete: $(BUILD_DIR)/image-sys-linux"

# Cross-compile for Windows (amd64)
build-windows: openapi
	@echo "Cross-compiling for Windows (amd64)..."
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@cd /d "$(ROOT)" && set GOOS=windows&& set GOARCH=amd64&& go build -o "$(BUILD_DIR)/image-sys-windows.exe" ./cmd/image-sys/
	@echo "Build complete: $(BUILD_DIR)/image-sys-windows.exe"

# Cross-compile for macOS (amd64)
build-darwin: openapi
	@echo "Cross-compiling for macOS (amd64)..."
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@cd /d "$(ROOT)" && set GOOS=darwin&& set GOARCH=amd64&& go build -o "$(BUILD_DIR)/image-sys-darwin" ./cmd/image-sys/
	@echo "Build complete: $(BUILD_DIR)/image-sys-darwin"

# Cross-compile for all platforms
build-all: build-linux build-windows build-darwin
	@echo "All cross-compilation builds complete"
	@echo "  - Linux:   $(BUILD_DIR)/image-sys-linux"
	@echo "  - Windows: $(BUILD_DIR)/image-sys-windows.exe"
	@echo "  - macOS:   $(BUILD_DIR)/image-sys-darwin"
