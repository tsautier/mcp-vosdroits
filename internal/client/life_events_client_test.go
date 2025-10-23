package client

import (
	"context"
	"testing"
	"time"
)

func TestListLifeEvents(t *testing.T) {
	// Create client with reasonable timeout
	c := New(30 * time.Second)
	ctx := context.Background()

	events, err := c.ListLifeEvents(ctx)
	if err != nil {
		t.Fatalf("ListLifeEvents failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected at least one life event, got none")
	}

	// Verify that we got some common life events
	expectedEvents := []string{
		"J'attends un enfant",
		"Je déménage",
		"Un proche est décédé",
	}

	foundEvents := make(map[string]bool)
	for _, event := range events {
		// Check that each event has required fields
		if event.Title == "" {
			t.Error("Event missing title")
		}
		if event.URL == "" {
			t.Error("Event missing URL")
		}

		// Track which expected events we found
		for _, expected := range expectedEvents {
			if event.Title == expected {
				foundEvents[expected] = true
			}
		}
	}

	// Verify we found at least some of the expected events
	if len(foundEvents) == 0 {
		t.Errorf("Did not find any of the expected life events: %v", expectedEvents)
		t.Logf("Found events: %v", events)
	}
}

func TestGetLifeEventDetails(t *testing.T) {
	c := New(30 * time.Second)
	ctx := context.Background()

	// Test with "J'attends un enfant" URL
	testURL := "https://www.service-public.gouv.fr/particuliers/vosdroits/F16225"

	details, err := c.GetLifeEventDetails(ctx, testURL)
	if err != nil {
		t.Fatalf("GetLifeEventDetails failed: %v", err)
	}

	// Verify basic fields
	if details.Title == "" {
		t.Error("Expected title, got empty string")
	}

	if details.URL != testURL {
		t.Errorf("Expected URL %s, got %s", testURL, details.URL)
	}

	// Should have either an introduction or sections
	if details.Introduction == "" && len(details.Sections) == 0 {
		t.Error("Expected either introduction or sections, got neither")
	}

	// If we have sections, verify they have content
	if len(details.Sections) > 0 {
		for i, section := range details.Sections {
			if section.Title == "" {
				t.Errorf("Section %d missing title", i)
			}
			if section.Content == "" {
				t.Errorf("Section %d (%s) missing content", i, section.Title)
			}
		}
	}

	t.Logf("Retrieved life event: %s", details.Title)
	t.Logf("Number of sections: %d", len(details.Sections))
	if len(details.Sections) > 0 {
		t.Logf("First section: %s", details.Sections[0].Title)
	}
}

func TestGetLifeEventDetailsInvalidURL(t *testing.T) {
	c := New(5 * time.Second)
	ctx := context.Background()

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
			name:    "invalid domain",
			url:     "https://example.com/test",
			wantErr: true,
		},
		{
			name:    "malformed URL",
			url:     "not-a-url",
			wantErr: true,
		},
		{
			name:    "category page (N-prefix)",
			url:     "https://www.service-public.gouv.fr/particuliers/vosdroits/N19808",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.GetLifeEventDetails(ctx, tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLifeEventDetails() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetLifeEventDetailsContextCancellation(t *testing.T) {
	c := New(30 * time.Second)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.GetLifeEventDetails(ctx, "https://www.service-public.gouv.fr/particuliers/vosdroits/F16225")
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}
}
