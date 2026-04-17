# Task 1.2 — DI registration and interface alias

## What was implemented

- **`internal/services/register.go`** — Replaced `NewAtlassianAccountsRepository` with `NewAccountsStoreWithDeps` and `di.ProvideAs[*AccountsStore, app.AtlassianAccountsRepository]` so the dig graph exposes the port from the concrete store (same pattern as Bitbucket in `internal/app/register.go`).
- **`internal/services/accounts_store.go`** — Added `var _ app.AtlassianAccountsRepository = (*AccountsStore)(nil)` as a compile-time port assertion.

The legacy `atlassian_accounts.go` implementation remains in the tree until Task 1.3; it is no longer registered in DI.

## Uncertainties / deviations

- **None.**
