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
- `internal/services/http/client_factory.go` - HTTP client factory with middleware composition (✓ created)
- `internal/services/http/send_request.go` - Shared SendRequest function with generic body and target types (✓ created)
- `doc/instructions/creating-http-clients.md` - Comprehensive instructions for API client implementation patterns (✓ created)
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
  - [x] 1.8 Create unit tests for client factory in `internal/services/http/client_factory_test.go` (TDD: write tests first)
  - [x] 1.9 Create `internal/services/http/client_factory.go` - Factory for composing middleware stack and creating configured `http.Client` instances
  - [x] 1.10 Register HTTP client infrastructure components in `internal/services/register.go`
- [x] 2.0 Client Generation Instructions & Patterns
  - [x] 2.1 Establish patterns and interfaces for API client implementation using standard Go HTTP types
  - [x] 2.2 Create `doc/instructions/` directory for client generation instruction documentation
  - [x] 2.3 Develop `doc/instructions/creating-http-clients.md` with comprehensive templates for struct definitions
  - [x] 2.4 Add method implementation patterns for converting OpenAPI endpoints to Go methods
  - [x] 2.5 Document error handling patterns and authentication integration guidelines
  - [x] 2.6 Create OpenAPI processing guidelines for schema mapping and parameter handling
  - [x] 2.7 Establish quality assurance instructions for testing, documentation, and validation
  - [x] 2.8 Define templates for converting OpenAPI schemas to Go structs with proper JSON tags
  - [x] 2.9 Document response processing patterns for different HTTP status codes and content types
  - [x] 2.10 Create shared doRequest function for common HTTP request handling  
  - [x] 2.11 Create shared mockClientFactory implementation for testing - not required. The test server will be used instead.
- [ ] 3.0 Atlassian HTTP Clients Implementation
  - [x] 3.1 Update `internal/config/default.json` with Atlassian configuration section including base URLs
  - [x] 3.2 Inject Atlassian client configuration into DI as per `internal/config/provide.go`
  - [ ] 3.3 Find official Bitbucket Cloud OpenAPI specification and add it to `internal/services/bitbucket/openapi.yaml`
  - [ ] 3.4 Find official Jira Cloud OpenAPI specification and add it to `internal/services/jira/openapi.yaml`
  - [ ] 3.5 Create Bitbucket client based on the openapi and `doc/instructions/creating-http-clients.md` instruction. Add methods for following Bitbucket API calls only: `CreatePR`, `GetPR`, `UpdatePR`, `ApprovePR`, `MergePR`
  - [ ] 3.6 Create Jira client based on the openapi and `doc/instructions/creating-http-clients.md` instruction. Add methods for following Jira API calls only: `GetTicket`, `TransitionTicket`, `ManageLabels`
  - [ ] 3.7 Add Atlassian-specific error response parsing and meaningful error messages (research if needed)
  - [ ] 3.8 Register Atlassian HTTP clients in `internal/services/register.go`
- [ ] 4.0 Accounts Management System
  - [ ] 4.1 Design accounts file JSON schema with multiple named accounts and default account specification
  - [ ] 4.2 Create `internal/services/atlassian_accounts.go` with accounts repository interface
  - [ ] 4.3 Implement `GetDefaultAccount()` method for retrieving configured default account
  - [ ] 4.4 Implement `GetAccountByName(name string)` method for retrieving specific named accounts
  - [ ] 4.5 Add support for separate Bitbucket and Jira credentials per account with dynamic URL parameters
  - [ ] 4.6 Implement file reading, parsing, and validation logic with proper error handling
  - [ ] 4.7 Handle workspace and domain parameters from account configuration (not main config)
  - [ ] 4.8 Create comprehensive unit tests in `internal/services/atlassian_accounts_test.go`
  - [ ] 4.9 Register accounts repository in `internal/services/register.go`
- [ ] 5.0 Bitbucket MCP Integration
  - [ ] 5.1 Create `internal/app/bitbucket.go` with business logic service interface and implementation
  - [ ] 5.2 Implement `bitbucket_create_pr` tool: Create pull requests with template support and reviewer assignment
  - [ ] 5.3 Implement `bitbucket_read_pr` tool: Retrieve comprehensive PR details including status and reviews
  - [ ] 5.4 Implement `bitbucket_update_pr` tool: Update PR titles and descriptions with validation
  - [ ] 5.5 Implement `bitbucket_approve_pr` tool: Approve pull requests with proper authorization
  - [ ] 5.6 Implement `bitbucket_merge_pr` tool: Merge PRs with strategy selection and pre-merge validation
  - [ ] 5.7 Create `internal/api/mcp/controllers/bitbucket.go` with MCP protocol handling for all tools
  - [ ] 5.8 Add account parameter support to all tools with default account resolution
  - [ ] 5.9 Create unit tests for business logic in `internal/app/bitbucket_test.go`
  - [ ] 5.10 Create unit tests for MCP controllers in `internal/api/mcp/controllers/bitbucket_test.go`
  - [ ] 5.11 Register Bitbucket services in `internal/app/register.go` and controllers in `internal/api/mcp/controllers/register.go`
- [ ] 6.0 Jira MCP Integration
  - [ ] 6.1 Create `internal/app/jira.go` with business logic service interface and implementation
  - [ ] 6.2 Implement `jira_read_ticket` tool: Retrieve complete ticket information including metadata and comments
  - [ ] 6.3 Implement `jira_transition_ticket` tool: Move tickets through workflow states with validation
  - [ ] 6.4 Implement `jira_manage_labels` tool: Add/remove labels with permission and format validation
  - [ ] 6.5 Create `internal/api/mcp/controllers/jira.go` with MCP protocol handling for all tools
  - [ ] 6.6 Add account parameter support to all tools with default account resolution
  - [ ] 6.7 Create unit tests for business logic in `internal/app/jira_test.go`
  - [ ] 6.8 Create unit tests for MCP controllers in `internal/api/mcp/controllers/jira_test.go`
  - [ ] 6.9 Register Jira services in `internal/app/register.go` and controllers in `internal/api/mcp/controllers/register.go`
- [ ] 9.0 Integration Testing and Validation
  - [ ] 9.1 Create sample accounts configuration file with multiple accounts for testing
  - [ ] 9.2 Test HTTP client infrastructure with middleware composition and authentication
  - [ ] 9.3 Validate client generation instructions by generating test client code
  - [ ] 9.4 Test end-to-end workflow: Create PR → Read PR → Update PR → Approve PR → Merge PR
  - [ ] 9.5 Test end-to-end workflow: Read Jira ticket → Transition ticket → Manage labels
  - [ ] 9.6 Validate multi-account functionality across all tools with dynamic URL parameter handling
  - [ ] 9.7 Test error handling scenarios: invalid credentials, missing permissions, API failures
  - [ ] 9.8 Verify MCP tool discovery and parameter validation
  - [ ] 9.9 Run full test suite and ensure all tests pass with proper coverage
  - [ ] 9.10 Validate configuration loading with base URLs and accounts file parsing
  - [ ] 9.11 Test default account resolution and named account selection
  - [ ] 9.12 Document configuration setup, usage examples, and client generation instruction usage 