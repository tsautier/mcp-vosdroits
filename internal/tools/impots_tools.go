// Package tools provides MCP tool implementations for searching French tax information.
package tools

import (
	"context"
	"fmt"

	"github.com/guigui42/mcp-vosdroits/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterImpotsTools registers all impots.gouv.fr MCP tools with the server.
func RegisterImpotsTools(server *mcp.Server, impotsClient *client.ImpotsClient) error {
	if err := registerSearchImpots(server, impotsClient); err != nil {
		return fmt.Errorf("failed to register search_impots: %w", err)
	}

	if err := registerGetImpotsArticle(server, impotsClient); err != nil {
		return fmt.Errorf("failed to register get_impots_article: %w", err)
	}

	if err := registerListImpotsCategories(server, impotsClient); err != nil {
		return fmt.Errorf("failed to register list_impots_categories: %w", err)
	}

	return nil
}

// SearchImpotsInput defines the input schema for search_impots.
type SearchImpotsInput struct {
	Query string `json:"query" jsonschema:"Search query for tax information and forms"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of results to return (1-100)"`
}

// SearchImpotsOutput defines the output schema for search_impots.
type SearchImpotsOutput struct {
	Results []ImpotsResult `json:"results" jsonschema:"List of matching tax documents and articles"`
}

// ImpotsResult represents a single search result from impots.gouv.fr.
type ImpotsResult struct {
	Title       string `json:"title" jsonschema:"Title of the tax document or article"`
	URL         string `json:"url" jsonschema:"URL to the document page"`
	Description string `json:"description,omitempty" jsonschema:"Brief description"`
	Type        string `json:"type,omitempty" jsonschema:"Type of document (Formulaire, Article, etc.)"`
	Date        string `json:"date,omitempty" jsonschema:"Publication or update date"`
}

func registerSearchImpots(server *mcp.Server, impotsClient *client.ImpotsClient) error {
	tool := &mcp.Tool{
		Name:        "search_impots",
		Description: "Search for tax forms, articles, and procedures on impots.gouv.fr",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input SearchImpotsInput) (*mcp.CallToolResult, SearchImpotsOutput, error) {
		if input.Limit == 0 {
			input.Limit = 10
		}

		if input.Query == "" {
			return nil, SearchImpotsOutput{}, fmt.Errorf("query cannot be empty")
		}

		results, err := impotsClient.SearchImpots(ctx, input.Query, input.Limit)
		if err != nil {
			return nil, SearchImpotsOutput{}, fmt.Errorf("search failed: %w", err)
		}

		output := SearchImpotsOutput{
			Results: make([]ImpotsResult, len(results)),
		}
		for i, r := range results {
			output.Results[i] = ImpotsResult{
				Title:       r.Title,
				URL:         r.URL,
				Description: r.Description,
				Type:        r.Type,
				Date:        r.Date,
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Found %d tax documents", len(results)),
				},
			},
		}, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}

// GetImpotsArticleInput defines the input schema for get_impots_article.
type GetImpotsArticleInput struct {
	URL string `json:"url" jsonschema:"URL of the tax article or form to retrieve from impots.gouv.fr"`
}

// GetImpotsArticleOutput defines the output schema for get_impots_article.
type GetImpotsArticleOutput struct {
	Title       string `json:"title" jsonschema:"Title of the tax document"`
	Content     string `json:"content" jsonschema:"Full content of the document"`
	URL         string `json:"url" jsonschema:"URL of the document"`
	Type        string `json:"type,omitempty" jsonschema:"Type of document"`
	Description string `json:"description,omitempty" jsonschema:"Brief description"`
}

func registerGetImpotsArticle(server *mcp.Server, impotsClient *client.ImpotsClient) error {
	tool := &mcp.Tool{
		Name:        "get_impots_article",
		Description: "Retrieve detailed information from a specific tax article or form URL on impots.gouv.fr",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input GetImpotsArticleInput) (*mcp.CallToolResult, GetImpotsArticleOutput, error) {
		if input.URL == "" {
			return nil, GetImpotsArticleOutput{}, fmt.Errorf("url cannot be empty")
		}

		article, err := impotsClient.GetImpotsArticle(ctx, input.URL)
		if err != nil {
			return nil, GetImpotsArticleOutput{}, fmt.Errorf("failed to get article: %w", err)
		}

		output := GetImpotsArticleOutput{
			Title:       article.Title,
			Content:     article.Content,
			URL:         article.URL,
			Type:        article.Type,
			Description: article.Description,
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Retrieved tax document: %s", article.Title),
				},
			},
		}, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}

// ListImpotsCategoriesOutput defines the output schema for list_impots_categories.
type ListImpotsCategoriesOutput struct {
	Categories []ImpotsCategory `json:"categories" jsonschema:"List of available tax service categories"`
}

// ImpotsCategory represents a tax service category.
type ImpotsCategory struct {
	Name        string `json:"name" jsonschema:"Name of the category"`
	Description string `json:"description" jsonschema:"Description of the category"`
	URL         string `json:"url" jsonschema:"URL to the category page"`
}

type ListImpotsCategoriesInput struct{}

func registerListImpotsCategories(server *mcp.Server, impotsClient *client.ImpotsClient) error {
	tool := &mcp.Tool{
		Name:        "list_impots_categories",
		Description: "List available categories of tax information on impots.gouv.fr",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input ListImpotsCategoriesInput) (*mcp.CallToolResult, ListImpotsCategoriesOutput, error) {
		categories, err := impotsClient.ListImpotsCategories(ctx)
		if err != nil {
			return nil, ListImpotsCategoriesOutput{}, fmt.Errorf("failed to list categories: %w", err)
		}

		output := ListImpotsCategoriesOutput{
			Categories: make([]ImpotsCategory, len(categories)),
		}
		for i, c := range categories {
			output.Categories[i] = ImpotsCategory{
				Name:        c.Name,
				Description: c.Description,
				URL:         c.URL,
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Found %d tax categories", len(categories)),
				},
			},
		}, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}
