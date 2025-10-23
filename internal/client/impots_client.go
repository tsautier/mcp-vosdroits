// Package client provides HTTP client functionality for impots.gouv.fr.
package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// ImpotsClient handles HTTP requests to impots.gouv.fr using Colly for web scraping.
type ImpotsClient struct {
	collector *colly.Collector
	baseURL   string
	timeout   time.Duration
}

// NewImpotsClient creates a new ImpotsClient with the specified timeout.
func NewImpotsClient(timeout time.Duration) *ImpotsClient {
	c := colly.NewCollector(
		colly.AllowedDomains("www.impots.gouv.fr", "impots.gouv.fr"),
		colly.UserAgent("VosDroits-MCP-Server/1.0"),
		colly.Async(false),
	)

	c.SetRequestTimeout(timeout)

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*.impots.gouv.fr",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})
	if err != nil {
		// Fallback without rate limiting if it fails
	}

	return &ImpotsClient{
		collector: c,
		baseURL:   "https://www.impots.gouv.fr",
		timeout:   timeout,
	}
}

// ImpotsSearchResult represents a search result from impots.gouv.fr.
type ImpotsSearchResult struct {
	Title       string
	URL         string
	Description string
	Type        string
	Date        string
}

// SearchImpots searches for tax information matching the query.
func (c *ImpotsClient) SearchImpots(ctx context.Context, query string, limit int) ([]ImpotsSearchResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	var results []ImpotsSearchResult
	errorChan := make(chan error, 1)

	scraper := c.collector.Clone()

	// Allow URL revisits to prevent "already visited" errors on repeated calls
	scraper.AllowURLRevisit = true

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			scraper = nil
		case <-done:
		}
	}()

	// Handle search results - impots.gouv.fr uses div.fr-card
	scraper.OnHTML("div.fr-card", func(e *colly.HTMLElement) {
		if len(results) >= limit {
			return
		}

		href := e.ChildAttr("a[href]", "href")
		if href == "" {
			return
		}

		fullURL := e.Request.AbsoluteURL(href)

		title := strings.TrimSpace(e.ChildText("h3.fr-card__title a"))
		if title == "" {
			title = strings.TrimSpace(e.ChildText("h3.fr-card__title"))
		}

		resultType := strings.TrimSpace(e.ChildText("div.fr-card__detail"))
		date := strings.TrimSpace(e.ChildText("p.fr-card__detail"))

		description := strings.TrimSpace(e.ChildText("p.fr-card__desc"))

		if title != "" && fullURL != "" {
			results = append(results, ImpotsSearchResult{
				Title:       title,
				URL:         fullURL,
				Description: description,
				Type:        resultType,
				Date:        date,
			})
		}
	})

	scraper.OnError(func(r *colly.Response, err error) {
		select {
		case errorChan <- fmt.Errorf("scraping error: %w", err):
		default:
		}
	})

	searchURL := fmt.Sprintf("%s/recherche/%s?origin[]=impots&search_filter=Filtrer",
		c.baseURL, url.QueryEscape(query))

	if err := scraper.Visit(searchURL); err != nil {
		return c.fallbackImpotsSearch(ctx, query, limit)
	}

	scraper.Wait()

	select {
	case err := <-errorChan:
		if len(results) == 0 {
			return nil, err
		}
	default:
	}

	if len(results) == 0 {
		return c.fallbackImpotsSearch(ctx, query, limit)
	}

	return results, nil
}

func (c *ImpotsClient) fallbackImpotsSearch(ctx context.Context, query string, limit int) ([]ImpotsSearchResult, error) {
	return []ImpotsSearchResult{
		{
			Title:       fmt.Sprintf("No results found for: %s", query),
			URL:         fmt.Sprintf("%s/recherche/%s?origin[]=impots&search_filter=Filtrer", c.baseURL, url.QueryEscape(query)),
			Description: "Try modifying your search terms or visit the website directly.",
			Type:        "Info",
		},
	}, nil
}

// ImpotsArticle represents an article from impots.gouv.fr.
type ImpotsArticle struct {
	Title       string
	Content     string
	URL         string
	Type        string
	Description string
}

// GetImpotsArticle retrieves an article from the specified URL.
func (c *ImpotsClient) GetImpotsArticle(ctx context.Context, articleURL string) (*ImpotsArticle, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Validate URL is not empty
	if articleURL == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(articleURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Ensure it's an impots.gouv.fr URL (check for empty host or relative URLs)
	if parsedURL.Host == "" {
		// Handle relative URLs by making them absolute
		articleURL = c.baseURL + articleURL
		parsedURL, err = url.Parse(articleURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL after making absolute: %w", err)
		}
	}

	// Check domain - accept both www.impots.gouv.fr and impots.gouv.fr
	host := strings.ToLower(parsedURL.Host)
	if host != "impots.gouv.fr" && host != "www.impots.gouv.fr" {
		return nil, fmt.Errorf("URL must be from impots.gouv.fr domain, got: %s", parsedURL.Host)
	}

	var article ImpotsArticle
	article.URL = articleURL
	errorChan := make(chan error, 1)

	scraper := c.collector.Clone()

	// Allow URL revisits to prevent "already visited" errors on repeated calls
	scraper.AllowURLRevisit = true

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			scraper = nil
		case <-done:
		}
	}()

	scraper.OnHTML("head", func(e *colly.HTMLElement) {
		if article.Title == "" {
			article.Title = strings.TrimSpace(e.ChildText("title"))
			article.Title = strings.Split(article.Title, " | ")[0]
		}
		if article.Description == "" {
			article.Description = e.ChildAttr("meta[property='og:title']", "content")
		}
	})

	scraper.OnHTML("main, article, div.main-content, div.content", func(e *colly.HTMLElement) {
		if article.Content == "" {
			var contentParts []string

			e.ForEach("h1, h2, h3, p, li, div.fr-callout, div.fr-card__desc", func(_ int, elem *colly.HTMLElement) {
				text := strings.TrimSpace(elem.Text)
				if text != "" &&
					!strings.Contains(text, "javascript") &&
					!strings.Contains(text, "Cookie") &&
					!strings.Contains(text, "Navigation") &&
					len(text) > 10 {
					contentParts = append(contentParts, text)
				}
			})

			article.Content = strings.Join(contentParts, "\n\n")
		}
	})

	scraper.OnHTML("div.fr-breadcrumb", func(e *colly.HTMLElement) {
		breadcrumb := strings.TrimSpace(e.Text)
		if strings.Contains(breadcrumb, "Formulaire") {
			article.Type = "Formulaire"
		} else if strings.Contains(breadcrumb, "Question") {
			article.Type = "Question-Réponse"
		} else {
			article.Type = "Article"
		}
	})

	scraper.OnError(func(r *colly.Response, err error) {
		select {
		case errorChan <- fmt.Errorf("failed to fetch article: %w", err):
		default:
		}
	})

	if err := scraper.Visit(articleURL); err != nil {
		return nil, fmt.Errorf("failed to visit article page: %w", err)
	}

	scraper.Wait()

	select {
	case err := <-errorChan:
		return nil, err
	default:
	}

	if article.Title == "" {
		article.Title = "Article from impots.gouv.fr"
	}
	if article.Content == "" {
		return nil, fmt.Errorf("no content found at URL: %s", articleURL)
	}

	return &article, nil
}

// ImpotsCategoryInfo represents a tax category.
type ImpotsCategoryInfo struct {
	Name        string
	Description string
	URL         string
}

// ListImpotsCategories retrieves available tax service categories.
func (c *ImpotsClient) ListImpotsCategories(ctx context.Context) ([]ImpotsCategoryInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var categories []ImpotsCategoryInfo
	errorChan := make(chan error, 1)

	scraper := c.collector.Clone()

	// Allow URL revisits to prevent "already visited" errors on repeated calls
	scraper.AllowURLRevisit = true

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			scraper = nil
		case <-done:
		}
	}()

	scraper.OnHTML("nav.fr-nav a.fr-nav__link", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.Text)
		href := e.Attr("href")

		if name != "" && href != "" && name != "Accueil" {
			fullURL := e.Request.AbsoluteURL(href)
			for _, cat := range categories {
				if cat.Name == name {
					return
				}
			}

			categories = append(categories, ImpotsCategoryInfo{
				Name:        name,
				Description: fmt.Sprintf("Information fiscale pour %s", strings.ToLower(name)),
				URL:         fullURL,
			})
		}
	})

	scraper.OnError(func(r *colly.Response, err error) {
		select {
		case errorChan <- fmt.Errorf("failed to fetch categories: %w", err):
		default:
		}
	})

	if err := scraper.Visit(c.baseURL + "/particulier"); err != nil {
		return c.getDefaultImpotsCategories(), nil
	}

	scraper.Wait()

	select {
	case <-errorChan:
		return c.getDefaultImpotsCategories(), nil
	default:
	}

	if len(categories) == 0 {
		return c.getDefaultImpotsCategories(), nil
	}

	return categories, nil
}

func (c *ImpotsClient) getDefaultImpotsCategories() []ImpotsCategoryInfo {
	return []ImpotsCategoryInfo{
		{
			Name:        "Particulier",
			Description: "Information fiscale pour les particuliers",
			URL:         c.baseURL + "/particulier",
		},
		{
			Name:        "Professionnel",
			Description: "Information fiscale pour les professionnels",
			URL:         c.baseURL + "/professionnel",
		},
		{
			Name:        "Partenaire",
			Description: "Information pour les partenaires",
			URL:         c.baseURL + "/partenaire",
		},
		{
			Name:        "Collectivité",
			Description: "Information pour les collectivités",
			URL:         c.baseURL + "/collectivite",
		},
		{
			Name:        "International",
			Description: "Information fiscale internationale",
			URL:         c.baseURL + "/international",
		},
	}
}
