# Later: stricter checks

Ideas to make the checks stricter, collected on 2026-11-10 while building the todo page. Done
already: svelte-check fails on warnings (`--fail-on-warnings` in the `check` script).

Already as strict as it gets: 100% coverage (Go and web, branches included), `check-i18n`
(including failing when a check can't run), Prettier and gofmt, the Tailwind lint rules as
errors, `go test -race`. `npm audit` stays at "high": "moderate" would fail on 12 known
advisories in dev tools that have no fix.

## Recommended (in this order: 1 and 3 first, then 4, then 2)

1. **ESLint warnings don't fail the lint.** Two rules in effect are only warnings:
   `svelte/no-at-debug-tags` (a leftover `{@debug}`) and `svelte/no-inspect` (a leftover
   `$inspect()`); both would pass CI. Fix: `eslint . --max-warnings 0` in the `lint` script.
   One flag; the lint passes today.
2. **Build warnings don't fail the build** (e.g. the CSS warnings from daisyUI's `join`). Vite
   has no fail-on-warnings flag, so `make build` would check its output for warnings, like
   `check-i18n` does for "incomplete".
3. **Runtime warnings during tests are ignored.** Svelte logs some problems only at runtime (e.g.
   mutating state it doesn't own), and libraries' `console.error`s go to the test output
   unnoticed. Fix: a few lines in `vitest-setup.ts` that fail a test on any `console.warn` or
   `console.error`.
4. **Go lint runs only golangci-lint's 5 defaults** (`errcheck`, `govet`, `ineffassign`,
   `staticcheck`, `unused`). A `.golangci.yml` could add e.g. `errorlint` (wrapping and comparing
   errors), `bodyclose` and `noctx` (HTTP mistakes), `gosec` (security), `revive` (doc comments,
   naming), `misspell`, `gocritic`. See what each flags in the current code before choosing.

## Possible, but costs code churn (measure first)

Turn each on temporarily, count what fails, and decide before changing anything.

5. **Stricter TypeScript** beyond `strict` (which `@tsconfig/svelte` sets):
   `noUncheckedIndexedAccess` (`array[i]` is "maybe undefined"), `noUnusedLocals` and
   `noUnusedParameters` (`tsconfig.node.json` already has them for `vite.config.ts`),
   `noImplicitReturns`, `noFallthroughCasesInSwitch`, `noImplicitOverride`. They catch real
   bugs; the first one adds checks in several places (`byStatus[...]`, `mock.calls[n]`).
6. **typescript-eslint `strictTypeChecked`** instead of `recommended`: many more type-aware rules
   (unnecessary conditions, unsafe `any`, non-null assertions). It would flag `result.data!` and
   `dropdown.querySelector(...)!`, so it means a cleanup pass.

## Open decision: daisyUI `join` and its build warnings

The add form (`routes/Home.svelte`) uses daisyUI's `join` for the input and the Add button. Its
CSS makes `vite build` print two "Invalid empty selector" warnings (a `>` at the start of a
selector inside `@scope`, which the CSS minifier doesn't understand). The CSS works.

- **(a)** Lay the form out with `flex gap-2` instead: slightly different look (a gap instead of
  joined), no warnings. *Recommended.*
- **(b)** Keep `join` and accept the warnings until daisyUI or the minifier fixes it. Not possible
  together with item 2 above.
