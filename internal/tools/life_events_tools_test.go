package tools

import (
	"testing"
	"time"

	"github.com/guigui42/mcp-vosdroits/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestRegisterListLifeEvents(t *testing.T) {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		},
		nil,
	)

	httpClient := client.New(10 * time.Second)

	err := registerListLifeEvents(server, httpClient)
	if err != nil {
		t.Fatalf("registerListLifeEvents failed: %v", err)
	}
}

func TestRegisterGetLifeEventDetails(t *testing.T) {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		},
		nil,
	)

	httpClient := client.New(10 * time.Second)

	err := registerGetLifeEventDetails(server, httpClient)
	if err != nil {
		t.Fatalf("registerGetLifeEventDetails failed: %v", err)
	}
}

func TestListLifeEventsInput(t *testing.T) {
	// Test that the input struct is properly defined
	input := ListLifeEventsInput{}
	_ = input // Just verify it compiles
}

func TestGetLifeEventDetailsInput(t *testing.T) {
	// Test that the input struct is properly defined
	input := GetLifeEventDetailsInput{
		URL: "https://www.service-public.gouv.fr/particuliers/vosdroits/F16225",
	}

	if input.URL == "" {
		t.Error("Expected URL to be set")
	}
}

func TestLifeEventStructs(t *testing.T) {
	// Test that output structs are properly defined
	event := LifeEventInfo{
		Title: "J'attends un enfant",
		URL:   "https://www.service-public.gouv.fr/particuliers/vosdroits/F16225",
	}

	if event.Title == "" {
		t.Error("Expected title to be set")
	}

	section := LifeEventSectionOutput{
		Title:   "Santé",
		Content: "Health information",
	}

	if section.Title == "" {
		t.Error("Expected section title to be set")
	}

	output := GetLifeEventDetailsOutput{
		Title:        "J'attends un enfant",
		URL:          "https://www.service-public.gouv.fr/particuliers/vosdroits/F16225",
		Introduction: "Introduction text",
		Sections:     []LifeEventSectionOutput{section},
	}

	if len(output.Sections) != 1 {
		t.Errorf("Expected 1 section, got %d", len(output.Sections))
	}
}

// Integration test - requires network access
func TestListLifeEventsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		},
		nil,
	)

	httpClient := client.New(30 * time.Second)

	err := registerListLifeEvents(server, httpClient)
	if err != nil {
		t.Fatalf("registerListLifeEvents failed: %v", err)
	}

	// Note: Actual MCP tool invocation would require a full server setup
	// This test just verifies registration doesn't error
}

// Integration test - requires network access
func TestGetLifeEventDetailsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		},
		nil,
	)

	httpClient := client.New(30 * time.Second)

	err := registerGetLifeEventDetails(server, httpClient)
	if err != nil {
		t.Fatalf("registerGetLifeEventDetails failed: %v", err)
	}

	// Note: Actual MCP tool invocation would require a full server setup
	// This test just verifies registration doesn't error
}

func TestLifeEventsToolsRegistration(t *testing.T) {
	// Test that both tools can be registered together
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		},
		nil,
	)

	httpClient := client.New(10 * time.Second)

	if err := registerListLifeEvents(server, httpClient); err != nil {
		t.Errorf("Failed to register list_life_events: %v", err)
	}

	if err := registerGetLifeEventDetails(server, httpClient); err != nil {
		t.Errorf("Failed to register get_life_event_details: %v", err)
	}
}
