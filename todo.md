# Todo

- [ ] **Update to Go 1.27.2 or newer.** `make security` fails on standard library vulnerabilities
      until then.
- [ ] **Update the README checkboxes.** Air, live reload and embedding `dist` are done.
- [ ] **Add GitHub CI.** A workflow running `make setup ci`, a separate one for `make security`,
      and `ci` as a required check on `main`.
- [ ] **Set up Renovate** for dependency updates, covering Go modules, `tools.mod` and npm.
- [ ] **Add a handler test** as a pattern to copy. `make test` has nothing to run yet.
- [ ] **Pin Node** with a `.nvmrc`.

## Good to know

- Run `make bindings` after changing the API. `schema.ts` does not update during `make dev`.
- Format by hand with `gofmt -w .`. `go fmt ./...` skips `web/embed.go` because of its build tag.
- Only `make build` embeds the UI. A plain `go build` reads `web/dist` from disk.
- `make check` and `make ci` fail until the intentional lint error is removed.
