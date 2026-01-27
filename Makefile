.PHONY: build test lint install clean help build-linux build-linux-arm64 build-darwin build-darwin-arm64 build-windows build-all

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
	@golangci-lint run --timeout=5m

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

# Cross-compilation targets
DIST_DIR=dist

build-linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .

build-linux-arm64:
	@echo "Building for Linux (arm64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .

build-darwin:
	@echo "Building for macOS (amd64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .

build-darwin-arm64:
	@echo "Building for macOS (arm64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .

build-windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .

build-all: build-linux build-linux-arm64 build-darwin build-darwin-arm64 build-windows
	@echo "All platform binaries built successfully!"

# Display help
help:
	@echo "Available targets:"
	@echo "  build              - Build the binary"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  lint               - Run linting"
	@echo "  install            - Build and install the binary"
	@echo "  clean              - Remove build artifacts"
	@echo "  build-linux        - Build for Linux (amd64)"
	@echo "  build-linux-arm64  - Build for Linux (arm64)"
	@echo "  build-darwin       - Build for macOS (amd64)"
	@echo "  build-darwin-arm64 - Build for macOS (arm64)"
	@echo "  build-windows      - Build for Windows (amd64)"
	@echo "  build-all          - Build for all platforms"
	@echo "  help               - Display this help message"
