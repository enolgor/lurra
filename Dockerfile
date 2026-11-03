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
COPY cmd cmd
COPY internal internal
COPY web/*.go web/
COPY --from=web /src/web/dist web/dist
# For CGO, set CGO_ENABLED=1 and switch the runtime to base-nossl below.
RUN CGO_ENABLED=0 go build -trimpath -tags embed -ldflags "-s -w -X main.version=${VERSION}" -o /app ./cmd/app

# Follows distroless's current Debian release; pin it with static-debian13 if needed.
FROM gcr.io/distroless/static:nonroot
# FROM gcr.io/distroless/base-nossl-debian13:nonroot
COPY --from=app /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
