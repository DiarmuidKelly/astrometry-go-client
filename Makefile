.PHONY: all build build-cli install test test-coverage test-integration test-all lint lint-fix clean help

# Variables
BINARY_NAME=astro-cli
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=gofmt
GOLINT=golangci-lint
VERSION=$(shell cat VERSION)
BUILD_DIR=bin
PKG=./pkg/...
CMD=./cmd/...

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

all: test build ## Run tests and build

build: build-cli ## Build all binaries

build-cli: ## Build the CLI tool
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/astro-cli

install: ## Install CLI to GOPATH
	@echo "Installing $(BINARY_NAME) to GOPATH..."
	$(GO) install $(LDFLAGS) ./cmd/astro-cli

test: ## Run unit tests
	@echo "Running unit tests..."
	$(GOTEST) -v -race -short $(PKG)

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic $(PKG)
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-integration: ## Run integration tests (requires Docker)
	@echo "Running integration tests..."
	$(GOTEST) -v -race -tags=integration $(PKG)

test-all: test test-integration ## Run all tests

lint: ## Run linter
	@echo "Running linter..."
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run $(PKG) $(CMD); \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

lint-fix: ## Run linter with auto-fix
	@echo "Running linter with auto-fix..."
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run --fix $(PKG) $(CMD); \
	else \
		echo "golangci-lint not installed."; \
		exit 1; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) $(PKG) $(CMD)

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f $(BINARY_NAME)
	$(GO) clean

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	$(GO) mod tidy

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
