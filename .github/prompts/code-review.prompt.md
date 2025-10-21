---
mode: 'agent'
model: Claude Sonnet 4
tools: ['codebase']
description: 'Review Go code for quality, security, and best practices'
---

# Go Code Review

Perform a comprehensive code review focusing on:

## Code Quality

### Style and Formatting
- [ ] Code is formatted with `gofmt`
- [ ] Imports are organized with `goimports`
- [ ] Names follow Go conventions (mixedCaps, descriptive)
- [ ] Package names are lowercase, single-word
- [ ] No stuttering in names

### Clarity and Simplicity
- [ ] Code is clear and self-documenting
- [ ] Happy path is left-aligned
- [ ] Early returns reduce nesting
- [ ] Functions are focused and single-purpose
- [ ] No unnecessary complexity

### Error Handling
- [ ] All errors are checked
- [ ] Errors are wrapped with context
- [ ] Error messages are clear and actionable
- [ ] Errors are returned as last value
- [ ] No `_` ignoring errors without justification

### Documentation
- [ ] All exported symbols are documented
- [ ] Comments start with the symbol name
- [ ] Documentation explains "why" not "what"
- [ ] Package has package comment
- [ ] Complex logic has explanatory comments

## MCP Server Specifics

### Tool Implementation
- [ ] Input structs have comprehensive `jsonschema` tags
- [ ] Handlers have correct signature
- [ ] Input validation is performed
- [ ] Context cancellation is checked
- [ ] Errors are informative for clients

### Type Safety
- [ ] Struct-based inputs/outputs are used
- [ ] No unnecessary use of `any` or `interface{}`
- [ ] Proper type conversions
- [ ] Zero values are considered

## Concurrency

- [ ] Goroutine lifecycle is clear
- [ ] No potential goroutine leaks
- [ ] Shared state is protected with mutexes
- [ ] Channels are used correctly
- [ ] `sync.WaitGroup` used for goroutine coordination

## Security

### Input Validation
- [ ] All external inputs are validated
- [ ] URLs and queries are sanitized
- [ ] Input size limits are enforced
- [ ] Type safety prevents invalid states

### External Communication
- [ ] HTTPS is used for external APIs
- [ ] Timeouts are set on HTTP requests
- [ ] Rate limiting is considered
- [ ] Error responses don't leak sensitive info

### Secrets
- [ ] No hardcoded secrets
- [ ] Environment variables used for config
- [ ] Sensitive data not logged

## Performance

- [ ] Unnecessary allocations are avoided
- [ ] Slices are preallocated when size is known
- [ ] String concatenation uses `strings.Builder`
- [ ] Resources are properly closed (`defer`)
- [ ] No obvious performance bottlenecks

## Testing

- [ ] Tests exist for new functionality
- [ ] Tests cover success and error cases
- [ ] Table-driven tests for multiple scenarios
- [ ] Tests are independent and isolated
- [ ] Mock external dependencies

## Common Issues to Check

- [ ] No race conditions (run with `-race`)
- [ ] Resources are cleaned up (files, connections)
- [ ] Maps not modified concurrently
- [ ] Nil pointer checks where needed
- [ ] Context passed to long-running operations
- [ ] HTTP response bodies are closed

## Review Feedback Format

Provide feedback as:
1. **Critical** - Must fix before merge (bugs, security issues)
2. **Important** - Should fix (best practices, potential issues)
3. **Nice to have** - Consider for improvement (style, clarity)
4. **Question** - Request clarification

Include specific line references and suggest improvements with code examples when possible.
