# Todo

- [ ] **Add GitHub CI.** A workflow running `make setup ci`, a separate one for `make security`,
      and `ci` as a required check on `main`. Report tests, coverage, etc. etc. for PR merge checks.
- [ ] **Install the Renovate app** when publishing. The config is in `.github/renovate.json`.
      Check that the onboarding PR detects `tools.mod`.
- [ ] devcontainers json?
- [ ] **Add i18n with Paraglide JS** (`@inlang/paraglide-js`, Vite plugin) and the ICU plugin
      (`@inlang/plugin-icu1`), so messages are standard ICU MessageFormat JSON (plurals, select)
      that any translation tool can sync. Not tied to one vendor.
      - `web/messages/en.json` is the master, edited by hand; Paraglide types `m.key()` from it.
      - Other languages sit next to it, written by `make pull-translations` and committed.
      - Exclude Paraglide's generated output from lint, Prettier and coverage.
- [ ] **Add `make push-translations` and `make pull-translations`.** Push uploads `en.json`, pull
      writes every other language (never `en.json`). Both run commands from Makefile variables
      (`TRANSLATIONS_PUSH`, `TRANSLATIONS_PULL`), Tolgee CLI by default; swap them for another
      tool. `make build` does not pull: builds need no network or credentials.
- [ ] **Add a translations workflow** that runs `make pull-translations` and opens a PR with the
      changes, reusing the open one (`peter-evans/create-pull-request`, fixed branch such as
      `translations/update`; no PR when nothing changed). Trigger it on `schedule` (daily) and
      `workflow_dispatch`, which works with any translation tool. Optionally add
      `repository_dispatch` for tools that can call it: most translation-tool webhooks can't send
      GitHub's required body and auth header without a relay. Needs the tool's API key as a repo
      secret, and the PR must trigger CI (a PAT or GitHub App token, not `GITHUB_TOKEN`).
- [ ] **Add dark mode switching** with `mode-watcher` (light/dark/system). The theme in
      `src/app.css` already has dark colors, enabled by a `.dark` class that nothing sets yet.

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
- shadcn-svelte copies components into `web/src/lib/components/ui/`; interactive ones get their
  behavior from `bits-ui`, an npm dependency that `add` installs with the first that needs it.
  Run these in `web/` with `npx shadcn-svelte@latest`, then `npm run format` and `git diff` to
  review and restore your own edits:
  - `add <name...>` adds components (`-a` for all of them).
  - `add <name> --overwrite` updates one to the latest upstream version.
  - `apply <preset>` switches the preset (style, theme, font) and reinstalls the components;
    `--only theme` or `--only font` changes just that, leaving the components alone.
  - `init --reinstall` reinstalls all components after a style change.

  CLI docs: https://www.shadcn-svelte.com/docs/cli
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
 - [ ] tolgee for translations
 - [x] openapi-fetch
