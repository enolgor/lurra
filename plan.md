# Plan: build the todo app (connect frontend and backend)

The template's backend has a full todo API, and the frontend has a typed client for it, but no
UI uses it yet: `/` shows demo content. This plan builds the real app on top of what exists.
GitHub CI comes after it (see the end), so CI checks the real app rather than the demo.

## Why the app first

- It shows what the template really needs: how API errors are shown, loading states, maybe a
  first test that goes from the UI to the API. CI should check that, not today's demo page.
- CI only wraps `make ci`, which already exists and passes, so waiting costs little.
- It replaces the demo leftovers from the last audit: the example messages in `Home.svelte`,
  the unused router exports (`navigate`, `isActive`, `preload`) and the unused API client.

## What exists

**Backend** (`internal/api`, under `/api/v1`):

| Method | Path | Body | Returns |
|---|---|---|---|
| GET | `/todos` | | `Todo[]` |
| POST | `/todos` | `{ title }` | `Todo` (201) |
| GET | `/todos/{id}` | | `Todo` |
| PATCH | `/todos/{id}` | `{ title?, done? }` | `Todo` |
| DELETE | `/todos/{id}` | | 204 |

- `Todo` is `{ id, title, done }`.
- Errors are `{ code, message, request_id }`. `code` is the HTTP status: 404 for an unknown todo,
  400 for an empty title *and* for an invalid id. The `message` is English server text.

**Frontend:**
- `$lib/api/client.ts` exports `api`, the typed openapi-fetch client (`api.GET('/todos')`, …).
- The page layout in `App.svelte`: a skip link, a `<header>` (daisyUI `navbar`: a home icon
  link, the theme and language menus), a `<main id="main">` where the router renders each page
  (each page renders its own `<h1>`), and a `<footer>` saying it's an example app built with the
  lurra template (linked to github.com/enolgor/lurra). The tab title follows the language
  (`app_title`, "Todos").
- sv-router with `/` and a catch-all, both to `routes/Home.svelte` (heading "Todos"; still the
  demo messages).
- daisyUI, Paraglide (en/de/es, live switching), the shared `Dropdown`, the language and theme
  menus in the header.

## Decided

- **Filter: all / done / undone,** in the frontend only (the API keeps returning every todo).
  The filter is part of the URL, so a reload or a shared link keeps it, as a query
  parameter: `/?filter=done` (`all` is the default and has no parameter), read and written with
  sv-router's reactive search params (`route.search`). Not the URL hash: `#…` already means
  "jump to this element" (the skip link uses `#main`), and the query keeps the path free for
  routes. The SPA fallback already serves `index.html` for `/?filter=…`, so reloads work.
- **One page,** `/`, with everything on it. The router still shows its value: the filter in
  `?filter=` (`route.search`), the typed `p('/')` link, and a **not-found page** for any other
  path (the catch-all route): "Page not found" with a link home, instead of showing the list on
  `/todos/1` or a typo.
- **Inline editing** of a todo's title, in the list (no detail page).
- **Errors are translated.** By HTTP status where that's enough (404 "not found", 500
  "something went wrong"); where one status covers different errors (400: an empty title and an
  invalid id), the backend adds a machine-readable error code to the response, for example
  `error: "title_invalid"`, and the frontend translates that code instead. The frontend maps
  codes and statuses to messages with a local map of message functions (never `m[key]`).
- **The list is updated locally** from each response (add, toggle, rename, delete), without
  reloading it.
- **A daisyUI `loading` spinner** while the list loads.
- **An empty list shows "Nothing to do yet"** in a different, quieter style than the list
  (for example centered, muted text), so it doesn't look like a todo.

- **Error display**, by where the error happens; always the translated message with the
  `request_id` in small muted text (for bug reports), and `role="alert"` so screen readers
  announce it:
  - **The list fails to load:** an inline daisyUI `alert alert-error` in place of the list, with
    a "Try again" button that reloads it.
  - **An action fails** (add, toggle, rename, delete): a daisyUI `toast` with an
    `alert alert-error` and a close button. It stays until closed (so the `request_id` can be
    copied), a new error replaces the previous one, and the list stays usable; the failed change
    is simply not applied, since the list only changes from responses.
- **No empty titles are sent:** the Add button is disabled while the field is blank, and
  blanking a title while editing cancels the edit. The server's `title_invalid` stays as a
  backup, shown in the toast like any other error.

## Steps

1. **Backend error codes**: an optional `error` field in the error response (`httpx`), set by
   `api.writeError` for errors whose status isn't specific enough (`title_invalid` for now),
   with tests and `make bindings`.
2. **Todo data module** (`$lib/todos.svelte.ts`): the list in `$state`, with `load`, `add`,
   `toggle`, `rename` and `remove` calling `api` and updating the list locally from each
   response, a loading flag, and the last error; plus the filtered list as a `$derived` of the
   list and the filter. Tests with a stubbed `fetch`, like `client.test.ts`.
3. **Error display**: a small component that turns an API error into a translated message (by
   `error` code, else by status) in a daisyUI alert with the `request_id`, used both inline
   (load failure, with "Try again") and in a toast (action failures, with a close button).
   Messages in en/de/es. Tests.
4. **The page**: replace `Home.svelte`'s demo with the list: an add form (Add disabled while
   blank), the all / done / undone filter, items with a done checkbox, inline title editing
   (blank cancels) and a delete button, the loading spinner, "Nothing to do yet" for an empty
   list, the load-error alert and the action-error toast. All text through `m.key()`, daisyUI
   classes only. Tests, including that the filter is read from and written to `?filter=` and
   survives a reload.
5. **Not-found page and demo leftovers**: a `routes/NotFound.svelte` ("Page not found", a link
   home) for the catch-all route, with tests; remove the example messages and the unused router
   exports.
6. **Docs**: CLAUDE.md layout, todo.md.
7. **Final check**: `make ci`, then the user tries it in `make dev` (Go API + Vite together).

Work it the usual way: one step at a time, each one reviewed, with a crossed-off checklist.

## After the app: GitHub CI

Drafted and postponed (also saved in Claude's memory as `github-ci-plan`):

1. `.github/workflows/ci.yml`: on PRs and pushes to `main`; Go from `go.mod`, Node from
   `web/package.json`; cache Go modules and npm; `make setup ci`.
2. `.github/workflows/security.yml`: `make security`, separate so a new advisory doesn't block
   unrelated PRs.
3. Results on the PR: Vitest annotations (automatic), Go test results and coverage in the job
   summary.
4. Docs; making `ci` a required check on `main` is a GitHub setting the user switches on.

Open questions for CI:
1. Coverage as a job summary (no extra permissions) or a PR comment? *Recommended: job summary.*
2. Run `security` weekly, on demand, and on PRs that change `go.mod`, `go.sum` or
   `package-lock.json`? *Recommended: yes.*
3. Pin actions by commit SHA with Renovate keeping them updated, or by tag? *Recommended: SHA.*
4. Check whether `actions/setup-node` reads `devEngines` from `package.json`; if not, read the
   version like the Makefile's `NODE_VERSION` does.
