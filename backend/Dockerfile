# Multi-stage build for the Go API
FROM golang:1.23.4-bullseye AS builder
WORKDIR /app

# Pre-fetch deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest and build
COPY . .
# Build a static binary to keep the final image small
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

# Minimal runtime image with certs for TLS DB connections (e.g., Supabase)
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/server /app/server

# Default port; override with SERVER_PORT if needed
ENV SERVER_HOST=0.0.0.0 \
    SERVER_PORT=8080

EXPOSE 8080
CMD ["./server"]
