# Task 1.1 — `DefaultAccountsFilePath` resolver

## What was implemented

- Added `internal/services/config_path.go` with `DefaultAccountsFilePath() (string, error)`:
  - On `darwin`, base dir is `filepath.Join(os.Getenv("HOME"), ".config")`.
  - Otherwise base dir comes from `os.UserConfigDir()`; errors propagate.
  - Final path is `filepath.Join(baseDir, "atlacp", "accounts.json")`.
- Unexported helpers `defaultAccountsBaseDir` and `defaultAccountsFilePath` take explicit `goos`, `home`, and a `userConfigDir` callback so behavior and error paths are testable on every platform.
- Added `internal/services/config_path_test.go`: suffix check, macOS integration with `HOME` override, non-macOS parity with real `UserConfigDir`, and table-style tests for the helpers (including “UserConfigDir not called on darwin”).

## Uncertainties / deviations from the plan

- **Deviation:** The plan only named `DefaultAccountsFilePath`. Extra unexported helpers were added so the repo’s **90% per-file** coverage gate passes on macOS (the real non-darwin and `UserConfigDir` error paths are not exercised when `runtime.GOOS` is `darwin`). Behavior matches the task spec; helpers are thin wrappers around the same rules.

## Verification

- `make lint`: pass  
- `make test`: pass (total coverage ~95.1% per project check)
