# VosDroits MCP Server

A Model Context Protocol (MCP) server that provides search and retrieval capabilities for French public service information from service-public.gouv.fr.

## Description

This MCP server enables AI assistants to search and retrieve official French administrative procedures and information. Built with Go and powered by intelligent web scraping, it provides three main capabilities:

- **Search Procedures**: Find relevant public service procedures and articles
- **Get Article Details**: Retrieve complete information from specific articles
- **List Categories**: Browse available categories of public service information

## Installation

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

## Available Tools

The server provides three MCP tools:

### 1. search_procedures

Search for procedures on service-public.gouv.fr.

### 2. get_article

Retrieve detailed information from a specific article URL.

### 3. list_categories

List available categories of public service information.

## Documentation

For developers and contributors:
- [Development Guide](docs/DEVELOPMENT.md) - Local development, testing, and contribution guidelines
- [Web Scraping Implementation](docs/SCRAPING.md) - Technical details on service-public.gouv.fr scraping
- [Colly Integration Guide](docs/COLLY_INTEGRATION.md) - Web scraping framework documentation

## License

MIT License - see LICENSE file for details
