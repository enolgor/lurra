# Makefile review

## To do

- [ ] **`dev` can leave a process running** if the other one crashes. It also uses `go run`, not
      air, so there is no Go live reload yet.

## Done

- [x] `build` and `dev` point at `cmd/app`.
- [x] `go:generate` removed, so `make bindings` is the only way to generate.
- [x] `tidy` tidies, and `tidy-check` verifies without changing files.
- [x] Added the `API_INFO_PATH`, `BIN_NAME` and `CMD_PATH` variables.
- [x] `-X main.version` now sets the `version` variable, logged at startup.
- [x] `build` runs `web-build` and embeds `dist` with `-tags embed`.
- [x] `ci` runs `build`, which compiles the frontend and the embed file.

## Good to know

- `tools.mod` also lists app dependencies like chi, because tidying it scans the app code.
- `setup` compiles the tools, but the first `go tool` run still has to link them.
- `make lint` runs untagged, so it does not check `web/embed.go`.
- Only `make build` produces a self-contained binary. A plain `go build` reads `web/dist` from
  disk. The server logs which source it uses at startup, and warns if `web/dist` is missing.
