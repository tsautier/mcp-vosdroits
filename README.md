# VosDroits MCP Server

A Model Context Protocol (MCP) server that provides search and retrieval capabilities for French public service and tax information from service-public.gouv.fr and impots.gouv.fr.

## Description

This MCP server enables AI assistants to search and retrieve official French administrative procedures and tax information. Built with Go and powered by intelligent web scraping, it provides eight main capabilities:

### Service-Public.gouv.fr Tools

- **search_procedures**: Find relevant public service procedures and articles
- **get_article**: Retrieve complete information from specific service-public.gouv.fr articles
- **list_categories**: Browse available categories of public service information
- **list_life_events**: List all available life events (événements de vie) guides
- **get_life_event_details**: Retrieve detailed information about specific life situations

### Impots.gouv.fr Tools

- **search_impots**: Search for tax forms, articles, and procedures on impots.gouv.fr
- **get_impots_article**: Retrieve detailed information from specific tax articles or forms
- **list_impots_categories**: List available tax service categories

## Installation

### Download Pre-Built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/guigui42/mcp-vosdroits/releases).

Available platforms:
- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64

```bash
# Example: Download and run Linux binary
curl -LO https://github.com/guigui42/mcp-vosdroits/releases/latest/download/mcp-vosdroits-linux-amd64
chmod +x mcp-vosdroits-linux-amd64
./mcp-vosdroits-linux-amd64
```

All binaries include SHA256 checksums for verification.

### Using Docker (Recommended)

Pull and run the official image from GitHub Container Registry:

```bash
docker pull ghcr.io/guigui42/mcp-vosdroits:latest
docker run -i ghcr.io/guigui42/mcp-vosdroits:latest
```

### VSCode with GitHub Copilot

To use this MCP server with GitHub Copilot in VSCode, you need to configure it in your MCP settings. See the [VSCode MCP documentation](https://code.visualstudio.com/docs/copilot/customization/mcp-servers) for detailed information.

#### One-Click Install:

[![Install in VS Code](https://img.shields.io/badge/VS_Code-VS_Code?style=flat-square&label=Install%20VosDroits%20MCP&color=0098FF)](https://insiders.vscode.dev/redirect?url=vscode%3Amcp%2Finstall%3F%257B%2522name%2522%253A%2522vosdroits%2522%252C%2522command%2522%253A%2522docker%2522%252C%2522args%2522%253A%255B%2522run%2522%252C%2522-i%2522%252C%2522--rm%2522%252C%2522ghcr.io%252Fguigui42%252Fmcp-vosdroits%253Alatest%2522%255D%257D) [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-VS_Code_Insiders?style=flat-square&label=Install%20VosDroits%20MCP&color=24bfa5)](https://insiders.vscode.dev/redirect?url=vscode-insiders%3Amcp%2Finstall%3F%257B%2522name%2522%253A%2522vosdroits%2522%252C%2522command%2522%253A%2522docker%2522%252C%2522args%2522%253A%255B%2522run%2522%252C%2522-i%2522%252C%2522--rm%2522%252C%2522ghcr.io%252Fguigui42%252Fmcp-vosdroits%253Alatest%2522%255D%257D)

#### Manual Installation:

Follow the [MCP install guide](https://code.visualstudio.com/docs/copilot/customization/mcp-servers), and use the standard config below. The configuration should be added to your MCP settings file (typically `~/Library/Application Support/Code/User/mcp.json` on macOS).

**Using Docker Image (Recommended):**

```json
{
  "servers": {
    "vosdroits": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "ghcr.io/guigui42/mcp-vosdroits:latest"
      ]
    }
  }
}
```

**Using Local Binary:**

If you've built the server from source:

```json
{
  "servers": {
    "vosdroits": {
      "command": "/absolute/path/to/mcp-vosdroits/bin/mcp-vosdroits"
    }
  }
}
```

**With Environment Variables:**

Configure custom environment variables for logging and timeout:

```json
{
  "servers": {
    "vosdroits": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e", "LOG_LEVEL=debug",
        "-e", "HTTP_TIMEOUT=60s",
        "ghcr.io/guigui42/mcp-vosdroits:latest"
      ]
    }
  }
}
```

After adding the configuration, restart VSCode or reload the window. The server will be available in GitHub Copilot Chat, and you can use the available tools to query French public service information.

### GitHub Copilot CLI

To use this MCP server with GitHub Copilot CLI, add the configuration to your MCP settings file (`~/.copilot/mcp-config.json`).

**Using Docker Image (Recommended):**

```json
{
  "mcpServers": {
    "vosdroits": {
      "type": "local",
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "ghcr.io/guigui42/mcp-vosdroits:latest"
      ],
      "tools": [
        "*"
      ]
    }
  }
}
```

**Using Local Binary:**

If you've built the server from source:

```json
{
  "mcpServers": {
    "vosdroits": {
      "type": "local",
      "command": "/absolute/path/to/mcp-vosdroits/bin/mcp-vosdroits",
      "args": [],
      "tools": [
        "*"
      ]
    }
  }
}
```

**With Environment Variables:**

Configure custom environment variables for logging and timeout:

```json
{
  "mcpServers": {
    "vosdroits": {
      "type": "local",
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "ghcr.io/guigui42/mcp-vosdroits:latest"
      ],
      "env": {
        "LOG_LEVEL": "debug",
        "HTTP_TIMEOUT": "60s"
      },
      "tools": [
        "*"
      ]
    }
  }
}
```

After adding the configuration, restart your terminal or run `gh copilot reload` to load the new MCP server.

### Claude Desktop

To use this MCP server with Claude Desktop, add the configuration to your MCP settings file:
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

**Using Docker Image (Recommended):**

```json
{
  "mcpServers": {
    "vosdroits": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "ghcr.io/guigui42/mcp-vosdroits:latest"
      ]
    }
  }
}
```

**Using Local Binary:**

If you've built the server from source:

```json
{
  "mcpServers": {
    "vosdroits": {
      "command": "/absolute/path/to/mcp-vosdroits/bin/mcp-vosdroits",
      "args": []
    }
  }
}
```

**With Environment Variables:**

Configure custom environment variables for logging and timeout:

```json
{
  "mcpServers": {
    "vosdroits": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e", "LOG_LEVEL=debug",
        "-e", "HTTP_TIMEOUT=60s",
        "ghcr.io/guigui42/mcp-vosdroits:latest"
      ]
    }
  }
}
```

After adding the configuration, restart Claude Desktop. The server will be available, and you can use the available tools to query French public service and tax information.

## Available Tools

The server provides eight MCP tools across two domains:

### Which Tool Should I Use?

**For major life situations** (buying a house, getting married, having a baby, death, moving, retirement, etc.):
- 🎯 **Start with `list_life_events`** - comprehensive guides organized by topic
- Then use **`get_life_event_details`** with a URL to get all procedures

**For specific administrative procedures** (passport renewal, driver's license, etc.):
- 🔍 **Use `search_procedures`** - targeted search for specific procedures
- Then use **`get_article`** to get full details

**For tax information**:
- 💰 **Use `search_impots`** for forms and tax procedures
- Then use **`get_impots_article`** for detailed information

### Service-Public.gouv.fr Tools

#### 1. search_procedures

Search for procedures on service-public.gouv.fr.

**Input:**
- `query` (string): Search query for procedures
- `limit` (int, optional): Maximum number of results to return (1-100, default: 10)

**Output:**
- `results`: Array of matching procedures with title, URL, and description

#### 2. get_article

Retrieve detailed information from a specific article URL on service-public.gouv.fr.

**Input:**
- `url` (string): URL of the article to retrieve

**Output:**
- `title`: Article title
- `content`: Full article content
- `url`: Article URL

#### 3. list_categories

List available categories of public service information.

**Output:**
- `categories`: Array of available categories with name and description

#### 4. list_life_events

List all available life events (événements de vie) from the "Comment faire si" section of service-public.gouv.fr. These are comprehensive practical guides for major life situations like expecting a child, moving, retirement, etc.

**Output:**
- `events`: Array of life events with title, URL, and description

**Example events:**
- "J'attends un enfant" (Expecting a child)
- "Je déménage en France" (Moving within France)
- "Un proche est décédé" (Death of a loved one)
- "Je prépare ma retraite" (Preparing for retirement)

#### 5. get_life_event_details

Retrieve detailed information about a specific life event, including all sections organized by topic (Health, Civil Status, Employment, etc.).

**Input:**
- `url` (string): URL of the life event to retrieve (from list_life_events results)

**Output:**
- `title`: Life event title
- `url`: Life event URL
- `introduction`: Overview text
- `sections`: Array of detailed sections with title and content

**See also:** [Life Events Documentation](docs/LIFE_EVENTS.md)

### Impots.gouv.fr Tools

#### 4. search_impots

Search for tax forms, articles, and procedures on impots.gouv.fr.

**Input:**
- `query` (string): Search query for tax information and forms (e.g., "formulaire 2042", "PEA")
- `limit` (int, optional): Maximum number of results to return (1-100, default: 10)

**Output:**
- `results`: Array of matching tax documents with title, URL, description, type, and date

**Example queries:**
- "formulaire 2042" - Find the income tax declaration form
- "PEA" - Find information about equity savings plans
- "crédit d'impôt" - Find information about tax credits

#### 5. get_impots_article

Retrieve detailed information from a specific tax article or form URL on impots.gouv.fr.

**Input:**
- `url` (string): URL of the tax article or form to retrieve

**Output:**
- `title`: Document title
- `content`: Full document content
- `url`: Document URL
- `type`: Type of document (Formulaire, Article, etc.)
- `description`: Brief description

#### 6. list_impots_categories

List available categories of tax information on impots.gouv.fr.

**Output:**
- `categories`: Array of tax categories (Particulier, Professionnel, Partenaire, Collectivité, International) with name, description, and URL

## Screenshots
<img width="1633" height="1292" alt="20251021212600" src="https://github.com/user-attachments/assets/12eb095f-37e6-4b18-89ad-767f1bf558a5" />



## Documentation

For developers and contributors:
- [Development Guide](docs/DEVELOPMENT.md) - Local development, testing, and contribution guidelines
- [Release Process](docs/RELEASE.md) - How releases are created and automated
- [Web Scraping Implementation](docs/SCRAPING.md) - Technical details on service-public.gouv.fr scraping
- [Colly Integration Guide](docs/COLLY_INTEGRATION.md) - Web scraping framework documentation

## License

MIT License - see LICENSE file for details
