# Task 1.2 summary: HTTP client — get / resolve comment

## What was done

- **`GetPRComment`**: Already implemented in `internal/services/bitbucket/get_pr_comment.go` with `get_pr_comment_test.go` (httptest: success with resolution JSON, HTTP error, token error) from Task 1.1 work.

- **`ResolvePRComment`**: Added `internal/services/bitbucket/resolve_pr_comment.go` — `POST` to `/repositories/{workspace}/{repo_slug}/pullrequests/{pr_id}/comments/{comment_id}/resolve` with path escaping, Bearer auth via `TokenProvider`, no request body, response unmarshaled into `CommentResolution`.

- **`CommentResolution`**: Added in `models.go` (`type`, `user`, `created_on`) for the resolve endpoint response body.

- **Tests**: `resolve_pr_comment_test.go` — success (method, path, `Authorization`, parsed `CommentResolution`), HTTP error path, token error.

## Verification

- `make lint`: pass  
- `make test`: pass (total coverage ~95.7% per last run)
