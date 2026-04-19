# Plan: Atlassian accounts store and management component

## 1. Introduction / Overview

**Problem:** Today the Atlassian accounts JSON file path is supplied via config/flags (`config.atlassian.accountsFilePath`), and `internal/services/atlassian_accounts.go` reads that file once at repository construction. The data is effectively **read-only** after startup; there is no first-class place to **persist** or **mutate** the account set (add/update/remove, change default, reload from disk).

**Goal:** Introduce a **standalone component** that can **hold**, **validate**, and **manage** a collection of `app.AtlassianAccount` values, including **optional persistence** to a JSON file (same on-disk shape as today: `{ "accounts": [...] }`). The component should be usable and testable on its own.

**Non-goals (this phase):**

- Wiring the component into DI, MCP, or replacing `NewAtlassianAccountsRepository`.
- Implementing the documented multi-path search order from `doc/atlassian-accounts-schema.md` (unless trivially useful inside the component as optional helpers—defer if ambiguous).
- UI or HTTP APIs for editing accounts.

## 2. Business Logic

- **Single source of truth in memory:** The component owns a slice (or map-by-name) of accounts that obeys the same rules as the existing file-backed repository.
- **Validation:** On load and after any mutation that changes the set, the full configuration must satisfy the same invariants already enforced in `validateAccountsConfig` and related helpers in `internal/services/atlassian_accounts.go` (at least one account, unique names, exactly one default, each account has Bitbucket and/or Jira with required token fields when present).
- **Persistence:** When a file path is configured, saving writes JSON atomically (write temp file in same directory, then rename) so partial writes do not corrupt the primary file.
- **Concurrency:** If the component can be called from multiple goroutines, exported methods should be safe (e.g. `sync.RWMutex` or document “caller must serialize” if explicitly single-threaded—prefer mutex for a reusable component).
- **Read accessors:** Support resolving the default account and lookup by name, aligned with `app.AtlassianAccountsRepository` behavior where it matters (e.g. `ErrNoDefaultAccount`, `ErrAccountNotFound`), so a future adapter can be thin.

## 3. High-Level Architecture

| Piece | Role |
|--------|------|
| **Accounts store (new)** | Holds accounts, validates, optional load/save to path, mutation API |
| **Existing `atlassianAccountsRepository`** | Unchanged in this phase; remains the DI-backed read-only implementation |
| **`app.AtlassianAccount` / errors in `ports.go`** | Reuse types and sentinel errors; no duplicate domain types |

Placement: implement under `internal/services/` (e.g. new file `accounts_store.go` or small subpackage `internal/services/accountsstore`) to stay next to JSON/file concerns and existing validation. Alternatively, extract **pure validation** into `internal/app` if both the repository and the store need it without import cycles—decide during detailed design (see Key Decisions).

## 4. Detailed Architecture

### 4.1 Public API (illustrative—finalize in implementation)

A concrete type (name TBD, e.g. `AccountsStore` or `ManagedAtlassianAccounts`) should expose, at minimum:

- `NewAccountsStore(...)` — options: logger (optional), initial path, maybe `WithAccounts` for tests without disk.
- `LoadFromFile(path string) error` — read JSON, validate, replace in-memory state.
- `SaveToFile(path string) error` or `Save() error` if path fixed at construction — persist current state.
- Mutations: e.g. `Upsert(account)`, `Remove(name string)`, `SetDefault(name string)` — each validates full result before committing.
- Queries: `GetDefault(ctx)` / `GetByName(ctx, name)` matching existing semantics, or unexported helpers if you prefer not to duplicate `context` on a non-port type—prefer consistency with existing repository.

Internal JSON envelope should match `atlassianAccountsConfig` (`{"accounts":[...]}`) for compatibility with existing files and `doc/atlassian-accounts-schema.md`.

### 4.2 Validation reuse

**Preferred:** Extract shared validation into a single place (e.g. `validateAtlassianAccountsConfig(accounts []app.AtlassianAccount) error` in `internal/app` or a small `internal/accounts` package) and call it from both:

- Existing `atlassianAccountsRepository` construction path, and  
- The new store.

If extraction risks a large refactor in one task, split: **Task 1** extract validation + update existing tests; **Task 2** add store using extracted validation.

### 4.3 Files

| Action | Path |
|--------|------|
| Likely **new** | `internal/services/accounts_store.go` (or `internal/services/accountsstore/*.go`) |
| Likely **new** | `internal/services/accounts_store_test.go` |
| **Maybe update** | `internal/services/atlassian_accounts.go` — delegate to shared validation |
| **Unchanged (this plan)** | `cmd/mcp/*`, `internal/services/register.go`, `internal/di/*` |

## 5. Key Architectural Decisions

1. **No integration in this PR:** The new type is not registered in `dig`; callers in tests (or temporary `main` snippets) construct it directly.
2. **File format unchanged:** Same JSON schema as today so future wiring can swap implementations without migrating files.
3. **Validation is shared:** Avoid duplicating `validateAccountsConfig` logic long-term; one function used by repository and store.
4. **Atomic writes** for `Save` to avoid corrupting user credential files on crash.

## 6. Uncertainties

- **Exact mutation API:** Whether to expose full CRUD or a minimal set (`Upsert` + `Remove` + `SetDefault`). The plan assumes minimal sufficient set; expand only if product needs list/rename batch operations.
- **Package layout:** Single file vs. `accountsstore` subpackage—choose based on line count and team preference (keep first iteration small).
- **Context on query methods:** If the store is not implementing `AtlassianAccountsRepository` yet, `context.Context` on getters may be omitted for simplicity; align with project conventions (`t.Context()` in tests).

## 7. Related Files

- `internal/services/atlassian_accounts.go` — current load/validate/read-only repository
- `internal/services/atlassian_accounts_test.go` — patterns for temp JSON files and validation cases
- `internal/app/atlassian_accounts.go` — `AtlassianAccount`, `AtlassianToken`
- `internal/app/ports.go` — `AtlassianAccountsRepository`, `ErrNoDefaultAccount`, `ErrAccountNotFound`, `ErrAccountConfigInvalid`
- `doc/atlassian-accounts-schema.md` — schema and validation rules (documentation)
- `cmd/mcp/root.go` — flag binding (reference only; no changes in this phase)

## 8. Task List

Follow **TDD** where practical: write failing tests for new behavior, then implement. After each task: `make lint` and `make test` (project completion protocol).

**Task 1.1: Extract shared validation for Atlassian accounts config**

- Move or duplicate-then-unify validation logic from `internal/services/atlassian_accounts.go` into a callable function that accepts `[]app.AtlassianAccount` (or a small struct with an `Accounts` field) and returns an error equivalent to current behavior.
- Update `NewAtlassianAccountsRepository` to use the shared validator (behavior unchanged).
- Write or adjust tests in `atlassian_accounts_test.go` so all existing cases still pass; add a test that calls the exported validator directly for one representative invalid case (if the function is exported) or test via repository only (if validation stays package-private in `services`).
- Success criteria: No behavior change for repository; `make lint` and `make test` clean.

**Task 1.2: Implement accounts store — load, validate, query**

- Add the new concrete type and constructor(s).
- Implement `LoadFromFile` (or equivalent) using the shared validator.
- Implement read-style methods consistent with existing semantics (`GetDefaultAccount`, `GetAccountByName`) for the in-memory state.
- Tests: happy path load from temp file; missing file; invalid JSON; validation failure (empty list, duplicate names, two defaults, etc.)—reuse scenarios from `TestAtlassianAccountsRepository` where applicable.
- Success criteria: Tests pass; component not yet required to persist mutations.

**Task 1.3: Implement accounts store — mutations and save**

- Implement mutations (`Upsert`, `Remove`, `SetDefault` or agreed minimal set) with full re-validation after each change.
- Implement `SaveToFile` / `Save` with atomic replace.
- Tests: round-trip save and reload; default switching; remove last default should error; concurrent calls if mutex is used (optional table-driven test).
- Success criteria: `make lint` and `make test` clean.

**Task 1.4: Documentation touchpoint (code comments only)**

- Package/file doc comment on the new type describing scope: in-memory management + optional file persistence; not yet wired to DI.
- No new markdown docs unless the team asks.

**Compress implementation summaries**

- Follow [.context/compress-implementation-summaries.md](/.context/compress-implementation-summaries.md) to compress the implementation summaries after all numbered tasks are done.
