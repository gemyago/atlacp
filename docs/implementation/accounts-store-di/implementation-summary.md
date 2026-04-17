# Implementation summary: AccountsStore DI and repository removal

**Plan:** [plan-accounts-store-di.md](./plan-accounts-store-di.md)

## Outcome

- **`AccountsStore` is the app port implementation:** `NewAccountsStoreWithDeps` loads from `config.atlassian.accountsFilePath`; `internal/services/register.go` provides `*AccountsStore` and `di.ProvideAs[*AccountsStore, app.AtlassianAccountsRepository]`.
- **Removed** the duplicate `atlassianAccountsRepository` implementation and merged its tests into `accounts_store_test.go`.
- **Code:** `internal/services/accounts_store.go` documents DI wiring at file level; compile-time `var _ app.AtlassianAccountsRepository = (*AccountsStore)(nil)`.

Prior work on validation and the store API is summarized in [../atlassian-accounts-store/implementation-summary.md](../atlassian-accounts-store/implementation-summary.md).
