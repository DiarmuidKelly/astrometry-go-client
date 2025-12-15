.PHONY: all build build-cli install test test-coverage test-integration test-integration-setup test-all lint lint-fix clean clean-indexes help

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
INDEX_DIR=astrometry-data

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

test-integration-setup: ## Download index files for integration tests
	@echo "Setting up integration test data..."
	@mkdir -p $(INDEX_DIR)
	@echo "Downloading astrometry index file for M42 test image (~6.6° FOV)..."
	@echo "  - index-4110.fits (3.00° - 4.20°, 24 MB)"
	@echo "Total download: 24 MB"
	@cd $(INDEX_DIR) && \
		if [ ! -f index-4110.fits ]; then \
			echo "Downloading index-4110.fits..." && \
			wget -q --show-progress http://data.astrometry.net/4100/index-4110.fits; \
		else \
			echo "index-4110.fits already exists"; \
		fi
	@echo "Index file ready:"
	@ls -lh $(INDEX_DIR)/index-4110.fits
	@echo "Converting test image to standard JPEG..."
	@if [ ! -f images/IMG_2820-converted.jpg ]; then \
		convert images/IMG_2820.JPG -quality 100 images/IMG_2820-converted.jpg && \
		echo "Created IMG_2820-converted.jpg (standard JPEG from MPO)"; \
	else \
		echo "IMG_2820-converted.jpg already exists"; \
	fi

test-integration: test-integration-setup ## Run integration tests (requires Docker and index files)
	@echo "Running integration tests..."
	ASTROMETRY_INDEX_PATH=$(PWD)/$(INDEX_DIR) $(GOTEST) -v -race -cover -tags=integration -timeout 10m $(PKG)

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

clean-indexes: ## Remove downloaded index files
	@echo "Removing index files..."
	rm -rf $(INDEX_DIR)
	@echo "Index files removed. Run 'make test-integration-setup' to re-download."

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
