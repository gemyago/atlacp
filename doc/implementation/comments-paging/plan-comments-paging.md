# Plan: PR Comments Pagination

## 1. Introduction/Overview

Currently, the `bitbucket_list_pr_comments` MCP tool fetches PR comments without any pagination support. The Bitbucket API returns paginated results by default (page 1, pagelen 10), but our implementation ignores pagination metadata (`page`, `pagelen`, `size`, `next`, `previous`) and doesn't allow callers to specify page or page size.

**Goal:** Extend the PR comments listing across all layers (Bitbucket client, app service, MCP controller) to support pagination parameters (`page`, `pagelen`) and return pagination metadata. Set the default page size to the maximum supported by the Bitbucket API (100). Also extend e2e tests to cover pagination scenarios.

## 2. Business Logic

- When listing PR comments, callers can optionally specify a `page` number and `pagelen` (page size).
- If no `pagelen` is specified, the system defaults to 100 (the maximum supported by the Bitbucket API for this endpoint).
- The response includes pagination metadata: `size` (total count), `page` (current page), `pagelen` (items per page), `next` (URL of next page if exists), `previous` (URL of previous page if exists).
- The MCP tool summary text should reflect pagination info (e.g., "Found 25 comments on pull request #42 (page 1, showing 25 of 150 total)").

## 3. High Level Architecture

The change touches three layers, following the same pattern as `ListPullRequestTasks`:

1. **Bitbucket Client** (`internal/services/bitbucket/`) - Add pagination query params to the API request and expand the response model.
2. **App Service** (`internal/app/`) - Pass pagination params through from controller to client.
3. **MCP Controller** (`internal/api/mcp/controllers/`) - Expose `page` and `pagelen` as optional tool parameters and include pagination info in the response.

## 4. Detailed Architecture

### 4.1 Bitbucket Client Layer

**File:** `internal/services/bitbucket/models.go`
- Add `Page` and `PageLen` fields to `ListPRCommentsParams`
- Add pagination metadata fields to `ListPRCommentsResponse`: `Size`, `Page`, `PageLen`, `Next`, `Previous` (matching `PaginatedTasks` pattern)

**File:** `internal/services/bitbucket/list_pr_comments.go`
- Build query parameters (`page`, `pagelen`) from params, following the pattern in `list_tasks.go`
- Append query string to the request URL

**File:** `internal/services/bitbucket/list_pr_comments_test.go`
- Update existing tests to verify query parameters are sent correctly
- Add test for pagination params being sent when specified
- Verify pagination metadata is parsed from response

### 4.2 App Service Layer

**File:** `internal/app/bitbucket.go`
- Add `Page` and `PageLen` fields to `BitbucketListPRCommentsParams`
- Pass these through to the Bitbucket client params in `ListPRComments`
- Set default `PageLen` to 100 if not specified (0)

**File:** `internal/app/bitbucket_test.go`
- Update tests to verify pagination params are passed through
- Test default pagelen behavior

### 4.3 MCP Controller Layer

**File:** `internal/api/mcp/controllers/bitbucket.go`
- Add `page` (number, optional) and `pagelen` (number, optional) parameters to the `bitbucket_list_pr_comments` tool definition
- Extract these from the request and pass to the service
- Update the summary text to include pagination info
- Include pagination metadata in the JSON response

**File:** `internal/api/mcp/controllers/bitbucket_test.go`
- Update tests to cover pagination parameters
- Verify pagination info appears in response

### 4.4 E2E Test Documentation

**File:** `doc/testing/bitbucket-mcp-integration-tests.md`
- Add a new test (Test 6) for PR comments pagination:
  - Create a PR
  - Post 20+ comments (mix of general and inline)
  - List comments with default pagelen, verify all returned
  - List with small pagelen (e.g., 5), verify only 5 returned and pagination metadata shows `next`
  - Fetch page 2, verify correct comments returned
  - Continue until all pages consumed and total matches expected count

## 5. Key Architectural Decisions

1. **Default pagelen = 100**: The Bitbucket API supports up to 100 items per page for most endpoints. Setting this as default ensures fewer round-trips while staying within API limits.
2. **Follow `ListPullRequestTasks` pattern**: Reuse the same pagination approach (query params, response structure) for consistency across the codebase.
3. **Default applied at app layer**: The app service sets the default pagelen (100) when the caller doesn't specify one. This keeps the client layer clean and the default centralized.
4. **Backward compatible MCP tool**: `page` and `pagelen` are optional parameters. Existing callers that don't specify them get the new default (100 items) instead of the old Bitbucket default (10).

## 6. Uncertainties

- The exact max pagelen for the Bitbucket PR comments endpoint may be 50 or 100. Bitbucket documentation typically states 100 for most paginated endpoints, but some are 50. This should be verified during implementation by testing against the live API. If it's 50, adjust the default accordingly.

## 7. Related Files

**Files to modify:**
- [internal/services/bitbucket/models.go](internal/services/bitbucket/models.go) - `ListPRCommentsParams`, `ListPRCommentsResponse`
- [internal/services/bitbucket/list_pr_comments.go](internal/services/bitbucket/list_pr_comments.go) - Add query params
- [internal/services/bitbucket/list_pr_comments_test.go](internal/services/bitbucket/list_pr_comments_test.go) - Update tests
- [internal/app/bitbucket.go](internal/app/bitbucket.go) - `BitbucketListPRCommentsParams`, `ListPRComments`
- [internal/app/bitbucket_test.go](internal/app/bitbucket_test.go) - Update tests
- [internal/app/ports.go](internal/app/ports.go) - No change needed (return type stays `*bitbucket.ListPRCommentsResponse`)
- [internal/api/mcp/controllers/bitbucket.go](internal/api/mcp/controllers/bitbucket.go) - `newListPRCommentsServerTool`
- [internal/api/mcp/controllers/bitbucket_test.go](internal/api/mcp/controllers/bitbucket_test.go) - Update tests
- [doc/testing/bitbucket-mcp-integration-tests.md](doc/testing/bitbucket-mcp-integration-tests.md) - Add pagination e2e test

**Reference files (pattern to follow):**
- [internal/services/bitbucket/list_tasks.go](internal/services/bitbucket/list_tasks.go) - Pagination pattern at client layer

## 8. Task List

TDD approach must be followed. Module-specific task completion protocol **must** be followed for each task.

**Task 1.1: Extend Bitbucket client models and implementation for pagination**
- Add `Page` and `PageLen` fields to `ListPRCommentsParams` in `internal/services/bitbucket/models.go`
- Add `Size`, `Page`, `PageLen`, `Next`, `Previous` fields to `ListPRCommentsResponse` in `internal/services/bitbucket/models.go`
- Write failing tests in `internal/services/bitbucket/list_pr_comments_test.go`:
  - Test that `page` and `pagelen` query params are sent when specified
  - Test that pagination metadata (`size`, `page`, `pagelen`, `next`, `previous`) is parsed from API response
  - Test that no query params are sent when pagination fields are zero-valued
- Run affected tests: `go test -v ./internal/services/bitbucket/... -run "^TestClient_ListPRComments$"`
  - Verify failure is expectation-based (not compilation errors)
- Update `ListPRComments` in `internal/services/bitbucket/list_pr_comments.go` to:
  - Build `url.Values` with `page` and `pagelen` when specified (following `list_tasks.go` pattern)
  - Append query string to request URL
- Run affected tests: `go test -v ./internal/services/bitbucket/... -run "^TestClient_ListPRComments$"`
  - Verify all tests pass
- Write summary to `doc/implementation/comments-paging/summary-task-1.1.md`
- All checks from completion protocol must be passed

**Task 1.2: Extend app service layer for pagination**
- Add `Page` and `PageLen` fields to `BitbucketListPRCommentsParams` in `internal/app/bitbucket.go`
- Write failing tests in `internal/app/bitbucket_test.go`:
  - Test that pagination params are passed through to the client
  - Test that default `PageLen` of 100 is set when caller specifies 0
- Run affected tests: `go test -v ./internal/app/... -run "ListPRComments"`
  - Verify failure is expectation-based (not compilation errors)
- Update `ListPRComments` in `internal/app/bitbucket.go`:
  - Map `Page` and `PageLen` to client params
  - Set default `PageLen = 100` when `params.PageLen == 0`
- Run affected tests: `go test -v ./internal/app/... -run "ListPRComments"`
  - Verify all tests pass
- Write summary to `doc/implementation/comments-paging/summary-task-1.2.md`
- All checks from completion protocol must be passed

**Task 1.3: Extend MCP controller for pagination**
- Add `page` (number, optional) and `pagelen` (number, optional) parameters to the `bitbucket_list_pr_comments` tool definition in `internal/api/mcp/controllers/bitbucket.go`
- Write failing tests in `internal/api/mcp/controllers/bitbucket_test.go`:
  - Test that `page` and `pagelen` are extracted from request and passed to service
  - Test that pagination info is included in the response summary and JSON
  - Test defaults (no page/pagelen specified)
- Run affected tests: `go test -v ./internal/api/mcp/controllers/... -run "ListPRComments"`
  - Verify failure is expectation-based (not compilation errors)
- Update the handler in `newListPRCommentsServerTool`:
  - Extract `page` and `pagelen` using `request.GetInt`
  - Pass to `BitbucketListPRCommentsParams`
  - Update summary text to include pagination info (e.g., page X, showing Y of Z total)
- Run affected tests: `go test -v ./internal/api/mcp/controllers/... -run "ListPRComments"`
  - Verify all tests pass
- Write summary to `doc/implementation/comments-paging/summary-task-1.3.md`
- All checks from completion protocol must be passed

**Task 1.4: Extend e2e test documentation with pagination test**
- Add a new Test 6 to `doc/testing/bitbucket-mcp-integration-tests.md`:
  - Title: "PR Comments Pagination"
  - Steps:
    1. Setup: Create branch, add test file, commit and push
    2. Create a PR
    3. Post 20+ comments using `mcp.bitbucket_add_pr_comment` (mix of general and inline comments)
    4. List comments with default pagelen - verify all 20+ returned in a single page (since default is 100)
    5. List comments with `pagelen: 5` - verify only 5 returned, pagination metadata shows `next` is present and `size` >= 20
    6. List comments with `pagelen: 5, page: 2` - verify second page of 5 comments
    7. Iterate through all pages with `pagelen: 5`, collecting all comment IDs - verify total matches expected count and all original comment IDs are present
    8. Clean up: Merge and delete branch
  - Follow same reporting format as other tests
- Write summary to `doc/implementation/comments-paging/summary-task-1.4.md`
- Non-coding task completion protocol applies (no lint/test needed)

**Task 1.5: Compress implementation summaries**
- Follow [compress-implementation-summaries.md](/.context/compress-implementation-summaries.md) to compress the implementation summaries.
