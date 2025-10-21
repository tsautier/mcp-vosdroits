# Colly Web Scraping - Quick Start

This guide shows how to use the Colly-powered web scraping features in the VosDroits MCP server.

## Installation

```bash
# Install dependencies
go get github.com/gocolly/colly/v2

# Build the server
make build
```

## Basic Usage

### Example 1: Search for Procedures

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/guigui42/mcp-vosdroits/internal/client"
)

func main() {
    // Create a client with 30 second timeout
    c := client.New(30 * time.Second)
    
    // Search for "carte d'identité"
    ctx := context.Background()
    results, err := c.SearchProcedures(ctx, "carte d'identité", 10)
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Display results
    for i, result := range results {
        fmt.Printf("%d. %s\n", i+1, result.Title)
        fmt.Printf("   URL: %s\n", result.URL)
        fmt.Printf("   %s\n\n", result.Description)
    }
}
```

### Example 2: Get Article Content

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/guigui42/mcp-vosdroits/internal/client"
)

func main() {
    c := client.New(30 * time.Second)
    ctx := context.Background()
    
    // Fetch a specific article
    articleURL := "https://www.service-public.gouv.fr/particuliers/vosdroits/F1234"
    article, err := c.GetArticle(ctx, articleURL)
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Display article
    fmt.Printf("Title: %s\n", article.Title)
    fmt.Printf("URL: %s\n", article.URL)
    fmt.Printf("\nContent:\n%s\n", article.Content)
}
```

### Example 3: List Categories

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/guigui42/mcp-vosdroits/internal/client"
)

func main() {
    c := client.New(30 * time.Second)
    ctx := context.Background()
    
    // Get available categories
    categories, err := c.ListCategories(ctx)
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Display categories
    fmt.Println("Available Categories:")
    for _, cat := range categories {
        fmt.Printf("- %s: %s\n", cat.Name, cat.Description)
    }
}
```

### Example 4: Using Context Cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/guigui42/mcp-vosdroits/internal/client"
)

func main() {
    c := client.New(30 * time.Second)
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // This will be cancelled if it takes more than 5 seconds
    results, err := c.SearchProcedures(ctx, "passeport", 10)
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Found %d results\n", len(results))
}
```

## Using with MCP Tools

The Colly client is integrated into the MCP tools automatically. When you use the MCP server:

### Via Stdio

```bash
./bin/mcp-vosdroits
```

Then use MCP clients to call the tools:

```json
{
  "tool": "search_procedures",
  "input": {
    "query": "carte d'identité",
    "limit": 10
  }
}
```

### Via HTTP

```bash
HTTP_PORT=8080 ./bin/mcp-vosdroits
```

Then make HTTP requests to the MCP endpoint.

## Advanced Configuration

### Custom Rate Limiting

Modify the client to adjust rate limiting:

```go
// In internal/client/client.go
c.Limit(&colly.LimitRule{
    DomainGlob:  "*.service-public.gouv.fr",
    Parallelism: 2,           // Allow 2 concurrent requests
    Delay:       500 * time.Millisecond,  // Wait 500ms between requests
})
```

### Custom User Agent

```go
c := colly.NewCollector(
    colly.AllowedDomains("www.service-public.gouv.fr"),
    colly.UserAgent("MyCustomBot/1.0"),
)
```

### Adding Callbacks

```go
scraper := c.collector.Clone()

// Before making a request
scraper.OnRequest(func(r *colly.Request) {
    fmt.Println("Visiting", r.URL)
})

// After receiving a response
scraper.OnResponse(func(r *colly.Response) {
    fmt.Println("Got response:", r.StatusCode)
})

// On error
scraper.OnError(func(r *colly.Response, err error) {
    fmt.Println("Error:", err)
})
```

## Performance Tips

### 1. Use Context Timeouts

Always set reasonable timeouts to prevent hanging:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

### 2. Limit Results

Request only what you need:

```go
// Only get top 5 results
results, err := c.SearchProcedures(ctx, "query", 5)
```

### 3. Reuse Client

Create the client once and reuse it:

```go
// Good: Reuse client
client := client.New(30 * time.Second)
results1, _ := client.SearchProcedures(ctx, "query1", 10)
results2, _ := client.SearchProcedures(ctx, "query2", 10)

// Bad: Creating new client each time
c1 := client.New(30 * time.Second)
results1, _ := c1.SearchProcedures(ctx, "query1", 10)
c2 := client.New(30 * time.Second)
results2, _ := c2.SearchProcedures(ctx, "query2", 10)
```

## Troubleshooting

### Issue: No Results Found

**Cause**: Search URL or selectors might have changed

**Solution**: Check fallback results, update selectors if needed

```go
// Fallback results are automatically returned
results, err := c.SearchProcedures(ctx, "query", 10)
// Even if scraping fails, you'll get a helpful fallback result
```

### Issue: Timeout Errors

**Cause**: Network slow or rate limiting too aggressive

**Solution**: Increase timeout

```go
// Increase timeout to 60 seconds
client := client.New(60 * time.Second)
```

### Issue: Context Cancelled

**Cause**: Context timeout or cancellation

**Solution**: Increase context timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()
```

## Testing

Run the test suite:

```bash
# Test client
go test ./internal/client -v

# Test all packages
go test ./... -v
```

## Next Steps

- Read the [Web Scraping Guide](web-scraping.md) for detailed information
- Check the [Colly Documentation](https://go-colly.org/docs/)
- Explore [Colly Examples](https://github.com/gocolly/colly/tree/master/_examples)

## Example Output

### Search Results

```
1. Carte nationale d'identité
   URL: https://www.service-public.gouv.fr/particuliers/vosdroits/F1234
   La carte nationale d'identité (CNI) est un document officiel d'identité...

2. Renouvellement de la carte d'identité
   URL: https://www.service-public.gouv.fr/particuliers/vosdroits/F21089
   Comment renouveler votre carte nationale d'identité...
```

### Article Content

```
Title: Carte nationale d'identité

Content:
La carte nationale d'identité (CNI) permet à son titulaire de certifier de 
son identité. Elle n'est pas obligatoire.

Qui peut avoir une carte d'identité ?
Toute personne de nationalité française peut être titulaire d'une carte 
nationale d'identité...
```

### Categories

```
Available Categories:
- Particuliers: Information and procedures for individuals - family, health, work, housing, etc.
- Professionnels: Information and procedures for professionals - business creation, taxes, employees, etc.
- Associations: Information and procedures for associations - creation, financing, management, etc.
```
