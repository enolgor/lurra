.DEFAULT_GOAL := help
.PHONY: help setup fmt lint test vuln audit security tidy bindings bindings-check build dev clean check ci

TOOLS_MOD       := tools.mod
GO_TOOL         := go tool -modfile=$(TOOLS_MOD)
WEB_DIR         := web
API_SCHEMA_PATH := $(WEB_DIR)/src/schema.ts
TMP_DIR         := .tmp
BIN_DIR         := bin
VERSION         := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

setup: ## Download deps and compile dev tools
	@go mod download
	@go build -modfile=$(TOOLS_MOD) tool
	@npm --prefix $(WEB_DIR) ci --silent --no-audit --no-fund
	@echo "setup: ok"

fmt: ## Check formatting
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
	@echo "fmt: ok"

lint: ## Run linters (includes go vet)
	@$(GO_TOOL) golangci-lint run
	@echo "lint: ok"

test: ## Run tests with race detector
	@go test -race -shuffle=on -cover ./...

vuln: ## Check Go dependencies for known vulnerabilities
	@$(GO_TOOL) govulncheck ./...
	@echo "vuln: ok"

audit: ## Check npm dependencies (fails on high/critical)
	@npm --prefix $(WEB_DIR) audit --audit-level=high
	@echo "audit: ok"

security: vuln audit ## Check Go and npm dependencies (use 'make -k security' to see both)

tidy: ## Fail if go.mod/go.sum/tools.mod/tools.sum are not tidy
	@go mod tidy
	@go mod tidy -modfile=$(TOOLS_MOD)
	@git diff --exit-code go.mod go.sum $(TOOLS_MOD) tools.sum
	@echo "tidy: ok"

bindings: ## Generate API bindings
	@trap 'rm -rf $(TMP_DIR)' EXIT; \
	$(GO_TOOL) swag init --v3.1 --outputTypes yaml -o $(TMP_DIR) && \
	npx --prefix $(WEB_DIR) --no-install openapi-typescript $(TMP_DIR)/swagger.yaml -o $(API_SCHEMA_PATH)

bindings-check: bindings ## Fail if bindings are out of date
	@git diff --exit-code -- $(API_SCHEMA_PATH) || (echo "bindings are stale, run 'make bindings'" && exit 1)
	@echo "bindings-check: ok"

build: ## Compile the binary
	@go build -ldflags "-X main.version=$(VERSION)" -o $(BIN_DIR)/app ./cmd/app
	@echo "build: ok ($(VERSION))"

dev: ## Run API and web in parallel (Ctrl+C stops both)
	@trap 'kill 0' INT TERM; go run ./cmd/app & npm --prefix $(WEB_DIR) run dev & wait

clean: ## Remove build artifacts
	@rm -rf $(BIN_DIR) $(TMP_DIR)
	@echo "clean: ok"

check: fmt lint test ## Quick local checks

ci: check tidy bindings-check ## Deterministic CI checks (security runs separately)