---
mode: 'agent'
model: Claude Sonnet 4
tools: ['codebase']
description: 'Add a new MCP tool to the VosDroits server'
---

# Add New MCP Tool

Your goal is to add a new MCP tool to the VosDroits server following the established patterns in the codebase.

## Information Needed

Ask the user for the following if not provided:
1. **Tool name** - The name of the tool (e.g., "get_procedure_details")
2. **Description** - What the tool does
3. **Input parameters** - What inputs the tool accepts
4. **Output format** - What the tool returns
5. **External API** - Which service-public.gouv.fr endpoint to call (if applicable)

## Implementation Steps

1. **Create Input/Output Structs**
   - Define typed structs in the appropriate file in `internal/tools/`
   - Add comprehensive `json` and `jsonschema` tags
   - Include field descriptions, constraints, and examples

2. **Implement Tool Handler**
   - Create handler function with signature:
     ```go
     func ToolHandler(ctx context.Context, req *mcp.CallToolRequest, input InputStruct) (*mcp.CallToolResult, OutputStruct, error)
     ```
   - Add input validation
   - Check context cancellation
   - Call external API through HTTP client
   - Handle errors appropriately
   - Return structured output

3. **Register Tool**
   - Add tool registration in `cmd/server/main.go` or appropriate setup file
   - Use `mcp.AddTool` with tool metadata:
     - Name
     - Description
     - Handler function

4. **Add Tests**
   - Create test file `internal/tools/toolname_test.go`
   - Write table-driven tests
   - Test success cases
   - Test error cases
   - Test input validation
   - Test context cancellation

5. **Update Documentation**
   - Add tool to README.md
   - Document input/output schema
   - Provide usage examples

## Code Quality Checklist

- [ ] Input struct has comprehensive jsonschema tags
- [ ] Output struct is well-typed
- [ ] Handler validates inputs
- [ ] Handler checks context cancellation
- [ ] Errors are wrapped with context
- [ ] Tests cover success and error paths
- [ ] Documentation is updated
- [ ] Code follows Go and MCP server guidelines

## Example Tool Structure

```go
// Input struct with JSON schema tags
type ToolInput struct {
    Field string `json:"field" jsonschema:"required,description=Description of field"`
}

// Output struct
type ToolOutput struct {
    Result string `json:"result" jsonschema:"description=Result description"`
}

// Handler function
func ToolHandler(ctx context.Context, req *mcp.CallToolRequest, input ToolInput) (*mcp.CallToolResult, ToolOutput, error) {
    // Check context
    if ctx.Err() != nil {
        return nil, ToolOutput{}, ctx.Err()
    }
    
    // Validate input
    if input.Field == "" {
        return nil, ToolOutput{}, fmt.Errorf("field cannot be empty")
    }
    
    // Perform operation
    result, err := performOperation(ctx, input.Field)
    if err != nil {
        return nil, ToolOutput{}, fmt.Errorf("operation failed: %w", err)
    }
    
    return nil, ToolOutput{Result: result}, nil
}
```

Follow the existing patterns in the codebase and ensure consistency with other tools.
