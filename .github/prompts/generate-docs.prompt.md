---
mode: 'agent'
model: Claude Sonnet 4
tools: ['codebase']
description: 'Generate comprehensive documentation for the project'
---

# Generate Documentation

Create or update documentation for the VosDroits MCP server project.

## Documentation Types

Ask the user which documentation to generate if not specified:

1. **README.md** - Main project documentation
2. **API Documentation** - MCP tool reference
3. **Code Comments** - Inline documentation
4. **CHANGELOG.md** - Version history
5. **CONTRIBUTING.md** - Contribution guidelines

## README.md Structure

Generate a comprehensive README with these sections:

### 1. Project Title and Description
- Clear, concise description
- Key features and capabilities
- Use case explanation

### 2. Installation

```markdown
## Installation

### Prerequisites
- Go 1.23 or later
- Docker (optional, for containerized deployment)

### Building from Source
\`\`\`bash
git clone https://github.com/yourusername/mcp-vosdroits.git
cd mcp-vosdroits
go mod download
go build -o bin/mcp-vosdroits ./cmd/server
\`\`\`

### Using Docker
\`\`\`bash
docker pull ghcr.io/yourusername/mcp-vosdroits:latest
\`\`\`
```

### 3. Configuration

Document environment variables:
```markdown
## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_NAME | MCP server name | vosdroits |
| SERVER_VERSION | Server version | v1.0.0 |
| LOG_LEVEL | Logging level (debug, info, warn, error) | info |
| HTTP_TIMEOUT | Timeout for HTTP requests | 30s |
```

### 4. Usage

Include quick start examples:
```markdown
## Usage

### Running with Stdio Transport
\`\`\`bash
./bin/mcp-vosdroits
\`\`\`

### Running with Docker
\`\`\`bash
docker run -it ghcr.io/yourusername/mcp-vosdroits:latest
\`\`\`
```

### 5. MCP Tools Reference

Document all available tools:

```markdown
## Available Tools

### search_procedures

Search for procedures on service-public.gouv.fr.

**Input:**
\`\`\`json
{
  "query": "carte d'identité",
  "limit": 10
}
\`\`\`

**Output:**
\`\`\`json
{
  "results": [
    {
      "title": "Carte nationale d'identité",
      "url": "https://...",
      "description": "..."
    }
  ],
  "count": 1
}
\`\`\`

### get_article

Retrieve detailed information from a specific article URL.

**Input:**
\`\`\`json
{
  "url": "https://www.service-public.gouv.fr/..."
}
\`\`\`

**Output:**
\`\`\`json
{
  "title": "Article Title",
  "content": "Full article content...",
  "sections": [...]
}
\`\`\`

### list_categories

List available categories of public service information.

**Input:** None

**Output:**
\`\`\`json
{
  "categories": [
    {
      "name": "Particuliers",
      "url": "..."
    }
  ]
}
\`\`\`
```

### 6. Development

```markdown
## Development

### Running Tests
\`\`\`bash
go test ./...
\`\`\`

### Running with Race Detector
\`\`\`bash
go test -race ./...
\`\`\`

### Code Coverage
\`\`\`bash
go test -cover ./...
\`\`\`

### Linting
\`\`\`bash
golangci-lint run
\`\`\`
```

### 7. Docker

```markdown
## Docker Deployment

### Building the Image
\`\`\`bash
docker build -t mcp-vosdroits .
\`\`\`

### Publishing to GitHub Packages
Images are automatically built and published via GitHub Actions when tags are pushed.
```

### 8. License

```markdown
## License

[Specify your license here]
```

## Code Comment Guidelines

For generating code comments:

### Package Comments
```go
// Package tools provides MCP tool implementations for searching
// French public service information from service-public.gouv.fr.
//
// The package includes tools for:
//   - Searching procedures
//   - Retrieving article content
//   - Listing categories
package tools
```

### Function Comments
```go
// SearchProcedures searches for procedures on service-public.gouv.fr
// matching the given query. It returns up to limit results.
//
// The function validates the query and limit parameters before making
// the HTTP request. It returns an error if the request fails or if
// the response cannot be parsed.
func SearchProcedures(ctx context.Context, query string, limit int) ([]Result, error)
```

## API Documentation Format

Use this format for API docs:

```markdown
# VosDroits MCP Server API Reference

## Tools

### Tool Name

**Description:** Brief description of what the tool does.

**Input Schema:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| field1 | string | Yes | Description |
| field2 | number | No | Description |

**Output Schema:**
| Field | Type | Description |
|-------|------|-------------|
| result1 | string | Description |

**Example Request:**
\`\`\`json
{
  "field1": "value"
}
\`\`\`

**Example Response:**
\`\`\`json
{
  "result1": "value"
}
\`\`\`

**Errors:**
- `empty query`: Query parameter is required
- `invalid limit`: Limit must be between 1 and 100
```

Keep documentation clear, up-to-date, and comprehensive.
