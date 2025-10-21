---
mode: 'agent'
model: Claude Sonnet 4
tools: ['codebase']
description: 'Debug issues in the Go MCP server'
---

# Debug Issue

Help debug and resolve issues in the VosDroits MCP server.

## Information Gathering

Ask the user for:
1. **Symptom** - What is the observed behavior?
2. **Expected Behavior** - What should happen?
3. **Error Messages** - Any error messages or stack traces?
4. **Steps to Reproduce** - How can the issue be reproduced?
5. **Environment** - Go version, OS, deployment method?
6. **Recent Changes** - What changed before the issue appeared?

## Debugging Approach

### 1. Reproduce the Issue
- Try to reproduce based on user's description
- Create a minimal reproduction case
- Isolate the problematic code

### 2. Analyze Error Messages
- Parse stack traces for root cause
- Identify the failing function
- Check error wrapping for context

### 3. Check Common Issues

#### MCP Tool Issues
- **Input validation failure**
  - Check jsonschema tags
  - Verify required fields
  - Validate input constraints
  
- **Context cancellation**
  - Check if context is cancelled
  - Verify timeout settings
  - Ensure context is passed correctly

- **External API errors**
  - Check HTTP status codes
  - Verify API endpoint URLs
  - Check rate limiting
  - Validate response parsing

#### Go-Specific Issues
- **Nil pointer dereference**
  - Check for nil before dereferencing
  - Validate return values
  
- **Race conditions**
  - Run with `-race` flag
  - Check concurrent map access
  - Verify mutex usage
  
- **Goroutine leaks**
  - Ensure goroutines exit
  - Check for blocked channels
  - Verify context cancellation

- **Resource leaks**
  - Ensure HTTP response bodies are closed
  - Check `defer` statements
  - Verify file handles are closed

### 4. Add Debugging

#### Logging
Add structured logging:
```go
log.Info("tool called",
    "name", req.Params.Name,
    "input", input,
)
```

#### Error Context
Add context to errors:
```go
if err != nil {
    return fmt.Errorf("failed to fetch article from %s: %w", url, err)
}
```

#### Debugging Tools
Use Go debugging tools:
```bash
# Race detector
go test -race ./...

# Memory profiling
go test -memprofile=mem.out

# CPU profiling
go test -cpuprofile=cpu.out

# Trace execution
go test -trace=trace.out
```

### 5. Common Solutions

#### HTTP Request Issues
```go
// Add timeout
client := &http.Client{
    Timeout: 30 * time.Second,
}

// Check response status
if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("unexpected status: %d", resp.StatusCode)
}

// Always close body
defer resp.Body.Close()
```

#### Context Handling
```go
// Check cancellation
if ctx.Err() != nil {
    return ctx.Err()
}

// Add timeout
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()
```

#### Input Validation
```go
// Validate required fields
if input.Query == "" {
    return fmt.Errorf("query cannot be empty")
}

// Validate ranges
if input.Limit < 1 || input.Limit > 100 {
    return fmt.Errorf("limit must be between 1 and 100")
}
```

### 6. Testing the Fix

After identifying and fixing the issue:
- Write a test that reproduces the bug
- Verify the fix resolves the issue
- Ensure no regressions
- Add edge case tests

Example test for a bug:
```go
func TestSearchProcedures_BugFix(t *testing.T) {
    // Test that previously caused the bug
    ctx := context.Background()
    input := ToolInput{
        Query: "", // Empty query that caused panic
    }
    
    _, _, err := SearchProcedures(ctx, nil, input)
    
    if err == nil {
        t.Error("expected error for empty query")
    }
}
```

## Debugging Checklist

- [ ] Understand the symptom and expected behavior
- [ ] Reproduce the issue
- [ ] Analyze error messages and stack traces
- [ ] Check for common issues (nil, race, leaks)
- [ ] Add logging and debugging output
- [ ] Identify root cause
- [ ] Implement fix
- [ ] Write test to prevent regression
- [ ] Verify fix resolves issue
- [ ] Document the issue and fix

Provide clear explanations of the issue, root cause, and solution with code examples.
