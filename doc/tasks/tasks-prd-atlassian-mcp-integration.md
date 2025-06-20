# Tasks: Atlassian MCP Integration

## Requirements

Please read referenced files to understand the problem:
- `doc/prd-atlassian-mcp-integration.md`
- `doc/erd-atlassian-mcp-integration.md`

## Relevant Files

### Source Files
- `internal/services/http/middleware/auth.go` - Authentication middleware for HTTP client (✓ created)
- `internal/services/http/middleware/logging.go` - Logging middleware for HTTP client (✓ created)
- `internal/services/http/middleware/error_handling.go` - Error handling middleware for HTTP client (✓ created)
- `internal/services/http/client_factory.go` - HTTP client factory with middleware composition
- `internal/services/atlassian_client.go` - Atlassian-specific HTTP client implementation
- `internal/services/atlassian_accounts.go` - Accounts repository for managing multiple named Atlassian accounts
- `internal/app/bitbucket.go` - Business logic for Bitbucket operations
- `internal/app/jira.go` - Business logic for Jira operations
- `internal/api/mcp/controllers/bitbucket.go` - MCP controller for Bitbucket tools
- `internal/api/mcp/controllers/jira.go` - MCP controller for Jira tools
- `internal/api/mcp/controllers/register.go` - Controller registration (modified)
- `internal/app/register.go` - App layer registration (modified)
- `internal/services/register.go` - Service registration (modified)
- `internal/config/load.go` - Configuration loading (modified)

### Test Files
- `internal/services/http/middleware/auth_test.go` - Unit tests for authentication middleware (✓ created)
- `internal/services/http/middleware/logging_test.go` - Unit tests for logging middleware (✓ created)
- `internal/services/http/middleware/error_handling_test.go` - Unit tests for error handling middleware (✓ created)
- `internal/services/http/client_factory_test.go` - Unit tests for client factory
- `internal/services/atlassian_client_test.go` - Unit tests for Atlassian HTTP client
- `internal/services/atlassian_accounts_test.go` - Unit tests for accounts repository
- `internal/app/bitbucket_test.go` - Unit tests for Bitbucket business logic
- `internal/app/jira_test.go` - Unit tests for Jira business logic
- `internal/api/mcp/controllers/bitbucket_test.go` - Unit tests for Bitbucket MCP controller
- `internal/api/mcp/controllers/jira_test.go` - Unit tests for Jira MCP controller

### Configuration Files
- `internal/config/default.json` - Default configuration (modified)
- `internal/config/local.json` - Local configuration example (modified)
- `doc/instructions/creating-http-clients.md` - Instructions for creating HTTP clients
- `.mockery.yaml` - Mockery configuration for generating mocks (if needed)
- `go.mod` - Go module dependencies (if new packages needed)

### Notes

- **Testing Framework:** Use testify for assertions and table-driven tests
- **Architecture:** Follow clean architecture with layers: MCP → App → Services
- **HTTP Client Infrastructure:** Built on standard Go `http.Client` and `http.RoundTripper` with middleware pattern
- **Code Organization:** Follow the existing structure:
  - `internal/api/mcp/` - MCP protocol layer (controllers)
  - `internal/app/` - Application layer (business logic)
  - `internal/services/` - Infrastructure layer (HTTP clients, repositories)
  - `internal/services/http/` - HTTP client infrastructure and middleware
- **Naming Conventions:** Use Go conventions (PascalCase for exported, camelCase for unexported)
- **Authentication:** API token-based authentication via middleware for MVP
- **Multi-Account Support:** Named accounts (e.g., "user", "merge-bot") with default account configuration
- **Error Handling:** Use Go's idiomatic error handling with clear, actionable error messages
- **MCP Tools:** Use service-based naming (`bitbucket_*`, `jira_*`)
- **OpenAPI Integration:** Use official Bitbucket and Jira OpenAPI specifications for client generation
- **Instructional Development:** Create comprehensive instructions for generating API clients
- **Testing Commands:**
  - `make test` - Run all tests with coverage
  - `go test -v ./internal/path/... --run TestName` - Run specific tests
  - `gow test -v ./internal/path/... --run TestName` - Watch mode for test development
- **Mock Generation:** Use `go generate` with mockery for interface mocks if needed
- **Dependency Injection:** Use uber/dig for dependency management
- **Configuration:** Extend existing viper-based configuration system with base URLs for Atlassian APIs

## Tasks

- [ ] 1.0 HTTP Client Infrastructure Foundation
  - [x] 1.1 Create `internal/services/http/middleware/` directory structure for HTTP middleware components
  - [x] 1.2 Create unit tests for authentication middleware in `internal/services/http/middleware/auth_test.go` (TDD: write tests first)
  - [x] 1.3 Implement `internal/services/http/middleware/auth.go` - Authentication middleware as `http.RoundTripper` wrapper
  - [x] 1.4 Create unit tests for logging middleware in `internal/services/http/middleware/logging_test.go` (TDD: write tests first)
  - [x] 1.5 Implement `internal/services/http/middleware/logging.go` - Structured logging middleware as `http.RoundTripper` wrapper
  - [x] 1.6 Create unit tests for error handling middleware in `internal/services/http/middleware/error_handling_test.go` (TDD: write tests first)
  - [x] 1.7 Implement `internal/services/http/middleware/error_handling.go` - Error handling middleware as `http.RoundTripper` wrapper
  - [ ] 1.8 Create unit tests for client factory in `internal/services/http/client_factory_test.go` (TDD: write tests first)
  - [ ] 1.9 Create `internal/services/http/client_factory.go`