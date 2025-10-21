package client

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	timeout := 30 * time.Second
	client := New(timeout)

	if client == nil {
		t.Fatal("New() returned nil")
	}

	if client.collector == nil {
		t.Error("collector should not be nil")
	}

	if client.timeout != timeout {
		t.Errorf("timeout = %v, want %v", client.timeout, timeout)
	}

	if client.baseURL == "" {
		t.Error("baseURL should not be empty")
	}
}

func TestSearchProcedures(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		limit   int
		wantErr bool
	}{
		{
			name:    "valid search",
			query:   "carte d'identité",
			limit:   10,
			wantErr: false, // Should return fallback results
		},
		{
			name:    "empty query",
			query:   "",
			limit:   10,
			wantErr: false, // Will use fallback search
		},
		{
			name:    "limit too high",
			query:   "passeport",
			limit:   200,
			wantErr: false, // Will be clamped to 10
		},
	}

	client := New(30 * time.Second)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := client.SearchProcedures(ctx, tt.query, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchProcedures() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && results == nil {
				t.Error("SearchProcedures() returned nil results")
			}

			// Should always return at least fallback results
			if len(results) == 0 {
				t.Error("SearchProcedures() returned empty results, expected at least fallback")
			}
		})
	}
}

func TestSearchProceduresContextCancellation(t *testing.T) {
	client := New(30 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.SearchProcedures(ctx, "test", 10)
	if err == nil {
		t.Error("SearchProcedures() should return error when context is cancelled")
	}
}

func TestGetArticle(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name: "valid service-public.gouv.fr url",
			url:  "https://www.service-public.gouv.fr/particuliers/vosdroits/F1234",
			// May or may not error depending on whether the page exists and has content
			wantErr: false,
		},
		{
			name:    "invalid url",
			url:     "not-a-url",
			wantErr: true,
		},
		{
			name:    "wrong domain",
			url:     "https://example.com/page",
			wantErr: true,
		},
	}

	client := New(5 * time.Second) // Shorter timeout for tests
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article, err := client.GetArticle(ctx, tt.url)
			if (err != nil) != tt.wantErr {
				// For the service-public test, we just log the result
				if tt.name == "valid service-public.gouv.fr url" {
					t.Logf("GetArticle() returned article=%v, err=%v", article != nil, err)
					return
				}
				t.Errorf("GetArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && article == nil {
				t.Error("GetArticle() returned nil article")
			}
		})
	}
}

func TestListCategories(t *testing.T) {
	client := New(30 * time.Second)
	ctx := context.Background()

	categories, err := client.ListCategories(ctx)
	if err != nil {
		t.Errorf("ListCategories() error = %v", err)
		return
	}

	if categories == nil {
		t.Error("ListCategories() returned nil")
	}

	if len(categories) == 0 {
		t.Error("ListCategories() returned empty list")
	}
}
