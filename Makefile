.PHONY: build test clean run docker-build docker-run fmt vet tidy release-build release-tag

# Build the binary
build:
	go build -o bin/mcp-vosdroits ./cmd/server

# Build release binaries for all platforms
release-build:
	@echo "Building release binaries..."
	@mkdir -p bin/release
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o bin/release/mcp-vosdroits-linux-amd64 ./cmd/server
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o bin/release/mcp-vosdroits-linux-arm64 ./cmd/server
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o bin/release/mcp-vosdroits-darwin-amd64 ./cmd/server
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o bin/release/mcp-vosdroits-darwin-arm64 ./cmd/server
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o bin/release/mcp-vosdroits-windows-amd64.exe ./cmd/server
	@echo "Release binaries built in bin/release/"

# Create and push a release tag
# Usage: make release-tag VERSION=v1.0.0
release-tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ ERROR: VERSION is required. Usage: make release-tag VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release tag $(VERSION)..."
	@if git rev-parse "$(VERSION)" >/dev/null 2>&1; then \
		echo "❌ ERROR: Tag $(VERSION) already exists"; \
		exit 1; \
	fi
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "✅ Created tag $(VERSION)"
	@echo "Pushing tag to origin..."
	@git push origin $(VERSION)
	@echo "✅ Tag pushed successfully"
	@echo ""
	@echo "🚀 Release workflow triggered for $(VERSION)"
	@echo "   Monitor: https://github.com/guigui42/mcp-vosdroits/actions"

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests with race detector
test-race:
	go test -race ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run the server
run: build
	./bin/mcp-vosdroits

# Build Docker image
docker-build:
	docker build -t mcp-vosdroits:latest .

# Run Docker container
docker-run:
	docker run -i mcp-vosdroits:latest

# Format code
fmt:
	go fmt ./...

# Run static analysis
vet:
	go vet ./...

# Tidy dependencies
tidy:
	go mod tidy

# Run all checks
check: fmt vet test

# Install dependencies
deps:
	go mod download
