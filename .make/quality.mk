# .make/quality.mk - Extended quality targets
# Included by main Makefile. Baseline `fmt` `lint` `check` live in the root Makefile.

.PHONY: vet fmt-check lint-check security tidy

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

fmt-check: ## Check formatting without writing (for CI)
	@echo "Checking code format..."
	@test -z "$$(gofmt -l .)" || { echo "Code is not formatted. Run: make fmt"; exit 1; }
	@echo "Code is properly formatted"

lint-check: ## Run lint without fixing (for CI)
	@echo "Checking lint..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed" >&2; exit 1; }
	golangci-lint run ./...

security: ## Run security scan (gosec)
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec -exclude-generated ./...; \
	else \
		echo "gosec not installed. Install: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
	fi

tidy: ## Tidy go modules
	@echo "Tidying modules..."
	$(GOMOD) tidy
