.DEFAULT_GOAL := help
MAKEFLAGS += --no-print-directory
.PHONY: help setup _fmt-check _lint test stress cover _vuln _audit security tidy _tidy-check _openapi bindings _bindings-check build docker dev _dev-api _dev-web clean check ci

TOOLS_MOD       := tools.mod
GO_TOOL         := go tool -modfile=$(TOOLS_MOD)
WEB_DIR         := web
API_SCHEMA_PATH := $(WEB_DIR)/src/schema.ts
API_INFO_PATH   := internal/api/api.go
TMP_DIR         := .tmp
COVER_PROFILE   := $(TMP_DIR)/cover.out
BIN_DIR         := bin
BIN_NAME        := app
CMD_PATH        := ./cmd/$(BIN_NAME)
IMAGE_NAME      := $(BIN_NAME)
VERSION         := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
GO_VERSION      := $(shell sed -n 's/^go //p' go.mod)
NODE_VERSION    := $(shell sed -n '/"runtime"/,/}/s/.*"version": "[^0-9]*\([0-9]*\).*/\1/p' $(WEB_DIR)/package.json)

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

setup: ## Download deps and compile dev tools
	@go mod download
	@go build -modfile=$(TOOLS_MOD) tool
	@npm --prefix $(WEB_DIR) ci --silent --no-audit --no-fund
	@echo "setup: ok"

_fmt-check: # Fail if Go code is not formatted
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
	@echo "fmt-check: ok"

_lint: # Run linters (includes go vet)
	@$(GO_TOOL) golangci-lint run
	@echo "lint: ok"

test: ## Run tests with race detector and write the coverage profile
	@mkdir -p $(TMP_DIR)
	@go test -race -shuffle=on -coverprofile=$(COVER_PROFILE) ./...

stress: ## Run concurrency tests repeatedly under the race detector
	@go test -race -count=20 -run Concurrent ./...

cover: ## Check coverage thresholds (.testcoverage.yml) on the profile written by make test
	@test -f $(COVER_PROFILE) || (echo "no coverage profile, run 'make test' first" && exit 1)
	@$(GO_TOOL) go-test-coverage -c .testcoverage.yml -p $(COVER_PROFILE)

_vuln: # Check Go dependencies for known vulnerabilities
	@$(GO_TOOL) govulncheck ./...
	@echo "vuln: ok"

_audit: # Check npm dependencies (fails on high/critical)
	@npm --prefix $(WEB_DIR) audit --audit-level=high
	@echo "audit: ok"

security: ## Check Go and npm dependencies
	@$(MAKE) -k _vuln _audit

tidy: ## Tidy go.mod/go.sum and tools.mod/tools.sum
	@go mod tidy
	@go mod tidy -modfile=$(TOOLS_MOD)
	@echo "tidy: ok"

_tidy-check: # Fail if go.mod/go.sum/tools.mod/tools.sum are not tidy
	@go mod tidy -diff
	@go mod tidy -diff -modfile=$(TOOLS_MOD)
	@echo "tidy-check: ok"

_openapi: # Generate the OpenAPI spec
	@$(GO_TOOL) swag init --v3.1 --outputTypes yaml -g $(API_INFO_PATH) -o $(TMP_DIR) --quiet

bindings: _openapi ## Generate API bindings
	@npx --prefix $(WEB_DIR) --no-install openapi-typescript $(TMP_DIR)/swagger.yaml -o $(API_SCHEMA_PATH)

_bindings-check: _openapi # Fail if bindings are out of date
	@npx --prefix $(WEB_DIR) --no-install openapi-typescript $(TMP_DIR)/swagger.yaml -o $(TMP_DIR)/schema.ts
	@diff -u $(API_SCHEMA_PATH) $(TMP_DIR)/schema.ts || (echo "bindings are stale, run 'make bindings'" && exit 1)
	@echo "bindings-check: ok"

build: ## Compile the binary with the SPA embedded
	@npm --prefix $(WEB_DIR) run build
	@CGO_ENABLED=0 go build -trimpath -tags embed -ldflags "-s -w -X main.version=$(VERSION)" -o $(BIN_DIR)/$(BIN_NAME) $(CMD_PATH)
	@echo "build: ok ($(VERSION))"

docker: ## Build the Docker image
	@docker build --build-arg GO_VERSION=$(GO_VERSION) --build-arg NODE_VERSION=$(NODE_VERSION) --build-arg VERSION=$(VERSION) -t $(IMAGE_NAME):$(VERSION) .

dev: ## Run air and vite with live reload
	@$(MAKE) -j2 _dev-api _dev-web

_dev-api:
	@$(GO_TOOL) air

_dev-web:
	@npm --prefix $(WEB_DIR) run dev

clean: ## Remove build artifacts
	@rm -rf $(BIN_DIR) $(TMP_DIR) $(WEB_DIR)/dist
	@echo "clean: ok"

check: ## Static checks (no tests)
	@$(MAKE) -k _fmt-check _tidy-check _bindings-check _lint

ci: check test build ## Deterministic CI checks (security runs separately)