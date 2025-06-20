# Tasks: Atlassian MCP Integration

## Requirements

Please read referenced files to understand the problem:
- `doc/prd-atlassian-mcp-integration.md`
- `doc/erd-atlassian-mcp-integration.md`

## Relevant Files

### Source Files
- `internal/services/http/middleware/auth.go` - Authentication middleware for HTTP client (✓ created)
- `internal/services/http/middleware/logging.go` - Logging middleware for HTTP client
- `internal/services/http/middleware/error_handling.go` - Error handling middleware for HTTP client
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
- `internal/services/http/middleware/logging_test.go` - Unit tests for logging middleware
- `internal/services/http/middleware/error_handling_test.go` - Unit tests for error handling middleware
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
  - [ ] 1.4 Create unit tests for logging middleware in `internal/services/http/middleware/logging_test.go` (TDD: write tests first)
  - [ ] 1.5 Implement `internal/services/http/middleware/logging.go` - Structured logging middleware as `http.RoundTripper` wrapper
  - [ ] 1.6 Create unit tests for error handling middleware in `internal/services/http/middleware/error_handling_test.go` (TDD: write tests first)
  - [ ] 1.7 Implement `internal/services/http/middleware/error_handling.go` - Error handling middleware as `http.RoundTripper` wrapper
  - [ ] 1.8 Create unit tests for client factory in `internal/services/http/client_factory_test.go` (TDD: write tests first)
  - [ ] 1.9 Create `internal/services/http/client_factory.go` - Factory for composing middleware stack and creating configured `http.Client` instances
  - [ ] 1.10 Establish patterns and interfaces for API client implementation using standard Go HTTP types
  - [ ] 1.11 Register HTTP client infrastructure components in `internal/services/register.go`
- [ ] 2.0 Client Generation Instructions
  - [ ] 2.1 Create `doc/instructions/` directory for client generation instruction documentation
  - [ ] 2.2 Develop `doc/instructions/creating-http-clients.md` with comprehensive templates for struct definitions
  - [ ] 2.3 Add method implementation patterns for converting OpenAPI endpoints to Go methods
  - [ ] 2.4 Document error handling patterns and authentication integration guidelines
  - [ ] 2.5 Create OpenAPI processing guidelines for schema mapping and parameter handling
  - [ ] 2.6 Establish quality assurance instructions for testing, documentation, and validation
  - [ ] 2.7 Define templates for converting OpenAPI schemas to Go structs with proper JSON tags
  - [ ] 2.8 Document response processing patterns for different HTTP status codes and content types
- [ ] 3.0 Configuration System Extension
  - [ ] 3.1 Extend `internal/config/load.go` to support Atlassian accounts file path configuration
  - [ ] 3.2 Add base URLs for Atlassian REST API endpoints in configuration schema
  - [ ] 3.3 Update `internal/config/default.json` with Atlassian configuration section including base URLs
  - [ ] 3.4 Update `internal/config/local.json` with example Atlassian configuration
  - [ ] 3.5 Add validation for Atlassian configuration parameters
  - [ ] 3.6 Create configuration struct types for Atlassian API settings
  - [ ] 3.7 Update `go.mod` with any required new dependencies for HTTP client functionality
- [ ] 4.0 Atlassian HTTP Client Implementation
  - [ ] 4.1 Use official Bitbucket Cloud OpenAPI specification to generate initial client structure types
  - [ ] 4.2 Use official Jira Cloud OpenAPI specification to generate initial client structure types
  - [ ] 4.3 Create `internal/services/atlassian_client.go` with HTTP client interface definition
  - [ ] 4.4 Implement Atlassian HTTP client using established middleware infrastructure
  - [ ] 4.5 Add methods for Bitbucket API calls: `CreatePR`, `GetPR`, `UpdatePR`, `ApprovePR`, `MergePR`
  - [ ] 4.6 Add methods for Jira API calls: `GetTicket`, `TransitionTicket`, `ManageLabels`
  - [ ] 4.7 Implement proper JSON marshaling/unmarshaling for Atlassian API request/response models
  - [ ] 4.8 Add Atlassian-specific error response parsing and meaningful error messages
  - [ ] 4.9 Create comprehensive unit tests in `internal/services/atlassian_client_test.go` with mocked HTTP responses
  - [ ] 4.10 Register Atlassian HTTP client in `internal/services/register.go`
- [ ] 5.0 Accounts Management System
  - [ ] 5.1 Design accounts file JSON schema with multiple named accounts and default account specification
  - [ ] 5.2 Create `internal/services/atlassian_accounts.go` with accounts repository interface
  - [ ] 5.3 Implement `GetDefaultAccount()` method for retrieving configured default account
  - [ ] 5.4 Implement `GetAccountByName(name string)` method for retrieving specific named accounts
  - [ ] 5.5 Add support for separate Bitbucket and Jira credentials per account with dynamic URL parameters
  - [ ] 5.6 Implement file reading, parsing, and validation logic with proper error handling
  - [ ] 5.7 Handle workspace and domain parameters from account configuration (not main config)
  - [ ] 5.8 Create comprehensive unit tests in `internal/services/atlassian_accounts_test.go`
  - [ ] 5.9 Register accounts repository in `internal/services/register.go`
- [ ] 6.0 Bitbucket MCP Integration
  - [ ] 6.1 Create `internal/app/bitbucket.go` with business logic service interface and implementation
  - [ ] 6.2 Implement `bitbucket_create_pr` tool: Create pull requests with template support and reviewer assignment
  - [ ] 6.3 Implement `bitbucket_read_pr` tool: Retrieve comprehensive PR details including status and reviews
  - [ ] 6.4 Implement `bitbucket_update_pr` tool: Update PR titles and descriptions with validation
  - [ ] 6.5 Implement `bitbucket_approve_pr` tool: Approve pull requests with proper authorization
  - [ ] 6.6 Implement `bitbucket_merge_pr` tool: Merge PRs with strategy selection and pre-merge validation
  - [ ] 6.7 Create `internal/api/mcp/controllers/bitbucket.go` with MCP protocol handling for all tools
  - [ ] 6.8 Add account parameter support to all tools with default account resolution
  - [ ] 6.9 Create unit tests for business logic in `internal/app/bitbucket_test.go`
  - [ ] 6.10 Create unit tests for MCP controllers in `internal/api/mcp/controllers/bitbucket_test.go`
  - [ ] 6.11 Register Bitbucket services in `internal/app/register.go` and controllers in `internal/api/mcp/controllers/register.go`
- [ ] 7.0 Jira MCP Integration
  - [ ] 7.1 Create `internal/app/jira.go` with business logic service interface and implementation
  - [ ] 7.2 Implement `jira_read_ticket` tool: Retrieve complete ticket information including metadata and comments
  - [ ] 7.3 Implement `jira_transition_ticket` tool: Move tickets through workflow states with validation
  - [ ] 7.4 Implement `jira_manage_labels` tool: Add/remove labels with permission and format validation
  - [ ] 7.5 Create `internal/api/mcp/controllers/jira.go` with MCP protocol handling for all tools
  - [ ] 7.6 Add account parameter support to all tools with default account resolution
  - [ ] 7.7 Create unit tests for business logic in `internal/app/jira_test.go`
  - [ ] 7.8 Create unit tests for MCP controllers in `internal/api/mcp/controllers/jira_test.go`
  - [ ] 7.9 Register Jira services in `internal/app/register.go` and controllers in `internal/api/mcp/controllers/register.go`
- [ ] 8.0 Integration Testing and Validation
  - [ ] 8.1 Create sample accounts configuration file with multiple accounts for testing
  - [ ] 8.2 Test HTTP client infrastructure with middleware composition and authentication
  - [ ] 8.3 Validate client generation instructions by generating test client code
  - [ ] 8.4 Test end-to-end workflow: Create PR → Read PR → Update PR → Approve PR → Merge PR
  - [ ] 8.5 Test end-to-end workflow: Read Jira ticket → Transition ticket → Manage labels
  - [ ] 8.6 Validate multi-account functionality across all tools with dynamic URL parameter handling
  - [ ] 8.7 Test error handling scenarios: invalid credentials, missing permissions, API failures
  - [ ] 8.8 Verify MCP tool discovery and parameter validation
  - [ ] 8.9 Run full test suite and ensure all tests pass with proper coverage
  - [ ] 8.10 Validate configuration loading with base URLs and accounts file parsing
  - [ ] 8.11 Test default account resolution and named account selection
  - [ ] 8.12 Document configuration setup, usage examples, and client generation instruction usage 