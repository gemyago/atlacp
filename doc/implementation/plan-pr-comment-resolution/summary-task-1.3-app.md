# Summary: App layer mapping and ListPRComments enrichment (PR comment resolution)

## What changed

- Added **`BitbucketPRComment`** (flat `resolved` + standard comment fields) and **`BitbucketListPRCommentsResult`** (pagination + `[]BitbucketPRComment`), plus **`BitbucketResolvePRCommentParams`**.
- **`ListPRComments`** now returns `*BitbucketListPRCommentsResult`: copies pagination from the Bitbucket list response, maps each `PRComment` with `prCommentToBitbucketPRComment`, and uses **`ResolvedStateFromResolutionJSON`**. When resolution is ambiguous (e.g. empty or `{}`), it calls **`GetPRComment`** once per comment and re-parses resolution; if still ambiguous, **`resolved`** is false.
- **`ResolvePRComment`** delegates to the Bitbucket client with validation on repo, PR id, and comment id.
- Extended **`bitbucketClient`** in `ports.go` with **`GetPRComment`** and **`ResolvePRComment`**; mocks regenerated with mockery v2.
- **`internal/api/mcp/controllers`**: `bitbucketService` **`ListPRComments`** return type updated to `*app.BitbucketListPRCommentsResult`; controller tests adjusted to return app DTOs.
- **`bitbucket_test.go`**: list tests use `Resolution: json.RawMessage("null")` where a single list call is enough; added a test that expects **`GetPRComment`** when list returns `{}` resolution; added **`ResolvePRComment`** tests.

## Verification

- `make lint` — pass  
- `make test` — pass (total coverage **95.3%** per project harness)
