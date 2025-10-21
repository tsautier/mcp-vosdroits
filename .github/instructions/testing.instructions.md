---
description: 'Testing standards and best practices for Go MCP server development'
applyTo: '**/*_test.go'
---

# Testing Standards

## General Principles

- Write tests for all public functions and methods
- Test both success and error paths
- Use table-driven tests for multiple scenarios
- Keep tests simple and focused
- Tests should be fast and deterministic
- Avoid testing implementation details

## Go Testing Patterns

### Table-Driven Tests

Use table-driven tests to test multiple scenarios efficiently:

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### Subtests

Use `t.Run` to organize related tests and enable selective test execution.

### Test Helpers

- Mark helper functions with `t.Helper()`
- Use `t.Cleanup()` for resource cleanup
- Create test fixtures for complex setup

## MCP Server Testing

### Testing Tool Handlers

- Test with valid and invalid inputs
- Verify JSON schema validation
- Test context cancellation
- Mock external dependencies
- Verify error messages are informative

### Testing Resources

- Test resource read operations
- Verify MIME types and content
- Test resource not found scenarios

### Testing Prompts

- Test prompt generation with various arguments
- Verify prompt message formatting

## Mocking

- Use interfaces for dependencies
- Create mock implementations for testing
- Consider using testify/mock for complex mocks
- Keep mocks simple and focused

## Coverage

- Aim for high test coverage but prioritize meaningful tests
- Use `go test -cover` to check coverage
- Focus on testing critical paths and edge cases

## Best Practices

- Name tests descriptively: `Test_FunctionName_Scenario`
- Don't test the Go standard library
- Avoid sleeps in tests; use channels or mocks
- Run tests with race detector: `go test -race`
- Keep tests isolated and independent
