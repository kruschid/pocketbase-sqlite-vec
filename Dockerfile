# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.25

FROM golang:${GO_VERSION}-bookworm AS build

ARG POCKETBASE_VERSION=0.39.10

RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 go test -tags no_default_driver ./...
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 go build \
      -tags no_default_driver \
      -trimpath \
      -ldflags="-s -w -X github.com/pocketbase/pocketbase.Version=v${POCKETBASE_VERSION}" \
      -o /out/pocketbase .

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates libsqlite3-0 tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --gid 1000 pocketbase \
    && useradd --uid 1000 --gid pocketbase --no-create-home --home-dir /pb --shell /usr/sbin/nologin pocketbase \
    && mkdir -p /pb/pb_data /pb/pb_hooks /pb/pb_migrations /pb/pb_public \
    && chown -R pocketbase:pocketbase /pb

WORKDIR /pb
COPY --from=build /out/pocketbase /usr/local/bin/pocketbase

USER pocketbase
EXPOSE 8090
VOLUME ["/pb/pb_data"]

ENTRYPOINT ["pocketbase"]
CMD ["serve", "--http=0.0.0.0:8090"]
