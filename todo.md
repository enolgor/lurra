# Todo

- [ ] **Add GitHub CI.** A workflow running `make setup ci`, a separate one for `make security`,
      and `ci` as a required check on `main`. Report tests, coverage, etc. etc. for PR merge checks.
- [ ] **Install the Renovate app** when publishing. The config is in `.github/renovate.json`.
      Check that the onboarding PR detects `tools.mod`. Add a rule so `@types/node` stays on the
      Node runtime's major (24, from `devEngines`) instead of jumping to the latest.
- [ ] devcontainers json?
- [ ] **Write the README** once the template is done, and merge this file into it. One README
      with everything: what the template is, requirements (Go, Node), quick start
      (`make setup`, `make dev`), the existing "Before deploying" notes, and this "Good to know"
      section. Drop the notes that repeat `CLAUDE.md` (migrations, `api.writeError`,
      `slog.*Context`, `make tidy`, Dockerfile `COPY` lines) or point to it, so each fact lives in
      one place.

## GitHub CI plan

1. `.github/workflows/ci.yml`: on PRs and pushes to `main`; Go from `go.mod`, Node from
   `web/package.json`; cache Go modules and npm; `make setup ci`.
2. `.github/workflows/security.yml`: `make security`, separate so a new advisory doesn't block
   unrelated PRs.
3. Results on the PR: Vitest annotations (automatic), Go test results and coverage in the job
   summary.
4. Docs; making `ci` a required check on `main` is a GitHub setting the user switches on.

Open questions:
1. Coverage as a job summary (no extra permissions) or a PR comment? *Recommended: job summary.*
2. Run `security` weekly, on demand, and on PRs that change `go.mod`, `go.sum` or
   `package-lock.json`? *Recommended: yes.*
3. Pin actions by commit SHA with Renovate keeping them updated, or by tag? *Recommended: SHA.*
4. Check whether `actions/setup-node` reads `devEngines` from `package.json`; if not, read the
   version like the Makefile's `NODE_VERSION` does.

## Good to know

- Run `make bindings` after changing the API. `schema.ts` does not update during `make dev`.
- Format by hand with `gofmt -w .`. `go fmt ./...` skips `web/embed.go` because of its build tag.
  Format the web code with `npm --prefix web run format`.
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

 - [ ] go get -tool github.com/getsops/sops/v3/cmd/sops@latest
 - [ ] jwt chi + oidc for auth
