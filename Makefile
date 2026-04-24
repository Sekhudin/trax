BINARY_NAME=trax
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
MODULE_PATH=github.com/Sekhudin/trax/cmd
LDFLAGS=-ldflags "-X $(MODULE_PATH).Version=$(VERSION)"

.PHONY: all build test clean run install help

all: build test ## Build binary and run tests (default)

build: ## Build the Trax binary with version injection
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) main.go

test: ## Run all unit tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Remove binary and build artifacts
	@echo "Cleaning up..."
	rm -rf bin/
	go clean

run: build ## Build and execute the binary
	./bin/$(BINARY_NAME)

install: ## Install binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME) to $(GOPATH)/bin..."
	go install $(LDFLAGS) ./...

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

