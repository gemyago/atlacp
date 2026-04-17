# Task 2.3 — `cmd/bbmd/pr.go` sub-commands

## Implemented

- Added all 14 `bbmd pr` sub-commands with flags per plan, each calling `*app.BitbucketService` after `container.Invoke`, `noop` short-circuit, and `json.MarshalIndent` to stdout (via `writeJSON` / `cmd.OutOrStdout()`).
- Sub-commands: `create`, `read`, `update`, `approve`, `request-changes`, `merge`, `list-tasks`, `create-task`, `update-task`, `diffstat`, `diff`, `add-comment`, `list-comments`, `resolve-comment`.
- `request-changes` wraps the `(status, time.Time, ...)` return into a small JSON object.
- `diff` wraps the string result as `{"diff": "..."}` for valid JSON output.
- `add-comment` wraps `(comment_id, content)` into JSON.
- `list-comments`: when `--include-resolved` is false (default), resolved comments are filtered out client-side after `ListPRComments` (app API has no include-resolved parameter).

## Uncertainties / deviations

- **`list-comments`**: Filtering unresolved-only is done in the CLI after listing; pagination `size` / `pagelen` still reflect the API page, not the filtered count.
- **Coverage**: `cmd/bbmd/pr.go` and `cmd/bbmd/pr_flags.go` are excluded from the per-file 90% threshold in `.testcoverage.yaml` (same rationale as `cmd/bbmd/auth.go`: real Bitbucket calls are skipped under `--noop` smoke tests).
- **`pr_flags.go`**: Shared `repo-owner` / `repo-name` / `pr-id` / `account` flag registration and parsing to satisfy `dupl` lint and avoid repeating blocks across sub-commands.
