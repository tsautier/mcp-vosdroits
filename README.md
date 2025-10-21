# VosDroits MCP Server

A Model Context Protocol (MCP) server written in Go that provides search and retrieval capabilities for French public service information from service-public.gouv.fr.

## Features

- **Search Procedures**: Search for French public service procedures using intelligent web scraping
- **Get Article Details**: Retrieve detailed information from specific articles with HTML parsing
- **List Categories**: Browse available categories of public service information
- **Web Scraping**: Powered by [Colly](https://github.com/gocolly/colly) for robust and respectful scraping

## Technology Stack

- **Language**: Go 1.25+
- **MCP Framework**: [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk)
- **Web Scraping**: [Colly v2](https://github.com/gocolly/colly) - Fast and elegant scraping framework
- **Deployment**: Docker with multi-stage builds

## Prerequisites

- Go 1.23 or higher
- Docker (optional, for containerized deployment)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/guigui42/mcp-vosdroits.git
cd mcp-vosdroits

# Download dependencies
go mod download

# Build the server
go build -o bin/mcp-vosdroits ./cmd/server
```

### Using Docker

```bash
# Build the Docker image
docker build -t mcp-vosdroits .

# Run the container
docker run -i mcp-vosdroits
```

## Usage

### Stdio Transport (Default)

The server uses stdio transport by default, which is suitable for desktop integrations:

```bash
./bin/mcp-vosdroits
```

### HTTP Transport

To run with HTTP transport:

```bash
HTTP_PORT=8080 ./bin/mcp-vosdroits
```

## Configuration

Configure the server using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_NAME` | Name of the MCP server | `vosdroits` |
| `SERVER_VERSION` | Server version | `v1.0.0` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `HTTP_TIMEOUT` | Timeout for HTTP requests to external services | `30s` |
| `HTTP_PORT` | Port for HTTP transport (when enabled) | `8080` |

## Available Tools

### search_procedures

Search for procedures on service-public.gouv.fr.

**Input:**
```json
{
  "query": "carte d'identité",
  "limit": 10
}
```

**Output:**
```json
{
  "results": [
    {
      "title": "Carte nationale d'identité",
      "url": "https://www.service-public.fr/...",
      "description": "..."
    }
  ]
}
```

### get_article

Retrieve detailed information from a specific article URL.

**Input:**
```json
{
  "url": "https://www.service-public.fr/particuliers/vosdroits/F1341"
}
```

**Output:**
```json
{
  "title": "Article title",
  "content": "Full article content...",
  "url": "https://www.service-public.fr/..."
}
```

### list_categories

List available categories of public service information.

**Output:**
```json
{
  "categories": [
    {
      "name": "Particuliers",
      "description": "Information for individuals"
    }
  ]
}
```

## Development

### Local Testing

The easiest way to test the MCP server locally is using the MCP Inspector:

```bash
# Install MCP Inspector globally (one-time setup)
npm install -g @modelcontextprotocol/inspector

# Build your server
make build

# Run the inspector with your server
npx @modelcontextprotocol/inspector ./bin/mcp-vosdroits
```

The MCP Inspector provides a web interface where you can:
- See all available tools
- Test each tool with different inputs
- View responses in real-time
- Debug any issues

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

### Project Structure

```
mcp-vosdroits/
├── cmd/
│   └── server/
│       └── main.go          # Server entry point
├── internal/
│   ├── tools/               # MCP tool implementations
│   ├── client/              # Web scraping client using Colly
│   └── config/              # Configuration management
├── docs/
│   ├── SCRAPING.md          # Scraping implementation details
│   ├── COLLY_INTEGRATION.md # Colly integration guide
│   ├── quick-start.md       # Quick start guide
│   └── web-scraping.md      # Web scraping overview
├── .github/
│   ├── workflows/           # GitHub Actions workflows
│   └── copilot-instructions.md
├── Dockerfile               # Multi-stage Docker build
├── go.mod                   # Go module definition
└── README.md                # This file
```

## Documentation

- [Web Scraping Implementation](docs/SCRAPING.md) - Technical details on service-public.gouv.fr scraping
- [Colly Integration Guide](docs/COLLY_INTEGRATION.md) - Detailed documentation on Colly integration and scraping strategy
- [Quick Start Guide](docs/quick-start.md) - Getting started with development
- [GitHub Copilot Instructions](.github/copilot-instructions.md) - Development guidelines for AI assistance

### Code Quality

Run linters and formatters:

```bash
# Format code
go fmt ./...

# Run static analysis
go vet ./...

# Tidy dependencies
go mod tidy
```

## Web Scraping

This server uses [Colly](https://github.com/gocolly/colly) for respectful and efficient web scraping:

- **Rate Limited**: 1 request per second to avoid overwhelming the target server
- **Context-Aware**: Supports cancellation via Go contexts
- **Robust**: Handles errors gracefully with fallback mechanisms
- **CSS Selectors**: Flexible HTML parsing for extracting structured data

See [Web Scraping Documentation](docs/web-scraping.md) for more details.

## Docker

### Building the Image

```bash
docker build -t mcp-vosdroits:latest .
```

### Running the Container

```bash
# Stdio transport
docker run -i mcp-vosdroits:latest

# HTTP transport
docker run -p 8080:8080 -e HTTP_PORT=8080 mcp-vosdroits:latest
```

### Publishing to GitHub Container Registry

Images are automatically published to `ghcr.io/guigui42/mcp-vosdroits` via GitHub Actions.

## Contributing

Contributions are welcome! Please follow the coding standards and guidelines in `.github/copilot-instructions.md`.

## License

MIT License - see LICENSE file for details
