BINARY_NAME=trax
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
MODULE_PATH=github.com/Sekhudin/trax/cmd
LDFLAGS=-ldflags "-X $(MODULE_PATH).Version=$(VERSION)"

MIN_COV=90

.PHONY: all build test clean run install help

all: build test ## Build binary and run tests (default)

build: ## Build the Trax binary with version injection
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p bin
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) main.go
	@echo "✔ Build complete: bin/$(BINARY_NAME)"

test: ## Run all unit tests
	@echo "Running tests..."
	@go test -v ./...

cov: ## Run all unit tests with coverage
	@go test -cover ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out

covc: cov ## Check if coverage is above ($(MIN_COV)%)
	@echo "Checking coverage threshold (Min: $(MIN_COV)%)"
	@total_cov=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$total_cov%"; \
	if [ $$(echo "$$total_cov < $(MIN_COV)" | bc -l) -ne 0 ]; then \
		echo "✖ Coverage is below $(MIN_COV)%!"; \
		exit 1; \
	else \
		echo "✔ Coverage is above $(MIN_COV)%. Solid work!!"; \
	fi

clean: ## Remove binary and build artifacts
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f coverage.out
	@go clean
	@echo "✔ Cleaned"

run: build ## Build and execute the binary
	@./bin/$(BINARY_NAME)

install: ## Install binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME) to $(GOPATH)/bin..."
	@go install $(LDFLAGS) ./...

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

