# ERD: Atlassian MCP Integration

## Introduction/Overview

This Engineering Requirements Document outlines the technical implementation of the Atlassian MCP Integration feature. The feature will extend the existing MCP server in `internal/api/mcp/` to provide tools for interacting with Bitbucket and Jira through natural language conversations with AI assistants.

**Goal:** Implement MCP tools that enable developers to manage Bitbucket pull requests and Jira tickets through conversational AI, reducing context switching and automating routine development workflows.

**Target Implementation:** The feature will be implemented with MCP controllers handling protocol communication and application layer services (`internal/app`) containing the actual business logic, following established patterns and maintaining consistency with the current codebase.

## Business Logic

### Bitbucket Operations
1. **Pull Request Creation**: Create PRs with auto-populated metadata, template compliance, and reviewer assignment
2. **Pull Request Reading**: Retrieve comprehensive PR details including status, reviews, and CI results
3. **Pull Request Updates**: Modify PR titles and descriptions while maintaining template standards
4. **Pull Request Approval**: Approve PRs with validation checks and proper authorization
5. **Pull Request Merging**: Execute merge operations with strategy selection and pre-merge validation
6. **Branch Operations**: Standard branch operations using PR creation with appropriate account selection

### Jira Operations
1. **Ticket Reading**: Retrieve complete ticket information including metadata, comments, and workflow status
2. **Ticket Transitions**: Move tickets through workflow states with proper validation
3. **Label Management**: Add/remove labels with permission and format validation

### Authentication & Configuration
- API token-based authentication for MVP implementation
- Multi-account support with user-defined account names (e.g., "user", "merge-bot", "review-bot")
- Each tool accepts optional `account` parameter to specify which configured account to use
- Default account configuration for tools when no specific account is specified
- Local credential storage using existing configuration system
- Graceful handling of authentication failures and permission issues

## High Level Architecture

```
MCP Server (existing)
├── Controllers (existing)
│   ├── math.go, time.go (existing)
│   ├── bitbucket.go (new - MCP protocol handling)
│   └── jira.go (new - MCP protocol handling)
├── Application Layer (internal/app)
│   ├── bitbucket.go (new - business logic)
│   └── jira.go (new - business logic)
├── Services Layer (internal/services)
│   ├── atlassian_client.go (new - HTTP client)
│   └── atlassian_accounts.go (new - accounts repository)
├── Configuration (existing, extended)
│   └── Atlassian account settings
└── External APIs
    ├── Bitbucket Cloud REST API
    └── Jira Cloud REST API
```

### Key Components
1. **MCP Controllers**: Handle MCP protocol communication and delegate to app layer
2. **Application Services**: Contain actual business logic for Bitbucket and Jira operations
3. **Atlassian HTTP Client**: Manages API communication and authentication
4. **Accounts Repository**: Manages multiple named Atlassian accounts with file-based storage
5. **Account Resolution**: Logic to determine which account to use for each operation

## Detailed Architecture

### New Files to Create

#### 1. MCP Controllers
**File: `internal/api/mcp/controllers/bitbucket.go`**
- Implements Bitbucket MCP tools following existing controller patterns
- Handles MCP protocol communication and parameter validation
- Delegates actual business logic to application layer
- Tools to implement:
  - `bitbucket_create_pr`: Create pull requests
  - `bitbucket_read_pr`: Get pull request details  
  - `bitbucket_update_pr`: Update PR title/description
  - `bitbucket_approve_pr`: Approve pull requests
  - `bitbucket_merge_pr`: Merge pull requests with strategy selection
- Each tool includes optional `account` parameter to specify which configured account to use
- Follows existing error handling and response patterns

**File: `internal/api/mcp/controllers/jira.go`**
- Implements Jira MCP tools following existing patterns
- Handles MCP protocol and delegates to app layer
- Tools to implement:
  - `jira_read_ticket`: Get ticket details by ID
  - `jira_transition_ticket`: Move tickets through workflow
  - `jira_manage_labels`: Add/remove ticket labels
- Consistent error handling and response formatting

#### 2. Application Layer Business Logic
**File: `internal/app/bitbucket.go`**
- Contains actual business logic for all Bitbucket operations
- Integrates with accounts repository for account resolution
- Integrates with Atlassian HTTP client for API calls
- Follows existing app layer patterns from `internal/app/`

**File: `internal/app/jira.go`**
- Contains business logic for Jira operations
- Handles ticket workflow validation and transitions
- Manages label operations and permissions
- Integrates with accounts repository and HTTP client

#### 3. HTTP Client Layer
**File: `internal/services/atlassian_client.go`**
- HTTP client for Bitbucket and Jira API communication using direct REST API calls
- Handles multi-account authentication and switching
- Implements rate limiting and error response handling
- Follows existing service patterns in `internal/services/`
- Uses standard Go HTTP client with custom authentication and request/response handling

**File: `internal/services/atlassian_accounts.go`**
- Implements accounts repository pattern (infrastructure layer)
- Provides `GetDefaultAccount()` method for default account retrieval
- Provides `GetAccountByName(name string)` method for named account retrieval
- Reads account configuration from file-based storage
- Handles account validation and error cases
- Supports separate credentials for Bitbucket and Jira services

#### 4. Test Files
**File: `internal/api/mcp/controllers/bitbucket_test.go`**
- Unit tests for MCP controller layer
- Tests parameter validation and MCP protocol handling
- Mocks application layer dependencies

**File: `internal/api/mcp/controllers/jira_test.go`**
- Unit tests for Jira MCP controller
- Tests MCP protocol compliance and error handling

**File: `internal/app/bitbucket_test.go`**
- Unit tests for Bitbucket business logic
- Tests integration with accounts repository and HTTP client
- Mocks both accounts repository and HTTP client responses

**File: `internal/app/jira_test.go`**
- Unit tests for Jira business logic
- Tests workflow validation and API operations
- Mocks accounts repository and HTTP client dependencies

**File: `internal/services/atlassian_accounts_test.go`**
- Unit tests for accounts repository
- Tests file reading, account resolution, and error handling
- Tests default account logic and named account retrieval
- Tests separate Bitbucket and Jira credential handling

**File: `internal/services/atlassian_client_test.go`**
- Unit tests for HTTP client with mocked responses
- Tests multi-account authentication and API communication

### Files to Modify

#### 1. Controller Registration
**File: `internal/api/mcp/controllers/register.go`**
- Add registration for new Bitbucket and Jira controllers
- Follow existing registration patterns

#### 2. Application Layer Registration
**File: `internal/app/register.go`**
- Register new Bitbucket and Jira application services
- Follow existing app layer registration patterns

#### 3. Service Registration  
**File: `internal/services/register.go`**
- Register Atlassian HTTP client service
- Add to dependency injection container

#### 4. Configuration Schema
**File: `internal/config/load.go`** and config JSON files
- Extend existing configuration. Include base URLs for Atlassian REST API endpoints (initially only base URLs, may be extended with additional options like timeouts, retries, etc. in future iterations)

#### 5. Accounts File Structure
**Separate accounts configuration file** (JSON format)
- Contains multiple named accounts with credentials
- Specifies default account name
- Each account includes separate credentials for Bitbucket and Jira services
- Account structure supports workspace ID, usernames, API tokens, and default repositories per service

### Configuration Design

#### Main Configuration
- Application configuration will include base URLs for Atlassian REST API endpoints
- Initially contains only base URLs (e.g., `https://api.bitbucket.org/2.0`, `https://{domain}.atlassian.net/rest/api/3`)
- Designed for future extension with additional HTTP client options such as:
  - Request timeouts
  - Retry policies and exponential backoff settings
  - Rate limiting configuration
  - Connection pooling options
  - Custom headers or user agents

#### Account Configuration
- Separate accounts file contains multiple named Atlassian accounts (e.g., "user", "merge-bot", "review-bot")
- Each account contains separate credentials for Bitbucket and Jira services
- Bitbucket credentials: workspace ID, username, API token, default repository
- Jira credentials: workspace/domain, username, API token, default project
- Accounts file specifies which account is the default
- Tools accept optional `account` parameter to specify which account to use
- If no account specified, accounts repository returns the configured default account
- Accounts repository (services layer) handles file reading, parsing, and account resolution

#### Tool Parameter Structure
- Each tool accepts relevant parameters for the operation (branches, repositories, ticket IDs, etc.)
- Optional `account` parameter for account selection
- Parameter validation handled at MCP controller layer
- Business logic parameters passed to application layer services

## Key Architectural Decisions

### 1. Atlassian API Integration Approach
**Decision**: Use direct REST API calls instead of third-party Go SDKs
**Rationale**: 
- **No Official SDK**: Atlassian does not provide official Go SDKs for Bitbucket Cloud or Jira Cloud APIs
- **Third-party Options**: While community SDKs exist (like `ktrysmt/go-bitbucket` and `ctreminiom/go-atlassian`), they add external dependencies and may not cover all required functionality
- **API Completeness**: Direct REST API access ensures we can use all features and stay current with API updates
- **Control**: Direct HTTP client implementation gives us full control over authentication, error handling, retry logic, and rate limiting
- **Maintainability**: Reduces dependency on external packages that may become outdated or unmaintained
- **Performance**: Custom implementation can be optimized for our specific use cases
- **Security**: Direct implementation allows for better security controls and audit trails

### 2. Layered Architecture
**Decision**: MCP controllers handle protocol, application layer contains business logic
**Rationale**: Follows existing codebase patterns, separates concerns properly, enables better testing and code reuse

### 2. Service-Based Tool Naming
**Decision**: Use `bitbucket_*` and `jira_*` prefixes for tool names
**Rationale**: Clear service separation, easy to understand, and follows logical grouping principles

### 3. API Token Authentication
**Decision**: Start with API token authentication for MVP
**Rationale**: Simplest implementation path, adequate security for MVP, can evolve to OAuth later without breaking changes

### 4. Accounts Repository Pattern
**Decision**: Separate accounts repository with file-based storage and dedicated methods (`GetDefaultAccount`, `GetAccountByName`)
**Rationale**: Clear separation of concerns, easier testing, file-based storage is simple and reliable, repository pattern provides clean interface

### 5. Multi-Account Configuration
**Decision**: Support multiple named accounts with flexible naming (e.g., "user", "merge-bot", "review-bot")
**Rationale**: More flexible than binary user/bot approach, allows teams to configure accounts based on their workflows, enables future expansion

### 6. Default Account Strategy
**Decision**: Accounts file specifies default account, tools accept optional account parameter
**Rationale**: Ensures system always has credentials to work with, while allowing per-operation account selection when needed

### 7. Parameter-Based Context
**Decision**: Tools receive all context (branch names, repositories) as required parameters
**Rationale**: Clear contract, predictable behavior, easier testing, and follows existing MCP tool patterns

### 8. Template Reading Approach
**Decision**: Simple template reading from repository settings via API
**Rationale**: Sufficient for MVP, avoids complex template management, can be enhanced later

## REST API Implementation Details

### Bitbucket Cloud REST API
- **Base URL**: `https://api.bitbucket.org/2.0`
- **Authentication**: HTTP Basic Auth with App Password or OAuth 2.0
- **Key Endpoints**:
  - Pull Requests: `/repositories/{workspace}/{repo_slug}/pullrequests`
  - Branches: `/repositories/{workspace}/{repo_slug}/refs/branches`
  - Repositories: `/repositories/{workspace}/{repo_slug}`
  - Users: `/user`, `/workspaces/{workspace}/members`
- **Dynamic URL Parameters**: The `{workspace}` parameter is taken from the user's account configuration (not from main config), as it varies per user account

### Jira Cloud REST API  
- **Base URL**: `https://{domain}.atlassian.net/rest/api/3`
- **Authentication**: HTTP Basic Auth with API Token or OAuth 2.0
- **Key Endpoints**:
  - Issues: `/issue/{issueIdOrKey}`
  - Transitions: `/issue/{issueIdOrKey}/transitions`
  - Labels: `/issue/{issueIdOrKey}`
  - Projects: `/project`
- **Dynamic URL Parameters**: The `{domain}` parameter is taken from the user's account configuration (not from main config), as it varies per user account

### HTTP Client Infrastructure Foundation
**Note**: The codebase currently lacks base infrastructure for external API calls. This implementation will establish foundational patterns for future API integrations.

#### Core HTTP Client Architecture
- **Standard Go Foundation**: All HTTP clients must be built on standard Go `http.Client` and `http.RoundTripper` types
- **Wrapper-Based Design**: Authentication, logging, and other cross-cutting concerns implemented as middleware wrappers around standard HTTP components
- **Composable Architecture**: Individual concerns (auth, logging, error handling) can be composed together via wrapping pattern
- **Standard Client Usage**: Actual API client implementations use standard `http.Client` instances with composed middleware stack

#### Infrastructure Components to Implement
1. **Authentication Middleware**: `http.RoundTripper` wrapper that adds authentication headers
2. **Logging Middleware**: `http.RoundTripper` wrapper for structured request/response logging
3. **Error Handling Middleware**: `http.RoundTripper` wrapper for consistent error response handling
4. **Client Factory**: Utility to compose middleware stack and create configured `http.Client` instances

**Note**: Rate limiting is out of scope for initial implementation and can be added later if needed.

#### API Client Implementation Pattern
- **Type Definitions**: Each API client defines Go structs corresponding to API request/response models
- **Method Definitions**: Each API client implements methods that correspond to specific REST API endpoints
- **Standard HTTP Usage**: Clients use standard `http.Client.Do()` with proper request construction and response parsing
- **JSON Marshaling**: Standard `encoding/json` for request/response serialization
- **Error Mapping**: Convert HTTP errors to domain-specific error types

### Atlassian HTTP Client Implementation
- **Request/Response Handling**: Custom JSON marshaling/unmarshaling for Atlassian API structures
- **Error Handling**: Atlassian-specific error response parsing and meaningful error messages
- **Authentication**: HTTP Basic Auth implementation using middleware wrapper
- **Logging**: Structured logging for API requests/responses for debugging

## Testing Strategy

### Unit Testing Approach
- **Follow Existing Patterns**: Use same testing structure as `math_test.go` and `time_test.go`
- **Mock HTTP Responses**: Create mock Atlassian API responses for all test scenarios
- **Test Coverage Requirements**:
  - All tool functions with success and error scenarios
  - Parameter validation for all tools
  - Authentication handling (valid/invalid tokens)
  - Bot account logic paths
  - Error response formatting

### Test Organization
```
internal/api/mcp/controllers/
├── bitbucket_test.go (MCP protocol tests)
├── jira_test.go (MCP protocol tests)
internal/app/
├── bitbucket_test.go (business logic tests)
├── jira_test.go (business logic tests)
internal/services/
├── atlassian_accounts_test.go (accounts repository tests)
└── atlassian_client_test.go (HTTP client tests)
```

### Test Data Management
- Create mock JSON responses for all Atlassian API endpoints
- Use table-driven tests for parameter validation scenarios
- Test multiple account configurations and account resolution logic
- Test separate Bitbucket and Jira credential handling in accounts repository
- Mock application layer dependencies in controller tests
- Mock accounts repository and HTTP client in application layer tests

## Implementation Phases

### Phase 1: Foundation (Week 1-2)
1. **HTTP Client Infrastructure**: Implement base HTTP client infrastructure with middleware pattern
   - Create authentication, logging, and error handling middleware
   - Implement client factory for composing middleware stack
   - Establish patterns for API client implementation
2. **AI Client Generation Instructions**: Develop comprehensive instructions for AI models to generate API clients
   - Create templates and patterns for type definitions and method implementations
   - Document how to use OpenAPI specifications to generate client code
   - Establish guidelines for error handling and response parsing
3. **Atlassian HTTP Client**: Implement Atlassian-specific HTTP client using established infrastructure
   - Use OpenAPI specifications for Bitbucket and Jira APIs to generate initial client structure
   - Implement multi-account authentication using middleware pattern
4. Extend existing configuration system for Atlassian accounts file path
5. Create accounts repository in services layer with `GetDefaultAccount` and `GetAccountByName` methods
6. Design account structure to support separate Bitbucket and Jira credentials
7. Create application layer services and registration infrastructure

### Phase 2: Bitbucket Core (Week 3-4)
1. Implement `bitbucket_create_pr` MCP controller and app level component
2. Implement `bitbucket_read_pr` MCP controller and app level component
3. Implement `bitbucket_update_pr` MCP controller and app level component
4. Add comprehensive unit tests for both layers

### Phase 3: Bitbucket Advanced (Week 5-6)
1. Implement `bitbucket_approve_pr` controller and app level component
2. Implement `bitbucket_merge_pr` controller and app level component
3. Test multi-account functionality across all tools
4. Complete Bitbucket integration testing

### Phase 4: Jira Integration (Week 7-8)
1. Implement `jira_read_ticket` controller and app level component
2. Implement `jira_transition_ticket` controller and app level component  
3. Implement `jira_manage_labels` controller and app level component
4. Complete test coverage and integration validation

## AI Model Client Generation Strategy

### OpenAPI Specification Usage
- **Bitbucket Cloud API**: Use official Bitbucket Cloud OpenAPI specification to generate client types and method signatures
- **Jira Cloud API**: Use official Jira Cloud OpenAPI specification to generate client types and method signatures
- **Code Generation Approach**: Develop instructions that allow AI models to generate client code by referencing OpenAPI specs
- **Type Safety**: Ensure generated clients provide strong typing for request/response models

### AI Model Instructions Development
As part of Phase 1 implementation, create comprehensive instructions for AI models that include:

#### Client Generation Templates
- **Struct Definition Patterns**: Templates for converting OpenAPI schemas to Go structs
- **Method Implementation Patterns**: Templates for implementing API endpoint methods
- **Error Handling Patterns**: Standard error handling and mapping approaches
- **Authentication Integration**: How to integrate with middleware authentication system

#### OpenAPI Processing Guidelines
- **Schema Mapping**: Instructions for converting OpenAPI schemas to Go types
- **Endpoint Mapping**: Instructions for converting OpenAPI paths to Go methods
- **Parameter Handling**: Guidelines for handling query parameters, path parameters, and request bodies
- **Response Processing**: Patterns for handling different response types and status codes

#### Quality Assurance Instructions
- **Testing Requirements**: Guidelines for generating comprehensive test coverage
- **Documentation Standards**: Requirements for code documentation and examples
- **Validation Logic**: Instructions for implementing request/response validation
- **Error Scenarios**: Guidelines for handling various error conditions

### Implementation Deliverables
1. **Base Infrastructure**: HTTP client infrastructure with middleware pattern
2. **AI Instructions Document**: Comprehensive guide for AI models to generate API clients
3. **Reference Implementation**: Atlassian clients generated using the AI instructions
4. **Validation Process**: Testing and validation of AI-generated client code

## Open Questions

### Technical Implementation
1. **Configuration Security**: Should API tokens be encrypted at rest, or is environment variable storage sufficient for MVP?

2. **Error Context**: How much detail should we include in error responses to help users troubleshoot issues without exposing sensitive information?

### User Experience
3. **Default Repository Detection**: If no repository is specified, should we attempt to detect it from git remote or require explicit specification?

4. **Reviewer Auto-Detection**: Should we implement simple reviewer detection based on git history or CODEOWNERS files, even though it's marked as non-goal?

5. **Branch Naming Conventions**: Should tools validate or enforce branch naming conventions, or accept any valid git branch name?

### Future Considerations
6. **Migration Path**: How should we design the authentication system to easily migrate from API tokens to OAuth 2.0 in future versions?

7. **Multi-Workspace Support**: Should the configuration support multiple Atlassian workspaces from the beginning, or add this later?

8. **Caching Strategy**: Should we implement any caching for frequently accessed data (PR details, ticket information) to improve performance?

---

**Document Status:** Ready for Development  
**Target Audience:** Junior Developer  
**Implementation Start:** Upon technical review approval  
**Estimated Timeline:** 8 weeks for complete implementation  
**Dependencies:** Existing MCP server architecture, Atlassian Cloud API access 