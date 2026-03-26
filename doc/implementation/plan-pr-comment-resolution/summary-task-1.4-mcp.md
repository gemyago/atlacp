# Task 1.3 — MCP + integration doc (summary)

## Changes

- **`internal/api/mcp/controllers/ports.go`**: Extended `bitbucketService` with `ResolvePRComment(ctx, app.BitbucketResolvePRCommentParams) (*bitbucket.CommentResolution, error)`.
- **`internal/api/mcp/controllers/bitbucket.go`**:
  - `bitbucket_list_pr_comments`: description now states JSON includes a per-comment `resolved` boolean (app-layer enrichment).
  - Refactored list handler into `makeListPRCommentsHandler` to satisfy `funlen` after the longer description.
  - Added `bitbucket_resolve_pr_comment` MCP tool (`pr_id`, `comment_id`, `repo_owner`, `repo_name`, optional `account`); handler returns summary text plus marshaled `CommentResolution` JSON.
  - Registered the new tool in `NewTools()`.
- **`internal/api/mcp/controllers/mock_bitbucket_service.go`**: Regenerated via mockery v2 (`go run github.com/vektra/mockery/v2@v2.53.4`).
- **`internal/api/mcp/controllers/bitbucket_test.go`**: Tool definition test, tool count 15, `resolved` assertion on list JSON, handler tests for resolve (success, service error, missing `comment_id`).
- **`doc/testing/bitbucket-mcp-integration-tests.md`**: Test 5 — new step 7 to resolve a comment and verify `resolved` in `bitbucket_list_pr_comments` JSON.

## Verification

- `make lint`: pass
- `make test`: pass (total coverage ~95.1% per project gate)
