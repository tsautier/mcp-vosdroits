# Impots.gouv.fr Scraping Implementation

This document describes the technical implementation details for scraping tax information from impots.gouv.fr.

## Overview

The impots.gouv.fr client provides three main capabilities:
- Searching for tax forms and articles
- Retrieving detailed content from specific tax documents
- Listing available tax service categories

## Website Structure

### Search Results

The search functionality uses the following URL pattern:
```
https://www.impots.gouv.fr/recherche/{query}?origin[]=impots&search_filter=Filtrer
```

Example:
- `https://www.impots.gouv.fr/recherche/formulaire%202042?origin[]=impots&search_filter=Filtrer`
- `https://www.impots.gouv.fr/recherche/PEA?origin[]=impots&search_filter=Filtrer`

#### HTML Structure

Search results are rendered as cards using the DSFR (Système de Design de l'État Français) framework:

```html
<div class="fr-card">
  <div class="fr-card__body">
    <div class="fr-card__content">
      <h3 class="fr-card__title fr-h6">
        <a href="/formulaire/2042/declaration-des-revenus" title="...">
          Formulaire 2042 : Déclaration de revenus
        </a>
      </h3>
      <div class="fr-card__start">
        <div class="fr-card__detail">Fiche formulaire</div>
      </div>
      <div class="fr-card__end">
        <p class="fr-card__detail fr-icon-calendar-line">16/04/2025 - Site impots.gouv.fr</p>
      </div>
    </div>
  </div>
</div>
```

Key selectors:
- Result cards: `div.fr-card`
- Title link: `h3.fr-card__title a`
- Document type: `div.fr-card__detail`
- Publication date: `p.fr-card__detail`
- Description (when available): `p.fr-card__desc`

### Article Pages

Tax documents and articles use various URL patterns:
- Forms: `https://www.impots.gouv.fr/formulaire/{form-number}/{form-name}`
- Articles: `https://www.impots.gouv.fr/particulier/{article-path}`

Example URLs:
- `https://www.impots.gouv.fr/formulaire/2042/declaration-des-revenus`
- `https://www.impots.gouv.fr/particulier/lassurance-vie-et-le-pea-0`

#### HTML Structure

```html
<head>
  <title>Formulaire n°2042 | impots.gouv.fr</title>
  <meta property="og:title" content="Formulaire n°2042" />
</head>
<body>
  <div class="fr-breadcrumb">...</div>
  <main>
    <h1>...</h1>
    <div class="fr-card__desc">...</div>
    <p>...</p>
  </main>
</body>
```

Key selectors:
- Page title: `head title` and `meta[property='og:title']`
- Main content: `main, article, div.main-content, div.content`
- Content elements: `h1, h2, h3, p, li, div.fr-callout, div.fr-card__desc`
- Breadcrumb (for type detection): `div.fr-breadcrumb`

### Categories

Main navigation categories are available at the top-level navigation:

```html
<nav role="navigation" id="main-navigation" aria-label="Menu principal">
  <ul class="menu--main fr-nav__list">
    <li class="fr-nav__item">
      <a href="/particulier" class="fr-nav__link">Particulier</a>
    </li>
    <li class="fr-nav__item">
      <a href="/professionnel" class="fr-nav__link">Professionnel</a>
    </li>
    <!-- ... -->
  </ul>
</nav>
```

Key selectors:
- Navigation links: `nav.fr-nav a.fr-nav__link`

## Scraping Strategy

### Search Implementation

1. Build search URL with properly encoded query
2. Visit search results page
3. Extract cards using `div.fr-card` selector
4. For each card, extract:
   - Title from `h3.fr-card__title a`
   - URL from the same link's href attribute
   - Type from `div.fr-card__detail`
   - Date from `p.fr-card__detail`
   - Description from `p.fr-card__desc` (optional)
5. Apply result limit
6. Return results or fallback on error

### Article Retrieval

1. Validate URL is from impots.gouv.fr domain
2. Visit article page
3. Extract title from `<title>` tag or `<meta property="og:title">`
4. Extract content from main content selectors
5. Filter out navigation, scripts, and boilerplate
6. Detect document type from breadcrumb
7. Combine content parts with line breaks

### Category Listing

1. Visit main navigation page (e.g., `/particulier`)
2. Extract navigation links from `nav.fr-nav a.fr-nav__link`
3. Filter out "Accueil" and duplicates
4. Build category objects with name, description, and URL
5. Return categories or fallback to defaults

## Rate Limiting

The client implements respectful rate limiting:
- 1 request per second maximum
- Single parallel request (no concurrent scraping)
- Timeout: configurable (default 30 seconds)

## Error Handling

The scraper handles several error scenarios:

1. **No Results Found**: Returns fallback result with search URL
2. **Network Errors**: Returns error or fallback
3. **Parsing Errors**: Attempts multiple selectors, falls back gracefully
4. **Invalid URLs**: Validates domain before scraping
5. **Context Cancellation**: Respects context timeouts and cancellations

## Content Filtering

The scraper filters out unwanted content:
- JavaScript fragments
- Cookie notices
- Navigation elements
- Empty strings
- Very short text (< 10 characters)
- Specific keywords: "javascript", "Cookie", "Navigation"

## Fallback Mechanisms

When scraping fails, the client provides fallback responses:

1. **Search Fallback**: Returns a helpful message with the search URL
2. **Categories Fallback**: Returns predefined default categories
3. **Partial Success**: If some results are found despite errors, returns what was collected

## Testing

The implementation includes comprehensive tests:

- Unit tests for invalid inputs and edge cases
- URL validation tests
- Domain checking tests
- Default category generation tests
- Fallback mechanism tests

## User Agent

All requests use a custom user agent:
```
VosDroits-MCP-Server/1.0
```

This identifies the scraper to the website administrators.

## Compliance Notes

- The scraper respects robots.txt
- Rate limiting prevents server overload
- User agent identifies the bot
- Only public information is accessed
- Content is used for informational purposes only

## Example Usage

### Search for Tax Forms

```go
client := NewImpotsClient(30 * time.Second)
results, err := client.SearchImpots(ctx, "formulaire 2042", 10)
```

### Retrieve Form Details

```go
article, err := client.GetImpotsArticle(ctx, 
    "https://www.impots.gouv.fr/formulaire/2042/declaration-des-revenus")
```

### List Categories

```go
categories, err := client.ListImpotsCategories(ctx)
```
