# Web Scraping with Colly

This document describes how the VosDroits MCP server uses [Colly](https://github.com/gocolly/colly) for web scraping French public service information from service-public.gouv.fr.

## Overview

Colly is a lightning-fast and elegant web scraping framework for Go. We use it to:

1. **Search for procedures** - Scrape search results from service-public.gouv.fr
2. **Extract article content** - Parse HTML to get structured content from articles
3. **List categories** - Discover available service categories

## Why Colly?

- **Fast**: Concurrent requests with rate limiting support
- **Clean API**: Simple callbacks for HTML element selection
- **Robust**: Handles common web scraping challenges (redirects, cookies, etc.)
- **Go-native**: Written in Go, integrates seamlessly with our MCP server
- **Respectful**: Built-in rate limiting to be polite to target servers

## Implementation Details

### Client Configuration

```go
// Create a new Colly collector with configuration
c := colly.NewCollector(
    colly.AllowedDomains("www.service-public.gouv.fr", "service-public.gouv.fr"),
    colly.UserAgent("VosDroits-MCP-Server/1.0"),
    colly.Async(false),
)

// Set timeout
c.SetRequestTimeout(timeout)

// Configure rate limiting to be respectful
c.Limit(&colly.LimitRule{
    DomainGlob:  "*.service-public.gouv.fr",
    Parallelism: 1,
    Delay:       1 * time.Second,
})
```

### Key Features Used

#### 1. HTML Element Selection

We use CSS selectors to extract data:

```go
scraper.OnHTML("div.search-result, article.item", func(e *colly.HTMLElement) {
    title := e.ChildText("h2, h3, .title")
    url := e.ChildAttr("a[href]", "href")
    description := e.ChildText("p, .description")
})
```

#### 2. Rate Limiting

Respectful scraping with delays:

- 1 request per second to service-public.gouv.fr
- Parallelism limited to 1 to avoid overwhelming the server
- Proper timeouts to prevent hanging

#### 3. Context Cancellation

Integration with Go's context for cancellability:

```go
go func() {
    select {
    case <-ctx.Done():
        // Context cancelled, stop scraping
        scraper = nil
    case <-done:
        // Operation completed normally
    }
}()
```

#### 4. Error Handling

Graceful error handling with fallbacks:

```go
scraper.OnError(func(r *colly.Response, err error) {
    // Log error and provide fallback results
})
```

## Search Implementation

The search functionality:

1. Builds a search URL with the query
2. Scrapes search results using CSS selectors
3. Extracts title, URL, and description for each result
4. Returns up to the specified limit
5. Falls back to helpful messages if no results found

### HTML Selectors

We use multiple selectors to handle different page structures:

- **Search results**: `div.search-result, article.item, li.result-item`
- **Titles**: `h2, h3, .title, a`
- **Descriptions**: `p, .description, .summary`

## Article Extraction

Article content extraction:

1. Validates the URL is from service-public.gouv.fr
2. Scrapes the article page
3. Extracts title from `<h1>` or `.page-title`
4. Extracts content from article body (paragraphs, headings, lists)
5. Returns structured Article object

### Content Filtering

We filter out unwanted content:

- JavaScript snippets
- Navigation elements
- Empty text nodes
- Duplicates

## Category Listing

Category discovery:

1. Visits the homepage
2. Scrapes navigation links
3. Filters for main categories (Particuliers, Professionnels, Associations)
4. Returns structured category information
5. Falls back to default categories if scraping fails

## Best Practices

### 1. Respectful Scraping

- **Rate limiting**: 1 second delay between requests
- **User agent**: Clear identification as "VosDroits-MCP-Server"
- **Domain restrictions**: Only scrape allowed domains
- **Timeouts**: Don't hang indefinitely on slow responses

### 2. Error Handling

- Always provide fallback results
- Don't fail completely if one element is missing
- Log errors for debugging but return partial results when possible

### 3. Context Awareness

- Respect context cancellation
- Clean up resources properly
- Use channels for error communication

### 4. Testing

- Test with real and mock data
- Handle edge cases (empty results, 404s, timeouts)
- Verify fallback mechanisms work

## Limitations

### Current Limitations

1. **Search URL Structure**: We rely on a specific URL pattern for search
2. **HTML Selectors**: Selectors may need updates if the website changes
3. **No JavaScript**: Colly doesn't execute JavaScript (use chromedp if needed)
4. **Rate Limiting**: Conservative to be polite (may be slow for bulk operations)

### Future Improvements

1. **Caching**: Cache frequent searches to reduce load
2. **Selector Discovery**: Automatically adapt to page structure changes
3. **JavaScript Support**: Add chromedp for dynamic content if needed
4. **Parallel Requests**: Increase parallelism for better performance
5. **Retry Logic**: Implement exponential backoff for failed requests

## Troubleshooting

### Common Issues

**Problem**: Search returns no results
- **Cause**: Search URL format changed or selectors don't match
- **Solution**: Check fallback results, update selectors if needed

**Problem**: Article content is empty
- **Cause**: HTML structure doesn't match selectors
- **Solution**: Inspect the page HTML and update selectors

**Problem**: Rate limiting errors
- **Cause**: Too many requests too quickly
- **Solution**: Increase delay in LimitRule

**Problem**: Timeout errors
- **Cause**: Slow network or server response
- **Solution**: Increase timeout or implement retry logic

## Resources

- [Colly Documentation](https://go-colly.org/docs/)
- [Colly GitHub](https://github.com/gocolly/colly)
- [CSS Selectors Reference](https://www.w3schools.com/cssref/css_selectors.asp)
- [service-public.gouv.fr](https://www.service-public.gouv.fr)

## Example Usage

```go
// Create a client
client := client.New(30 * time.Second)

// Search for procedures
results, err := client.SearchProcedures(ctx, "carte d'identité", 10)

// Get an article
article, err := client.GetArticle(ctx, "https://www.service-public.gouv.fr/...")

// List categories
categories, err := client.ListCategories(ctx)
```
