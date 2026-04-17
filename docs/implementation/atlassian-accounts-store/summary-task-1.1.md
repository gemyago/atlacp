# Summary: Task 1.1 — Extract shared validation for Atlassian accounts config

## What was implemented

- Added `app.ValidateAtlassianAccounts([]AtlassianAccount) error` in `internal/app/atlassian_accounts_validate.go` with the same rules as the former `validateAccountsConfig` path (non-empty list, unique names, exactly one default, Bitbucket/Jira token fields when those services are present).
- `NewAtlassianAccountsRepository` now calls `app.ValidateAtlassianAccounts(config.Accounts)` after JSON parse; error wrapping and messages for the repository path are unchanged.
- Validation unit tests live in `internal/app/atlassian_accounts_validate_test.go`. `internal/services/atlassian_accounts_test.go` keeps repository-focused tests and adds one case that rejects invalid config on disk via `NewAtlassianAccountsRepository` (no default among two accounts).

## Uncertainties / deviations

- Removed an obsolete `//nolint:ireturn` on `internal/di/dig.go` `ProvideAs` (triggered `nolintlint`) so `make lint` passes; behavior unchanged.
