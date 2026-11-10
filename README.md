# lurra

A full-stack project template: a Go API and a Svelte 5 single-page app, shipped together as one
binary. It comes with a small todo app that shows how the pieces fit, and with checks strict
enough that you can trust a green build.

**It's opinionated.** These are my preferred tools and conventions, wired together the way I like
to work: plain Go with few dependencies, a Svelte SPA without SvelteKit, SQLite, a Makefile, and
100% test coverage on both sides. If that matches how you like to build things, it should save
you a few days of setup. If it doesn't, it's probably not the template for you, and that's fine.

**Built with Claude, under my control.** I made this template with
[Claude Code](https://claude.com/claude-code), one small, reviewed step at a time: nothing went in
without my go-ahead, and nothing was done until `make ci` passed and I'd tried it myself. The
choices are mine; Claude did much of the typing and checking.

## What you get

- **A Go API** with chi, structured logs (slog) that carry a request ID, SQLite with migrations,
  security headers including a strict Content Security Policy, and graceful shutdown.
- **A typed API client for free.** Handlers are annotated for OpenAPI; `make bindings` turns
  that into TypeScript types, so the frontend breaks at compile time when the API changes.
- **A Svelte 5 app** styled with daisyUI and Tailwind, with client-side routing, live language
  switching (English, German and Spanish to start with) and a light, dark or system theme.
- **One binary to deploy.** The built SPA is embedded in the Go binary, with an SPA fallback and
  sensible cache headers. There's also a small distroless Docker image.
- **Live reload while you work:** `make dev` runs the API and the frontend together.
- **Strict checks:** 100% test coverage for Go and web, golangci-lint with a broad set of
  linters, strict TypeScript and typescript-eslint, translation checks that fail on anything
  missing or unused, and a build that fails on warnings.
- **CI on GitHub Actions** that runs every check as its own PR check, plus dependency security
  scans and Renovate for updates.

## The stack

| Area | Choice |
|---|---|
| Language and server | Go 1.27, [chi](https://go-chi.io) v5, [httplog](https://github.com/go-chi/httplog) + `log/slog` |
| Configuration | environment variables via [caarlos0/env](https://github.com/caarlos0/env), `.env` via [godotenv](https://github.com/joho/godotenv) |
| Database | SQLite ([ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3), no cgo), [goose](https://github.com/pressly/goose) migrations |
| API description | [swag](https://github.com/swaggo/swag) annotations → OpenAPI → [openapi-typescript](https://openapi-ts.dev) |
| Frontend | [Svelte 5](https://svelte.dev) (runes, no SvelteKit), [Vite](https://vite.dev), [sv-router](https://github.com/colinlienard/sv-router) |
| API client | [openapi-fetch](https://openapi-ts.dev/openapi-fetch/), typed from the generated schema |
| UI | [daisyUI](https://daisyui.com) 5 on [Tailwind CSS](https://tailwindcss.com) 4, [Lucide](https://lucide.dev) icons, self-hosted Inter |
| Translations | [Paraglide JS](https://github.com/opral/paraglide-js) and [inlang](https://inlang.com) |
| Tests | Go's `testing` with [gotestsum](https://github.com/gotestyourself/gotestsum) and [go-test-coverage](https://github.com/vladopajic/go-test-coverage); [Vitest](https://vitest.dev) with Testing Library |
| Linting and formatting | gofmt, [golangci-lint](https://golangci-lint.run), svelte-check, ESLint, Prettier |
| Dev tooling | Make, [air](https://github.com/air-verse/air) (Go live reload), [Delve](https://github.com/go-delve/delve) (debugger) |
| CI and updates | GitHub Actions, [actionlint](https://github.com/rhysd/actionlint), govulncheck, npm audit, [Renovate](https://docs.renovatebot.com) |
| Deployment | a single binary, or a [distroless](https://github.com/GoogleContainerTools/distroless) Docker image |

A few things are left out on purpose: no SvelteKit (the API is Go, so a plain SPA is enough), no
UI component library (daisyUI is just CSS classes), no task runner beyond Make, and no ORM
(plain `database/sql`).

## Getting started

You need **Go 1.27** and **Node 24** (an LTS release). Go fetches the exact toolchain from
`go.mod` by itself; the Node version is checked from `web/package.json`. Docker is only needed for
`make docker`.

1. Create your repository from the template (the **Use this template** button on GitHub), and
   clone it.
2. Install everything:

   ```sh
   make setup
   ```

   This downloads the Go and npm dependencies and builds the dev tools (golangci-lint, swag, air,
   …). They're pinned in `tools.mod`, so nothing is installed globally.
3. Optionally, copy the example settings: `cp .env.example .env`.
4. Start the app:

   ```sh
   make dev
   ```

   Open **http://localhost:5173**. Vite serves the frontend with hot reload and forwards `/api`
   to the Go server on port 8080, which air rebuilds whenever a `.go` or `.sql` file changes.

Run `make` (or `make help`) any time to see all targets.

## Everyday commands

| Command | What it does |
|---|---|
| `make dev` | API and frontend with live reload |
| `make format` | Format all Go and web code (gofmt, Prettier) |
| `make ci` | Everything the PR checks run: run it before you push |
| `make bindings` | Regenerate the TypeScript API types after changing the API |
| `make build` | The production binary, with the SPA embedded, in `bin/app` |
| `make docker` | The Docker image |
| `make release` | Binaries for Linux, macOS and Windows, as archives with checksums, in `bin/release/` |
| `make debug` | Run the API under Delve, with the built SPA; attach your editor |
| `make tidy` | Tidy `go.mod` and `tools.mod` after changing Go dependencies |
| `make security` | Known vulnerabilities in Go and npm dependencies |
| `make translate` | Machine-translate missing messages (see [Translations](#translations)) |

`make ci` is made of smaller targets you can also run on their own: `check-go`, `check-web`,
`check-i18n`, `check-workflows`, `test-go`, `test-web` and `build`. `make stress` runs the
concurrency tests repeatedly under the race detector, and `make clean` removes build output.

## Project layout

```
cmd/app/            main: wiring only (config, logging, store, server)
internal/
  api/              HTTP handlers and the OpenAPI annotations
  config/           settings from the environment and .env
  httpx/            JSON responses and errors
  logging/          the slog handler that adds the request ID
  model/            domain types and errors
  server/           router, middleware, security headers, SPA fallback
  store/sqlite/     the SQLite store and its migrations/
web/
  src/routes/       pages (one per route)
  src/lib/          shared code: API client, state, router, i18n, theme
  src/lib/components/  your own components
  src/lib/test/     test helpers and fixtures
  messages/         translations, one JSON file per language
  *.go              serves the built SPA (embedded or from disk)
.github/            CI workflows, their scripts and the Renovate config
tools.mod           the Go dev tools, pinned separately from the app's dependencies
```

## Configuration

Settings come from environment variables. For local work, put them in `.env` (see
`.env.example`); real environment variables always win over `.env`.

| Variable | Default | |
|---|---|---|
| `PORT` | `8080` | HTTP port |
| `LOG_LEVEL` | `INFO` | `DEBUG`, `INFO`, `WARN` or `ERROR` |
| `DB_PATH` | `data/app.db` | SQLite file, created on first start (`/data/app.db` in Docker) |

## How it works

### API and types

The API lives under `/api/v1` (the example has the usual todo endpoints), and `/ping` answers
health checks. Handlers carry swag annotations; `make bindings` generates the OpenAPI spec from
them and turns it into `web/src/lib/api/schema.ts`. The frontend uses one typed client built on
it, so a renamed field or a changed response is a compile error, not a runtime surprise.
`schema.ts` is generated: never edit it by hand, and run `make bindings` after changing the API
(it doesn't update during `make dev`). CI fails if it's out of date.

Errors are JSON with a status `code`, a `message` and the `request_id`. When one status can mean
different things (a 400 for an empty title and a 400 for a bad ID), the response also carries a
machine-readable `error` code, which the frontend translates. A 500 never shows the real error to
the client; find it in the logs by its `request_id`.

### Frontend

A plain Svelte 5 SPA: pages in `src/routes/`, routing with sv-router, state in small
`*.svelte.ts` modules. Styling is daisyUI classes and its theme colors, so dark mode works without
any `dark:` variants. The theme choice is saved in `localStorage` and applied by a tiny script
before the first paint, so there's no flash of the wrong theme.

The filter in the example (all, done, to do) lives in the URL (`/?filter=done`), so a reload or a
shared link keeps it.

### Translations

Every piece of text the user sees comes from `web/messages/{locale}.json` through Paraglide's
`m.key()` functions. Switching languages updates the page live, without a reload.

- **To add a language,** add it to `locales` in `web/project.inlang/settings.json` and create its
  messages file; the language menu picks it up by itself.
- **`make check-i18n`** is strict and part of `make ci`: a missing, empty or broken translation, or
  a message no code uses, fails it. To loosen it, change the `i18n:check` script in
  `web/package.json` (`--no-fail` only reports).
- **`make translate`** machine-translates new messages from English into the other languages,
  using your own provider (Google Translate, DeepL or OpenAI) through the inlang CLI's
  environment variables. Review the result before committing.

### Building and deploying

`make build` produces `bin/app`, with the SPA embedded: copy that one file to a server and run it.
A plain `go build` (without the `embed` build tag) serves `web/dist` from disk instead.

`make docker` builds a distroless image that runs as a non-root user, listens on 8080 and keeps
its database in `/data`. It's the `static` variant, without libc, since the binary is pure Go; if
you ever enable cgo, the Dockerfile says which image to switch to. Mount a volume at `/data` to
keep your data:

```sh
docker run -p 8080:8080 -v lurra-data:/data app:<version>
```

The server sends `index.html` for any path that isn't a file, so client-side routes survive a
reload. Hashed files in `/assets/` are cached forever; everything else is revalidated.

### Releases

Push a version tag, and `release.yml` does the rest:

```sh
git tag v1.2.3 && git push origin v1.2.3
```

- **A GitHub Release** with generated notes (one line per merged PR) and the binaries for Linux,
  macOS (amd64 and arm64 each) and Windows (amd64), as `.tar.gz` archives with the license and a
  `SHA256SUMS` file. `make release` builds the same files locally.
- **A Docker image** for amd64 and arm64 on GitHub's registry, as
  `ghcr.io/<owner>/<repo>:v1.2.3` and `:latest`. The Dockerfile cross-compiles, so building both
  needs no emulation.
- **Pre-releases:** a tag with a suffix, like `v1.0.0-rc.1`, becomes a pre-release and leaves
  `latest` alone.

The version comes from the tag: the binary reports it, and every build in between is named after
`git describe` (`v1.2.3-4-gabc1234`). No checks run on release, since a tag points at a commit
that already passed `ci` on main. Protect `v*` tags with a ruleset, so only you can create them.
The first image push creates a private package; make it public under the package's settings if you
want anyone to pull it.

## Quality checks and CI

Everything that guards the code runs through `make ci`, and CI runs exactly the same targets:

- **Go:** gofmt, tidy modules, golangci-lint (`.golangci.yml`), and the tests with the race
  detector and shuffled order. Coverage must be 100%; a branch that isn't worth testing can be
  marked `// coverage-ignore: <reason>`.
- **Web:** svelte-check with strict TypeScript, ESLint (`strictTypeChecked`, Svelte and Tailwind
  rules), Prettier, and Vitest at 100% coverage. A test fails if anything logs a warning or error.
- **Translations, workflows and the build:** `check-i18n`, actionlint, and a build that fails on
  any warning.

On GitHub, `ci.yml` runs each of those as its own check, in parallel, and also builds the Docker
image. A final `ci` job passes only when all of them did; make that the one required check on
your main branch. Each run's summary shows the Vitest report and Go test and coverage tables.
`security.yml` scans dependencies weekly and on dependency changes, without blocking PRs, and
`release.yml` publishes a release for every version tag ([Releases](#releases)).

Renovate (`.github/renovate.json`) keeps npm and Go dependencies, the dev tools and the pinned
GitHub Actions up to date; install the Renovate app once the config is on your main branch.

## Making it yours

The todo app is there to show the patterns; replace it with your own domain:

1. **Rename the module:** change `github.com/enolgor/lurra` in `go.mod` and in the imports
   (`go mod edit -module <your/module>`, then a search and replace).
2. **Replace the todo domain:** the types in `internal/model`, the store and a new migration in
   `internal/store/sqlite`, the handlers in `internal/api`, then `make bindings` and the pages in
   `web/src/routes`.
3. **Change the names and texts:** the app title and footer in `web/messages/`, and the binary and
   image name (`BIN_NAME`) in the Makefile.
4. **Choose your license,** and replace this README.

If you use Claude Code, `CLAUDE.md` holds the project's rules (layout, conventions, how to run an
audit), so it follows the same conventions you do.

## Good to know

### Backend

- **Migrations are append-only.** Never edit one that has run; add a new `NNNNN_name.sql`. To
  start over locally, delete `data/`.
- **New domain errors need a case in `api.writeError`,** or they come back as a generic 500.
- **Log with `slog.InfoContext(ctx, …)`,** not `slog.Info`, so your lines carry the request ID
  (golangci-lint enforces it).
- **Timeouts depend on each other:** `writeTimeout` must stay above `handlerTimeout`
  (`internal/server/server.go`), and air's `kill_delay` (`.air.toml`) above the app's 5 s
  `shutdownTimeout` (`cmd/app/main.go`). The request body limit is 1 MB.
- **New Go directories must be added to the Dockerfile,** which copies `cmd/`, `internal/` and
  `web/*.go` explicitly. CI's Docker build catches it if you forget.
- **Use `make tidy`, not just `go mod tidy`:** the dev tools have their own `tools.mod`, which
  shares some dependencies with the app.

### Frontend

- **The Content Security Policy only allows the app's own files.** Loading anything from another
  origin (a CDN, analytics, an external API or login provider) needs it added in
  `internal/server/headers.go`. `make dev` runs through Vite without the CSP, so check those
  problems with `make debug` or the built binary.
- **Client routes can't contain dots.** The SPA fallback treats `/users/john.doe` as a missing
  file and returns 404. Anything under `/api/` returns 404 too, never the SPA.
- **Translation keys must be called directly,** as `m.key()`, never `m[key]`, and not stored in a
  plain `const` (it wouldn't update when the language changes).
- **The unused-messages check needs plain imports.** A re-export (`export { x } from '…'`) or
  dynamic code (`eval`, `new Function`) makes it unable to run, and `make check-i18n` fails with
  "a check could not run". Import, then export instead.
- **The inlang plugins load from `node_modules`.** If a tool rewrites the `modules` in
  `web/project.inlang/settings.json` to `https://cdn…` URLs, revert it, or every build downloads
  them.
- **daisyUI's `join`** currently makes the build print CSS warnings (which fail `make build`), so
  the example's add form uses `flex gap-2`. Try it again after a daisyUI or Tailwind update.

### Tooling

- **Format with `make format`.** Not `go fmt ./...` (it skips `web/embed.go`, which has a build
  tag), and not `gofmt -w .` (it would rewrite Go files that some npm packages ship).
- **`make build` checks the web build's output for warnings,** because Vite has no flag for it.
  `npm run build` on its own doesn't fail on them.
- **Debugging:** `make debug` runs the API under Delve on port 2345, serving the built SPA with
  source maps on 8080. In VS Code, attach with **Attach to make debug** (F5). Ctrl+C stops it
  without a graceful shutdown.
- **golangci-lint only runs in `make check-go`,** not in the editor: it's pinned in `tools.mod`,
  and VS Code's Go extension needs an executable path. To lint on save, add a small script that
  runs `go tool -modfile=tools.mod golangci-lint "$@"`, and point `go.lintTool` and
  `go.alternateTools` at it.
- **Node is pinned** with `devEngines` in `web/package.json`, and `make setup` refuses other
  versions. When you move to a new Node major, update `@types/node` with it (Renovate leaves that
  major alone on purpose).
- **No end-to-end tests,** on purpose: each piece is tested on its own, but nothing drives a real
  browser against the binary. What only shows up there (the embedded SPA, the CSP, the SPA
  fallback) is worth a quick manual check.

### CI

- **Only one workflow caches Go.** setup-go keys its cache on `go.sum` and `tools.sum` and saves
  it once, from the first job that succeeds. A second workflow with the same key could save a
  partial cache that every `ci` job then restores, so `security.yml` runs without one. If tool
  downloads show up in every job, delete the `setup-go-…` cache under **Actions → Caches**.
- **actionlint runs shellcheck on the workflows in CI** (GitHub's runners have it), so CI can be
  stricter than your machine if you don't have shellcheck installed.
- **actionlint isn't in `tools.mod`** (one of its dependencies conflicts with another tool's), so
  the Makefile runs a pinned version with `go run`, and Renovate updates it through the comment
  above it.

## Before deploying

- **Log the client IP.** Add the chi `ClientIPFrom*` middleware that matches your deployment
  (see [Client IP](https://go-chi.io/#/pages/client_ip)), before the request logger:
  - exposed directly: `middleware.ClientIPFromRemoteAddr`
  - behind nginx, Cloudflare, etc.: `middleware.ClientIPFromHeader("X-Real-IP", <proxy CIDRs>)`
  - behind N load balancers: `middleware.ClientIPFromXFFTrustedProxies(n)`

  httplog's concise schema drops the IP, so also add a middleware after the logger that calls
  `httplog.SetAttrs(ctx, slog.String("client.ip", middleware.GetClientIP(ctx)))`.
- **Back up the database.** It's a single SQLite file (`DB_PATH`); copy it with the SQLite
  backup API or `sqlite3 app.db ".backup backup.db"`, not with `cp` while the app is running.

## Maybe later

Ideas that aren't in the template yet: authentication (JWT and OIDC with chi), secrets in the repo
with [sops](https://github.com/getsops/sops), and a devcontainer for VS Code and Codespaces.

## License

[MIT](LICENSE)
