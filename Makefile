.PHONY: build test clean run docker-build docker-run fmt vet tidy

# Build the binary
build:
	go build -o bin/mcp-vosdroits ./cmd/server

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
