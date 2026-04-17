# Plan: `bbmd` — Bitbucket CLI Tool

## 1. Introduction / Overview

**Goal:** Implement a new CLI binary called `bbmd` that exposes all currently supported Bitbucket operations and auth (account) management as interactive command-line tools. The binary will follow the same structural pattern as the existing `cmd/mcp` binary (cobra commands, DI via dig, layered architecture), but instead of running an MCP server it will execute individual operations directly, print results as JSON to stdout, and send logs to a file by default (to avoid mixing log noise with machine-readable output).

**Problem solved:** Developers and scripts need a lightweight CLI wrapper around the Bitbucket service layer, without having to run or connect to an MCP server. `bbmd` provides direct access to every Bitbucket operation as a sub-command.

---

## 2. Business Logic

### PR sub-commands
Each Bitbucket operation becomes a named sub-command under `bbmd pr <operation>`. The operation accepts flags corresponding to its parameters, calls the `app.BitbucketService` method, and writes the result as pretty-printed JSON to `os.Stdout`. Errors are written to `stderr`.

Supported operations (mirrors existing `app.BitbucketService` methods):
- `create` — create a pull request
- `read` — get pull request details
- `update` — update title / description / draft state
- `approve` — approve a PR
- `request-changes` — remove approval / request changes
- `merge` — merge a PR
- `list-tasks` — list tasks on a PR
- `create-task` — create a task
- `update-task` — update a task
- `diffstat` — list changed files summary
- `diff` — raw diff text
- `add-comment` — post a comment (general or inline)
- `list-comments` — list all PR comments
- `resolve-comment` — resolve a comment thread

There is also one non-PR-scoped Bitbucket command:
- `bbmd file content` — get file content at a commit (under a separate `file` sub-command group)

### Auth / account management sub-commands
All operations on the `AccountsStore` are exposed under `bbmd auth <operation>`. Mutating commands persist changes by calling `AccountsStore.SaveToFile` before exiting.

Supported auth operations:
- `status` — print all configured accounts (name, default flag; partially redact token values — show first 4 and last 4 characters, mask middle with `***`)
- `add` — add or replace an account (`--name` required, `--default` bool, `--token-type` (default `Bearer`), `--token-value`)
- `remove` — remove an account by name (`--name` required)
- `set-default` — set a named account as the default (`--name` required)

### Accounts file resolution
Both `bbmd` and `cmd/mcp` resolve the accounts file path as follows (in priority order):
1. Explicit `--atlassian-accounts-file` / `-a` flag.
2. Platform-specific config directory: `~/.config/atlacp/accounts.json` on all platforms (XDG-style). On macOS `os.UserConfigDir()` returns `~/Library/Application Support`, which is inconvenient to work with from the shell; the resolver must override this to `$HOME/.config` on macOS (and fall back to `os.UserConfigDir()` on other platforms). This matches the convention used by most developer tools (e.g. gh, kubectl).

A new `internal/services/config_path.go` component (`AccountsFilePathResolver`) encapsulates this resolution. Both `cmd/bbmd/root.go` and `cmd/mcp/root.go` will call it and bind the resolved path to the `atlassian.accountsFilePath` viper key before DI setup.

The `AccountsStore` currently fails at startup when the accounts file does not exist. To allow working without a pre-existing file, the store must be changed: if the resolved path does not yet exist, it starts with an empty in-memory store (no accounts). The first mutating `auth add` / `SaveToFile` call will create the file.

### Logging
By default logs go to `bbmd.log` in the current working directory. The user may override this with `--logs-file`. This prevents log output from polluting the JSON printed to stdout. `bbmd.log` must be added to `.gitignore`.

---

## 3. High-Level Architecture

```
cmd/bbmd/
  main.go      – entry point; setupCommands(); same DI pattern as cmd/mcp
  root.go      – root cobra command; persistent flags; DI wiring in PersistentPreRunE
  pr.go        – "pr" sub-command group + all PR operation sub-commands
  file.go      – "file" sub-command group (get-content sub-command)
  auth.go      – "auth" sub-command group + account management sub-commands
  main_test.go – DI wiring smoke tests with --noop

internal/services/
  config_path.go       – AccountsFilePathResolver (platform-specific default path logic)
  config_path_test.go  – tests

internal/app/         – no changes needed
internal/config/      – minor: bind resolved accounts path before Provide()
```

The existing layers (`internal/app`, `internal/services`, `internal/config`, `internal/diag`, `internal/di`) are largely reused. Cobra commands talk directly to `app.BitbucketService` and `services.AccountsStore` — no intermediate CLI controller layer is needed; cobra itself serves as that boundary.

---

## 4. Detailed Architecture

### `internal/services/config_path.go` — `AccountsFilePathResolver`
- Exported function: `DefaultAccountsFilePath() (string, error)` — returns `filepath.Join(os.UserConfigDir(), "atlacp", "accounts.json")`.
- Used in both `cmd/bbmd/root.go` and `cmd/mcp/root.go` to set the default value of the `--atlassian-accounts-file` flag (via `cobra` flag default or pre-resolution before `cfg.BindPFlag`).

### `AccountsStore` — tolerate missing file
- Change `NewAccountsStoreWithDeps`: if `configPath` resolves and the file does not exist (`os.IsNotExist`), start with an empty store instead of returning an error. Log a warning.
- Keep the existing "path not specified" error unchanged (path must always be known so `SaveToFile` knows where to write).

### `cmd/bbmd/root.go`
- Mirrors `cmd/mcp/root.go`.
- Binary name: `bbmd`.
- Defines the **entire command hierarchy**: registers `pr`, `file`, and `auth` group commands (constructed in `pr.go`, `file.go`, `auth.go`) and adds them to the root command.
- Persistent flags bound on the root command: `--env`/`-e`, `--log-level`/`-l`, `--json-logs`, `--logs-file` (default `"bbmd.log"`), `--atlassian-accounts-file`/`-a`, `--noop` (bool, default false).
- `PersistentPreRunE`:
  1. Resolve accounts file path (flag value or `DefaultAccountsFilePath()`)
  2. Set viper key `atlassian.accountsFilePath`
  3. Load config, setup logger
  4. Register DI: `config.Provide → app.Register → services.Register`
- The `noop` bool variable is declared at root scope and passed as a parameter into each sub-command executor function.

### `cmd/mcp/root.go` — update
- Same accounts file path resolution as `bbmd`: if `--atlassian-accounts-file` is not provided, fall back to `DefaultAccountsFilePath()`.

### `cmd/bbmd/pr.go`
- Returns a `"pr"` cobra group command (called from `root.go` which adds it to the root).
- All 14 PR-scoped sub-commands are defined here as `*cobra.Command` values.
- Each sub-command:
  - Declares flags using cobra typed binding (e.g. `cmd.Flags().StringVar`, `cmd.Flags().IntVar`).
  - Marks required flags with `cmd.MarkFlagRequired`.
  - `RunE` calls a local executor function (e.g. `runPRCreate(ctx, container, noop, params...)`).
  - The executor calls `container.Invoke` to get `*app.BitbucketService`, then checks `if noop { return nil }` before calling the service, marshalling the result to `json.MarshalIndent`, and writing to `os.Stdout`.
  - On error, returns the error (cobra prints to stderr, exits non-zero).

### `cmd/bbmd/file.go`
- `"file"` cobra group with a `"content"` sub-command for `GetFileContent`.

### `cmd/bbmd/auth.go`
- Returns an `"auth"` cobra group command (added to root by `root.go`).
- Sub-commands:
  - `status`: executor calls `container.Invoke` → `*services.AccountsStore`, checks `noop`, then `store.ListAccounts()` → redact tokens → `json.MarshalIndent` → `os.Stdout`.
  - `add`: flags `--name` (required), `--default` (bool), `--token-type` (string, default `"Bearer"`), `--token-value` (string); executor builds `app.AtlassianAccount{Bitbucket: &app.AtlassianToken{...}}`, checks `noop`, calls `store.Upsert`, then `store.SaveToFile(resolvedPath)`.
  - `remove`: `--name` (required); executor checks `noop`, calls `store.Remove`, then `store.SaveToFile`.
  - `set-default`: `--name` (required); executor checks `noop`, calls `store.SetDefault`, then `store.SaveToFile`.
- The resolved accounts file path (set in `PersistentPreRunE`) is available to executors via closure or a variable shared in the same `root.go` scope.

### `cmd/bbmd/main_test.go`
- Tests call `setupCommands()` with `--noop` (no accounts file — missing-file tolerance means this is fine) and verify no error. The `--noop` flag causes each executor to exit after `container.Invoke` succeeds, so the full DI graph is exercised without real side effects. Coverage of `main()` itself should use surgical `// coverage-ignore` comment (matching the pattern used in `cmd/mcp`) rather than file-level suppression.

### `internal/services/accounts_store.go` — `ListAccounts`
- Add `ListAccounts() []app.AtlassianAccount` method: returns a copy of the current accounts slice (thread-safe).

### Output format
All output is `json.MarshalIndent(result, "", "  ")` written to `os.Stdout`. For `auth status` token values are partially redacted: show first 4 and last 4 characters with `***` in the middle (e.g. `"abcd***wxyz"`). If the value is ≤ 8 characters, replace entirely with `"***"`.

### `.gitignore`
Add `bbmd.log`.

---

## 5. Key Architectural Decisions

1. **Reuse existing app and services layers unchanged** (except `ListAccounts` and missing-file tolerance). `bbmd` is purely a command layer.
2. **No intermediate CLI controller layer** — cobra sub-commands call `app.BitbucketService` and `services.AccountsStore` directly, using cobra's typed flag binding for params. This avoids an unnecessary indirection layer.
3. **`pr` sub-command group** for all PR-scoped operations; `file` for file operations; `auth` for account management.
4. **`DefaultAccountsFilePath()` in `internal/services/config_path.go`** — centralises platform config dir logic; on macOS overrides `~/Library/Application Support` with `~/.config` for shell-friendliness; consumed by both `cmd/bbmd` and `cmd/mcp`.
5. **`AccountsStore` tolerates missing file** — allows bootstrapping without a pre-existing config file; the file is created on first `SaveToFile` call.
6. **Default log file `bbmd.log`** — ensures stdout is clean for JSON output piping/scripting.
7. **`--noop` pattern mirrors `cmd/mcp/stdio.go`** — `noop` is a local bool bound to a root persistent flag; each sub-command executor calls `container.Invoke` (fully exercising DI) then checks `if noop { return nil }` before the real operation. This ensures `main_test.go` tests DI wiring end-to-end without side effects.
8. **Token partial redaction** — first 4 + last 4 characters visible, middle replaced with `***`, so users can identify which token is configured without exposing the full value.

---

## 6. Uncertainties

- `os.UserConfigDir()` returns an error on systems where the config directory cannot be determined. `DefaultAccountsFilePath()` should propagate this error; callers fall back to requiring an explicit `--atlassian-accounts-file`.
- The `--noop` flag follows the same pattern as `cmd/mcp/stdio.go`: it is a local `bool` bound to a `--noop` flag on the root command (persistent). Each sub-command's executor function (e.g. `runPRCreate`, `runAuthStatus`) receives `noop` as a parameter alongside the bound flag values. The executor calls `container.Invoke` to wire all dependencies, then checks `if noop { return nil }` before performing the actual operation. This means DI is always fully exercised (the goal of the test) while the real side-effecting call is skipped. The root command (`root.go`) defines the entire command hierarchy and binds all persistent flags; per-operation logic lives in the respective sub-command files (`pr.go`, `auth.go`, `file.go`).

---

## 7. Related Files

### New files to create
| File | Purpose |
|------|---------|
| `cmd/bbmd/main.go` | Entry point |
| `cmd/bbmd/root.go` | Root cobra command + DI wiring |
| `cmd/bbmd/pr.go` | `pr` sub-command group (14 operations) |
| `cmd/bbmd/file.go` | `file` sub-command group (`content`) |
| `cmd/bbmd/auth.go` | `auth` sub-command group |
| `cmd/bbmd/main_test.go` | DI wiring smoke tests |
| `internal/services/config_path.go` | `DefaultAccountsFilePath()` resolver |
| `internal/services/config_path_test.go` | Tests for path resolver |

### Existing files to modify
| File | Change |
|------|--------|
| `.gitignore` | Add `bbmd.log` |
| `internal/services/accounts_store.go` | Add `ListAccounts()` method; tolerate missing file on init |
| `internal/services/accounts_store_test.go` | Tests for `ListAccounts`; tests for missing-file behaviour |
| `cmd/mcp/root.go` | Use `DefaultAccountsFilePath()` as default when flag not provided |
| `AGENTS.md` | Update run section with `bbmd` command examples |

---

## 8. Task List

> TDD approach must be followed: write failing tests first, verify failure is of expected kind (not compilation errors), implement, then verify all tests pass. Each task must leave the codebase in a green state as per the Task Completion Protocol.

---

**Task 1.1: Add `DefaultAccountsFilePath` resolver**
- Create `internal/services/config_path.go`:
  - `DefaultAccountsFilePath() (string, error)` — determines the base config dir:
    - On macOS (`runtime.GOOS == "darwin"`): use `filepath.Join(os.Getenv("HOME"), ".config")`
    - On all other platforms: use `os.UserConfigDir()`
  - Returns `filepath.Join(baseDir, "atlacp", "accounts.json")`
- Add stub so it compiles
- Write failing tests in `internal/services/config_path_test.go`:
  - Returns a non-empty path ending in `atlacp/accounts.json`
  - Path contains `.config/atlacp/accounts.json` (XDG style) on macOS
- Run tests, verify failure is "not equal" style (not compilation)
- Implement `DefaultAccountsFilePath`
- Run tests, verify all pass
- Write summary to `docs/implementation/bbmd/summary-task-1.1.md`
- Success criteria: `make lint && make test` pass

---

**Task 1.2: Update `AccountsStore` — `ListAccounts` and tolerate missing file**
- Add `ListAccounts() []app.AtlassianAccount` method to `AccountsStore` in `internal/services/accounts_store.go`
  - Returns a snapshot (copy) of the current accounts slice (thread-safe via `mu.RLock`)
- Modify `NewAccountsStoreWithDeps`: if the config file does not exist (`os.IsNotExist`), log a warning and return an empty store rather than an error
- Add stubs so it compiles
- Write failing tests in `internal/services/accounts_store_test.go`:
  - `ListAccounts` returns all accounts after `LoadFromFile`
  - `ListAccounts` returns updated list after `Upsert`
  - `ListAccounts` returns empty slice when no accounts loaded
  - `NewAccountsStoreWithDeps` with a non-existent file returns a non-nil empty store (no error)
  - `NewAccountsStoreWithDeps` with an empty path still returns an error
- Run tests, verify failures are expected (not compilation)
- Implement changes
- Run tests, verify all pass
- Write summary to `docs/implementation/bbmd/summary-task-1.2.md`
- Success criteria: `make lint && make test` pass

---

**Task 1.3: Update `cmd/mcp/root.go` to use `DefaultAccountsFilePath`**
- In `cmd/mcp/root.go`, after parsing the `--atlassian-accounts-file` flag and before `config.Provide`, if the flag value is empty resolve via `services.DefaultAccountsFilePath()` and set the viper key `atlassian.accountsFilePath`
- Update `cmd/mcp/main_test.go` if necessary (the test currently passes an explicit file; it should still pass)
- Run `make lint && make test`
- Write summary to `docs/implementation/bbmd/summary-task-1.3.md`
- Success criteria: `make lint && make test` pass

---

**Task 2.1: Create `cmd/bbmd` scaffold**
- Create `cmd/bbmd/main.go`:
  - `func main()` calls `setupCommands().Execute()`
  - `setupCommands()` creates `dig.Container`, builds root command via `newRootCmd(container)` which internally adds all sub-command groups
- Create `cmd/bbmd/root.go`:
  - Root command name `bbmd`
  - Declares `noop bool` at package scope (accessible to all executor functions)
  - Declares `resolvedAccountsFilePath string` at package scope (set in `PersistentPreRunE`, used by auth sub-commands)
  - Persistent flags: `--env`/`-e`, `--log-level`/`-l`, `--json-logs`, `--logs-file` (default `"bbmd.log"`), `--atlassian-accounts-file`/`-a`, `--noop` bound to `noop`
  - `newRootCmd` creates the root command, constructs `pr`, `file`, and `auth` group commands and adds them via `AddCommand`
  - `PersistentPreRunE`:
    1. Resolve accounts file path (flag value or `DefaultAccountsFilePath()`) → store in `resolvedAccountsFilePath`
    2. Set viper key `atlassian.accountsFilePath`
    3. Load config, setup logger
    4. Register DI: `config.Provide → app.Register → services.Register`
- Create `cmd/bbmd/pr.go`, `cmd/bbmd/file.go`, `cmd/bbmd/auth.go` as empty group command constructors (return `*cobra.Command`)
- Create `cmd/bbmd/main_test.go` with a test that runs `bbmd pr read --noop` (with required flags satisfied by dummy values) and verifies no error; use `// coverage-ignore` on `main()` body
- Run `make lint && make test`
- Write summary to `docs/implementation/bbmd/summary-task-2.1.md`
- Success criteria: `make lint && make test` pass

---

**Task 2.2: Implement `cmd/bbmd/auth.go` sub-commands**
- Add sub-commands to the `auth` group. Each sub-command's `RunE` calls a named executor (e.g. `runAuthStatus`). Each executor:
  1. Calls `container.Invoke` to wire dependencies
  2. Checks `if noop { return nil }` before performing the real operation
  - `status`: invokes `*services.AccountsStore`, calls `store.ListAccounts()`, redacts tokens (first 4 + `***` + last 4; full `"***"` if ≤ 8 chars), writes `json.MarshalIndent` result to `os.Stdout`
  - `add`: flags `--name` (required), `--default` (bool), `--token-type` (string, default `"Bearer"`), `--token-value` (string); invokes `*services.AccountsStore`, builds `app.AtlassianAccount`, calls `store.Upsert`, then `store.SaveToFile(resolvedAccountsFilePath)`
  - `remove`: `--name` (required); invokes store, calls `store.Remove`, then `store.SaveToFile`
  - `set-default`: `--name` (required); invokes store, calls `store.SetDefault`, then `store.SaveToFile`
- Expand `main_test.go` to run `auth status --noop` and verify no error
- Run `make lint && make test`
- Write summary to `docs/implementation/bbmd/summary-task-2.2.md`
- Success criteria: `make lint && make test` pass

---

**Task 2.3: Implement `cmd/bbmd/pr.go` sub-commands**
- Add all 14 PR sub-commands to the `pr` group. Each sub-command's `RunE` calls a named executor. Each executor calls `container.Invoke` to get `*app.BitbucketService`, then checks `if noop { return nil }` before the real call:
  - `create`: `--title` (req), `--source-branch` (req), `--target-branch` (req), `--repo-owner` (req), `--repo-name` (req), `--description`, `--account`
  - `read`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req, int)
  - `update`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req), `--title`, `--description`, `--draft` (bool)
  - `approve`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req), `--account`
  - `request-changes`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req), `--account`
  - `merge`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req), `--strategy` (merge_commit/squash/fast_forward), `--account`
  - `list-tasks`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req)
  - `create-task`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req), `--content` (req), `--comment-id` (int)
  - `update-task`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req), `--task-id` (req, int), `--content`, `--state` (RESOLVED/UNRESOLVED)
  - `diffstat`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req)
  - `diff`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req), `--path`, `--context` (int)
  - `add-comment`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req), `--content` (req), `--file-path`, `--line` (int), `--account`
  - `list-comments`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req), `--include-resolved` (bool)
  - `resolve-comment`: `--repo-owner` (req), `--repo-name` (req), `--pr-id` (req), `--comment-id` (req, int)
- After the noop check: call the service method, marshal result to `json.MarshalIndent`, write to `os.Stdout`
- Expand `main_test.go` to smoke-test `pr read --noop --repo-owner x --repo-name x --pr-id 1` to confirm DI is valid
- Run `make lint && make test`
- Write summary to `docs/implementation/bbmd/summary-task-2.3.md`
- Success criteria: `make lint && make test` pass

---

**Task 2.4: Implement `cmd/bbmd/file.go` and finalize**
- Add `file` group with `content` sub-command:
  - Flags: `--repo-owner` (req), `--repo-name` (req), `--commit` (req), `--path` (req), `--account`
  - Executor: `container.Invoke` → `*app.BitbucketService`, check `if noop { return nil }`, call `GetFileContent`, `json.MarshalIndent` → `os.Stdout`
- Update `.gitignore`: add `bbmd.log`
- Update `AGENTS.md` run section:
  - `go run ./cmd/bbmd --help`
  - `go run ./cmd/bbmd auth status`
  - `go run ./cmd/bbmd pr read --repo-owner <owner> --repo-name <repo> --pr-id <id>`
  - Note: `--noop` skips execution (dry-run / startup checks)
- Run `make lint && make test`
- Write summary to `docs/implementation/bbmd/summary-task-2.4.md`
- Success criteria: `make lint && make test` pass

---

**Compress implementation summaries**
- Follow [compress-implementation-summaries.md](/.context/compress-implementation-summaries.md) to compress the implementation summaries.
