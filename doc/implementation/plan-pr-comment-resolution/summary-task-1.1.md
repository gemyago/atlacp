# Task 1.1 summary: Service models and resolution parsing

## What was done

- **`PRComment`**: `Resolution json.RawMessage` with `json:"resolution,omitempty"` was already present in `internal/services/bitbucket/models.go` (captures list `resolution: {}` vs full GET payloads).

- **`ResolvedStateFromResolutionJSON`**: Logic was already in `internal/services/bitbucket/pr_comment_resolution.go`. Adjusted the signature to use unnamed `(bool, bool)` returns to satisfy the `nonamedreturns` linter.

- **Unit tests**: Added `pr_comment_resolution_test.go` with `TestResolvedStateFromResolutionJSON` and nested cases for nil/empty, `null`, `{}` (ambiguous), `resolved` true/false, `resolved_by` / `resolved_on`, whitespace trimming, invalid JSON, and unrelated object keys.

- **Coverage**: `get_pr_comment.go` had no tests and failed the repo’s 90% per-file coverage gate. Added `get_pr_comment_test.go` (success + HTTP error + token error) so `make test` passes; this overlaps with plan Task 1.2 but is required for CI green on this branch.

## Verification

- `make lint`: pass  
- `make test`: pass (total coverage ~95.7% per last run)
