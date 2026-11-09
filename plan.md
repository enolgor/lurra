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

## Open questions

1. **Scope.** Just the basics (list todos, add one, mark done/not done, delete), or more: edit
   the title, filter (all / open / done), a page per todo?
2. **Pages.** One page with everything on `/`, or the list on `/` plus a detail page like
   `/todos/1`? One page is enough for the basics; a detail page also shows off the router and the
   SPA fallback (and is where editing the title would live).
3. **Errors.** Show API errors as a daisyUI alert with a translated message and the `request_id`
   in small text, for support? *Recommended.*
4. **Error ids.** Translating by error *type* needs the backend to say which error it is, since
   `code` is only the HTTP status (two different 400s). Add a machine-readable field, for example
   `error: "title_invalid" | "not_found" | "bad_request" | "internal"`, to the error response, or
   translate by status only ("Invalid input", "Not found", "Something went wrong")? An `error`
   field is the cleaner template, at the cost of a small backend change (`httpx`, `writeError`,
   `make bindings`).
5. **Data loading.** Plain `await api.GET(...)` into `$state` in the page, with a small loading
   state, and reload the list after each change? Or update the list locally from each response
   (no extra request)? *Recommended: update locally from the responses; simple and no flicker.*
6. **Empty and loading states.** A daisyUI `loading` spinner while the list loads, and an
   "Nothing to do yet" message for an empty list? *Recommended.*

## Draft steps (to refine once the questions are answered)

1. **Backend error ids** (if question 4 says so): add the field, map it in `api.writeError`,
   tests, `make bindings`.
2. **Todo data module** (`$lib/todos.svelte.ts`): the list in `$state`, with `load`, `add`,
   `toggle`, `remove` (and `rename` if in scope) calling `api`, and error handling. Tests with a
   stubbed `fetch`, like `client.test.ts`.
3. **Error display**: a small component that turns an API error into a translated daisyUI alert
   with the `request_id`. Messages in en/de/es. Tests.
4. **The page(s)**: replace `Home.svelte`'s demo with the list (add form, items with a done
   checkbox and delete button, empty and loading states); add the detail page and its route if
   in scope. All text through `m.key()`, daisyUI classes only. Tests.
5. **Remove the demo leftovers**: the example messages, unused router exports; replace the
   catch-all with a not-found page if the app has more than one page.
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
