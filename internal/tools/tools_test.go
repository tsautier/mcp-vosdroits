package tools

import (
	"testing"
	"time"

	"github.com/guigui42/mcp-vosdroits/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestRegisterTools(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful registration",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mcp.NewServer(
				&mcp.Implementation{
					Name:    "test",
					Version: "v0.0.0",
				},
				nil,
			)

			httpClient := client.New(30 * time.Second)

			if err := registerSearchProcedures(server, httpClient); (err != nil) != tt.wantErr {
				t.Errorf("registerSearchProcedures() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := registerGetArticle(server, httpClient); (err != nil) != tt.wantErr {
				t.Errorf("registerGetArticle() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := registerListCategories(server, httpClient); (err != nil) != tt.wantErr {
				t.Errorf("registerListCategories() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchProcedures(t *testing.T) {
	tests := []struct {
		name    string
		input   SearchProceduresInput
		wantErr bool
	}{
		{
			name: "valid search",
			input: SearchProceduresInput{
				Query: "carte d'identité",
				Limit: 10,
			},
			wantErr: false,
		},
		{
			name: "empty query",
			input: SearchProceduresInput{
				Query: "",
				Limit: 10,
			},
			wantErr: true,
		},
		{
			name: "default limit",
			input: SearchProceduresInput{
				Query: "passeport",
			},
			wantErr: false,
		},
	}

	httpClient := client.New(30 * time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mcp.NewServer(
				&mcp.Implementation{
					Name:    "test",
					Version: "v0.0.0",
				},
				nil,
			)

			if err := registerSearchProcedures(server, httpClient); err != nil {
				t.Fatalf("Failed to register tool: %v", err)
			}

			// Note: This is a simplified test. In practice, you'd invoke the tool
			// through the MCP server's tool handling mechanism.
			// For now, we're just testing registration succeeds.
		})
	}
}
