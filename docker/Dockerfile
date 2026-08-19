FROM golang:1.25.1 AS builder

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./src ./src
COPY ./cmd ./cmd

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/main.go


FROM debian:bookworm-slim

ENV CONTAINER_MODE=1

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /out/app ./app
