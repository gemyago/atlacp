# Task 1.1 — Constructor with deps and startup parity

## What was implemented

- **`AccountsStoreDeps`** — `dig.In`, `RootLogger *slog.Logger`, `ConfigPath` with `name:"config.atlassian.accountsFilePath"` (same shape as `AtlassianAccountsRepositoryDeps`).
- **`NewAccountsStoreWithDeps`** — Rejects empty path with `errors.New("accounts configuration path not specified")` (same string as the old repository). Uses `os.Stat` for missing-file detection and returns `fmt.Errorf("accounts configuration file not found at %s", path)` without wrapping, matching the former repository. Then builds `&AccountsStore{logger: ...}`, calls `LoadFromFile(configPath)`, and returns the store or the load error.
- **`AccountsStore.logger`** — Set to `deps.RootLogger.WithGroup("atlassian-accounts")` in the constructor for parity with the removed type (field reserved for future use; getters unchanged).

Tests in `accounts_store_test.go` under `TestAccountsStore/NewAccountsStoreWithDeps`: empty path, missing file, happy path with temp JSON file and `GetDefaultAccount` assertion.

## Uncertainties / deviations

- **None.** Behavior for the covered cases matches the old repository startup (empty path, missing file via `Stat`, load via existing `LoadFromFile`).
