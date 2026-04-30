BINARY_NAME=trax
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
MODULE_PATH=github.com/Sekhudin/trax/cmd
LDFLAGS=-ldflags "-X $(MODULE_PATH).Version=$(VERSION)"

PACKAGES := $(shell go list ./... | grep -v /internal/testutil)
PACKAGES_DEV := $(shell go list ./... | grep -v /internal/testutil)
MIN_COV=100

.PHONY: all build test clean run install help

all: build covc ## Build binary and run covc (default)

build: ## Build the Trax binary with version injection
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p bin
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) main.go
	@echo -e "\033[32m✔\033[0m Build complete: \033[32mbin/$(BINARY_NAME)\033[0m"

cov: ## Run all unit tests with coverage
	@go test -cover $(PACKAGES) -coverprofile=coverage.out
	@go tool cover -func=coverage.out

covd: ## Run all unit tests with coverage (development)
	@go test -cover $(PACKAGES_DEV) -coverprofile=coverage.out
	@go tool cover -func=coverage.out

covc: cov ## Check if coverage is above $(MIN_COV)%
	@echo ""
	@echo -e "\033[34mChecking coverage threshold\033[0m (\033[1mMin: $(MIN_COV)%\033[0m)"
	@echo ""
	@total_cov=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total_cov < $(MIN_COV)" | bc -l) -ne 0 ]; then \
		echo -e "\033[31m✖\033[0m Target coverage not reached. \033[1mKeep grinding\033[0m!"; \
		echo -e "  ↳ Coverage: \033[31m$$total_cov%\033[0m"; \
		exit 1; \
	else \
		echo -e "\033[32m✔\033[0m Minimum threshold satisfied. \033[1mKeep it up\033[0m!"; \
		echo -e "  ↳ Coverage: \033[32m$$total_cov%\033[0m"; \
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
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; { \
		help_msg = $$2; \
		gsub(/\$$\(MIN_COV\)/, "$(MIN_COV)", help_msg); \
		printf "\033[36m%-15s\033[0m %s\n", $$1, help_msg \
	}'

