# lurra

Go 1.27 todo app template: chi + httplog/slog API, SQLite (ncruces driver, goose migrations),
swag → OpenAPI → `web/src/schema.ts`, Vite SPA embedded in the binary.

## Layout

- `cmd/app`: wiring only. `internal/config` env/.env, `logging` slog handler, `server` router and
  middleware, `api` handlers, `httpx` JSON/errors, `model` types, errors and `TodoStore`,
  `store/sqlite` store and `migrations/`.

## Workflow

- Answer questions and wait; edit files only after an explicit go-ahead.
- After every change run `make check-go test-go` (Go) or `make check-web test-web` (web); run
  `make ci` before finishing. Fix all failures.
- Never run `docker build`, containers or the app to verify; the user does that.
- Run `make bindings` after changing API annotations or request/response types.
  `web/src/api/schema.ts` is generated; never edit it by hand.
- Fix formatting with `gofmt -w .` (Go) and `npm --prefix web run format` (web); make only
  checks it.
- Add UI components with `npx shadcn-svelte@latest add <name>` in `web/`, then format: the
  generated files don't follow Prettier.
- Run `make tidy` (not just `go mod tidy`) after dependency changes.
- Prefer the simplest option; keep Makefile targets few (private helpers start with `_`).
- Keep this file current: when a new rule, convention or repeated correction appears, propose
  adding it here, in as few words as possible.

## Rules

- Coverage target is 100% (checked by `make test-go`). New Go code ships with tests; when a
  branch isn't worth a test (near-impossible or too costly), mark it `// coverage-ignore: <reason>`.
  Plain `if err != nil { return ..., err }` bodies are excluded automatically.
- Test behavior, not just lines: new features get tests for their intended behavior and edge
  cases, and bug fixes get a regression test, even when coverage is already 100%.
- New packages need tests or an exclude in `.testcoverage.yml`.
- Never edit an applied migration; add a new `NNNNN_name.sql` with Up and Down sections.
- New domain errors need a case in `api.writeError`, otherwise they become a 500.
- Log with `slog.*Context(ctx, ...)` so lines carry `request_id`.
- New Go source dirs must be added to the Dockerfile's explicit `COPY` lines.
