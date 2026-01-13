.PHONY: all build build-en build-ar test clean install help

# Binary names
BINARY_EN=daad
BINARY_AR=ض

# Build directory
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod

# Default target
all: test build

# Build both English and Arabic binaries
build: build-en build-ar

# Build English binary
build-en:
	@echo "Building English binary: $(BINARY_EN)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_EN) .
	@echo "✓ Built $(BUILD_DIR)/$(BINARY_EN)"

# Build Arabic binary
build-ar:
	@echo "Building Arabic binary: $(BINARY_AR)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_AR) .
	@echo "✓ Built $(BUILD_DIR)/$(BINARY_AR)"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "✓ Tests passed"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Run tests for specific package
test-lexer:
	@echo "Running lexer tests..."
	$(GOTEST) -v ./tests/lexer
	@echo "✓ Lexer tests passed"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_EN)
	rm -f $(BUILD_DIR)/$(BINARY_AR)
	rm -f coverage.out coverage.html
	@echo "✓ Cleaned"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✓ Dependencies installed"

# Install binaries to system (requires sudo on Linux/Mac)
install: build
	@echo "Installing binaries..."
	install -m 755 $(BUILD_DIR)/$(BINARY_EN) /usr/local/bin/$(BINARY_EN)
	install -m 755 $(BUILD_DIR)/$(BINARY_AR) /usr/local/bin/$(BINARY_AR)
	@echo "✓ Installed to /usr/local/bin"

# Uninstall binaries from system
uninstall:
	@echo "Uninstalling binaries..."
	rm -f /usr/local/bin/$(BINARY_EN)
	rm -f /usr/local/bin/$(BINARY_AR)
	@echo "✓ Uninstalled"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "✓ Code formatted"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run ./...
	@echo "✓ Linting complete"

# Run the English binary
run-en: build-en
	$(BUILD_DIR)/$(BINARY_EN)

# Run the Arabic binary
run-ar: build-ar
	$(BUILD_DIR)/$(BINARY_AR)

# Display help
help:
	@echo "Daad (ض) - Arabic Programming Language"
	@echo ""
	@echo "Available targets:"
	@echo "  make              - Run tests and build both binaries"
	@echo "  make build        - Build both English and Arabic binaries"
	@echo "  make build-en     - Build English binary (daad)"
	@echo "  make build-ar     - Build Arabic binary (ض)"
	@echo "  make test         - Run all tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-lexer   - Run only lexer tests"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make deps         - Install dependencies"
	@echo "  make install      - Install binaries to /usr/local/bin"
	@echo "  make uninstall    - Remove binaries from /usr/local/bin"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
	@echo "  make run-en       - Build and run English binary"
	@echo "  make run-ar       - Build and run Arabic binary"
	@echo "  make help         - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make test"
	@echo "  ./bin/daad tokenize tests/examples/basic_arithmetic.daad"
	@echo "  ./bin/ض رمز tests/examples/basic_arithmetic.daad"
