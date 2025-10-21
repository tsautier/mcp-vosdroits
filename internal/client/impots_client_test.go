package client

import (
	"context"
	"testing"
	"time"
)

func TestNewImpotsClient(t *testing.T) {
	client := NewImpotsClient(30 * time.Second)
	if client == nil {
		t.Fatal("NewImpotsClient returned nil")
	}
	if client.baseURL != "https://www.impots.gouv.fr" {
		t.Errorf("expected baseURL to be https://www.impots.gouv.fr, got %s", client.baseURL)
	}
	if client.timeout != 30*time.Second {
		t.Errorf("expected timeout to be 30s, got %v", client.timeout)
	}
}

func TestSearchImpots_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		query string
		limit int
	}{
		{
			name:  "empty query",
			query: "",
			limit: 10,
		},
		{
			name:  "negative limit",
			query: "formulaire 2042",
			limit: -1,
		},
		{
			name:  "excessive limit",
			query: "formulaire 2042",
			limit: 200,
		},
	}

	client := NewImpotsClient(30 * time.Second)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.SearchImpots(ctx, tt.query, tt.limit)
			// Empty query should return fallback results, not error
			// Invalid limits should be normalized, not error
			if err != nil && tt.query != "" {
				t.Errorf("SearchImpots() unexpected error: %v", err)
			}
		})
	}
}

func TestGetImpotsArticle_InvalidURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			url:     "not a url",
			wantErr: true,
		},
		{
			name:    "wrong domain",
			url:     "https://www.google.com",
			wantErr: true,
		},
		{
			name:    "valid domain",
			url:     "https://www.impots.gouv.fr/formulaire/2042/declaration-des-revenus",
			wantErr: false,
		},
	}

	client := NewImpotsClient(30 * time.Second)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetImpotsArticle(ctx, tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImpotsArticle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListImpotsCategories_DefaultCategories(t *testing.T) {
	client := NewImpotsClient(30 * time.Second)
	categories := client.getDefaultImpotsCategories()

	if len(categories) == 0 {
		t.Fatal("getDefaultImpotsCategories returned empty slice")
	}

	expectedCategories := []string{"Particulier", "Professionnel", "Partenaire", "Collectivité", "International"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("expected %d categories, got %d", len(expectedCategories), len(categories))
	}

	for i, cat := range categories {
		if cat.Name != expectedCategories[i] {
			t.Errorf("expected category %s, got %s", expectedCategories[i], cat.Name)
		}
		if cat.Description == "" {
			t.Errorf("category %s has empty description", cat.Name)
		}
		if cat.URL == "" {
			t.Errorf("category %s has empty URL", cat.Name)
		}
	}
}

func TestFallbackImpotsSearch(t *testing.T) {
	client := NewImpotsClient(30 * time.Second)
	ctx := context.Background()

	results, err := client.fallbackImpotsSearch(ctx, "test query", 10)
	if err != nil {
		t.Fatalf("fallbackImpotsSearch returned error: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("fallbackImpotsSearch returned empty results")
	}

	if results[0].Title == "" {
		t.Error("fallback result has empty title")
	}

	if results[0].URL == "" {
		t.Error("fallback result has empty URL")
	}
}
