---
mode: 'agent'
model: Claude Sonnet 4
tools: ['codebase']
description: 'Write comprehensive tests for MCP tools'
---

# Write MCP Tool Tests

Generate comprehensive tests for MCP tool handlers following Go testing best practices.

## Test Requirements

For the specified tool, create tests that cover:

1. **Success Cases**
   - Valid inputs return expected outputs
   - Different input variations
   - Edge cases within valid range

2. **Error Cases**
   - Invalid or missing required inputs
   - External API failures
   - Network timeouts
   - Invalid response formats

3. **Context Handling**
   - Context cancellation is respected
   - Timeouts are handled properly

## Test Structure

Use table-driven tests with `t.Run` for organization:

```go
func TestToolName(t *testing.T) {
    tests := []struct {
        name    string
        input   ToolInput
        want    ToolOutput
        wantErr bool
        errMsg  string
    }{
        {
            name: "success case description",
            input: ToolInput{
                // valid input
            },
            want: ToolOutput{
                // expected output
            },
            wantErr: false,
        },
        {
            name: "error case description",
            input: ToolInput{
                // invalid input
            },
            wantErr: true,
            errMsg: "expected error message",
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()
            
            result, output, err := ToolHandler(ctx, nil, tt.input)
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error, got nil")
                }
                if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
                    t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
                }
                return
            }
            
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            
            // Compare output with expected
            if !reflect.DeepEqual(output, tt.want) {
                t.Errorf("got %+v, want %+v", output, tt.want)
            }
        })
    }
}
```

## Mock External Dependencies

If the tool calls external APIs:

1. Create a mock HTTP client or server
2. Use `httptest.NewServer` for HTTP mocking
3. Test with various response scenarios:
   - Success responses
   - Error responses (4xx, 5xx)
   - Network errors
   - Timeout scenarios

## Context Cancellation Test

Always include a test for context cancellation:

```go
func TestToolName_ContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel() // Cancel immediately
    
    input := ToolInput{ /* valid input */ }
    _, _, err := ToolHandler(ctx, nil, input)
    
    if err == nil {
        t.Error("expected error due to cancelled context")
    }
    if !errors.Is(err, context.Canceled) {
        t.Errorf("expected context.Canceled error, got %v", err)
    }
}
```

## Test Helpers

Use test helpers to reduce duplication:

```go
func TestToolName_helper(t *testing.T) {
    t.Helper()
    // Common setup code
}
```

## Best Practices

- Keep tests focused and independent
- Use meaningful test names that describe the scenario
- Don't test implementation details
- Test behavior, not internal structure
- Clean up resources with `t.Cleanup()`
- Run tests with race detector: `go test -race`
