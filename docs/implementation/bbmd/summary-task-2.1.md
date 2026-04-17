# Task 2.1 — `cmd/bbmd` scaffold

## Implemented

- Added `cmd/bbmd` Cobra binary: `main.go` (`setupCommands`, `main` with `// coverage-ignore`), `root.go` with persistent flags (`--env/-e`, `--log-level/-l`, `--json-logs`, `--logs-file` default `bbmd.log`, `--atlassian-accounts-file/-a`, `--noop`), package-level `noop` and `resolvedAccountsFilePath`, `PersistentPreRunE` resolving accounts path (default via `services.DefaultAccountsFilePath()` when `-a` is empty), `config.Load`, logger setup, and DI (`config.Provide` → `app.Register` → `services.Register` + root logger).
- `pr.go`: `pr` group with a minimal `read` sub-command that invokes `*app.BitbucketService` under DI and returns early when `noop` (otherwise stub error until task 2.3).
- `file.go` / `auth.go`: empty group commands only.
- `main_test.go`: smoke and DI tests (`pr read --noop`, explicit accounts file, bad log level / env, non-noop stub error).

## Uncertainties / deviations

- **Order of operations vs plan bullet list:** `config.Load` runs before applying the default accounts file path, matching `cmd/mcp/root.go`, so embedded `local.json` is merged first and the resolved default path is applied afterward without being overwritten by merge order.
- **`pr` is not strictly “empty”:** the plan asked for empty group constructors, but Task 2.1 also requires `bbmd pr read --noop` with flags; `read` is included as a minimal stub so the smoke test can run.
