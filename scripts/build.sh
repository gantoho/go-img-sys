#!/bin/bash

# Build and Run script for Linux/Mac

set -e

usage() {
    echo "Usage: ./build.sh [command]"
    echo ""
    echo "Commands:"
    echo "  build    - Build the application"
    echo "  run      - Build and run the application"
    echo "  clean    - Clean build artifacts"
    echo "  test     - Run tests"
    echo "  help     - Show this help message"
}

build() {
    echo "Building application..."
    go build -o image-sys main.go
    echo "Build successful: image-sys"
}

run() {
    build
    echo "Starting application on http://localhost:3128"
    ./image-sys
}

clean() {
    echo "Cleaning build artifacts..."
    rm -f image-sys
    go clean
    echo "Clean complete"
}

test() {
    echo "Running tests..."
    go test -v ./...
}

if [ $# -eq 0 ]; then
    usage
    exit 0
fi

case "$1" in
    build)
        build
        ;;
    run)
        run
        ;;
    clean)
        clean
        ;;
    test)
        test
        ;;
    help)
        usage
        ;;
    *)
        echo "Unknown command: $1"
        usage
        exit 1
        ;;
esac
