# Task 1.2 — `AccountsStore`: `ListAccounts` and missing-file tolerance

## What was implemented

- **`ListAccounts() []app.AtlassianAccount`**: returns a copy of the current accounts under `RWMutex` read lock (empty slice is non-nil with length 0).
- **`NewAccountsStoreWithDeps`**: if `os.Stat` reports `NotExist` for `ConfigPath`, logs a **warning** (path in attrs) and returns a non-nil store with no accounts instead of an error. Empty `ConfigPath` still errors unchanged.
- **`accounts_store_test.go`**: updated “missing file” constructor test; added `ListAccounts` tests (empty store, after `LoadFromFile`, after `Upsert`, snapshot independence).

## Uncertainties / deviations from the plan

- none

## Verification

- `make lint`: pass  
- `make test`: pass
