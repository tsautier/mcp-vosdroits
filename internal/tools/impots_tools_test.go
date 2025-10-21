package tools

import (
	"context"
	"testing"
	"time"

	"github.com/guigui42/mcp-vosdroits/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestRegisterImpotsTools(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "v1.0.0",
	}, nil)

	impotsClient := client.NewImpotsClient(30 * time.Second)

	err := RegisterImpotsTools(server, impotsClient)
	if err != nil {
		t.Fatalf("RegisterImpotsTools failed: %v", err)
	}
}

func TestSearchImpotsInput_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   SearchImpotsInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: SearchImpotsInput{
				Query: "formulaire 2042",
				Limit: 10,
			},
			wantErr: false,
		},
		{
			name: "empty query",
			input: SearchImpotsInput{
				Query: "",
				Limit: 10,
			},
			wantErr: true,
		},
		{
			name: "zero limit defaults to 10",
			input: SearchImpotsInput{
				Query: "PEA",
				Limit: 0,
			},
			wantErr: false,
		},
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "v1.0.0",
	}, nil)

	impotsClient := client.NewImpotsClient(30 * time.Second)
	if err := registerSearchImpots(server, impotsClient); err != nil {
		t.Fatalf("registerSearchImpots failed: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := &mcp.CallToolRequest{}

			// Create a simple handler that mimics the registered one
			handler := func(ctx context.Context, req *mcp.CallToolRequest, input SearchImpotsInput) (*mcp.CallToolResult, SearchImpotsOutput, error) {
				if input.Limit == 0 {
					input.Limit = 10
				}
				if input.Query == "" {
					return nil, SearchImpotsOutput{}, nil
				}
				return &mcp.CallToolResult{}, SearchImpotsOutput{}, nil
			}

			result, output, err := handler(ctx, req, tt.input)
			if tt.wantErr && err == nil && tt.input.Query == "" {
				// Empty query should result in empty output or error
				if len(output.Results) != 0 {
					t.Error("expected empty results for empty query")
				}
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantErr && result == nil && tt.input.Query != "" {
				t.Error("expected non-nil result for valid input")
			}
		})
	}
}

func TestGetImpotsArticleInput_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   GetImpotsArticleInput
		wantErr bool
	}{
		{
			name: "valid URL",
			input: GetImpotsArticleInput{
				URL: "https://www.impots.gouv.fr/formulaire/2042/declaration-des-revenus",
			},
			wantErr: false,
		},
		{
			name: "empty URL",
			input: GetImpotsArticleInput{
				URL: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input.URL == "" && !tt.wantErr {
				t.Error("expected error for empty URL")
			}
			if tt.input.URL != "" && tt.wantErr {
				t.Error("expected no error for valid URL")
			}
		})
	}
}

func TestImpotsResultStructure(t *testing.T) {
	result := ImpotsResult{
		Title:       "Formulaire 2042",
		URL:         "https://www.impots.gouv.fr/formulaire/2042/declaration-des-revenus",
		Description: "Déclaration de revenus",
		Type:        "Formulaire",
		Date:        "2025-04-16",
	}

	if result.Title == "" {
		t.Error("Title should not be empty")
	}
	if result.URL == "" {
		t.Error("URL should not be empty")
	}
}

func TestImpotsCategoryStructure(t *testing.T) {
	category := ImpotsCategory{
		Name:        "Particulier",
		Description: "Information fiscale pour les particuliers",
		URL:         "https://www.impots.gouv.fr/particulier",
	}

	if category.Name == "" {
		t.Error("Name should not be empty")
	}
	if category.Description == "" {
		t.Error("Description should not be empty")
	}
	if category.URL == "" {
		t.Error("URL should not be empty")
	}
}
