# Plan: Wire AccountsStore into DI and remove file-backed repository

## 1. Introduction / Overview

**Problem:** `AccountsStore` already implements the same read surface as `app.AtlassianAccountsRepository`, but production still registers `NewAtlassianAccountsRepository`, duplicating load/validation logic and maintaining two types.

**Goal (phase 1):** Add a **constructor with dependencies** mirroring `AtlassianAccountsRepositoryDeps` (logger + named config path), call existing `LoadFromFile` from that constructor so load stays a separate, testable method, register `*AccountsStore` in the services DI graph, and expose it as `app.AtlassianAccountsRepository` via `di.ProvideAs`. Then **delete** `atlassianAccountsRepository` and its tests, leaving `AccountsStore` as the single implementation.

**Non-goals (this phase):**

- MCP/HTTP APIs for mutating accounts at runtime (mutations exist on the type but remain unused by app wiring unless already called elsewhere).
- Changing JSON schema or validation rules.
- Implementing multi-path account file search from `doc/atlassian-accounts-schema.md`.

## 2. Business Logic

- **Startup:** If `config.atlassian.accountsFilePath` is empty, fail (same as today). If the file is missing, unreadable, invalid JSON, or fails `app.ValidateAtlassianAccounts`, fail construction so the process does not start with bad credentials—same operational contract as `NewAtlassianAccountsRepository`.
- **Runtime:** Consumers continue to use only `GetDefaultAccount` / `GetAccountByName` through the port; behavior matches validated in-memory data as today.

## 3. High-Level Architecture

| Piece | Role |
|--------|------|
| **`AccountsStoreDeps` + `New…(*AccountsStore, error)`** | DI constructor: validate deps, `LoadFromFile(path)`, return store |
| **`LoadFromFile`** | Unchanged public method; used by constructor and tests |
| **`di.ProvideAs[*AccountsStore, app.AtlassianAccountsRepository]`** | Satisfies the app port with the concrete store (same pattern as `ProvideAs[*bitbucket.Client, bitbucketClient]` in `internal/app/register.go`) |
| **Remove `atlassianAccountsRepository`** | No second implementation; tests consolidated onto store + constructor |

## 4. Detailed Architecture

### 4.1 Constructor and deps

- Add a struct `AccountsStoreDeps` (name may match team preference) with `dig.In`, `RootLogger *slog.Logger`, and `ConfigPath string` with tag ``name:"config.atlassian.accountsFilePath"`` — **same shape as** `AtlassianAccountsRepositoryDeps` in `internal/services/atlassian_accounts.go`.
- Add a constructor, e.g. `NewAccountsStoreWithDeps(deps AccountsStoreDeps) (*AccountsStore, error)` (exact name TBD; avoid clashing with the current parameterless `NewAccountsStore()` — see below).
- Constructor behavior should mirror the old repository startup **in order**: reject empty path; optionally keep `os.Stat` before read if you want identical “file not found” behavior to today (repository stat’d first); then `store.LoadFromFile(configPath)` on a new `AccountsStore` (or empty `&AccountsStore{}`).
- **Logger:** Repository stored `logger := deps.RootLogger.WithGroup("atlassian-accounts")` but did not use it in getters. For parity, either attach the same child logger on `AccountsStore` for future use or omit until needed—prefer attaching for consistency with the removed type.

### 4.2 DI registration

- Update `internal/services/register.go`:
  - Import `github.com/gemyago/atlacp/internal/app` for `app.AtlassianAccountsRepository` in `ProvideAs`.
  - Replace `NewAtlassianAccountsRepository` with the new constructor plus:
    - `di.ProvideAs[*AccountsStore, app.AtlassianAccountsRepository]`
- Confirm `dig` resolves: constructor produces `*AccountsStore`; `ProvideAs` consumes it and produces the interface (same pattern as app layer Bitbucket client).

### 4.3 Remove repository

- Delete `internal/services/atlassian_accounts.go`.
- Fold scenarios from `internal/services/atlassian_accounts_test.go` into `accounts_store_test.go` (or a dedicated test function): empty path, missing file, invalid JSON, invalid validation—assert against **`NewAccountsStoreWithDeps`** (or the chosen constructor name) and/or `LoadFromFile`, preserving expectations where they match product behavior.
- Remove `//nolint:ireturn` from the deleted repository; the new provider returns `*AccountsStore`; `ProvideAs` yields the interface.

### 4.4 Empty store for tests

- Today tests use `NewAccountsStore()` to get an empty store. **Either:**
  - Replace call sites with `&AccountsStore{}` (valid zero value), **or**
  - Keep a tiny helper only for tests (same package) if repetition is noisy.
- Remove or repurpose the parameterless `NewAccountsStore()` so there is a single clear DI entry point and no duplicate “empty” API unless intentionally kept as a test convenience.

### 4.5 Compatibility and docs

- Update the file doc comment on `accounts_store.go`: it is **registered** in DI as `AtlassianAccountsRepository` (remove “not registered” wording).
- If `docs/implementation/atlassian-accounts-store/implementation-summary.md` or related notes claim the store is not wired, add a short note or link to this plan’s outcome (optional; avoid large doc churn).

## 5. Key Architectural Decisions

1. **Keep `LoadFromFile` public** — Constructor delegates to it; tests can still target load failures without going through DI.
2. **One implementation behind the port** — Drop `atlassianAccountsRepository` after DI wiring and tests pass.
3. **`ProvideAs` for the port** — Matches existing project convention; no new adapter type unless `dig`/`ProvideAs` constraints force one (unlikely).

## 6. Uncertainties

- **Exact constructor name:** `NewAccountsStoreWithDeps` vs renaming the zero-arg factory; choose one exported DI constructor and minimize public surface.
- **Error text parity:** `LoadFromFile` wraps missing-file errors with `%w`; old repository used a plain message without wrapping. Decide whether to align messages in the constructor path for identical log output (ops/debugging), or accept small differences.
- **`Stat` vs read-only missing file:** Repository used `Stat` then `ReadFile`; `LoadFromFile` only `ReadFile`. Behavior is effectively the same for normal files; edge cases (path is directory) may differ slightly—only address if tests covered it.

## 7. Related Files

| Action | Path |
|--------|------|
| Update | `internal/services/accounts_store.go` — deps struct, constructor calling `LoadFromFile`, optional logger field |
| Update | `internal/services/register.go` — new provider + `ProvideAs` |
| Update | `internal/services/accounts_store_test.go` — empty-store construction, constructor/DI-equivalent tests |
| Delete | `internal/services/atlassian_accounts.go` |
| Delete or merge | `internal/services/atlassian_accounts_test.go` → into `accounts_store_test.go` |
| Maybe update | `docs/implementation/atlassian-accounts-store/implementation-summary.md` (if it states store is unwired) |
| No change expected | `internal/app/ports.go` — interface unchanged; `internal/app/bitbucket_auth.go` — still depends on `AtlassianAccountsRepository` |

## 8. Task List

Follow **TDD** where practical: adjust or add failing tests for the new constructor and DI-visible behavior, then implement. After each task: `make lint` and `make test` (project completion protocol).

**Task 1.1: Constructor with deps and startup parity**

- Add `AccountsStoreDeps` (`dig.In`, `RootLogger`, named `config.atlassian.accountsFilePath`).
- Add constructor that enforces non-empty path (same error style as old repository), builds store, calls `LoadFromFile`, returns `*AccountsStore` or error.
- Optionally attach `RootLogger.WithGroup("atlassian-accounts")` to the store if adding a field.
- Write tests: empty path fails; missing file fails; happy path loads (reuse temp-file patterns from former `TestAtlassianAccountsRepository` where applicable).
- Success criteria: Tests pass; behavior matches previous repository startup for the covered cases.

**Task 1.2: DI registration and interface alias**

- Update `internal/services/register.go`: provide new constructor; add `di.ProvideAs[*AccountsStore, app.AtlassianAccountsRepository]`.
- Remove `NewAtlassianAccountsRepository` registration.
- Add compile-time check in `accounts_store.go` if desired: `var _ app.AtlassianAccountsRepository = (*AccountsStore)(nil)`.
- Smoke: `go run ./cmd/mcp stdio --env local --noop` (or project-standard noop startup) still resolves the container.
- Success criteria: Lint/tests green; DI builds without missing types.

**Task 1.3: Remove repository and consolidate tests**

- Delete `internal/services/atlassian_accounts.go`.
- Merge remaining repository tests into `accounts_store_test.go`; remove `atlassian_accounts_test.go`.
- Replace `NewAccountsStore()` usages in tests with `&AccountsStore{}` or an agreed test helper; remove redundant zero-arg constructor if obsolete.
- Success criteria: No references to `atlassianAccountsRepository` / `NewAtlassianAccountsRepository`; full `make test` green.

**Task 1.4: Comments and implementation summary**

- Refresh `accounts_store.go` package/file comments (DI wiring, no “unwired” claim).
- Add a short bullet to implementation summary under `docs/implementation/accounts-store-di/` or update existing atlassian-accounts-store summary per team practice.

**Compress implementation summaries**

- Follow [.context/compress-implementation-summaries.md](/.context/compress-implementation-summaries.md) to compress the implementation summaries after all numbered tasks are done.
