# Todo

- [ ] **Make `ci` a required check on `main`** once the repo is on GitHub (Settings → Rules →
      require status checks: `ci`). It's the job that passes only when every other check did.
- [ ] **Install the Renovate app** after the first PR is merged into `main` and `ci` is required
      there: Renovate reads `.github/renovate.json` from `main`, so it starts right away, with no
      onboarding PR. Check its "Dependency Dashboard" issue lists `tools.mod`, the pinned actions
      and `ACTIONLINT_VERSION` from the Makefile.
- [ ] **Write the README** once the template is done, and merge this file into it. One README
      with everything: what the template is, requirements (Go, Node), quick start
      (`make setup`, `make dev`), the existing "Before deploying" notes, and this "Good to know"
      section. Drop the notes that repeat `CLAUDE.md` (migrations, `api.writeError`,
      `slog.*Context`, `make tidy`, Dockerfile `COPY` lines) or point to it, so each fact lives in
      one place.

## Good to know

- Run `make bindings` after changing the API. `schema.ts` does not update during `make dev`.
- Format with `make format` (gofmt and Prettier). Not `go fmt ./...`: it skips `web/embed.go`
  because of its build tag. Not `gofmt -w .` either: it would rewrite the `.go` files some npm
  packages ship in `web/node_modules`.
- golangci-lint (`.golangci.yml`) runs only in `make check-go`, not in VS Code: it's a Go tool in
  `tools.mod`, and the Go extension's `go.alternateTools` needs an executable, not a command. To
  lint on save, add a script that runs `go tool -modfile=tools.mod golangci-lint "$@"` and point
  `"go.lintTool": "golangci-lint-v2"` and `go.alternateTools` at it.
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
  Renovate bumps it, keep it on an LTS release. `@types/node` majors are off in Renovate (they'd
  describe a newer Node than the one we run): bump that major by hand together with Node's.
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
- The unused-messages check can't run if any source file in `web/` has a re-export
  (`export { x } from '…'` or `export * from '…'`) or dynamic code (`new Function`, `eval`), and
  `make check-i18n` then fails with "a check could not run". Import, then export instead
  (`import { x } from 'y'; export { x };`), and avoid dynamic code.
- i18n is Paraglide JS with inlang's default message format, in `web/messages/{locale}.json`.
  To add a language, list it in `locales` in `web/project.inlang/settings.json` and add its
  file; the language selector picks it up. The inlang plugins load from `node_modules`
  (`settings.json` `modules`); if a tool writes `https://cdn…` URLs back, revert them, or every
  build downloads them.
- Theme: daisyUI's `light`/`dark` themes, chosen by `data-theme` on `<html>`; without it, dark
  follows the system. `src/lib/theme.svelte.ts` saves the choice in `localStorage` under the
  `theme` key, and `public/theme.js` (a plain `<script>` in `index.html`, allowed by the CSP)
  applies it before the first paint. Keep the key in sync; a test checks both use it.
- `make build` fails if the web build prints any warning (it checks the output; Vite has no flag
  for it, and Tailwind prints CSS warnings with a bare `console.warn`). `npm run build` alone
  doesn't fail. daisyUI's `join` triggers two "Invalid empty selector" warnings, so the add form
  in `routes/Home.svelte` uses `flex gap-2`; try `join` again after daisyUI or Tailwind updates.
- No end-to-end tests (e.g. Playwright), on purpose: each piece is tested on its own (Go
  handlers and server, web components and state), but nothing tests them together in a browser.
  Check what only shows up there (the embedded SPA, the CSP, the SPA fallback) by hand with
  `make debug` or the built binary.
- CI (`.github/workflows/ci.yml`) runs each `make ci` target as its own job, so they show as
  separate PR checks and run in parallel, plus `docker`, which only CI runs (it needs Docker).
  Locally, `make ci` runs the same checks in one go. Require only the `ci` job on `main`. Test
  results show in each run's summary: Vitest's report by itself; for Go, a test results table
  with the output of each failure (`.github/scripts/go-tests.jq`, from gotestsum's JSON) and a
  coverage table per package (`.github/scripts/go-coverage.awk`, from go-test-coverage's
  breakdown, its own numbers, with the thresholds from `.testcoverage.yml`). On GitHub, files
  below the threshold are also annotated on the PR. To preview a summary locally:
  `sort .tmp/coverage-go.txt | awk -F';' -f .github/scripts/go-coverage.awk`.
- setup-go's cache key comes from `go.sum` and `tools.sum`, and a cache is saved only on a miss,
  by the first job that succeeds; later runs restore it and never save again. So only `ci.yml`
  caches Go: another workflow with the same files (`security.yml`, or a new one) has
  `cache: false`, or it could save a partial cache that every `ci` job then restores. If tool
  downloads show up in every job, delete the `setup-go-…` entry under Actions → Caches.
- In CI, actionlint also runs shellcheck on the workflows' `run:` scripts (GitHub's runners have
  it); locally only if shellcheck is installed, so CI's `check-workflows` can be stricter.
- actionlint isn't in `tools.mod` (its YAML dependency conflicts with another tool's): the
  Makefile runs it with `go run …@$(ACTIONLINT_VERSION)`, and Renovate updates that version
  through the `# renovate:` comment above it.
- About half of `make test-web` is jsdom setup, one per test file. Vitest suggests
  `pool: 'vmThreads'` or `isolate: false`; both are left off on purpose: `isolate: false` leaks
  state between files (locale, `localStorage`, URL), and `vmThreads` has known memory issues.
  Revisit if the suite gets slow.

## Stack

 - [x] go + chi v5
 - [x] godotenv + caarlos0/env
 - [x] svelte (no sveltekit)
 - [x] tool github.com/swaggo/swag/v2/cmd/swag
 - [x] tool github.com/go-task/task/v3/cmd/task@latest (nope - use makefile)
 - [x] tool github.com/air-verse/air
 - [x] sqlite
 - [x] air + vite/proxy for live-reload
 - [x] embed dist + serve
 - [x] paraglide / inlang for translations
 - [x] openapi-fetch

## Future ideas

 - [ ] devcontainer (`.devcontainer/devcontainer.json`: Go and Node, the VS Code extensions,
   `make setup` on create, ports 8080 and 5173)
 - [ ] go get -tool github.com/getsops/sops/v3/cmd/sops@latest
 - [ ] jwt chi + oidc for auth
