# Tasks: Atlassian MCP Integration

## Requirements

Please read referenced files to understand the problem:
- `doc/prd-atlassian-mcp-integration.md`
- `doc/erd-atlassian-mcp-integration.md`

## Relevant Files

### Source Files
- `internal/services/atlassian_client.go` - HTTP client for Bitbucket and Jira API communication
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
- `internal/services/atlassian_client_test.go` - Unit tests for HTTP client
- `internal/services/atlassian_accounts_test.go` - Unit tests for accounts repository
- `internal/app/bitbucket_test.go` - Unit tests for Bitbucket business logic
- `internal/app/jira_test.go` - Unit tests for Jira business logic
- `internal/api/mcp/controllers/bitbucket_test.go` - Unit tests for Bitbucket MCP controller
- `internal/api/mcp/controllers/jira_test.go` - Unit tests for Jira MCP controller

### Configuration Files
- `internal/config/default.json` - Default configuration (modified)
- `internal/config/local.json` - Local configuration example (modified)
- `.mockery.yaml` - Mockery configuration for generating mocks (if needed)
- `go.mod` - Go module dependencies (if new packages needed)

### Notes

- **Testing Framework:** Use testify for assertions and table-driven tests
- **Architecture:** Follow clean architecture with layers: MCP → App → Services
- **Code Organization:** Follow the existing structure:
  - `internal/api/mcp/` - MCP protocol layer (controllers)
  - `internal/app/` - Application layer (business logic)
  - `internal/services/` - Infrastructure layer (HTTP clients, repositories)
- **Naming Conventions:** Use Go conventions (PascalCase for exported, camelCase for unexported)
- **Authentication:** API token-based authentication for MVP
- **Multi-Account Support:** Named accounts (e.g., "user", "merge-bot") with default account configuration
- **Error Handling:** Use Go's idiomatic error handling with clear, actionable error messages
- **MCP Tools:** Use service-based naming (`bitbucket_*`, `jira_*`)
- **Testing Commands:**
  - `make test` - Run all tests with coverage
  - `go test -v ./internal/path/... --run TestName` - Run specific tests
  - `gow test -v ./internal/path/... --run TestName` - Watch mode for test development
- **Mock Generation:** Use `go generate` with mockery for interface mocks if needed
- **Dependency Injection:** Use uber/dig for dependency management
- **Configuration:** Extend existing viper-based configuration system

## Tasks

- [ ] 1.0 Foundation Infrastructure Setup
  - [ ] 1.1 Extend configuration system to support Atlassian accounts file path in `internal/config/load.go`
  - [ ] 1.2 Update `internal/config/default.json` and `internal/config/local.json` with Atlassian configuration section
  - [ ] 1.3 Create directory structure for new services and app layer components
  - [ ] 1.4 Update `go.mod` with any required new dependencies for HTTP client and JSON handling
  - [ ] 1.5 Set up basic error types and constants for Atlassian integration
- [ ] 2.0 Atlassian HTTP Client and Authentication
  - [ ] 2.1 Create `internal/services/atlassian_client.go` with HTTP client interface definition
  - [ ] 2.2 Implement basic HTTP client struct with authentication support (API token based)
  - [ ] 2.3 Add methods for Bitbucket API calls: `CreatePR`, `GetPR`, `UpdatePR`, `ApprovePR`, `MergePR`
  - [ ] 2.4 Add methods for Jira API calls: `GetTicket`, `TransitionTicket`, `ManageLabels`
  - [ ] 2.5 Implement rate limiting and error handling for API responses
  - [ ] 2.6 Create comprehensive unit tests in `internal/services/atlassian_client_test.go`
  - [ ] 2.7 Add HTTP client to service registration in `internal/services/register.go`
- [ ] 3.0 Accounts Management System
  - [ ] 3.1 Design accounts file JSON schema with multiple named accounts and default account specification
  - [ ] 3.2 Create `internal/services/atlassian_accounts.go` with accounts repository interface
  - [ ] 3.3 Implement `GetDefaultAccount()` method for retrieving configured default account
  - [ ] 3.4 Implement `GetAccountByName(name string)` method for retrieving specific named accounts
  - [ ] 3.5 Add support for separate Bitbucket and Jira credentials per account
  - [ ] 3.6 Implement file reading, parsing, and validation logic with proper error handling
  - [ ] 3.7 Create comprehensive unit tests in `internal/services/atlassian_accounts_test.go`
  - [ ] 3.8 Add accounts repository to service registration in `internal/services/register.go`
- [ ] 4.0 Bitbucket MCP Integration
  - [ ] 4.1 Create `internal/app/bitbucket.go` with business logic service interface and implementation
  - [ ] 4.2 Implement `bitbucket_create_pr` tool: Create pull requests with template support and reviewer assignment
  - [ ] 4.3 Implement `bitbucket_read_pr` tool: Retrieve comprehensive PR details including status and reviews
  - [ ] 4.4 Implement `bitbucket_update_pr` tool: Update PR titles and descriptions with validation
  - [ ] 4.5 Implement `bitbucket_approve_pr` tool: Approve pull requests with proper authorization
  - [ ] 4.6 Implement `bitbucket_merge_pr` tool: Merge PRs with strategy selection and pre-merge validation
  - [ ] 4.7 Create `internal/api/mcp/controllers/bitbucket.go` with MCP protocol handling for all tools
  - [ ] 4.8 Add account parameter support to all tools with default account resolution
  - [ ] 4.9 Create unit tests for business logic in `internal/app/bitbucket_test.go`
  - [ ] 4.10 Create unit tests for MCP controllers in `internal/api/mcp/controllers/bitbucket_test.go`
  - [ ] 4.11 Register Bitbucket services in `internal/app/register.go` and controllers in `internal/api/mcp/controllers/register.go`
- [ ] 5.0 Jira MCP Integration
  - [ ] 5.1 Create `internal/app/jira.go` with business logic service interface and implementation
  - [ ] 5.2 Implement `jira_read_ticket` tool: Retrieve complete ticket information including metadata and comments
  - [ ] 5.3 Implement `jira_transition_ticket` tool: Move tickets through workflow states with validation
  - [ ] 5.4 Implement `jira_manage_labels` tool: Add/remove labels with permission and format validation
  - [ ] 5.5 Create `internal/api/mcp/controllers/jira.go` with MCP protocol handling for all tools
  - [ ] 5.6 Add account parameter support to all tools with default account resolution
  - [ ] 5.7 Create unit tests for business logic in `internal/app/jira_test.go`
  - [ ] 5.8 Create unit tests for MCP controllers in `internal/api/mcp/controllers/jira_test.go`
  - [ ] 5.9 Register Jira services in `internal/app/register.go` and controllers in `internal/api/mcp/controllers/register.go`
- [ ] 6.0 Integration Testing and Validation
  - [ ] 6.1 Create sample accounts configuration file with multiple accounts for testing
  - [ ] 6.2 Test end-to-end workflow: Create PR → Read PR → Update PR → Approve PR → Merge PR
  - [ ] 6.3 Test end-to-end workflow: Read Jira ticket → Transition ticket → Manage labels
  - [ ] 6.4 Validate multi-account functionality across all tools
  - [ ] 6.5 Test error handling scenarios: invalid credentials, missing permissions, API failures
  - [ ] 6.6 Verify MCP tool discovery and parameter validation
  - [ ] 6.7 Run full test suite and ensure all tests pass with proper coverage
  - [ ] 6.8 Validate configuration loading and accounts file parsing
  - [ ] 6.9 Test default account resolution and named account selection
  - [ ] 6.10 Document configuration setup and usage examples 