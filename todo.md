# Todo

- [ ] **Add GitHub CI.** A workflow running `make setup ci`, a separate one for `make security`,
      and `ci` as a required check on `main`. Report tests, coverage, etc. etc. for PR merge checks.
- [ ] **Install the Renovate app** when publishing. The config is in `.github/renovate.json`.
      Check that the onboarding PR detects `tools.mod`.
- [x] **Merge `make test` and `make cover`** before coverage goes into CI, so tests run once: let
      `test` write the profile (`-coverpkg=./internal/... -coverprofile=.tmp/cover.out`, works with
      `-race`) and drop `cover` or make it only print the report from that profile.
- [x] **Pin Node** with `devEngines.runtime` in `web/package.json`, using `"onFail": "error"`.
      Renovate then picks it up, so check its first Node PR stays on LTS.
- [x] **Add an index.html fallback for the SPA** once the Svelte version has client-side routes.
      Test that an unknown route (`/todos/1`) returns `index.html` with 200, a real asset returns
      the file, a missing asset (`/assets/nope.js`) returns 404, and an unknown API route
      (`/api/v1/nope`) returns 404, not the SPA.
- [x] Configure using dlv easily
- [ ] devcontainers json?
- [x] **Configure ESLint** (with the Svelte plugin) for `web/` and add it to `_lint-web` in the
      Makefile, so `check-web` and `ci` run it.
- [ ] **Add an ICU message format i18n library** for Svelte (plurals, gender, etc.). Check how it
      fits with Tolgee from the stack list.
- [x] **Add a client-side routing library** for Svelte. It needs the index.html fallback above.
- [x] **Add shadcn-svelte** as the UI component library (with Tailwind, which it requires).
- [ ] **Add dark mode switching** with `mode-watcher` (light/dark/system). The theme in
      `src/app.css` already has dark colors, enabled by a `.dark` class that nothing sets yet.
- [x] **Add frontend tests** with Vitest and `@testing-library/svelte`, replacing the `test-web`
      placeholder in the Makefile.
- [x] **Add Prettier** (with the Svelte plugin) for `web/` formatting, checked in `_lint-web`.
- [x] **Decide on frontend coverage:** Vitest thresholds in `test-web`, or Go-only. State it in
      `CLAUDE.md` either way. Exclude `src/lib/components/ui/` (copied shadcn-svelte code).
- [x] **Remove the create-vite demo content:** `src/lib/Counter.svelte`, `src/assets/`, the demo
      `App.svelte`, `web/README.md` and the placeholder icons in `public/`.
- [ ] **Describe `web/` in `CLAUDE.md`'s Layout:** Svelte 5 + Vite, API client in `src/api/`, and
      `make check-web test-web` after web changes.
- [x] **Add cache headers for the SPA:** `Cache-Control: public, max-age=31536000, immutable` for
      hashed `/assets/*`, `no-cache` for `index.html`. Embedded files have no modtime, so
      `http.FileServer` sends no `Last-Modified` and browsers can't revalidate them.
- [x] **Add security headers** (chi middleware): at least `X-Content-Type-Options: nosniff` and a
      Content Security Policy for the SPA.
- [x] **Add a shared API client** (`src/api/client.ts`, `openapi-fetch` with `baseUrl: "/api/v1"`),
      showing errors with their `request_id`.


## Good to know

- Run `make bindings` after changing the API. `schema.ts` does not update during `make dev`.
- Format by hand with `gofmt -w .`. `go fmt ./...` skips `web/embed.go` because of its build tag.
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
- shadcn-svelte copies components into `web/src/lib/components/ui/`; the behavior comes from
  `bits-ui` (an npm dependency). Run these in `web/` with `npx shadcn-svelte@latest`, then
  `npm run format` and `git diff` to review and restore your own edits:
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
