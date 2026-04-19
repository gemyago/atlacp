# Implementation Summary: Wire AccountsStore into DI and remove file-backed repository

**Plan:** [plan-accounts-store-di.md](./plan-accounts-store-di.md)

## Overview

`AccountsStore` is now the single implementation behind `app.AtlassianAccountsRepository`: `AccountsStoreDeps` and `NewAccountsStoreWithDeps` mirror the old repository startup (path checks, `os.Stat`, `LoadFromFile`, grouped logger), DI registers the store with `ProvideAs`, and the legacy `atlassian_accounts` implementation and tests were removed or folded into `accounts_store_test.go`. Documentation and file comments were updated so production wiring and test-only construction are accurately described.

## Tasks

### Task 1.1: Constructor with deps and startup parity

Introduced `AccountsStoreDeps` and `NewAccountsStoreWithDeps` so construction matches the old repository: same error strings, `os.Stat` for a missing file, `LoadFromFile`, and a logger with `WithGroup("atlassian-accounts")`. Tests cover empty path, missing file, and a happy path with a temp JSON file and `GetDefaultAccount`.

### Task 1.2: DI registration and interface alias

DI now wires `AccountsStore` via `NewAccountsStoreWithDeps` and `di.ProvideAs[*AccountsStore, app.AtlassianAccountsRepository]` (same pattern as Bitbucket), with a compile-time port assertion on `AccountsStore`. The legacy `atlassian_accounts.go` file was still in the tree for Task 1.3 but was no longer registered in DI.

### Task 1.3: Remove repository and consolidate tests

The Atlassian accounts repository files were removed; `atlassianAccountsConfig` lives in `accounts_store.go`, parameterless `NewAccountsStore()` was dropped in favor of `&AccountsStore{}` in tests, and former repository scenarios plus unique cases were folded into `accounts_store_test.go` via `NewAccountsStoreWithDeps` and same-package struct initialization.

### Task 1.4: Comments and implementation summary

Refreshed the file-level comment on `internal/services/accounts_store.go` for production DI and test-only construction, added `docs/implementation/accounts-store-di/implementation-summary.md`, and updated `docs/implementation/atlassian-accounts-store/implementation-summary.md` for DI wiring links and to correct the Task 1.4 bullet that previously claimed the store was unwired.

## Completion

- Lint: ✓
- Type check: ✓
- Tests: ✓
