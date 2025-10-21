package client

import (
	"context"
	"testing"
	"time"
)

// Integration tests for impots.gouv.fr scraping

func TestSearchImpotsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewImpotsClient(30 * time.Second)
	ctx := context.Background()

	tests := []struct {
		name          string
		query         string
		limit         int
		expectResults bool
	}{
		{
			name:          "Search for formulaire 2042",
			query:         "formulaire 2042",
			limit:         5,
			expectResults: true,
		},
		{
			name:          "Search for PEA",
			query:         "PEA",
			limit:         5,
			expectResults: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := client.SearchImpots(ctx, tt.query, tt.limit)
			if err != nil {
				t.Fatalf("SearchImpots failed: %v", err)
			}

			if tt.expectResults && len(results) == 0 {
				t.Error("Expected results but got none")
			}

			if len(results) > 0 {
				t.Logf("Found %d results for '%s':", len(results), tt.query)
				for i, result := range results {
					t.Logf("  %d. %s", i+1, result.Title)
					t.Logf("     URL: %s", result.URL)
					if result.Type != "" {
						t.Logf("     Type: %s", result.Type)
					}
					if result.Date != "" {
						t.Logf("     Date: %s", result.Date)
					}
				}

				// Verify result structure
				firstResult := results[0]
				if firstResult.Title == "" {
					t.Error("First result has empty title")
				}
				if firstResult.URL == "" {
					t.Error("First result has empty URL")
				}
			}
		})
	}
}

func TestGetImpotsArticleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewImpotsClient(30 * time.Second)
	ctx := context.Background()

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "Formulaire 2042",
			url:  "https://www.impots.gouv.fr/formulaire/2042/declaration-des-revenus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article, err := client.GetImpotsArticle(ctx, tt.url)
			if err != nil {
				t.Fatalf("GetImpotsArticle failed: %v", err)
			}

			if article.Title == "" {
				t.Error("Article has empty title")
			}
			if article.Content == "" {
				t.Error("Article has empty content")
			}
			if article.URL != tt.url {
				t.Errorf("Expected URL %s, got %s", tt.url, article.URL)
			}

			t.Logf("Title: %s", article.Title)
			t.Logf("URL: %s", article.URL)
			t.Logf("Type: %s", article.Type)
			t.Logf("Content length: %d characters", len(article.Content))
			if len(article.Content) > 200 {
				t.Logf("Content preview: %s...", article.Content[:200])
			}
		})
	}
}

func TestListImpotsCategoriesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewImpotsClient(30 * time.Second)
	ctx := context.Background()

	categories, err := client.ListImpotsCategories(ctx)
	if err != nil {
		t.Fatalf("ListImpotsCategories failed: %v", err)
	}

	if len(categories) == 0 {
		t.Fatal("Expected at least one category")
	}

	t.Logf("Found %d categories:", len(categories))
	for i, cat := range categories {
		t.Logf("  %d. %s: %s", i+1, cat.Name, cat.Description)
		t.Logf("     URL: %s", cat.URL)

		if cat.Name == "" {
			t.Error("Category has empty name")
		}
		if cat.Description == "" {
			t.Error("Category has empty description")
		}
		if cat.URL == "" {
			t.Error("Category has empty URL")
		}
	}
}
