# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically-linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o mcp-vosdroits ./cmd/server

# Production stage
FROM scratch

# Copy binary from builder
COPY --from=builder /build/mcp-vosdroits .

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/mcp-vosdroits .

# Run as non-root user
USER nobody

# Set default environment variables
ENV SERVER_NAME=vosdroits \
    SERVER_VERSION=v1.0.0 \
    LOG_LEVEL=info

ENTRYPOINT ["/app/mcp-vosdroits"]
