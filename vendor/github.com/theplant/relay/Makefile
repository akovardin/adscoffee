.PHONY: test lint ci

# Default target
all: test lint

# Run all tests
test:
	go test -v ./...

# Run code linting
lint:
	golangci-lint run ./... --timeout=5m

# CI pipeline target
ci: test lint

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Clean build artifacts
clean:
	go clean
	rm -rf ./bin

# Build project
build:
	go build -o bin/ ./...

# Help information
help:
	@echo "Available targets:"
	@echo "  all    - Run tests and linters"
	@echo "  test   - Run tests"
	@echo "  lint   - Run linters"
	@echo "  ci     - Run CI pipeline locally"
	@echo "  fmt    - Format code"
	@echo "  clean  - Clean build artifacts"
	@echo "  build  - Build the project"
	@echo "  help   - Show this help message" 