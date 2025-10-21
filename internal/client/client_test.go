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

	if client.httpClient.Timeout != timeout {
		t.Errorf("timeout = %v, want %v", client.httpClient.Timeout, timeout)
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
			wantErr: false,
		},
		{
			name:    "empty query",
			query:   "",
			limit:   10,
			wantErr: false, // Current implementation doesn't validate
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
			name:    "valid url",
			url:     "https://www.service-public.gouv.fr/particuliers/vosdroits/F1234",
			wantErr: false,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: false, // Current implementation doesn't validate
		},
	}

	client := New(30 * time.Second)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article, err := client.GetArticle(ctx, tt.url)
			if (err != nil) != tt.wantErr {
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
