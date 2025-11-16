.PHONY: help build test lint fmt clean install all

# Variables
BINARY_NAME=my-context
BUILD_DIR=bin
CMD_DIR=cmd/my-context
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Default target
help: ## Display this help message
	@echo "My-Context Copilot - Makefile Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

all: clean lint test build ## Run all checks and build

build: ## Build binary for current platform
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)/
	@echo "✓ Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build binaries for all platforms
	@echo "Building for all platforms..."
	@./scripts/build-all.sh
	@echo "✓ All binaries built in $(BUILD_DIR)/"

test: ## Run all tests
	@echo "Running tests..."
	@go test ./... -v -race -coverprofile=coverage.out
	@echo "✓ Tests complete"

test-short: ## Run tests without race detector
	@echo "Running tests (short)..."
	@go test ./... -v
	@echo "✓ Tests complete"

test-integration: ## Run only integration tests
	@echo "Running integration tests..."
	@go test ./tests/integration/... -v
	@echo "✓ Integration tests complete"

test-unit: ## Run only unit tests
	@echo "Running unit tests..."
	@go test ./tests/unit/... -v
	@echo "✓ Unit tests complete"

coverage: test ## Generate coverage report
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

benchmark: ## Run benchmark tests
	@echo "Running benchmarks..."
	@go test ./tests/benchmarks/... -bench=. -benchmem
	@echo "✓ Benchmarks complete"

lint: ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --timeout=5m; \
		echo "✓ Linting complete"; \
	else \
		echo "⚠ golangci-lint not installed. Install: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -w -s .
	@goimports -w . 2>/dev/null || echo "⚠ goimports not installed (optional)"
	@echo "✓ Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet complete"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

install: build ## Install binary to local bin directory
	@echo "Installing $(BINARY_NAME)..."
	@./scripts/install.sh
	@echo "✓ Installation complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify
	@echo "✓ Dependencies downloaded"

tidy: ## Tidy go.mod and go.sum
	@echo "Tidying modules..."
	@go mod tidy
	@echo "✓ Modules tidied"

check: lint vet test ## Run all checks (lint, vet, test)
	@echo "✓ All checks passed"

ci: deps check build ## CI pipeline (deps, checks, build)
	@echo "✓ CI pipeline complete"

release-check: ## Verify release readiness
	@echo "Checking release readiness..."
	@./scripts/pre-release-check.sh || echo "⚠ Create scripts/pre-release-check.sh for release verification"
	@echo "✓ Release check complete"

# Development helpers
dev-setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	@echo "Installing development tools..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "✓ Development environment ready"

# Watch for changes and run tests (requires entr or similar)
watch: ## Watch for changes and run tests
	@echo "Watching for changes..."
	@find . -name "*.go" | entr -c make test-short

.DEFAULT_GOAL := help
