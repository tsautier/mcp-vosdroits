# Service-Public.gouv.fr Scraping Implementation

This document describes the web scraping implementation for the VosDroits MCP server.

## Overview

The scraper is built using [Colly](https://github.com/gocolly/colly), a fast and elegant web scraping framework for Go. It extracts information from service-public.gouv.fr with three main capabilities:

1. **Search Procedures** - Search for procedures and articles
2. **Get Article** - Retrieve detailed article content
3. **List Categories** - List available service categories

## Search Implementation

### URL Pattern
```
https://www.service-public.gouv.fr/particuliers/recherche?keyword={query}
```

### HTML Structure
Search results are contained in `<li>` elements with IDs matching the pattern `result_*`:

```html
<li id="result_fichePratique_1">
  <div class="sp-link">
    <a href="/particuliers/vosdroits/F2726" class="fr-link">
      <span><span>Nationalité française par mariage</span></span>
    </a>
  </div>
</li>
```

### Extraction Logic
- **Title**: Extracted from innermost `<span>` to avoid duplication
- **URL**: From `a.fr-link` href attribute, converted to absolute URL
- **Description**: Currently not available in search results (field kept for future use)

### Example Results
Query: `"nationalité française"`

```
1. Nationalité française
   URL: https://www.service-public.gouv.fr/particuliers/vosdroits/N111

2. Nationalité française par mariage
   URL: https://www.service-public.gouv.fr/particuliers/vosdroits/F2726
```

## Article Retrieval

### URL Pattern
```
https://www.service-public.gouv.fr/particuliers/vosdroits/{article_id}
```

Examples:
- `https://www.service-public.gouv.fr/particuliers/vosdroits/F2726` (Fiche pratique)
- `https://www.service-public.gouv.fr/particuliers/vosdroits/N111` (Navigation/theme)

### HTML Structure
Articles use a consistent structure:

```html
<h1 id="titlePage">Nationalité française par mariage</h1>

<div id="intro">
  <p class="fr-text--lg">Introduction text...</p>
</div>

<article class="article">
  <h2>Section heading</h2>
  <p data-test="contenu-texte">Content paragraph...</p>
  ...
</article>
```

### Extraction Logic
- **Title**: Extracted from `h1#titlePage`
- **Content**: Combines:
  - Introduction from `div#intro p.fr-text--lg`
  - Main content from paragraphs with `data-test="contenu-texte"`
  - Headings (`h2`, `h3`) for structure
  - Callout boxes and sections

### Content Filtering
The scraper filters out:
- JavaScript content
- Navigation elements ("Votre situation", "Abonnement")
- Very short text (< 10 characters)

### Example Output
```
Title: Nationalité française par mariage
Content length: ~34,000 characters
Content preview: Vous êtes marié(e) avec un(e) français(e) et vous 
                 voulez avoir la nationalité française ? Vous pouvez 
                 faire une déclaration de nationalité française par 
                 mariage...
```

## Categories Implementation

### URL Pattern
```
https://www.service-public.gouv.fr/particuliers
```

### HTML Structure
Categories are listed in the footer theme list:

```html
<ul class="sp-theme-list">
  <li>
    <a href="/particuliers/vosdroits/N19810" class="fr-footer__top-link">
      Papiers - Citoyenneté - Élections
    </a>
  </li>
  ...
</ul>
```

### Extraction Logic
- **Name**: Text content of the link
- **URL Pattern**: Links matching `/particuliers/vosdroits/N*`
- **Description**: Generated as "Information and procedures for {name}"

### Available Categories
The scraper typically finds 11 main categories:

1. Papiers - Citoyenneté - Élections
2. Famille - Scolarité
3. Social - Santé
4. Travail - Formation
5. Logement
6. Transports - Mobilité
7. Argent - Impôts - Consommation
8. Justice
9. Étranger - Europe
10. Loisirs - Sports - Culture
11. Associations, fondations et fonds de dotation

### Fallback
If scraping fails, default categories are returned:
- Particuliers
- Professionnels
- Associations

## Technical Details

### Rate Limiting
The scraper implements respectful rate limiting:
- 1 request per second to `*.service-public.gouv.fr`
- 30-second timeout for requests
- Sequential requests (no parallelism)

### User Agent
```
VosDroits-MCP-Server/1.0
```

### Error Handling
- Context cancellation support
- Graceful fallback to default results on errors
- Partial results returned if some data is extracted

### Domain Restrictions
Only allows scraping from:
- `www.service-public.gouv.fr`
- `service-public.gouv.fr`

## Testing

Integration tests verify real-world scraping:

```bash
# Run all integration tests
go test -v ./internal/client

# Run specific test
go test -v -run TestSearchProceduresIntegration ./internal/client

# Skip integration tests (short mode)
go test -short ./...
```

### Test Coverage
- ✅ Search with multiple queries
- ✅ Article retrieval with real URLs
- ✅ Category listing
- ✅ Title extraction (no duplication)
- ✅ Content extraction with proper filtering

## Maintenance Notes

### Potential Breaking Changes
If service-public.gouv.fr updates their HTML structure, these selectors may need adjustment:

1. **Search results**: `li[id^='result_']` and `a.fr-link span span`
2. **Article title**: `h1#titlePage`
3. **Article content**: `article.article`, `p[data-test='contenu-texte']`
4. **Categories**: `ul.sp-theme-list` and `a.fr-footer__top-link`

### Debugging
To debug scraping issues, use curl to inspect the HTML:

```bash
# Check search results
curl -s "https://www.service-public.gouv.fr/particuliers/recherche?keyword=carte" | grep "result_"

# Check article structure
curl -s "https://www.service-public.gouv.fr/particuliers/vosdroits/F2726" | grep -E "titlePage|contenu-texte"
```

## Dependencies

- [Colly v2](https://github.com/gocolly/colly) - Web scraping framework
- Go 1.23+ - Programming language

## Related Files

- `/internal/client/client.go` - Main scraping implementation
- `/internal/client/integration_test.go` - Integration tests
- `/internal/tools/tools.go` - MCP tool handlers
