// Package client provides HTTP client functionality for service-public.gouv.fr.
package client

import (
	"context"
	"testing"
	"time"
)

// TestSearchProceduresIntegration tests actual search against service-public.gouv.fr
// Run with: go test -v -run TestSearchProceduresIntegration
func TestSearchProceduresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := New(30 * time.Second)
	ctx := context.Background()

	tests := []struct {
		name  string
		query string
		limit int
	}{
		{
			name:  "Search for carte identité",
			query: "carte identité",
			limit: 5,
		},
		{
			name:  "Search for nationalité française",
			query: "nationalité française",
			limit: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := client.SearchProcedures(ctx, tt.query, tt.limit)
			if err != nil {
				t.Fatalf("SearchProcedures() error = %v", err)
			}

			if len(results) == 0 {
				t.Error("Expected at least one result, got none")
			}

			// Log results for inspection
			t.Logf("Found %d results for '%s':", len(results), tt.query)
			for i, r := range results {
				t.Logf("  %d. %s", i+1, r.Title)
				t.Logf("     URL: %s", r.URL)
				if r.Description != "" {
					t.Logf("     Description: %s", r.Description)
				}

				// Validate result fields
				if r.Title == "" {
					t.Errorf("Result %d has empty title", i+1)
				}
				if r.URL == "" {
					t.Errorf("Result %d has empty URL", i+1)
				}
			}
		})
	}
}

// TestGetArticleIntegration tests actual article retrieval
// Run with: go test -v -run TestGetArticleIntegration
func TestGetArticleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := New(30 * time.Second)
	ctx := context.Background()

	tests := []struct {
		name       string
		articleURL string
	}{
		{
			name:       "Article on French nationality by marriage",
			articleURL: "https://www.service-public.gouv.fr/particuliers/vosdroits/F2726",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article, err := client.GetArticle(ctx, tt.articleURL)
			if err != nil {
				t.Fatalf("GetArticle() error = %v", err)
			}

			// Log article for inspection
			t.Logf("Title: %s", article.Title)
			t.Logf("URL: %s", article.URL)
			t.Logf("Content length: %d characters", len(article.Content))
			t.Logf("Content preview: %s...", truncate(article.Content, 200))

			// Validate article fields
			if article.Title == "" {
				t.Error("Article has empty title")
			}
			if article.Content == "" {
				t.Error("Article has empty content")
			}
			if article.URL != tt.articleURL {
				t.Errorf("Article URL = %s, want %s", article.URL, tt.articleURL)
			}
		})
	}
}

// TestListCategoriesIntegration tests category listing
// Run with: go test -v -run TestListCategoriesIntegration
func TestListCategoriesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := New(30 * time.Second)
	ctx := context.Background()

	categories, err := client.ListCategories(ctx)
	if err != nil {
		t.Fatalf("ListCategories() error = %v", err)
	}

	if len(categories) == 0 {
		t.Error("Expected at least one category, got none")
	}

	// Log categories for inspection
	t.Logf("Found %d categories:", len(categories))
	for i, c := range categories {
		t.Logf("  %d. %s: %s", i+1, c.Name, c.Description)

		// Validate category fields
		if c.Name == "" {
			t.Errorf("Category %d has empty name", i+1)
		}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
