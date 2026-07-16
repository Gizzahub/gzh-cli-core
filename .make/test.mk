# .make/test.mk - Extended testing targets
# Included by main Makefile. Baseline `test` lives in the root Makefile.

.PHONY: test-unit test-verbose test-coverage bench

test-unit: ## Run unit tests only (skip integration)
	@echo "Running unit tests..."
	$(GOTEST) -short -timeout $(TEST_TIMEOUT) ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	$(GOTEST) -v $(RACE_FLAG) -timeout $(TEST_TIMEOUT) ./...

test-coverage: test ## Generate HTML coverage report
	@echo "Generating coverage report..."
	$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report: $(COVERAGE_HTML)"
	@$(GO) tool cover -func=$(COVERAGE_OUT) | tail -1

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem -run=^$$ ./...
