# Targets (run `make help` for the same list from the ## comments below)
#
#   help             Show this help (the default: plain `make`)
#   setup            Download deps and compile dev tools (checks go and node first)
#   format           Format Go (gofmt) and web code (Prettier); the checks only report it
#   check-go         Go static checks: gofmt, go.mod tidiness, golangci-lint
#   test-go          Go tests with the race detector, then coverage thresholds
#   stress           Concurrency tests, repeatedly, under the race detector
#   check-web        Web static checks: API bindings up to date, svelte-check, Prettier, ESLint
#   check-i18n       Translation checks: missing, empty or broken translations, unused messages
#   translate        Machine-translate messages from en into the other locales (your own provider)
#   test-web         Web tests (Vitest), then coverage thresholds
#   bindings         Regenerate web/src/lib/api/schema.ts from the API annotations
#   tidy             Tidy go.mod/go.sum and tools.mod/tools.sum
#   security         Vulnerability checks for Go and npm dependencies
#   dev              Run air and vite with live reload
#   debug            Run the API under Delve with the built SPA; attach your editor
#   build            Compile the binary with the SPA embedded
#   release          Binaries for every release platform, as archives with checksums, in bin/release
#   docker           Build the Docker image
#   clean            Remove build artifacts
#   check-workflows  Lint the GitHub Actions workflows (actionlint)
#   ci               Everything the PR checks run: the checks above, the tests and the build
#
# Targets starting with _ are helpers: hidden from help, but callable (e.g. `make _lint`).

.DEFAULT_GOAL := help
MAKEFLAGS += --no-print-directory
.PHONY: help _check-env setup format \
	check-go _fmt-check _tidy-check _lint test-go stress \
	check-web _bindings-check _lint-web check-i18n translate _openapi bindings test-web \
	tidy security _vuln _audit \
	dev _dev-api _dev-web debug \
	build release _build-web docker clean check-workflows ci

# --- Variables ----------------------------------------------------------------

TOOLS_MOD       := tools.mod
GO_TOOL         := go tool -modfile=$(TOOLS_MOD)
WEB_DIR         := web
API_SCHEMA_PATH := $(WEB_DIR)/src/lib/api/schema.ts
API_INFO_PATH   := internal/api/api.go
TMP_DIR         := .tmp
COVER_PROFILE   := $(TMP_DIR)/cover.out
# Go test results (gotestsum's JSON) and coverage per file (go-test-coverage's breakdown,
# `file;statements;covered`); CI's job summary reads both (ci.yml)
GO_TEST_JSON    := $(TMP_DIR)/go-test.json
COVER_BREAKDOWN := $(TMP_DIR)/coverage-go.txt
BIN_DIR         := bin
BIN_NAME        := app
CMD_PATH        := ./cmd/$(BIN_NAME)
# The release binary: pure Go (no cgo), the SPA embedded, the version stamped into main.version
GO_BUILD        := CGO_ENABLED=0 go build -trimpath -tags embed -ldflags "-s -w -X main.version=$(VERSION)"
RELEASE_DIR     := $(BIN_DIR)/release
RELEASE_PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64
IMAGE_NAME      := $(BIN_NAME)
# Extra `docker build` flags, e.g. the release's --platform, more -t tags and --push (release.yml)
DOCKER_FLAGS    :=
VERSION         := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
GO_VERSION      := $(shell sed -n 's/^go //p' go.mod)
NODE_VERSION    := $(shell sed -n '/"runtime"/,/}/s/.*"version": "[^0-9]*\([0-9]*\).*/\1/p' $(WEB_DIR)/package.json)
# Our Go files: not the ones some npm packages ship in node_modules. `=`, so only targets that use
# it run the find; find (not git ls-files) also sees files that aren't committed yet.
GO_FILES         = $(shell find . \( -path ./web/node_modules -o -path ./.git \) -prune -o -name '*.go' -print)
# Not in tools.mod: its YAML dependency conflicts with another tool's. `go run pkg@version` builds
# it with its own go.mod; Renovate updates the version through the comment (renovate.json).
# renovate: datasource=go depName=github.com/rhysd/actionlint
ACTIONLINT_VERSION := v1.7.12

# Header for each stage `make ci` runs, so its output says which stage is running.
STAGE = @printf '\n\033[1m==> %s\033[0m\n' $@

# --- Help and setup -----------------------------------------------------------

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

_check-env: # Fail early if go is missing, or node is missing or has the wrong major version
	@command -v go >/dev/null || (echo "go not found: install Go (go.mod's toolchain is fetched automatically)" && exit 1)
	@node -v 2>/dev/null | grep -q '^v$(NODE_VERSION)\.' || (echo "need node $(NODE_VERSION).x (web/package.json)" && exit 1)

setup: _check-env ## Download deps and compile dev tools
	@go mod download
	@# fetch what `go mod tidy` loads (dependencies' test imports) so check-go doesn't download it
	@# every run, and CI's cache has it; -diff writes nothing, check-go does the real check
	@go mod tidy -diff >/dev/null 2>&1 || true
	@go mod tidy -diff -modfile=$(TOOLS_MOD) >/dev/null 2>&1 || true
	@go build -modfile=$(TOOLS_MOD) tool
	@npm --prefix $(WEB_DIR) ci --silent --no-audit --no-fund
	@echo "setup: ok"

format: ## Format Go (gofmt) and web code (Prettier); the checks only report it
	@gofmt -w $(GO_FILES)
	@npm --prefix $(WEB_DIR) run format --silent -- --log-level warn
	@echo "format: ok"

# --- Go -----------------------------------------------------------------------

check-go: ## Go static checks: gofmt, go.mod tidiness, golangci-lint
	$(STAGE)
	@$(MAKE) -k _fmt-check _tidy-check _lint

_fmt-check: # Fail if Go code is not formatted (fix with make format)
	@test -z "$$(gofmt -l $(GO_FILES))" || (gofmt -l $(GO_FILES) && exit 1)
	@echo "fmt-check: ok"

_tidy-check: # Fail if go.mod/go.sum/tools.mod/tools.sum are not tidy
	@go mod tidy -diff
	@go mod tidy -diff -modfile=$(TOOLS_MOD)
	@echo "tidy-check: ok"

_lint: # Run linters (includes go vet)
	@$(GO_TOOL) golangci-lint run
	@echo "lint: ok"

test-go: ## Go tests with the race detector, then coverage thresholds (.testcoverage.yml)
	$(STAGE)
	@mkdir -p $(TMP_DIR)
	@# gotestsum runs go test with the flags after --; GOTESTSUM_FORMAT picks the output (CI: github-actions)
	@$(GO_TOOL) gotestsum --jsonfile $(GO_TEST_JSON) -- -race -shuffle=on -coverprofile=$(COVER_PROFILE) ./...
	@# on GitHub, files and packages below the thresholds also become annotations on the PR
	@$(GO_TOOL) go-test-coverage -c .testcoverage.yml -p $(COVER_PROFILE) \
		--breakdown-file-name $(COVER_BREAKDOWN) $(if $(GITHUB_ACTIONS),--github-action-output)

stress: ## Run concurrency tests repeatedly under the race detector
	@go test -race -count=20 -run Concurrent ./...

# --- Web ----------------------------------------------------------------------

check-web: ## Web static checks: API bindings up to date, svelte-check, Prettier, ESLint
	$(STAGE)
	@$(MAKE) -k _bindings-check _lint-web

_lint-web: # Type-check (svelte-check), then formatting (Prettier) and lint (ESLint); fix formatting with make format
	@npm --prefix $(WEB_DIR) run check --silent
	@npm --prefix $(WEB_DIR) run lint --silent
	@echo "lint-web: ok"

check-i18n: ## Translation checks: missing, empty or broken translations, unused messages
	$(STAGE)
	@# inlang exits 0 when a check can't run (e.g. unused messages after a re-export); fail on that too
	@out=$$(npm --prefix $(WEB_DIR) run i18n:check --silent 2>&1); status=$$?; echo "$$out"; \
	[ $$status -eq 0 ] || exit $$status; \
	! echo "$$out" | grep -qE 'incomplete|not checked' || { echo "check-i18n: a check could not run, see above"; exit 1; }

translate: ## Machine-translate messages from en into the other locales (your own provider)
	@npm --prefix $(WEB_DIR) run i18n:translate --silent

_openapi: # Generate the OpenAPI spec
	@$(GO_TOOL) swag init --v3.1 --outputTypes yaml -g $(API_INFO_PATH) -o $(TMP_DIR) --quiet

bindings: _openapi ## Regenerate web/src/lib/api/schema.ts from the API annotations
	@npx --prefix $(WEB_DIR) --no-install openapi-typescript $(TMP_DIR)/swagger.yaml -o $(API_SCHEMA_PATH)

_bindings-check: _openapi # Fail if bindings are out of date
	@npx --prefix $(WEB_DIR) --no-install openapi-typescript $(TMP_DIR)/swagger.yaml -o $(TMP_DIR)/schema.ts
	@diff -u $(API_SCHEMA_PATH) $(TMP_DIR)/schema.ts || (echo "bindings are stale, run 'make bindings'" && exit 1)
	@echo "bindings-check: ok"

test-web: ## Web tests (Vitest), then coverage thresholds
	$(STAGE)
	@npm --prefix $(WEB_DIR) test --silent

# --- Dependencies -------------------------------------------------------------

tidy: ## Tidy go.mod/go.sum and tools.mod/tools.sum
	@go mod tidy
	@go mod tidy -modfile=$(TOOLS_MOD)
	@echo "tidy: ok"

security: ## Check Go and npm dependencies
	@$(MAKE) -k _vuln _audit

_vuln: # Check Go dependencies for known vulnerabilities
	@$(GO_TOOL) govulncheck ./...
	@echo "vuln: ok"

_audit: # Check npm dependencies (fails on high/critical)
	@npm --prefix $(WEB_DIR) audit --audit-level=high
	@echo "audit: ok"

# --- Run ----------------------------------------------------------------------

dev: ## Run air and vite with live reload
	@$(MAKE) -j2 _dev-api _dev-web

_dev-api:
	@$(GO_TOOL) air

_dev-web:
	@npm --prefix $(WEB_DIR) run dev

# Remove --continue to start paused until a debugger attaches (e.g. to debug startup code).
# stdin from /dev/null keeps Delve from handing the terminal to the app, so Ctrl+C stops Delve too.
debug: ## Build the SPA with source maps, then run the API under Delve (headless on :2345); attach your editor
	@npm --prefix $(WEB_DIR) run build -- --sourcemap
	@$(GO_TOOL) dlv debug $(CMD_PATH) --output $(BIN_DIR)/debug --headless --listen=127.0.0.1:2345 --api-version=2 --accept-multiclient --continue < /dev/null

# --- Build --------------------------------------------------------------------

build: ## Compile the binary with the SPA embedded
	$(STAGE)
	@$(MAKE) _build-web
	@$(GO_BUILD) -v -o $(BIN_DIR)/$(BIN_NAME) $(CMD_PATH)
	@echo "build: ok ($(VERSION))"

release: ## Binaries for every release platform, as archives with checksums, in bin/release
	$(STAGE)
	@$(MAKE) _build-web
	@rm -rf $(RELEASE_DIR) && mkdir -p $(RELEASE_DIR)
	@for platform in $(RELEASE_PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; ext=; [ $$os = windows ] && ext=.exe; \
		name=$(BIN_NAME)-$(VERSION)-$$os-$$arch; \
		mkdir -p $(RELEASE_DIR)/$$name && cp LICENSE README.md $(RELEASE_DIR)/$$name/ && \
		GOOS=$$os GOARCH=$$arch $(GO_BUILD) -o $(RELEASE_DIR)/$$name/$(BIN_NAME)$$ext $(CMD_PATH) && \
		tar -czf $(RELEASE_DIR)/$$name.tar.gz -C $(RELEASE_DIR) $$name && rm -r $(RELEASE_DIR)/$$name && \
		echo "  $$name.tar.gz" || exit 1; \
	done
	@# shasum (not sha256sum) is on macOS too; `sha256sum -c SHA256SUMS` checks the result
	@cd $(RELEASE_DIR) && shasum -a 256 *.tar.gz > SHA256SUMS
	@echo "release: ok ($(VERSION), $(RELEASE_DIR))"

_build-web: # Build the SPA (web/dist), failing on warnings
	@# vite has no fail-on-warnings flag, and Tailwind prints CSS warnings with a bare console.warn
	@out=$$(npm --prefix $(WEB_DIR) run build 2>&1); status=$$?; echo "$$out"; \
	[ $$status -eq 0 ] || exit $$status; \
	! echo "$$out" | grep -qi 'warning' || { echo "build: the web build printed warnings, see above"; exit 1; }

docker: ## Build the Docker image
	@docker build --build-arg GO_VERSION=$(GO_VERSION) --build-arg NODE_VERSION=$(NODE_VERSION) --build-arg VERSION=$(VERSION) -t $(IMAGE_NAME):$(VERSION) $(DOCKER_FLAGS) .

clean: ## Remove build artifacts
	@rm -rf $(BIN_DIR) $(TMP_DIR) $(WEB_DIR)/dist $(WEB_DIR)/src/lib/paraglide
	@echo "clean: ok"

# --- CI -----------------------------------------------------------------------

check-workflows: ## Lint the GitHub Actions workflows (actionlint)
	$(STAGE)
	@go run github.com/rhysd/actionlint/cmd/actionlint@$(ACTIONLINT_VERSION)
	@echo "check-workflows: ok"

# Keep in sync with the jobs in .github/workflows/ci.yml. Only CI also runs `docker` (it needs
# Docker and takes minutes); security runs separately.
ci: check-go check-web check-i18n check-workflows test-go test-web build ## Everything the PR checks run
