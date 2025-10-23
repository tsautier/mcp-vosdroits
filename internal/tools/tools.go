// Package tools provides MCP tool implementations for searching French public service information.
package tools

import (
	"context"
	"fmt"

	"github.com/guigui42/mcp-vosdroits/internal/client"
	"github.com/guigui42/mcp-vosdroits/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all available MCP tools with the server.
func RegisterTools(server *mcp.Server, cfg *config.Config) error {
	// Create HTTP client for service-public.gouv.fr
	httpClient := client.New(cfg.HTTPTimeout)

	// Register search_procedures tool
	if err := registerSearchProcedures(server, httpClient); err != nil {
		return fmt.Errorf("failed to register search_procedures: %w", err)
	}

	// Register get_article tool
	if err := registerGetArticle(server, httpClient); err != nil {
		return fmt.Errorf("failed to register get_article: %w", err)
	}

	// Register list_categories tool
	if err := registerListCategories(server, httpClient); err != nil {
		return fmt.Errorf("failed to register list_categories: %w", err)
	}

	// Create HTTP client for impots.gouv.fr
	impotsClient := client.NewImpotsClient(cfg.HTTPTimeout)

	// Register impots.gouv.fr tools
	if err := RegisterImpotsTools(server, impotsClient); err != nil {
		return fmt.Errorf("failed to register impots tools: %w", err)
	}

	return nil
}

// SearchProceduresInput defines the input schema for search_procedures.
type SearchProceduresInput struct {
	Query string `json:"query" jsonschema:"Search query for procedures (e.g. 'carte d'identité' or 'passport renewal')"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of results to return (1-100), default 10"`
}

// SearchProceduresOutput defines the output schema for search_procedures.
type SearchProceduresOutput struct {
	Results []ProcedureResult `json:"results" jsonschema:"List of matching procedures. Each result includes a URL that can be used with the get_article tool to retrieve full details."`
}

// ProcedureResult represents a single procedure search result.
type ProcedureResult struct {
	Title       string `json:"title" jsonschema:"Title of the procedure"`
	URL         string `json:"url" jsonschema:"URL to the procedure page. Use this URL with the get_article tool to retrieve complete information."`
	Description string `json:"description" jsonschema:"Brief summary. For full content, requirements, and instructions, use get_article with this URL."`
}

func registerSearchProcedures(server *mcp.Server, httpClient *client.Client) error {
	tool := &mcp.Tool{
		Name:        "search_procedures",
		Description: "Search for procedures on service-public.gouv.fr. Returns URLs and brief descriptions. For detailed information about any procedure, use the get_article tool with the returned URLs.",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input SearchProceduresInput) (*mcp.CallToolResult, SearchProceduresOutput, error) {
		// Set default limit
		if input.Limit == 0 {
			input.Limit = 10
		}

		// Validate input
		if input.Query == "" {
			return nil, SearchProceduresOutput{}, fmt.Errorf("query cannot be empty")
		}

		// TODO: Implement actual search logic using client
		results, err := httpClient.SearchProcedures(ctx, input.Query, input.Limit)
		if err != nil {
			return nil, SearchProceduresOutput{}, fmt.Errorf("search failed: %w", err)
		}

		// Convert client results to output format
		output := SearchProceduresOutput{
			Results: make([]ProcedureResult, len(results)),
		}
		for i, r := range results {
			output.Results[i] = ProcedureResult{
				Title:       r.Title,
				URL:         r.URL,
				Description: r.Description,
			}
		}

		// Create helpful message that encourages follow-up
		message := fmt.Sprintf("Found %d procedures matching '%s'. ", len(results), input.Query)
		if len(results) > 0 {
			message += "Use the get_article tool with any of the returned URLs to retrieve complete details about a specific procedure."
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: message,
				},
			},
		}, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}

// GetArticleInput defines the input schema for get_article.
type GetArticleInput struct {
	URL string `json:"url" jsonschema:"URL of the article to retrieve (typically from search_procedures results)"`
}

// GetArticleOutput defines the output schema for get_article.
type GetArticleOutput struct {
	Title   string `json:"title" jsonschema:"Title of the article"`
	Content string `json:"content" jsonschema:"Full content of the article"`
	URL     string `json:"url" jsonschema:"URL of the article"`
}

func registerGetArticle(server *mcp.Server, httpClient *client.Client) error {
	tool := &mcp.Tool{
		Name:        "get_article",
		Description: "Retrieve detailed article content from service-public.gouv.fr URLs. Use this after search_procedures to get complete information about a specific procedure, including full text, requirements, and step-by-step instructions.",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input GetArticleInput) (*mcp.CallToolResult, GetArticleOutput, error) {
		if input.URL == "" {
			return nil, GetArticleOutput{}, fmt.Errorf("url cannot be empty")
		}

		// TODO: Implement actual article retrieval using client
		article, err := httpClient.GetArticle(ctx, input.URL)
		if err != nil {
			return nil, GetArticleOutput{}, fmt.Errorf("failed to get article: %w", err)
		}

		output := GetArticleOutput{
			Title:   article.Title,
			Content: article.Content,
			URL:     article.URL,
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Retrieved article: %s", article.Title),
				},
			},
		}, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}

// ListCategoriesOutput defines the output schema for list_categories.
type ListCategoriesOutput struct {
	Categories []Category `json:"categories" jsonschema:"List of available categories"`
}

// Category represents a service category.
type Category struct {
	Name        string `json:"name" jsonschema:"Name of the category"`
	Description string `json:"description" jsonschema:"Description of the category"`
}

type ListCategoriesInput struct{}

func registerListCategories(server *mcp.Server, httpClient *client.Client) error {
	tool := &mcp.Tool{
		Name:        "list_categories",
		Description: "List available categories of public service information",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input ListCategoriesInput) (*mcp.CallToolResult, ListCategoriesOutput, error) {
		// TODO: Implement actual category listing using client
		categories, err := httpClient.ListCategories(ctx)
		if err != nil {
			return nil, ListCategoriesOutput{}, fmt.Errorf("failed to list categories: %w", err)
		}

		output := ListCategoriesOutput{
			Categories: make([]Category, len(categories)),
		}
		for i, c := range categories {
			output.Categories[i] = Category{
				Name:        c.Name,
				Description: c.Description,
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Found %d categories", len(categories)),
				},
			},
		}, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}
