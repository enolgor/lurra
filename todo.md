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

## Good to know

- Run `make bindings` after changing the API. `schema.ts` does not update during `make dev`.
- Format by hand with `gofmt -w .`. `go fmt ./...` skips `web/embed.go` because of its build tag.
- Only `make build` embeds the UI. A plain `go build` reads `web/dist` from disk.
- `make check` and `make ci` fail until the intentional lint error is removed.
