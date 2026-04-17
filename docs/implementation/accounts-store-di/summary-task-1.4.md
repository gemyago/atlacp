# Task 1.4 — Comments and implementation summary

## What was implemented

- **Refreshed** the file-level comment at the top of `internal/services/accounts_store.go`: describes production DI registration (`NewAccountsStoreWithDeps`, `app.AtlassianAccountsRepository`, `register.go`) and test-only construction; removed the obsolete “not registered in DI” wording.
- **Added** `docs/implementation/accounts-store-di/implementation-summary.md` with a short outcome summary (port wiring, repository removal, pointer to prior atlassian-accounts-store doc).
- **Updated** `docs/implementation/atlassian-accounts-store/implementation-summary.md`: overview now notes DI wiring with links to the `accounts-store-di` plan and implementation summary; Task 1.4 bullet adjusted so it no longer claims the store is unwired.

## Uncertainties / deviations

- **None.**
