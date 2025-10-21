<!-- Based on: https://github.com/github/awesome-copilot/blob/main/instructions/go.instructions.md -->
---
description: 'Instructions for writing Go code following idiomatic Go practices and community standards'
applyTo: '**/*.go,**/go.mod,**/go.sum'
---

# Go Development Instructions

Follow idiomatic Go practices and community standards when writing Go code. These instructions are based on [Effective Go](https://go.dev/doc/effective_go), [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments), and [Google's Go Style Guide](https://google.github.io/styleguide/go/).

## General Instructions

- Write simple, clear, and idiomatic Go code
- Favor clarity and simplicity over cleverness
- Follow the principle of least surprise
- Keep the happy path left-aligned (minimize indentation)
- Return early to reduce nesting
- Prefer early return over if-else chains
- Make the zero value useful
- Write self-documenting code with clear, descriptive names
- Document exported types, functions, methods, and packages
- Use Go modules for dependency management
- Leverage the Go standard library instead of reinventing the wheel
- Write comments in English

## Naming Conventions

### Packages
- Use lowercase, single-word package names
- Avoid underscores, hyphens, or mixedCaps
- Choose names that describe what the package provides
- Avoid generic names like `util`, `common`, or `base`

### Variables and Functions
- Use mixedCaps or MixedCaps (camelCase) rather than underscores
- Keep names short but descriptive
- Exported names start with a capital letter
- Unexported names start with a lowercase letter
- Avoid stuttering (e.g., avoid `http.HTTPServer`, prefer `http.Server`)

### Interfaces
- Name interfaces with -er suffix when possible (e.g., `Reader`, `Writer`)
- Single-method interfaces should be named after the method
- Keep interfaces small and focused

## Code Style and Formatting

- Always use `gofmt` to format code
- Use `goimports` to manage imports automatically
- Keep line length reasonable (no hard limit, but consider readability)
- Add blank lines to separate logical groups of code

## Error Handling

- Check errors immediately after the function call
- Don't ignore errors using `_` unless documented why
- Wrap errors with context using `fmt.Errorf` with `%w` verb
- Create custom error types when needed
- Place error returns as the last return value
- Name error variables `err`
- Keep error messages lowercase and don't end with punctuation

## Concurrency

### Goroutines
- Always know how a goroutine will exit
- Use `sync.WaitGroup` or channels to wait for goroutines
- Avoid goroutine leaks by ensuring cleanup

### Channels
- Use channels to communicate between goroutines
- Close channels from the sender side
- Use buffered channels when you know the capacity
- Use `select` for non-blocking operations

### Synchronization
- Use `sync.Mutex` for protecting shared state
- Keep critical sections small
- Use `sync.RWMutex` when you have many readers
- Use `sync.Once` for one-time initialization
- For Go 1.25+, use `WaitGroup.Go` method for cleaner goroutine management

## Testing

- Use table-driven tests for multiple test cases
- Name tests descriptively using `Test_functionName_scenario`
- Use subtests with `t.Run` for better organization
- Test both success and error cases
- Mark helper functions with `t.Helper()`
- Clean up resources using `t.Cleanup()`

## Performance

- Minimize allocations in hot paths
- Reuse objects when possible (consider `sync.Pool`)
- Preallocate slices when size is known
- Avoid unnecessary string conversions
- Profile before optimizing
- Use built-in profiling tools (`pprof`)

## Common Pitfalls to Avoid

- Not checking errors
- Ignoring race conditions
- Creating goroutine leaks
- Not using defer for cleanup
- Modifying maps concurrently
- Forgetting to close resources (files, connections)
- Using global variables unnecessarily
- Not considering the zero value of types
