// Package client provides HTTP client functionality for service-public.gouv.fr.
package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// Client handles HTTP requests to service-public.gouv.fr using Colly for web scraping.
type Client struct {
	collector *colly.Collector
	baseURL   string
	timeout   time.Duration
}

// New creates a new Client with the specified timeout.
func New(timeout time.Duration) *Client {
	// Create a new Colly collector with configuration
	c := colly.NewCollector(
		colly.AllowedDomains(
			"www.service-public.gouv.fr",
			"service-public.gouv.fr",
			"www.service-public.fr",
			"service-public.fr",
		),
		colly.UserAgent("VosDroits-MCP-Server/1.0"),
		colly.Async(false),
	)

	// Set timeout
	c.SetRequestTimeout(timeout)

	// Configure rate limiting to be respectful
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*.service-public.*",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})
	if err != nil {
		// Fallback without rate limiting if it fails
		// This shouldn't happen, but we handle it gracefully
	}

	return &Client{
		collector: c,
		baseURL:   "https://www.service-public.gouv.fr",
		timeout:   timeout,
	}
}

// SearchResult represents a search result.
type SearchResult struct {
	Title       string
	URL         string
	Description string
}

// SearchProcedures searches for procedures matching the query.
func (c *Client) SearchProcedures(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Validate limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	var results []SearchResult
	errorChan := make(chan error, 1)

	// Clone the collector to avoid conflicts
	scraper := c.collector.Clone()

	// Set up context cancellation
	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			// Context cancelled, close the collector
			scraper = nil
		case <-done:
			// Operation completed normally
		}
	}()

	// Handle search results - service-public.gouv.fr uses <li> with id pattern "result_*"
	scraper.OnHTML("li[id^='result_']", func(e *colly.HTMLElement) {
		if len(results) >= limit {
			return
		}

		// Extract URL and title from the link
		href := e.ChildAttr("a.fr-link", "href")
		if href == "" {
			return
		}

		// Make URL absolute
		fullURL := e.Request.AbsoluteURL(href)

		// Extract title - get the innermost span text to avoid duplication
		var title string
		e.ForEach("a.fr-link span span", func(_ int, elem *colly.HTMLElement) {
			if title == "" {
				title = strings.TrimSpace(elem.Text)
			}
		})
		// Fallback if no nested span found
		if title == "" {
			title = strings.TrimSpace(e.ChildText("a.fr-link"))
		}

		// Extract description if available (some results may have descriptions)
		description := strings.TrimSpace(e.ChildText(".sp-description, .description"))

		if title != "" && fullURL != "" {
			results = append(results, SearchResult{
				Title:       title,
				URL:         fullURL,
				Description: description,
			})
		}
	})

	// Handle errors
	scraper.OnError(func(r *colly.Response, err error) {
		select {
		case errorChan <- fmt.Errorf("scraping error: %w", err):
		default:
		}
	})

	// Build search URL - service-public.gouv.fr uses /particuliers/recherche
	searchURL := fmt.Sprintf("%s/particuliers/recherche?keyword=%s", c.baseURL, url.QueryEscape(query))

	// Visit the search page
	if err := scraper.Visit(searchURL); err != nil {
		// If the search fails, return fallback results
		return c.fallbackSearch(ctx, query, limit)
	}

	// Wait for scraping to complete
	scraper.Wait()

	// Check for errors
	select {
	case err := <-errorChan:
		if len(results) == 0 {
			return nil, err
		}
		// If we got some results, return them despite the error
	default:
	}

	// If no results found, try alternative search
	if len(results) == 0 {
		return c.fallbackSearch(ctx, query, limit)
	}

	return results, nil
}

// fallbackSearch provides a fallback search method
func (c *Client) fallbackSearch(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	// For now, return a helpful message
	// In a real implementation, we could try different search strategies
	return []SearchResult{
		{
			Title:       fmt.Sprintf("No results found for: %s", query),
			URL:         fmt.Sprintf("%s/particuliers/recherche?keyword=%s", c.baseURL, url.QueryEscape(query)),
			Description: "Try modifying your search terms or visit the website directly.",
		},
	}, nil
}

// Article represents an article from service-public.gouv.fr.
type Article struct {
	Title   string
	Content string
	URL     string
}

// GetArticle retrieves an article from the specified URL.
func (c *Client) GetArticle(ctx context.Context, articleURL string) (*Article, error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Validate URL is not empty
	if articleURL == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}

	// Validate URL
	parsedURL, err := url.Parse(articleURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Ensure it's a service-public.gouv.fr URL (check for empty host or relative URLs)
	if parsedURL.Host == "" {
		// Handle relative URLs by making them absolute
		articleURL = c.baseURL + articleURL
		parsedURL, err = url.Parse(articleURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL after making absolute: %w", err)
		}
	}

	// Check domain - accept service-public.gouv.fr and service-public.fr (both with and without www)
	host := strings.ToLower(parsedURL.Host)
	validDomains := []string{
		"service-public.gouv.fr",
		"www.service-public.gouv.fr",
		"service-public.fr",
		"www.service-public.fr",
	}

	domainValid := false
	for _, validDomain := range validDomains {
		if host == validDomain {
			domainValid = true
			break
		}
	}

	if !domainValid {
		return nil, fmt.Errorf("URL must be from service-public.gouv.fr or service-public.fr domain, got: %s", parsedURL.Host)
	}

	var article Article
	article.URL = articleURL
	errorChan := make(chan error, 1)

	// Clone the collector
	scraper := c.collector.Clone()

	// Set up context cancellation
	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			scraper = nil
		case <-done:
		}
	}()

	// Extract article title - service-public.gouv.fr uses h1#titlePage
	scraper.OnHTML("h1#titlePage", func(e *colly.HTMLElement) {
		if article.Title == "" {
			article.Title = strings.TrimSpace(e.Text)
		}
	})

	// Extract main content - service-public.gouv.fr has content in article.article
	scraper.OnHTML("article.article", func(e *colly.HTMLElement) {
		if article.Content == "" {
			// Get all relevant content elements
			var contentParts []string

			// Add introduction text
			intro := strings.TrimSpace(e.ChildText("div#intro p.fr-text--lg"))
			if intro != "" {
				contentParts = append(contentParts, intro)
			}

			// Add main content sections
			e.ForEach("h2, h3, p[data-test='contenu-texte'], .fr-text--lg, div.sp-section p, div.fr-callout p", func(_ int, elem *colly.HTMLElement) {
				text := strings.TrimSpace(elem.Text)
				// Filter out navigation, scripts, and empty content
				if text != "" &&
					!strings.Contains(text, "javascript") &&
					!strings.Contains(text, "Votre situation") &&
					!strings.Contains(text, "Abonnement") &&
					len(text) > 10 {
					contentParts = append(contentParts, text)
				}
			})

			article.Content = strings.Join(contentParts, "\n\n")
		}
	})

	// Handle HTTP responses to detect 404s and other errors
	scraper.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 404 {
			select {
			case errorChan <- fmt.Errorf("article not found (404) at URL: %s", r.Request.URL):
			default:
			}
		} else if r.StatusCode >= 400 {
			select {
			case errorChan <- fmt.Errorf("HTTP error %d at URL: %s", r.StatusCode, r.Request.URL):
			default:
			}
		}
	})

	// Handle errors
	scraper.OnError(func(r *colly.Response, err error) {
		statusMsg := ""
		if r != nil {
			statusMsg = fmt.Sprintf(" (HTTP %d)", r.StatusCode)
		}
		select {
		case errorChan <- fmt.Errorf("failed to fetch article%s: %w", statusMsg, err):
		default:
		}
	})

	// Visit the article page
	if err := scraper.Visit(articleURL); err != nil {
		return nil, fmt.Errorf("failed to visit article page: %w", err)
	}

	// Wait for scraping to complete
	scraper.Wait()

	// Check for errors
	select {
	case err := <-errorChan:
		return nil, err
	default:
	}

	// Validate that we got content
	if article.Title == "" {
		article.Title = "Article from service-public.gouv.fr"
	}
	if article.Content == "" {
		return nil, fmt.Errorf("no content found at URL: %s", articleURL)
	}

	return &article, nil
}

// CategoryInfo represents a service category.
type CategoryInfo struct {
	Name        string
	Description string
}

// ListCategories retrieves available service categories.
func (c *Client) ListCategories(ctx context.Context) ([]CategoryInfo, error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var categories []CategoryInfo
	errorChan := make(chan error, 1)

	// Clone the collector
	scraper := c.collector.Clone()

	// Set up context cancellation
	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			scraper = nil
		case <-done:
		}
	}()

	// Extract main category sections from the footer theme list
	scraper.OnHTML("ul.sp-theme-list li a.fr-footer__top-link", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.Text)
		href := e.Attr("href")

		// Filter for main categories (these are the primary themes)
		if name != "" && len(name) > 3 && strings.Contains(href, "/particuliers/vosdroits/N") {
			// Avoid duplicates
			for _, cat := range categories {
				if cat.Name == name {
					return
				}
			}

			categories = append(categories, CategoryInfo{
				Name:        name,
				Description: fmt.Sprintf("Information and procedures for %s", strings.ToLower(name)),
			})
		}
	})

	// Handle errors
	scraper.OnError(func(r *colly.Response, err error) {
		select {
		case errorChan <- fmt.Errorf("failed to fetch categories: %w", err):
		default:
		}
	})

	// Visit the home page for particuliers
	if err := scraper.Visit(c.baseURL + "/particuliers"); err != nil {
		// Fallback to default categories if scraping fails
		return c.getDefaultCategories(), nil
	}

	// Wait for scraping to complete
	scraper.Wait()

	// Check for errors
	select {
	case <-errorChan:
		// If there was an error, return default categories
		return c.getDefaultCategories(), nil
	default:
	}

	// If no categories found, return default ones
	if len(categories) == 0 {
		return c.getDefaultCategories(), nil
	}

	return categories, nil
}

// getDefaultCategories returns the standard categories
func (c *Client) getDefaultCategories() []CategoryInfo {
	return []CategoryInfo{
		{
			Name:        "Particuliers",
			Description: "Information and procedures for individuals - family, health, work, housing, etc.",
		},
		{
			Name:        "Professionnels",
			Description: "Information and procedures for professionals - business creation, taxes, employees, etc.",
		},
		{
			Name:        "Associations",
			Description: "Information and procedures for associations - creation, financing, management, etc.",
		},
	}
}
