# Task 1.3 — Remove repository and consolidate tests

## What was implemented

- **Removed** `internal/services/atlassian_accounts.go` and `internal/services/atlassian_accounts_test.go`.
- **Moved** `atlassianAccountsConfig` into `accounts_store.go` so JSON load/save keeps a single definition.
- **Removed** parameterless `NewAccountsStore()`; tests use `&AccountsStore{}` for an empty store.
- **Extended** `accounts_store_test.go`: former repository scenarios are covered on `NewAccountsStoreWithDeps` (read failure on directory path, invalid JSON, shared validation failure), plus `GetDefaultAccount` / `GetAccountByName` cases that were unique to the old tests (no default among in-memory accounts via same-package field init; not-found name when other accounts are loaded).

## Uncertainties / deviations

- **None.** The “no default among accounts” case uses `&AccountsStore{accounts: ...}` in tests (same package) instead of the deleted `atlassianAccountsRepository` struct literal; behavior under test is unchanged.
