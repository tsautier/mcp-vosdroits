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

	// Register life events tools
	if err := registerListLifeEvents(server, httpClient); err != nil {
		return fmt.Errorf("failed to register list_life_events: %w", err)
	}

	if err := registerGetLifeEventDetails(server, httpClient); err != nil {
		return fmt.Errorf("failed to register get_life_event_details: %w", err)
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
		Description: "Search for SPECIFIC administrative procedures on service-public.gouv.fr (e.g., 'passport renewal', 'driver's license'). WARNING: For MAJOR LIFE SITUATIONS (buying house, marriage, birth, death, moving, job change, retirement), you MUST use list_life_events + get_life_event_details FIRST before trying this tool. Only use search_procedures if the life event details are insufficient or for simple administrative tasks only. Returns URLs and brief descriptions. For detailed information about any procedure, use the get_article tool with the returned URLs.",
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
			// Return a clear error message that discourages retrying the same URL
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("ERROR: Unable to retrieve article from %s. Reason: %v\n\nDo NOT retry this same URL. Instead, inform the user that this specific article could not be retrieved and suggest they visit the URL directly in their browser, or try searching for alternative procedures.", input.URL, err),
					},
				},
				IsError: true,
			}, GetArticleOutput{}, fmt.Errorf("failed to get article from %s: %w", input.URL, err)
		}

		output := GetArticleOutput{
			Title:   article.Title,
			Content: article.Content,
			URL:     article.URL,
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Retrieved article: %s\n\nSource: %s\n\nIMPORTANT: Always provide this source URL to the user so they can access the original article.", article.Title, article.URL),
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

// ListLifeEventsOutput defines the output schema for list_life_events.
type ListLifeEventsOutput struct {
	Events []LifeEventInfo `json:"events" jsonschema:"List of available life events - THESE ARE ONLY TITLES AND URLS. You MUST call get_life_event_details with a URL to get actual procedures and detailed information."`
}

// LifeEventInfo represents a life event.
type LifeEventInfo struct {
	Title string `json:"title" jsonschema:"Title of the life event - this is ONLY the title, not the actual content"`
	URL   string `json:"url" jsonschema:"IMPORTANT: URL to pass to get_life_event_details tool to retrieve actual procedures and information. Must be used exactly as-is - do not modify."`
}

type ListLifeEventsInput struct{}

func registerListLifeEvents(server *mcp.Server, httpClient *client.Client) error {
	tool := &mcp.Tool{
		Name:        "list_life_events",
		Description: "PRIMARY TOOL for major life situations in France (e.g., 'acheter une maison', 'buying a house', 'getting married', 'having a baby', 'death of relative', 'moving', 'changing jobs', 'retirement', 'divorce'). Use THIS tool first for comprehensive guidance on major life events. Returns list of available life events with titles and URLs. MANDATORY NEXT STEP: You MUST call get_life_event_details with one of the returned URLs to get actual procedures - DO NOT skip this step and jump to search_procedures. Only use search_procedures if get_life_event_details doesn't provide enough information.",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input ListLifeEventsInput) (*mcp.CallToolResult, ListLifeEventsOutput, error) {
		events, err := httpClient.ListLifeEvents(ctx)
		if err != nil {
			return nil, ListLifeEventsOutput{}, fmt.Errorf("failed to list life events: %w", err)
		}

		output := ListLifeEventsOutput{
			Events: make([]LifeEventInfo, len(events)),
		}
		for i, e := range events {
			output.Events[i] = LifeEventInfo{
				Title: e.Title,
				URL:   e.URL,
			}
		}

		// Build a formatted list of events for the LLM to see
		eventsList := "\n\nAvailable life events:\n"
		for i, e := range events {
			eventsList += fmt.Sprintf("%d. %s\n   URL: %s\n\n", i+1, e.Title, e.URL)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Found %d life events (titles only).\n\nMANDATORY NEXT STEP: You MUST immediately call get_life_event_details with the most relevant URL from the list below. DO NOT call search_procedures yet - the life event details contain comprehensive information organized by topic that will likely answer the user's question.\n\nOnly use search_procedures if get_life_event_details doesn't provide sufficient information.\n\nIMPORTANT: Use the EXACT URL from the results - these are fiche pratique URLs (F-prefix). Do NOT modify or substitute with category URLs (N-prefix).%s", len(events), eventsList),
				},
			},
		}, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}

// GetLifeEventDetailsInput defines the input schema for get_life_event_details.
type GetLifeEventDetailsInput struct {
	URL string `json:"url" jsonschema:"required,EXACT URL from list_life_events results. Must be a fiche pratique URL with F-prefix like https://www.service-public.gouv.fr/particuliers/vosdroits/F16225. Do NOT use category URLs with N-prefix or modify the URL."`
}

// GetLifeEventDetailsOutput defines the output schema for get_life_event_details.
type GetLifeEventDetailsOutput struct {
	Title        string                   `json:"title" jsonschema:"Title of the life event"`
	URL          string                   `json:"url" jsonschema:"URL of the life event page"`
	Introduction string                   `json:"introduction" jsonschema:"Introduction text explaining the life event"`
	Sections     []LifeEventSectionOutput `json:"sections" jsonschema:"Detailed sections organized by topic (Health, Civil Status, Employment, etc.)"`
}

// LifeEventSectionOutput represents a section within a life event.
type LifeEventSectionOutput struct {
	Title   string `json:"title" jsonschema:"Section title (e.g., Santé, État civil, Emploi-Travail)"`
	Content string `json:"content" jsonschema:"Detailed content for this section"`
}

func registerGetLifeEventDetails(server *mcp.Server, httpClient *client.Client) error {
	tool := &mcp.Tool{
		Name:        "get_life_event_details",
		Description: "Retrieve detailed information about a specific life event from service-public.gouv.fr. CRITICAL: Use the EXACT URL from list_life_events results (fiche pratique F-URLs only, like /F16225). Do NOT use category URLs (N-prefix like /N20020). Provides comprehensive guidance organized by topic (Health, Civil Status, Employment, etc.) for major life situations.",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input GetLifeEventDetailsInput) (*mcp.CallToolResult, GetLifeEventDetailsOutput, error) {
		if input.URL == "" {
			return nil, GetLifeEventDetailsOutput{}, fmt.Errorf("url cannot be empty")
		}

		details, err := httpClient.GetLifeEventDetails(ctx, input.URL)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("ERROR: Unable to retrieve life event details from %s. Reason: %v\n\nDo NOT retry this same URL. Instead, inform the user that this specific page could not be retrieved and suggest they visit the URL directly in their browser, or try searching for alternative life events.", input.URL, err),
					},
				},
				IsError: true,
			}, GetLifeEventDetailsOutput{}, fmt.Errorf("failed to get life event details from %s: %w", input.URL, err)
		}

		output := GetLifeEventDetailsOutput{
			Title:        details.Title,
			URL:          details.URL,
			Introduction: details.Introduction,
			Sections:     make([]LifeEventSectionOutput, len(details.Sections)),
		}
		for i, s := range details.Sections {
			output.Sections[i] = LifeEventSectionOutput{
				Title:   s.Title,
				Content: s.Content,
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Retrieved life event: %s\n\nThis guide contains %d sections with detailed information.\n\nSource: %s\n\nIMPORTANT: Always provide this source URL to the user so they can access the original page.", details.Title, len(details.Sections), details.URL),
				},
			},
		}, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}
