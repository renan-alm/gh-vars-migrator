.PHONY: build test lint install clean help

# Binary name
BINARY_NAME=gh-vars-migrator
BINARY_DIR=bin

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Run linting
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR)
	@rm -rf dist
	@rm -f coverage.out coverage.html
	@rm -f *.test
	@rm -f *.out

# Display help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linting"
	@echo "  install        - Build and install the binary"
	@echo "  clean          - Remove build artifacts"
	@echo "  help           - Display this help message"
