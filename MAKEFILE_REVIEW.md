# Makefile review

Some of these are expected while the repo and Makefile are still being tied together.
Checked items are done and were run, except `make dev`, which was only dry-run. The fixes
suggested for open items have not been tried.

## Broken now

- [x] **`build` and `dev` pointed at a missing path.** `main.go` now lives in `cmd/app`, which
      is what both targets expect.
- [ ] **`make check` and `make ci` fail on lint.** `errcheck` flags the ignored error from
      `json.NewEncoder(w).Encode(v)` at `httputil/error.go:12`. `make test` never runs because
      lint stops first.
- [ ] **`-X main.version` does nothing.** `cmd/app/main.go` has no `version` variable, so the
      flag is silently ignored and `build: ok (...)` overstates what happened. Add
      `var version = "dev"`.

## Gaps

- [ ] **CI never compiles the frontend.** `bindings-check` proves `schema.ts` is current, not
      that `web/src/main.ts` still works with it, so an API change can break the client while CI
      stays green. Add a target that runs `npm --prefix web run build` and include it in `ci`.
- [ ] **`build` does not build the SPA.** This works today because `dist` is read from disk, but
      the binary is only usable from the repo root. Once `dist` is embedded (planned in the
      README), `build` must depend on the web build or a fresh checkout will not compile.
- [x] **`go:generate` was a second, outdated path.** Both directives are gone, so
      `make bindings` is the only way to generate. It reads the general API info from
      `controller/doc.go` via `swag init -g`.

## Smaller things

- [x] **`tidy` changed files while acting as a check.** Now split: `tidy` tidies both modfiles,
      and `tidy-check` uses `go mod tidy -diff`, which verifies without modifying. `ci` runs
      `tidy-check`.
- [x] **Repeated and hardcoded paths.** Added `API_INFO_PATH`, `BIN_NAME` and `CMD_PATH`.
- [ ] **`dev` can leave a process running.** If one process crashes, the other keeps going. It
      also uses `go run`, while the README plans on air, so there is no Go live reload yet.
- **`tools.mod` picks up app dependencies.** Tidying it also scans the app code, which is why chi
  is in there. Harmless, but it shares the module path, so it is not fully isolated.
- **`setup` does not warm `go tool`.** `go build -modfile=tools.mod tool` works as a compile
  check, but `go tool` caches binaries separately, so its first run still links.

## Looks right

- `npx --prefix web --no-install` works.
- `web/src/schema.ts` is tracked, so `bindings-check` is meaningful.
- `--audit-level=high` will not fail on the moderate `minimatch` issues.
