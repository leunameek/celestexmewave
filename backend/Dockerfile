# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-bookworm AS builder
WORKDIR /app

# Faster, reproducible builds
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Only copy what is needed to download deps
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/. .

# Build the API binary
RUN go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/api


FROM alpine:3.20
WORKDIR /app

# Add minimal runtime deps and a non-root user
RUN apk add --no-cache ca-certificates tzdata curl \
  && addgroup -S app && adduser -S -G app app

# Copy binary and static assets (build from repo root so assets/ is in the context)
COPY --from=builder /out/server /app/server
COPY assets /app/assets

# Entry script to honor DigitalOcean's PORT env var
COPY backend/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh && chown -R app:app /app

USER app

# Sensible defaults; override in your deployment
ENV GIN_MODE=release \
    SERVER_HOST=0.0.0.0 \
    SERVER_PORT=8080 \
    UPLOAD_DIR=/app/assets/images

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD curl -fsS http://127.0.0.1:${SERVER_PORT}/health || exit 1

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/server"]
