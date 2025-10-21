// Package client provides HTTP client functionality for service-public.gouv.fr.
package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Client handles HTTP requests to service-public.gouv.fr.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// New creates a new Client with the specified timeout.
func New(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://www.service-public.gouv.fr",
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

	// TODO: Implement actual HTTP request to service-public.gouv.fr search API
	// This is a placeholder implementation
	return []SearchResult{
		{
			Title:       fmt.Sprintf("Example result for: %s", query),
			URL:         fmt.Sprintf("%s/particuliers/vosdroits/F1234", c.baseURL),
			Description: "This is a placeholder result. Implement actual search logic.",
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
func (c *Client) GetArticle(ctx context.Context, url string) (*Article, error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// TODO: Implement actual HTTP request and HTML parsing
	// This is a placeholder implementation
	return &Article{
		Title:   "Example Article",
		Content: "This is placeholder content. Implement actual article retrieval and parsing.",
		URL:     url,
	}, nil
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

	// TODO: Implement actual category listing
	// This is a placeholder implementation
	return []CategoryInfo{
		{
			Name:        "Particuliers",
			Description: "Information and procedures for individuals",
		},
		{
			Name:        "Professionnels",
			Description: "Information and procedures for professionals",
		},
		{
			Name:        "Associations",
			Description: "Information and procedures for associations",
		},
	}, nil
}
