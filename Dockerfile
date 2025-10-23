# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies with cache mount for faster rebuilds
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build statically-linked binary with build cache
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-w -s" -o mcp-vosdroits ./cmd/server

# Production stage
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D -u 1000 appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/mcp-vosdroits .

# Run as non-root user
USER appuser

# Set default environment variables
ENV SERVER_NAME=vosdroits \
    SERVER_VERSION=v1.0.0 \
    LOG_LEVEL=info

ENTRYPOINT ["/app/mcp-vosdroits"]
