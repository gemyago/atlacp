# Draft Pull Request Support Implementation Tasks

## Overview
This document outlines the tasks required to fully implement the draft pull request functionality in the Atlassian MCP integration. While the `Draft` field is already defined in the `BitbucketCreatePRParams` struct and used in the `CreatePR` method, we need to extend this functionality to other PR operations and ensure proper test coverage.

## Tasks

### 1. Extend PR Update Functionality
- [ ] Update `BitbucketUpdatePRParams` struct in `internal/app/bitbucket.go` to include the `Draft` field
- [ ] Update `UpdatePR` method to include draft status in the update request
- [ ] Ensure validation logic handles the draft parameter appropriately

### 2. Update Client Layer
- [ ] Verify that `internal/services/bitbucket/update_pr.go` passes draft status to the API
- [ ] Make sure `internal/services/bitbucket/pullrequest.go` includes draft field in appropriate structs

### 3. API Controller Layer
- [ ] Update MCP controllers in `internal/api/mcp/controllers/bitbucket.go` to handle draft parameter
- [ ] Ensure JSON payload parsing includes draft status for both create and update operations

### 4. Run mcp integration tests
- [ ] Ask user to start the MCP server
- [ ] Run `make bitbucket-mcp-integration-integration-tests` to ensure the draft PR functionality works as expected

## Implementation Approach
1. Follow TDD: Write failing tests first to validate draft PR behavior
2. Implement the necessary changes to make tests pass
3. Refactor if needed while maintaining test coverage

## Files to Modify
- `internal/app/bitbucket.go`
- `internal/app/bitbucket_test.go`
- `internal/services/bitbucket/update_pr.go` (maybe)
- `internal/services/bitbucket/update_pr_test.go` (maybe)
- `internal/api/mcp/controllers/bitbucket.go`
- `internal/api/mcp/controllers/bitbucket_test.go`

## Completion Criteria
- All tests pass
- `make lint test` passes without errors
- `bitbucket-mcp-integration-integration-tests.md` passes without errors