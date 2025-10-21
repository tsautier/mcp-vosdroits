---
description: 'Docker containerization best practices for Go MCP servers'
applyTo: '**/Dockerfile,**/*.dockerfile,**/docker-compose.yml'
---

# Docker Guidelines

## Multi-Stage Builds

Use multi-stage builds to create minimal production images:

1. **Build stage**: Compile the Go binary with all dependencies
2. **Production stage**: Copy only the binary to a minimal base image

## Go-Specific Best Practices

### Static Compilation

Build statically-linked binaries for minimal images:

```dockerfile
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
```

### Base Images

- Use `golang:1.23-alpine` for build stage (smaller than full golang image)
- Use `scratch` or `alpine` for production stage
- For `scratch` images, ensure binary is statically linked
- For `alpine`, include necessary CA certificates

### Binary Optimization

- Use `-ldflags="-w -s"` to strip debug information and reduce binary size
- Consider using UPX for further compression (if acceptable)

## Security

### User Permissions

- Don't run as root in production
- Create a non-root user in the image
- Use `USER` directive to switch to non-root user

### Image Scanning

- Regularly scan images for vulnerabilities
- Keep base images updated
- Minimize the number of layers

### Secrets Management

- Never hardcode secrets in Dockerfile
- Use build arguments for build-time secrets
- Use environment variables or secret management for runtime secrets
- Don't commit sensitive files

## Image Optimization

### Layer Caching

- Order commands from least to most frequently changing
- Copy `go.mod` and `go.sum` first, then download dependencies
- Copy source code last

### Size Reduction

- Remove unnecessary files
- Use `.dockerignore` to exclude files from build context
- Combine RUN commands to reduce layers
- Clean up package manager caches

## Health Checks

- Add `HEALTHCHECK` instruction for container health monitoring
- Keep health checks lightweight
- Set appropriate timeout and interval

## Labels

- Use `LABEL` instructions for metadata
- Include version, description, maintainer
- Follow OCI image spec annotations

## Example Structure

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o app ./cmd/server

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/app .
USER nobody
ENTRYPOINT ["/app/app"]
```

## Docker Compose

- Use for local development and testing
- Define service dependencies clearly
- Use environment files for configuration
- Set resource limits appropriately
