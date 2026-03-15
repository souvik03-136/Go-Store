# Dockerfile

# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

# Install CA certificates and git (needed for private module fetching).
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Download dependencies first so Docker caches the layer when only source
# code changes (not go.mod / go.sum).
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a fully static binary (-tags netgo) so it runs in a scratch image.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags netgo -ldflags="-w -s" -o bin/api ./cmd/api

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# Copy the binary and nothing else.
COPY --from=builder /app/bin/api .

EXPOSE 8080

# Run as non-root (distroless provides uid 65532 "nonroot" by default).
USER nonroot:nonroot

ENTRYPOINT ["./api"]