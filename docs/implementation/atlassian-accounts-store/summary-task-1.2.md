# Summary: Task 1.2 — Implement accounts store — load, validate, query

## What was implemented

- Added `AccountsStore` in `internal/services/accounts_store.go` with `NewAccountsStore`, `LoadFromFile` (read JSON envelope `{"accounts":[...]}`, `app.ValidateAtlassianAccounts`, replace state only on success), and `GetDefaultAccount` / `GetAccountByName` matching repository semantics (`ErrNoDefaultAccount`, `ErrAccountNotFound` with name). `sync.RWMutex` protects the in-memory slice for concurrent read access after load.
- Tests in `internal/services/accounts_store_test.go`: happy path; missing file; invalid JSON; validation failures (empty list, duplicate names, two defaults, no default); failed load preserves prior loaded state; empty-store getters.

## Uncertainties / deviations

- Adjusted `internal/app/bitbucket_auth.go`: multi-line `getTokenProvider` signature for `golines`, and `//nolint:ireturn` on the line immediately above `func` so `ireturn` and `nolintlint` agree (replacing the previous block-style nolint that failed `nolintlint`). Behavior unchanged.
