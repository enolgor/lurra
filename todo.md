# Todo

- [x] **Update to Go 1.27.2 or newer.** `make security` fails on standard library vulnerabilities
      until then.
- [x] **Update the README checkboxes.** Air, live reload and embedding `dist` are done.
- [ ] **Add GitHub CI.** A workflow running `make setup ci`, a separate one for `make security`,
      and `ci` as a required check on `main`.
- [ ] **Install the Renovate app** when publishing. The config is in `.github/renovate.json`.
      Check that the onboarding PR detects `tools.mod`.
- [ ] **Add a handler test** as a pattern to copy. `make test` has nothing to run yet.
- [ ] **Pin Node** with `devEngines.runtime` in `web/package.json`, using `"onFail": "error"`.
      Renovate then picks it up, so check its first Node PR stays on LTS.
- [ ] **Enable graceful shutdown in air** when SQLite arrives: `send_interrupt = true` and
      `kill_delay = "1s"` under `[build]` in `.air.toml`, so `make dev` exercises the shutdown path.

## Code review

- [ ] **Replace the global model with a `model.Store`.** Create it in `main` and pass it with
      `controller.New(store)`, so handler tests don't share state. Do it before the handler test.
      Its methods take `ctx context.Context` first, so their logs include the request ID.
- [x] **Add a sentinel `model.ErrNotFound`.** Controllers map every model error to 404; use
      `errors.Is` and return 500 for anything else.
- [ ] **Rewrite the swag-example doc comments** (`// Controller example`, `// TodosAll example`, ...).
- [x] **Return 413 for oversized bodies.** `RequestSize` makes `Decode` fail with
      `*http.MaxBytesError`, which currently returns 400.
- [ ] **Add an index.html fallback for the SPA** once the Svelte version has client-side routes.
- [x] **Reject unknown fields** when decoding JSON (`json.RejectUnknownMembers` in `httputil.Decode`).
- [x] **Add graceful shutdown** with `signal.NotifyContext` and `srv.Shutdown`, treating
      `http.ErrServerClosed` as a normal exit.
- [x] **Log with the request context in `httputil.JSON`** so its errors include the request ID.

## Good to know

- Run `make bindings` after changing the API. `schema.ts` does not update during `make dev`.
- Format by hand with `gofmt -w .`. `go fmt ./...` skips `web/embed.go` because of its build tag.
- Only `make build` embeds the UI. A plain `go build` reads `web/dist` from disk.

## Stack

 - [x] go + chi v5
 - [x] godotenv + caarlos0/env
 - [ ] svelte (no sveltekit)
 - [x] tool github.com/swaggo/swag/v2/cmd/swag
 - [x] tool github.com/go-task/task/v3/cmd/task@latest (nope - use makefile)
 - [x] tool github.com/air-verse/air
 - [ ] sqlite
 - [ ] go get -tool github.com/getsops/sops/v3/cmd/sops@latest
 - [x] air + vite/proxy for live-reload (vite proxy done, air pending)
 - [x] embed dist + serve (serving from disk, embed pending)
 - [ ] jwt chi + oidc for auth
 - [ ] tolgee for translations
 - [ ] openapi-fetch and https://github.com/tanstack/query (svelte) (maybe) (openapi-fetch done, query pending)
