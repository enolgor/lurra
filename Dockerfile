# check=skip=InvalidDefaultArgInFrom
# Build with `make docker`, which reads the versions from go.mod and web/package.json.
# Multi-platform (`--platform linux/amd64,linux/arm64`, as the release does): the build stages run
# on the build machine's platform and Go cross-compiles, so only the final stage is per platform,
# and nothing runs under emulation.
ARG GO_VERSION
ARG NODE_VERSION

FROM --platform=$BUILDPLATFORM node:${NODE_VERSION} AS web
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY web/ ./
RUN npm run build

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS app
ARG VERSION=dev
# set by Docker for each platform being built: the binary is compiled for it
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
# Copy only what the build needs, so edits to other files (Makefile, docs, this file) don't rebuild.
# New Go source dirs or files the build needs must be added here.
# Alternative: COPY . . and let .dockerignore exclude the rest, which needs no upkeep but rebuilds more.
COPY cmd cmd
COPY internal internal
COPY web/*.go web/
COPY --from=web /src/web/dist web/dist
# For cgo, set CGO_ENABLED=1, add a C cross-compiler for the target platforms, and switch the
# runtime image below (see its comment).
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -v -trimpath -tags embed -ldflags "-s -w -X main.version=${VERSION}" -o /app ./cmd/app
# Empty dir for the SQLite database; distroless has no shell to create it in the final stage.
RUN mkdir /data

# static: no libc, enough because the binary is pure Go (the SQLite driver needs no cgo). With
# cgo, use gcr.io/distroless/base-nossl-debian13:nonroot instead, which has glibc. Follows
# distroless's current Debian release; pin it with static-debian13 if needed.
FROM gcr.io/distroless/static:nonroot
COPY --from=app /app /app
COPY --from=app --chown=nonroot:nonroot --chmod=1777 /data /data
ENV DB_PATH=/data/app.db
VOLUME /data
EXPOSE 8080
ENTRYPOINT ["/app"]
