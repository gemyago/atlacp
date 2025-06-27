# Tasks: Atlassian MCP Integration

## Requirements

Please read referenced files to understand the problem:
- `doc/prd-atlassian-mcp-integration.md`
- `doc/erd-atlassian-mcp-integration.md`

## Relevant Files

### Source Files
- `internal/services/http/middleware/auth.go` - Authentication middleware for HTTP client (✓ created)
- `internal/services/http/middleware/logging.go` - Logging middleware for HTTP client (✓ created)
- `internal/services/http/middleware/error_handling.go` - Error handling middleware for HTTP client with response body logging (✓ created)
- `internal/services/http/client_factory.go` - HTTP client factory with middleware composition (✓ created)
- `internal/services/http/send_request.go` - Shared SendRequest function with generic body and target types (✓ created)
- `doc/instructions/creating-http-clients.md` - Comprehensive instructions for API client implementation patterns (✓ created)
- `internal/services/bitbucket/client.go` - Bitbucket API client implementation (✓ created)
- `internal/services/bitbucket/pullrequest.go` - Bitbucket pull request models (✓ created)
- `internal/services/bitbucket/create_pr.go` - CreatePR operation implementation (✓ created)
- `internal/services/bitbucket/get_pr.go` - GetPR operation implementation (✓ created)
- `internal/services/bitbucket/update_pr.go` - UpdatePR operation implementation (✓ created)
- `internal/services/bitbucket/approve_pr.go` - ApprovePR operation implementation (✓ created)
- `internal/services/bitbucket/merge_pr.go` - MergePR operation implementation (✓ created)
- `internal/services/jira/client.go` - Jira API client implementation (✓ created)
- `internal/services/jira/ticket.go` - Jira ticket models (✓ created)
- `internal/services/jira/get_ticket.go` - GetTicket operation implementation (✓ created)
- `internal/services/jira/transition_ticket.go` - TransitionTicket operation implementation (✓ created)
- `internal/services/jira/manage_labels.go` - ManageLabels operation implementation (✓ created)
- `internal/app/atlassian_accounts.go` - Atlassian account models in application layer (✓ created)
- `internal/app/ports.go` - Application layer ports (interfaces) for repositories (✓ created)
- `internal/app/auth.go` - Token provider implementation for API authentication (✓ created)
- `doc/atlassian-accounts-schema.md` - JSON schema documentation for Atlassian accounts (✓ created)
- `examples/atlassian-accounts-stub.json` - Example Atlassian accounts configuration (✓ created)
- `internal/services/atlassian_client.go` - Atlassian-specific HTTP client implementation
- `internal/services/atlassian_accounts.go` - Accounts repository for managing multiple named Atlassian accounts (✓ created)
- `internal/app/bitbucket.go` - Business logic for Bitbucket operations (✓ created with CreatePR implemented)
- `internal/app/jira.go` - Business logic for Jira operations
- `internal/api/mcp/controllers/bitbucket.go` - MCP controller for Bitbucket tools
- `internal/api/mcp/controllers/jira.go` - MCP controller for Jira tools
- `internal/api/mcp/controllers/register.go` - Controller registration (modified)
- `internal/app/register.go` - App layer registration (modified)
- `internal/services/register.go` - Service registration with Atlassian HTTP clients (✓ modified)
- `internal/config/load.go` - Configuration loading (modified)

### Test Files
- `internal/services/http/middleware/auth_test.go` - Unit tests for authentication middleware (✓ created)
- `internal/services/http/middleware/logging_test.go` - Unit tests for logging middleware (✓ created)
- `internal/services/http/middleware/error_handling_test.go` - Unit tests for error handling middleware (✓ created)
- `internal/services/http/client_factory_test.go` - Unit tests for client factory
- `internal/services/bitbucket/client_test.go` - Unit tests for Bitbucket client (✓ created)
- `internal/services/bitbucket/create_pr_test.go` - Unit tests for CreatePR operation (✓ created)
- `internal/services/bitbucket/get_pr_test.go` - Unit tests for GetPR operation (✓ created)
- `internal/services/bitbucket/update_pr_test.go` - Unit tests for UpdatePR operation (✓ created)
- `internal/services/bitbucket/approve_pr_test.go` - Unit tests for ApprovePR operation (✓ created)
- `internal/services/bitbucket/merge_pr_test.go` - Unit tests for MergePR operation (✓ created)
- `internal/services/jira/client_test.go` - Unit tests for Jira client (✓ created)
- `internal/services/jira/get_ticket_test.go` - Unit tests for GetTicket operation (✓ created)
- `internal/services/jira/transition_ticket_test.go` - Unit tests for TransitionTicket operation (✓ created)
- `internal/services/jira/manage_labels_test.go` - Unit tests for ManageLabels operation (✓ created)
- `internal/services/atlassian_client_test.go` - Unit tests for Atlassian HTTP client
- `internal/services/atlassian_accounts_test.go` - Unit tests for accounts repository (✓ created)
- `internal/app/bitbucket_test.go` - Unit tests for Bitbucket business logic (✓ created with tests for CreatePR)
- `internal/app/auth_test.go` - Unit tests for token provider (✓ created)
- `internal/app/jira_test.go` - Unit tests for Jira business logic
- `internal/api/mcp/controllers/bitbucket_test.go` - Unit tests for Bitbucket MCP controller
- `internal/api/mcp/controllers/jira_test.go` - Unit tests for Jira MCP controller

### Configuration Files
- `internal/config/default.json` - Default configuration (modified)
- `internal/config/local.json` - Local configuration example (modified)
- `doc/instructions/creating-http-clients.md` - Instructions for creating HTTP clients
- `.mockery.yaml` - Mockery configuration for generating mocks (if needed)
- `go.mod` - Go module dependencies (if new packages needed)

### Example and Testing Files
- `examples/docker-compose.yml` - Docker compose configuration for testing MCP integration (✓ created)
- `examples/test-client/package.json` - Node.js package for test client (✓ created)
- `examples/test-client/test-bitbucket-workflow.js` - Complete Bitbucket workflow test script (✓ created)
- `examples/test-client/test-multi-account.js` - Multi-account functionality test script (✓ created)
- `examples/README.md` - Comprehensive documentation for examples and testing (✓ created)

## Tasks

- [x] 1.0 HTTP Client Infrastructure Foundation
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
- [x] 3.0 Atlassian HTTP Clients Implementation
  - [x] 3.1 Update `internal/config/default.json` with Atlassian configuration section including base URLs
  - [x] 3.2 Inject Atlassian client configuration into DI as per `internal/config/provide.go`
  - [x] 3.3 Find official Bitbucket Cloud OpenAPI specification and add it to `internal/services/bitbucket/openapi.yaml` (Used community-maintained spec from magmax/atlassian-openapi)
  - [x] 3.4 Find official Jira Cloud OpenAPI specification and add it to `internal/services/jira/openapi.yaml`
  - [x] 3.5 Create Bitbucket client based on the openapi and `doc/instructions/creating-http-clients.md` instruction. Add methods for following Bitbucket API calls only: `CreatePR`, `GetPR`, `UpdatePR`, `ApprovePR`, `MergePR`
  - [x] 3.6 Create Jira client based on the openapi and `doc/instructions/creating-http-clients.md` instruction. Add methods for following Jira API calls only: `GetTicket`, `TransitionTicket`, `ManageLabels`
  - [x] 3.7 Add Atlassian-specific error response parsing and meaningful error messages (research if needed)
  - [x] 3.8 Register Atlassian HTTP clients in `internal/services/register.go`
- [x] 4.0 Accounts Management System
  - [x] 4.1 Design accounts file JSON schema with multiple named accounts and default account specification
  - [x] 4.2 Create `internal/app/ports.go` with Atlassian accounts repository interface
  - [x] 4.3 Implement `GetDefaultAccount()` method for retrieving configured default account
  - [x] 4.4 Implement `GetAccountByName(name string)` method for retrieving specific named accounts
  - [x] 4.5 Add support for separate Bitbucket and Jira credentials per account with dynamic URL parameters
  - [x] 4.6 Implement file reading, parsing, and validation logic with proper error handling
  - [x] 4.7 Handle workspace and domain parameters from account configuration (not main config)
  - [x] 4.8 Create comprehensive unit tests in `internal/services/atlassian_accounts_test.go`
  - [x] 4.9 Register accounts repository in `internal/services/register.go`
- [x] 5.0 Bitbucket MCP Integration
  - [x] 5.1 Create `internal/app/bitbucket.go` with business logic service with all method stubs and initial skeleton for unit tests
  - [x] 5.2 Implement `bitbucket_create_pr` tool: Create pull requests with template support and reviewer assignment
  - [x] 5.3 Implement `bitbucket_read_pr` tool: Retrieve comprehensive PR details including status and reviews
  - [x] 5.4 Implement `bitbucket_update_pr` tool: Update PR titles and descriptions with validation
  - [x] 5.5 Implement `bitbucket_approve_pr` tool: Approve pull requests with proper authorization
  - [x] 5.6 Implement `bitbucket_merge_pr` tool: Merge PRs with strategy selection and pre-merge validation
  - [x] 5.7 Create `internal/api/mcp/controllers/bitbucket.go` with MCP protocol handling for all tools
  - [x] 5.8 Add account parameter support to all tools with default account resolution
  - [x] 5.9 Create unit tests for business logic in `internal/app/bitbucket_test.go` (partial - CreatePR tests implemented)
  - [x] 5.10 Create unit tests for MCP controllers in `internal/api/mcp/controllers/bitbucket_test.go`
  - [x] 5.11 Register Bitbucket services in `internal/app/register.go` and controllers in `internal/api/mcp/controllers/register.go`
- [ ] 6.0 Bitbucket Integration Testing and Usage Documentation
  - [x] 6.1 Create example docker compose
  - [x] 6.2 Support user token authentication (Basic token)
  - [x] 6.3 Allow draft PRs to be created
  - [x] 6.4 Test end-to-end bitbucket workflow (from AI code editor): Create PR → Read PR → Update PR → Approve PR → Merge PR. Automate this with instructions for AI code editor.
  - [x] 6.5 Validate multi-account functionality: Create PR as user -> Approve PR as bot -> Merge PR as user. Automate this with instructions for AI code editor.
  - [x] 6.6 Test default account resolution and named account selection: Create PR as default account -> Approve PR as bot. Automate this with instructions for AI code editor.
  - [x] 6.7 Created comprehensive testing document at `doc/testing/bitbucket-mcp-integration-tests.md`.
  - [ ] 6.8 Update README with quick start and usage instruction
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
  - [ ] 7.10 Prepare Jira Integration Testing and Usage Documentation task