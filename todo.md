# Todo

## Makefile

- [ ] **`bindings-check` should not change files.** It rewrites `schema.ts` before diffing.
      Generate to `.tmp` and compare instead, like `tidy-check`. This also catches a stale
      `schema.ts` that is already staged.
- [ ] **Split `fmt` and `fmt-check`.** `fmt` only checks today. Make it format, add `fmt-check`,
      and use that in `check`.
- [ ] **`clean` should remove `web/dist`.** A stale `dist` can keep being served from disk on `:8080`.

## Good to know

- CI needs `make setup` first, because `ci` assumes `web/node_modules` exists. Run `make setup ci`.
- `schema.ts` does not update during `make dev`. Run `make bindings` after changing the API.
- If air or vite crashes in `dev`, the other keeps running until you press Ctrl+C.
- Air runs the app in its own process group. If air is force-killed with `kill -9`, the app
  survives and keeps port 8080.
- `tools.mod` also lists app dependencies like chi, because tidying it scans the app code.
- `setup` compiles the tools, but the first `go tool` run still has to link them.
- `make lint` runs untagged, so it does not check `web/embed.go`.
- Only `make build` produces a self-contained binary. A plain `go build` reads `web/dist` from
  disk. The server logs which source it uses at startup, and warns if `web/dist` is missing.
- `-trimpath` and static-build flags can wait until binaries or Docker images are shipped.
