# Todo

- [ ] **Add GitHub CI.** A workflow running `make setup ci`, a separate one for `make security`,
      and `ci` as a required check on `main`. Report tests, coverage, etc. etc. for PR merge checks.
- [ ] **Install the Renovate app** when publishing. The config is in `.github/renovate.json`.
      Check that the onboarding PR detects `tools.mod`.
- [ ] devcontainers json?
- [ ] **Add a theme switcher** (light/dark/system) next to the language selector, with daisyUI's
      `theme-controller`. Dark mode already follows the system (daisyUI's `dark` theme is
      `--prefersdark`); the switcher only adds a manual choice, remembered for the next visit.

## Good to know

- Run `make bindings` after changing the API. `schema.ts` does not update during `make dev`.
- Format by hand with `gofmt -w .`. `go fmt ./...` skips `web/embed.go` because of its build tag.
  Format the web code with `npm --prefix web run format`.
- Only `make build` embeds the UI. A plain `go build` reads `web/dist` from disk.
- Copy `.env.example` to `.env` for local settings. Real environment variables take precedence.
- Never edit an applied migration; add a new one. Goose skips versions it has already applied.
  To start over locally, delete `data/`.
- New domain errors need a case in `api.writeError`, otherwise they are returned as a generic 500.
- 500 responses hide the real error; find it in the logs by the `request_id` from the response.
  Log with `slog.InfoContext(ctx, ...)` (not `slog.Info`) so your own lines carry it too.
- Run `make tidy`, not just `go mod tidy`, after changing dependencies; `tools.sum` shares some.
- The Dockerfile copies `cmd/`, `internal/` and `web/*.go` explicitly; add new Go dirs there.
- In Docker the database is `/data/app.db`; mount a volume at `/data` to keep it.
- `kill_delay` in `.air.toml` must stay above `shutdownTimeout` in `main.go`.
- `writeTimeout` must exceed `handlerTimeout` in `server.go`.
- `make debug` serves the built SPA (with source maps) at `:8080` under Delve; attach with
  "Attach to make debug" (F5) in VS Code. Ctrl+C kills the app without a graceful shutdown.
- The Content Security Policy (`internal/server/headers.go`) only allows the app's own
  files. Loading anything from another origin (CDN, analytics, an external API or OIDC
  provider) needs it added there. `make dev` serves through Vite without the CSP, so check
  CSP problems with `make debug` or the built binary.
- Node is pinned with `devEngines.runtime` in `web/package.json` (`"onFail": "error"`). When
  Renovate bumps it, keep it on an LTS release.
- The SPA fallback (`internal/server/spa.go`) serves `index.html` for any path without a file
  extension, so client routes must not contain dots (`/users/john.doe` gets a 404). Missing
  files and anything under `/api/` return 404, never the SPA.
- Hashed `/assets/*` are cached forever; `index.html` and other files get `no-cache`. Embedded
  files have no modtime, so there is no `Last-Modified` and `index.html` is downloaded again on
  every load (it's tiny).
- `make check-i18n` runs `inlang check` strictly, as part of `make ci`: a missing, empty or
  broken translation, or a message no code uses, fails it. To loosen it, change the `i18n:check`
  script in `web/package.json`: `--no-fail` reports without failing, listed check flags (such as
  `--missing-variables`) run only those, and `--locales de` limits the locales.
- `make translate` machine-translates the messages in `web/messages/en.json` into the other
  locales and writes their files (`de.json`, …); review the diff and run `make check-i18n`. It
  uses your own provider (Google Translate, DeepL or OpenAI), set through the inlang CLI's
  environment variables. Not part of any check or build.
- The unused-messages check gives up if any source file in `web/` has a re-export
  (`export { x } from '…'` or `export * from '…'`). Import, then export instead
  (`import { x } from 'y'; export { x };`).
- i18n is Paraglide JS with inlang's default message format, in `web/messages/{locale}.json`.
  To add a language, list it in `locales` in `web/project.inlang/settings.json` and add its
  file; the language selector picks it up. The inlang plugins load from `node_modules`
  (`settings.json` `modules`); if a tool writes `https://cdn…` URLs back, revert them, or every
  build downloads them.

## Stack

 - [x] go + chi v5
 - [x] godotenv + caarlos0/env
 - [x] svelte (no sveltekit)
 - [x] tool github.com/swaggo/swag/v2/cmd/swag
 - [x] tool github.com/go-task/task/v3/cmd/task@latest (nope - use makefile)
 - [x] tool github.com/air-verse/air
 - [x] sqlite
 - [ ] go get -tool github.com/getsops/sops/v3/cmd/sops@latest
 - [x] air + vite/proxy for live-reload (vite proxy done, air pending)
 - [x] embed dist + serve (serving from disk, embed pending)
 - [ ] jwt chi + oidc for auth
 - [ ] paraglide / inlang for translations
 - [x] openapi-fetch
