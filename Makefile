# gzh-cli-core Makefile
# Modular build automation for gzh-cli-* projects
#
# gzh-cli-core is a library: `build` compiles all packages, it produces no binary.

include .make/vars.mk
include .make/test.mk
include .make/quality.mk

.PHONY: help build test fmt lint check clean

CYAN  := \033[36m
RESET := \033[0m

help: ## Show this help
	@echo 'Usage: make [target]'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(CYAN)%-16s$(RESET) %s\n", $$1, $$2}' \
		$(MAKEFILE_LIST) | sort

build: ## Compile all packages
	@echo "Building packages..."
	$(GOBUILD) ./...

test: ## Run all tests with race detection
	@echo "Running tests..."
	$(GOTEST) $(RACE_FLAG) -timeout $(TEST_TIMEOUT) -coverprofile=$(COVERAGE_OUT) ./...

fmt: ## Format code with gofmt and gofumpt
	@echo "Formatting code..."
	$(GOFMT) ./...
	@if command -v gofumpt >/dev/null 2>&1; then \
		gofumpt -w .; \
	else \
		echo "gofumpt not installed, using go fmt only"; \
	fi

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed" >&2; exit 1; }
	golangci-lint run ./...

check: fmt lint test ## Run all quality checks (fmt + lint + test)
	@echo "All quality checks passed"

clean: ## Remove test and coverage artifacts
	@echo "Cleaning..."
	@rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)
	$(GO) clean ./...
	@echo "Cleaned"

.DEFAULT_GOAL := help
