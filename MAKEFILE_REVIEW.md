# Makefile review

## To do

- [ ] **CI never compiles the frontend.** Add a target that runs `npm --prefix web run build`
      and include it in `ci`.
- [ ] **`build` does not build the SPA.** It must depend on the web build once `dist` is embedded.
- [ ] **`dev` can leave a process running** if the other one crashes. It also uses `go run`, not
      air, so there is no Go live reload yet.

## Done

- [x] `build` and `dev` point at `cmd/app`.
- [x] `go:generate` removed, so `make bindings` is the only way to generate.
- [x] `tidy` tidies, and `tidy-check` verifies without changing files.
- [x] Added the `API_INFO_PATH`, `BIN_NAME` and `CMD_PATH` variables.
- [x] `-X main.version` now sets the `version` variable, logged at startup.

## Good to know

- `tools.mod` also lists app dependencies like chi, because tidying it scans the app code.
- `setup` compiles the tools, but the first `go tool` run still has to link them.
