author:	gemyago
association:	owner
edited:	true
status:	none
--
# Engineering Requirements Document: Bitbucket PR Review Automation

## Introduction/Overview

This document outlines the requirements for implementing automated Bitbucket pull request (PR) review functionality within the existing MCP (Model Context Protocol) server architecture. The feature will enable AI coding assistants (Cursor, VSCode Copilot) to automatically analyze PR diffs, check for coding best practices, and suggest improvements through automated review comments.

**Problem Statement:** Development teams spend significant time reviewing obvious issues in pull requests (code style, basic logic errors, common anti-patterns) rather than focusing on complex business logic and architectural decisions.

**Goal:** Reduce developer time spent on reviewing obvious issues by providing AI-powered automated PR analysis that can identify and comment on common problems, allowing human reviewers to focus on higher-level concerns.

## Business Logic

### Core Workflow
1. **PR Analysis Request**: AI coding assistant receives user request to review a specific Bitbucket PR
2. **PR Data Retrieval**: AI uses existing `bitbucket_read_pr` tool to get PR metadata
3. **File Discovery**: AI uses new `bitbucket_get_pr_diffstat` tool to get list of changed files
4. **Diff Analysis**: AI uses new `bitbucket_get_pr_diff` tool to retrieve changes for analysis
5. **File Content Review**: AI uses new `bitbucket_get_file_content` tool for full context when needed
6. **Issue Detection**: AI analyzes changes and identifies problems based on coding standards
7. **Comment Generation**: AI uses new `bitbucket_add_pr_comment` tool to add review comments (both general and inline)
8. **Review Decision**: AI uses existing `bitbucket_approve_pr` or new `bitbucket_request_pr_changes` based on findings

### Tool Integration Pattern
The implementation follows the established pattern where:
1. **MCP tools provide data access**: Each tool handles a specific Bitbucket API operation
2. **AI assistants orchestrate workflow**: AI combines multiple tool calls to complete review process
3. **Service layer handles business logic**: App services contain the actual API integration logic
4. **Error handling is consistent**: All tools use established error patterns from existing Bitbucket tools

### Analysis Capabilities Enabled
- **Diff-based review**: Direct access to line-by-line changes for targeted analysis  
- **Full file context**: Ability to retrieve complete file content when diff context isn't sufficient
- **Precise commenting**: Line-specific comments tied to actual problematic code
- **Approval workflow**: Proper approval/request changes flow using existing mechanisms

## High Level Architecture

### Components Involved
1. **Existing BitbucketController** (Enhanced): Add 4 new MCP tools for PR review
2. **Existing Bitbucket Services** (Extended): Add new client methods for diff, comments, and file content
3. **AI Coding Assistant** (External): Orchestrates review workflow using combination of existing and new tools

### Integration Approach
- **Extend existing patterns**: Add new tools to current `BitbucketController` following established conventions
- **Reuse authentication**: Leverage existing Atlassian account authentication system
- **Consistent error handling**: Use same error patterns as existing Bitbucket tools

### Tool Workflow Example
```
1. AI calls bitbucket_read_pr (existing) → Get PR metadata
2. AI calls bitbucket_get_pr_diffstat (new) → Get list of changed files
3. AI calls bitbucket_get_pr_diff (new) → Get detailed code changes
4. AI analyzes diff and identifies issues
5. AI calls bitbucket_add_pr_comment (new) → Add review comments (general or inline)
6. AI calls bitbucket_approve_pr (existing) OR bitbucket_request_pr_changes (new)
```

## Detailed Architecture

### New MCP Tools Required

**Note**: The existing `bitbucket_read_pr` tool already provides PR metadata and the `bitbucket_approve_pr` tool handles PR approval. The following tools add the missing functionality needed specifically for automated PR review.

#### 1. `bitbucket_get_pr_diffstat`
**Purpose**: Get a list of files changed in the pull request with basic statistics

**API Endpoint**: `GET /repositories/{workspace}/{repo_slug}/diffstat/{spec}` (where spec is `source_commit..destination_commit`)

**Parameters**:
- `pr_id`: PR ID number (number, required)
- `repo_owner`: Repository owner (username/workspace) (string, required)
- `repo_name`: Repository name (slug) (string, required)
- `file_paths` (optional): Array of specific file paths to limit diffstat to. Maps to the `path` query parameter.
- `context` (optional): Number of context lines around changes (number, default: 3)
- `account` (optional): Atlassian account name to use (string)

**Implementation Notes**:
- Use PR's source and destination commit hashes to construct the spec parameter
- The `path` query parameter can be repeated for multiple specific files: `?path=file1.js&path=file2.py`
- This allows LLM to get stats for all files or filter to specific files of interest

**Returns**:
- JSON array of diffstat objects with file statistics
- Each file includes: `status` (modified/added/removed), `lines_added`, `lines_removed`, `old` and `new` file paths
- If results are paginated, the MCP tool will automatically fetch all pages. It will feed the results to the LLM as it receives them using stream.

**Reference**: [Bitbucket API Documentation](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-diffstat-spec-get)

#### 2. `bitbucket_get_pr_diff`
**Purpose**: Retrieve detailed diff information for specific files or entire PR

**API Endpoint**: `GET /repositories/{workspace}/{repo_slug}/diff/{spec}` (where spec is `source_commit..destination_commit`)

**Parameters**:
- `pr_id`: PR ID number (number, required)
- `repo_owner`: Repository owner (username/workspace) (string, required)
- `repo_name`: Repository name (slug) (string, required)
- `file_paths` (optional): Array of specific file paths to get diff for. Maps to the `path` query parameter.
- `context` (optional): Number of context lines around changes (number, default: 3)
- `account` (optional): Atlassian account name to use (string)

**Implementation Notes**:
- Use PR's source and destination commit hashes to construct the spec parameter: `source_commit..destination_commit`
- The `path` query parameter can be repeated for multiple specific files: `?path=file1.js&path=file2.py`
- The `context` parameter controls the number of context lines shown around changes
- This allows LLM to get diff for all files or filter to specific files of interest

**Returns**:
- Raw diff content (unified diff format)
- When `file_paths` is omitted: complete diff for all files in the PR
- When `file_paths` is specified: diff for only those specific files

**Note**: The LLM can decide whether to fetch all files at once or filter to specific files using the path parameters.

#### 3. `bitbucket_get_file_content`
**Purpose**: Retrieve full file content from specific branch/commit

**Parameters**:
- `repo_owner`: Repository owner (username/workspace) (string, required)
- `repo_name`: Repository name (slug) (string, required)
- `commit_hash`: Commit hash or branch name (string, required)
- `file_path`: Path to the file (string, required)
- `account` (optional): Atlassian account name to use (string)

**Returns**:
- Complete file content as text
- File metadata (size, type, encoding)

#### 4. `bitbucket_add_pr_comment`
**Purpose**: Add comments to pull requests (both general and inline comments)

**API Endpoint**: `POST /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments`

**Parameters**:
- `pr_id`: PR ID number (number, required)
- `repo_owner`: Repository owner (username/workspace) (string, required)
- `repo_name`: Repository name (slug) (string, required)
- `comment_text`: The comment content (string, required)
- `file_path` (optional): Path to the file for inline comments, relative to repository root (string)
- `line_number_from` (optional): The comment's anchor line in the old version of the file (number)
- `line_number_to` (optional): The comment's anchor line in the new version of the file (number)
- `account` (optional): Atlassian account name to use (string)

**API Request Body Structure**:
For general PR comments:
```json
{
  "content": {
    "raw": "[comment_text]"
  }
}
```

For inline comments:
```json
{
  "content": {
    "raw": "[comment_text]"
  },
  "inline": {
    "path": "[file_path]", 
    "from": [line_number_from],
    "to": [line_number_to]
  }
}
```

**Implementation Notes**:
- If `file_path` is provided, creates an inline comment with the `inline` object
- If `file_path` is omitted, creates a general PR comment
- This extends functionality rather than creating a new service since no comment API is currently implemented

**Returns**:
- Comment ID for tracking
- Success/failure status

#### 5. `bitbucket_request_pr_changes`
**Purpose**: Request changes by removing approval from the pull request

**Parameters**:
- `pr_id`: PR ID number (number, required)
- `repo_owner`: Repository owner (username/workspace) (string, required)
- `repo_name`: Repository name (slug) (string, required)
- `account` (optional): Atlassian account name to use (string)

**Returns**:
- Request status
- Timestamp of approval removal

### File Structure Updates

**Note**: Following existing codebase patterns, new tools will be added to the existing `BitbucketController` rather than creating a separate PR review controller.

#### New Files to Create:
```
internal/services/bitbucket/
├── get_diffstat.go           # Get diff statistics using commits diffstat API
├── get_diffstat_test.go      # Tests for diffstat
├── get_diff.go               # Get diff using commits diff API
├── get_diff_test.go          # Tests for diff
├── get_file_content.go       # Get file content from branches/commits
├── get_file_content_test.go  # Tests for file content
├── add_pr_comment.go         # Add comments to PRs (general and inline)
├── add_pr_comment_test.go    # Tests for PR comments
└── request_pr_changes.go     # Request changes by removing approval
└── request_pr_changes_test.go # Tests for request changes
```

#### Files to Update:
```
internal/api/mcp/controllers/
├── bitbucket.go              # Add new MCP tools to existing BitbucketController
├── bitbucket_test.go         # Add tests for new tools
└── mock_bitbucket_service.go # Add new service methods to mock

internal/services/bitbucket/
├── models.go                 # Add new model structures (Comment, Diff, etc.)
└── ENDPOINTS.md              # Document new API endpoints

internal/app/
├── bitbucket.go              # Add new service methods for PR review
├── bitbucket_test.go         # Add tests for new service methods
└── ports.go                  # Add new methods to bitbucketClient interface
```

### API Integration Details

#### New Bitbucket Cloud API Endpoints to Implement:
- ✅ `GET /2.0/repositories/{workspace}/{repo_slug}/diffstat/{spec}` - Get diff statistics for commits (with optional `context` parameter)
- ✅ `GET /2.0/repositories/{workspace}/{repo_slug}/diff/{spec}` - Get diff for commits (with optional `context` parameter, supports file-specific diffs)
- ✅ `GET /2.0/repositories/{workspace}/{repo_slug}/src/{commit}/{path}` - Get file content from commit/branch
- ✅ `POST /2.0/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments` - Add PR comments (supports optional `inline` object for line-specific comments)
- ✅ `DELETE /2.0/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/approve` - Remove approval (request changes)

#### Existing Endpoints (Already Implemented):
- `GET /2.0/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}` - PR details (`bitbucket_read_pr`)
- `POST /2.0/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/approve` - Approve PR (`bitbucket_approve_pr`)

### Error Handling
- **Authentication Errors**: Proper error messages for invalid/expired tokens
- **Permission Errors**: Clear messages when user lacks PR review permissions  
- **Rate Limiting**: Implement backoff strategies for API rate limits
- **Large PR Handling**: Chunking strategy for PRs with many files or large diffs
- **Network Failures**: Retry mechanisms with exponential backoff

## Key Architectural Decisions

### 1. Extend Existing BitbucketController
**Decision**: Add new PR review tools to the existing `BitbucketController` rather than creating a new controller
**Rationale**: Maintains consistency with existing codebase patterns and avoids duplication of shared functionality like authentication and error handling

### 2. Reuse Existing Infrastructure
**Decision**: Leverage existing `bitbucket_read_pr` and `bitbucket_approve_pr` tools instead of reimplementing
**Rationale**: Avoids code duplication and maintains consistency with established tool interfaces

### 3. Focus on Missing Functionality Only
**Decision**: Implement only the 5 missing tools needed for PR review automation (diffstat, diff, file content, inline comments, request changes)
**Rationale**: Minimal implementation approach that adds only what's necessary for the PR review workflow, with proper file discovery capability

### 4. Follow Existing Naming Conventions
**Decision**: Use `bitbucket_` prefix and existing parameter patterns (`pr_id`, `repo_owner`, `repo_name`, `account`)
**Rationale**: Ensures consistency with existing tools and familiar interface for users

### 5. Granular Tool Design
**Decision**: Separate tools for file discovery (diffstat), diff retrieval, file content access, and inline comment creation
**Rationale**: Allows AI assistants to compose complex review workflows from simple, focused operations with clear separation of concerns

### 6. Line-Level Comment Focus
**Decision**: Provide dedicated inline comment tool with required file path and line number parameters
**Rationale**: Matches Bitbucket API requirements and focuses on code review use case where line-specific feedback is most valuable

## Testing Strategy

### Test-Driven Development (TDD) Approach
1. **Write tests first** for each new service function
2. **Implement minimum code** to make tests pass
3. **Refactor** while maintaining test coverage

### Testing Levels

#### Unit Tests
- **Service Layer**: Test each new Bitbucket API client method following existing patterns
- **Controller Layer**: Test new MCP tool request/response handling in existing `BitbucketController`
- **Mock Dependencies**: Extend existing `mock_bitbucket_service.go` with new methods

#### Integration Tests
- **API Integration**: Test new Bitbucket API endpoints with test repositories
- **Tool Registration**: Verify new tools are properly registered with existing `BitbucketController`
- **Workflow Testing**: Test complete PR review scenarios using combination of existing and new tools

#### Test Coverage Requirements
- **Follow existing patterns**: Match the test coverage approach used by existing Bitbucket tools
- **Critical path testing**: Ensure diff parsing, comment creation, and approval removal work correctly
- **Error scenario coverage**: Test API failures, malformed responses, and permission errors

### Test Data Strategy
- **Reuse existing test infrastructure**: Use existing test repository and mock patterns
- **Diff test cases**: Create PRs with various diff scenarios (small/large changes, binary files, etc.)
- **Mock API responses**: Add new Bitbucket API response mocks for diff, comments, and file content endpoints

## Implementation Decisions

Based on requirements analysis and AI assistant optimization considerations, the following decisions have been made:

1. **Diff Context Lines**: Default to 3 context lines to provide sufficient context for AI analysis while maintaining manageable response sizes.

2. **Comment API Implementation**: Create a unified comment tool that supports both general and inline comments using the same API endpoint. The `inline` object is optional - when provided with file path and line numbers, creates inline comments; when omitted, creates general PR comments.

3. **File Content Caching**: No caching or size limits will be implemented initially to keep the design simple and avoid premature optimization.

4. **Error Response Formatting**: Follow existing Bitbucket tool patterns for consistent error handling across all MCP tools in the controller.

5. **Diff Format**: The Bitbucket API returns unified diff format (not JSON). This will be returned as-is to the LLM, which can parse unified diff effectively.

6. **File Discovery**: Added `bitbucket_get_pr_diffstat` tool using the commits diffstat API with PR source and destination commits to get the list of changed files with statistics.

--
author:	gemyago
association:	owner
edited:	true
status:	none
--
# Task list based on the requirements

## Requirements

Please read referenced files to understand the problem:
- `https://github.com/gemyago/atlacp/issues/21#issuecomment-3193509026`

## Relevant Files

- internal/api/mcp/controllers/bitbucket.go
- internal/api/mcp/controllers/bitbucket_test.go
- internal/api/mcp/controllers/ports.go

### Source Files
- `internal/api/mcp/controllers/bitbucket.go` - Add 5 new MCP tools to existing BitbucketController
- `internal/api/mcp/controllers/ports.go` - Add new service methods to bitbucketService interface
- `internal/services/bitbucket/get_diffstat.go` - Service for getting PR diff statistics
- `internal/services/bitbucket/get_diff.go` - Service for getting PR diff content
- `internal/services/bitbucket/get_file_content.go` - Service for getting file content from commits
- `internal/services/bitbucket/add_pr_comment.go` - Service for adding PR comments (general and inline)
- `internal/services/bitbucket/request_pr_changes.go` - Service for requesting PR changes
- `internal/services/bitbucket/models.go` - Add new model structures for diff, comment, file content
- `internal/app/bitbucket.go` - Add new service methods for PR review functionality
- `internal/app/ports.go` - Add new methods to bitbucketClient interface
- `internal/services/bitbucket/ENDPOINTS.md` - Document new API endpoints

### Test Files
- `internal/api/mcp/controllers/bitbucket_test.go` - Unit tests for new MCP tools
- `internal/services/bitbucket/get_diffstat_test.go` - Unit tests for diffstat service
- `internal/services/bitbucket/get_diff_test.go` - Unit tests for diff service
- `internal/services/bitbucket/get_file_content_test.go` - Unit tests for file content service
- `internal/services/bitbucket/add_pr_comment_test.go` - Unit tests for comment service
- `internal/services/bitbucket/request_pr_changes_test.go` - Unit tests for request changes service
- `internal/app/bitbucket_test.go` - Unit tests for new application layer methods

## Tasks

- [x] 1.0 Implement New Bitbucket Service Layer Methods and Models
    - [x] 1.1 Add new model structures to `internal/services/bitbucket/models.go`
    - [x] 1.2 Implement `get_diffstat.go` service following HTTP client generation pattern
    - [x] 1.3 Implement `get_diff.go` service following HTTP client generation pattern
    - [x] 1.4 Implement `get_file_content.go` service following HTTP client generation pattern
    - [x] 1.5 Implement `add_pr_comment.go` service following HTTP client generation pattern
    - [x] 1.6 Implement `request_pr_changes.go` service following HTTP client generation pattern

- [ ] 2.0 Update Application Layer Integration
    - [x] 2.1 Extend `bitbucketClient` interface in `internal/app/ports.go` with new methods copied from bitbucket client (implemented in scope of 1.0 task):
        - GetPRDiffstat
        - GetPRDiff
        - GetFileContent
        - AddPRComment
        - RequestPRChanges
    - [x] 2.2 Regenerate mocks (using @mockery.mdc)
    - [x] 2.3 Implement new service method in `internal/app/bitbucket.go` for GetPRDiffstat
    - [x] 2.4 Implement new service method in `internal/app/bitbucket.go` for GetPRDiff ✓
    - [x] 2.5 Implement new service method in `internal/app/bitbucket.go` for GetFileContent
    - [ ] 2.6 Implement new service method in `internal/app/bitbucket.go` for AddPRComment
    - [ ] 2.7 Implement new service method in `internal/app/bitbucket.go` for RequestPRChanges

- [ ] 3.0 Extend BitbucketController with New MCP Tools
    - [ ] 3.1 Extend `bitbucketService` in `internal/api/mcp/controllers/ports.go` to include new methods (implemented in scope of 2.0 task). Regenerate mocks (using @mockery.mdc)
    - [ ] 3.1 Add `bitbucket_get_pr_diffstat` tool to existing controller
    - [ ] 3.2 Add `bitbucket_get_pr_diff` tool to existing controller
    - [ ] 3.3 Add `bitbucket_get_file_content` tool to existing controller
    - [ ] 3.4 Add `bitbucket_add_pr_comment` tool to existing controller
    - [ ] 3.5 Add `bitbucket_request_pr_changes` tool to existing controller

- [ ] 4.0 Update Documentation and API References
    - [ ] 4.1 Update README or other relevant documentation if needed

--
