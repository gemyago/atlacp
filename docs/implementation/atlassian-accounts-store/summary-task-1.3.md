# Task 1.3 summary: accounts store — mutations and save

## Implemented

- **Mutations:** `Upsert`, `Remove`, and `SetDefault` on `AccountsStore`, each building the next slice, running `app.ValidateAtlassianAccounts` on the full result, and committing only on success (same error wrapping as `LoadFromFile` for validation failures). `Remove` / `SetDefault` use `ErrAccountNotFound` when the name is absent.
- **Persistence:** `SaveToFile(path)` marshals `{"accounts":[...]}`, rejects empty path and invalid state, and writes via a temp file in the destination directory + `Sync` + `Rename` for atomic replace.
- **Tests:** Round-trip load → save → reload; default switching with `SetDefault`; removing the sole account fails with validation error and leaves state unchanged; `ErrAccountNotFound` for missing names; failed invalid `Upsert` preserves state; concurrent `Upsert` stress with per-goroutine unique names. Additional cases: non-file `LoadFromFile`, `SaveToFile` empty path / invalid empty store / missing parent dir for atomic write, `Upsert` replace-by-name, successful `Remove` of a non-default account (coverage toward the 90% per-file gate).

## Uncertainties / deviations

- None. API matches the plan’s minimal set (`Upsert`, `Remove`, `SetDefault`, `SaveToFile` only; no separate `Save()` because no path is stored on the struct).
