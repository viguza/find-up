.PHONY: build test clean deps lint vet fmt install

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Binary names
BINARY_NAME=find-up
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Test
test:
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Run benchmarks
bench:
	$(GOTEST) -bench=. -benchmem ./...

# Clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out

# Dependencies (no external dependencies)
deps:
	$(GOMOD) tidy
	$(GOMOD) verify

# Install dependencies (no external dependencies)
deps-install:
	$(GOMOD) tidy

# Lint
lint:
	$(GOLINT) run

# Vet
vet:
	$(GOCMD) vet ./...

# Format
fmt:
	$(GOFMT) -s -w .

# Install
install:
	$(GOCMD) install ./...

# Run examples
example:
	$(GOCMD) run examples/main.go

# Run all checks
check: fmt vet lint test

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-linux   - Build for Linux"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  bench         - Run benchmarks"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  deps-install  - Install dependencies"
	@echo "  lint          - Run linter"
	@echo "  vet           - Run go vet"
	@echo "  fmt           - Format code"
	@echo "  install       - Install the package"
	@echo "  example       - Run example"
	@echo "  check         - Run all checks (fmt, vet, lint, test)"
	@echo "  help          - Show this help"
