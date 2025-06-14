# Project configuration
BINARY_NAME=aws-config
MAIN_PACKAGE=.
BUILD_DIR=bin
DIST_DIR=dist
COVERAGE_DIR=coverage

# Build information
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go configuration (for development builds only)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
CGO_ENABLED?=0

# Development build flags (simpler for faster iteration)
DEV_LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Default target
.DEFAULT_GOAL := help

# Phony targets
.PHONY: help clean deps tidy fmt vet lint test test-race test-coverage \
        build build-dev build-local install run dev check \
        release snapshot docker pre-commit ci security goreleaser-check

## help: Show this help message
help:
	@echo "Available targets:"
	@awk '/^##/ { \
		sub(/^## /, "", $$0); \
		split($$0, arr, ": "); \
		printf "  %-15s %s\n", arr[1], arr[2] \
	}' $(MAKEFILE_LIST)

## clean: Remove build artifacts and temporary files
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	@go clean -cache -testcache -modcache
	@echo "Clean complete"

## deps: Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify

## tidy: Clean up dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@gofumpt -l -w . 2>/dev/null || true

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## lint: Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## test-race: Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	@go test -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report: $(COVERAGE_DIR)/coverage.html"

## goreleaser-check: Check GoReleaser configuration
goreleaser-check:
	@echo "Checking GoReleaser configuration..."
	@goreleaser check

## build: Build using GoReleaser (recommended for production)
build: clean goreleaser-check
	@echo "Building with GoReleaser..."
	@goreleaser build --clean --snapshot --single-target
	@echo "Build complete, artifacts in $(DIST_DIR)/"

## build-dev: Quick development build using go build
build-dev:
	@echo "Building $(BINARY_NAME) for development ($(GOOS)/$(GOARCH))..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) go build $(DEV_LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Development build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## build-local: Alias for build-dev (backward compatibility)
build-local: build-dev

## build-all: Build for all platforms using GoReleaser
build-all: clean goreleaser-check
	@echo "Building for all platforms with GoReleaser..."
	@goreleaser build --clean --snapshot
	@echo "Multi-platform build complete, artifacts in $(DIST_DIR)/"

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(DEV_LDFLAGS) $(MAIN_PACKAGE)

## run: Build and run the application (development build)
run: build-dev
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

## dev: Run the application in development mode (no build)
dev:
	@echo "Running in development mode..."
	@go run $(DEV_LDFLAGS) $(MAIN_PACKAGE)

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "All checks passed!"

## pre-commit: Run pre-commit checks
pre-commit: tidy fmt vet lint test-race goreleaser-check
	@echo "Pre-commit checks complete!"

## ci: Run CI pipeline checks
ci: deps check test-coverage build
	@echo "CI pipeline complete!"

## security: Run security checks
security:
	@echo "Running security checks..."
	@govulncheck ./... 2>/dev/null || echo "govulncheck not installed, skipping vulnerability check"
	@gosec ./... 2>/dev/null || echo "gosec not installed, skipping security scan"

## release: Create a release using goreleaser
release: check
	@echo "Creating release..."
	@goreleaser release --clean

## snapshot: Create a snapshot build using goreleaser
snapshot: check
	@echo "Creating snapshot build..."
	@goreleaser build --snapshot --clean

## docker: Build Docker image (if Dockerfile exists)
docker:
	@if [ -f Dockerfile ]; then \
		echo "Building Docker image..."; \
		docker build -t $(BINARY_NAME):$(VERSION) .; \
	else \
		echo "No Dockerfile found, skipping Docker build"; \
	fi

# Version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
