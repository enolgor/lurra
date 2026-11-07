# check=skip=InvalidDefaultArgInFrom
# Build with `make docker`, which reads the versions from go.mod and web/package.json.
ARG GO_VERSION
ARG NODE_VERSION

FROM node:${NODE_VERSION} AS web
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY web/ ./
RUN npm run build

FROM golang:${GO_VERSION} AS app
ARG VERSION=dev
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
# For CGO, set CGO_ENABLED=1 and switch the runtime to base-nossl below.
RUN CGO_ENABLED=0 go build -v -trimpath -tags embed -ldflags "-s -w -X main.version=${VERSION}" -o /app ./cmd/app
# Empty dir for the SQLite database; distroless has no shell to create it in the final stage.
RUN mkdir /data

# Follows distroless's current Debian release; pin it with static-debian13 if needed.
FROM gcr.io/distroless/static:nonroot
# FROM gcr.io/distroless/base-nossl-debian13:nonroot
COPY --from=app /app /app
COPY --from=app --chown=nonroot:nonroot --chmod=1777 /data /data
ENV DB_PATH=/data/app.db
VOLUME /data
EXPOSE 8080
ENTRYPOINT ["/app"]
