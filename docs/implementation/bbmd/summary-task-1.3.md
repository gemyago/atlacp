# Summary: Task 1.3 — `cmd/mcp/root.go` uses `DefaultAccountsFilePath`

## What was implemented

- In `PersistentPreRunE`, after `config.Load` and before DI (`config.Provide`), if `--atlassian-accounts-file` / `-a` is empty, the viper key `atlassian.accountsFilePath` is set from `services.DefaultAccountsFilePath()`.
- `cmd/mcp/main_test.go`: added sub-tests that run `http` and `stdio` with `--noop` and `--logs-file` but **without** `--atlassian-accounts-file`, asserting initialization succeeds when the default path is used (missing file tolerated by `AccountsStore`).

## Uncertainties / deviations

- `internal/services/time.go`: moved the `//nolint:ireturn` directive onto the `NewTimeProvider` function line so `nolintlint` accepts it (standalone directive above the function was reported unused while `ireturn` still failed after removal). No behavioral change.
