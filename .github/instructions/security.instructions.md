---
description: 'Security best practices for Go MCP server development'
applyTo: '**/*.go'
---

# Security Best Practices

## Input Validation

- Validate all external inputs (tool parameters, resource URIs)
- Use strong typing to prevent invalid states
- Sanitize user-provided URLs and queries
- Validate data formats (emails, URLs, etc.)
- Set reasonable limits on input sizes

## HTTP Security

### Client Configuration

- Always use HTTPS for external API calls
- Set reasonable timeouts to prevent hanging requests
- Validate SSL/TLS certificates
- Handle redirects carefully

### Rate Limiting

- Implement rate limiting for external API calls
- Handle 429 (Too Many Requests) responses gracefully
- Use exponential backoff for retries

## Error Handling

- Don't expose sensitive information in error messages
- Log errors securely without leaking secrets
- Return generic error messages to clients
- Log detailed errors for debugging

## Secrets Management

- Never hardcode secrets in source code
- Use environment variables for configuration
- Rotate secrets regularly
- Don't log sensitive information
- Use secret management services when available

## Data Privacy

- Don't log personally identifiable information (PII)
- Sanitize logs before storage
- Respect user privacy in caching
- Follow data retention policies

## Dependency Security

- Regularly update dependencies
- Use `go mod tidy` to remove unused dependencies
- Scan for known vulnerabilities
- Pin dependency versions in production

## Cryptography

- Use Go's standard library crypto packages
- Don't implement custom cryptography
- Use `crypto/rand` for random number generation
- Use TLS for network communication

## Context Security

- Set timeouts on contexts to prevent resource exhaustion
- Cancel contexts when operations complete
- Don't store sensitive data in context values

## Common Vulnerabilities

### SQL Injection
Not applicable for this project, but if using databases, always use parameterized queries.

### Path Traversal
- Validate file paths from user input
- Use `filepath.Clean` to sanitize paths
- Restrict file access to allowed directories

### Command Injection
- Avoid executing shell commands with user input
- If necessary, use Go's `os/exec` with separate arguments
- Sanitize all inputs

## Logging Security

- Use structured logging
- Don't log secrets or sensitive data
- Sanitize URLs before logging
- Log security events (authentication failures, etc.)

## Docker Security

- Run containers as non-root user
- Use minimal base images
- Scan images for vulnerabilities
- Keep images updated
- Don't include secrets in images
