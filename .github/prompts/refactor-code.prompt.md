---
mode: 'agent'
model: Claude Sonnet 4
tools: ['codebase']
description: 'Refactor Go code to improve quality and maintainability'
---

# Refactor Code

Refactor the specified code to improve quality, maintainability, and adherence to best practices.

## Refactoring Goals

Ask the user which areas to focus on if not specified:

1. **Simplify complexity** - Reduce nested logic, improve readability
2. **Improve error handling** - Better error messages, proper wrapping
3. **Extract functions** - Break down large functions
4. **Reduce duplication** - DRY principle
5. **Improve naming** - More descriptive names
6. **Add type safety** - Replace `any` with specific types
7. **Optimize performance** - Reduce allocations, improve efficiency
8. **Improve testability** - Better separation of concerns

## Refactoring Process

### 1. Analyze Current Code
- Identify code smells
- Find duplication
- Locate complex functions
- Check error handling patterns
- Review naming conventions

### 2. Plan Refactoring
- Ensure tests exist before refactoring
- Refactor in small, incremental steps
- Run tests after each change
- Keep commits focused

### 3. Common Refactorings

#### Extract Function
Break down large functions:
```go
// Before
func ProcessData(data []byte) error {
    // 50 lines of code
}

// After
func ProcessData(data []byte) error {
    validated, err := validateData(data)
    if err != nil {
        return err
    }
    
    transformed := transformData(validated)
    return saveData(transformed)
}
```

#### Simplify Conditionals
Use early returns:
```go
// Before
func Process(input string) error {
    if input != "" {
        // lots of nested code
    } else {
        return errors.New("empty input")
    }
}

// After
func Process(input string) error {
    if input == "" {
        return errors.New("empty input")
    }
    // un-nested code
}
```

#### Reduce Duplication
Extract common code:
```go
// Before - duplication in multiple functions
func HandleA() { /* setup code */ /* operation A */ /* cleanup */ }
func HandleB() { /* setup code */ /* operation B */ /* cleanup */ }

// After
func withSetup(operation func() error) error {
    // setup code
    defer cleanup()
    return operation()
}
```

#### Improve Error Handling
Add context to errors:
```go
// Before
if err != nil {
    return err
}

// After
if err != nil {
    return fmt.Errorf("failed to process data: %w", err)
}
```

#### Interface Extraction
For better testability:
```go
// Before - hard to test
type Service struct {
    httpClient *http.Client
}

// After - easy to mock
type HTTPClient interface {
    Do(*http.Request) (*http.Response, error)
}

type Service struct {
    client HTTPClient
}
```

### 4. Verify Refactoring

After each refactoring:
- [ ] Run tests: `go test ./...`
- [ ] Check for races: `go test -race ./...`
- [ ] Verify formatting: `gofmt -s -w .`
- [ ] Run linter: `golangci-lint run`
- [ ] Ensure behavior unchanged

## Refactoring Principles

- **Don't change behavior** - Refactoring should not change functionality
- **Test first** - Ensure good test coverage before refactoring
- **Small steps** - Make incremental changes
- **One thing at a time** - Focus on one improvement per refactoring
- **Keep it simple** - Don't over-engineer

## Common Code Smells to Address

1. **Long functions** - Extract smaller functions
2. **Deep nesting** - Use early returns
3. **Magic numbers** - Use named constants
4. **Poor names** - Rename to be descriptive
5. **Duplicate code** - Extract common functionality
6. **Large structs** - Split into smaller types
7. **Too many parameters** - Use struct or context
8. **Global variables** - Use dependency injection

Provide clear explanations for each refactoring decision and show before/after comparisons.
