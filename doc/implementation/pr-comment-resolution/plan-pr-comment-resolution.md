# Plan: PR comment resolution (list + resolve MCP)

## Introduction / Overview

Bitbucket’s list-comments API often returns `resolution: {}` (empty object), while fetching a single comment returns a full `resolution` object. We expose a flat `resolved` boolean on MCP list output, map it in the application service (not in the MCP controller), add a `POST .../comments/{id}/resolve` client and `bitbucket_resolve_pr_comment` MCP tool, and extend Test 5 in the Bitbucket MCP integration doc.

## Business Logic

- Parse `resolution` JSON when it carries meaningful fields (`resolved`, `resolved_by`, `resolved_on`, etc.).
- When resolution is absent, `null`, or empty `{}`, treat as not resolved (list payload only; no per-comment GET).
- Resolve: call Bitbucket resolve endpoint; no request body.

## High-Level Architecture

- **Services (`internal/services/bitbucket`)**: Raw API types (`PRComment` + `resolution` as `json.RawMessage`), `ResolvePRComment`, `ResolvedStateFromResolutionJSON`.
- **App (`internal/app`)**: `BitbucketListPRCommentsResult` / `BitbucketPRComment` with `resolved`; `ListPRComments` enriches when needed; `ResolvePRComment`.
- **MCP**: List tool returns app DTO JSON; new resolve tool.

## Key Architectural Decisions

- Application layer returns `BitbucketListPRCommentsResult` (not `*bitbucket.ListPRCommentsResponse`) so MCP stays flat without Bitbucket’s nested `resolution` object.
- No N+1 GETs; resolution state comes only from the list payload.

## Uncertainties

- Exact shape of Bitbucket’s full `resolution` object may vary; parser accepts common fields and treats non-empty objects with resolver metadata as resolved.

## Related Files

- `internal/services/bitbucket/models.go`, new `pr_comment_resolution.go`, `resolve_pr_comment.go`, tests
- `internal/app/bitbucket.go`, `ports.go`, mocks
- `internal/api/mcp/controllers/bitbucket.go`, `ports.go`, mocks, tests
- `doc/testing/bitbucket-mcp-integration-tests.md`

## Task List

**Task 1.1: Service models and resolution parsing**

- Add `Resolution json.RawMessage` to `PRComment`.
- Implement `ResolvedStateFromResolutionJSON` and unit tests (TDD).

**Task 1.2: HTTP client — resolve comment**

- Add `ResolvePRComment` with tests (httptest).

**Task 1.2: App layer mapping and ListPRComments enrichment**

- Add `BitbucketPRComment`, `BitbucketListPRCommentsResult`.
- Implement `ListPRComments` enrichment + `ResolvePRComment`; update `bitbucketClient` port and regenerate mocks.
- Update `bitbucket_test.go` (set `resolution` in list mocks as needed).

**Task 1.3: MCP + integration doc**

- Update `bitbucketService` port, list tool JSON, add `bitbucket_resolve_pr_comment`.
- Regenerate controller mocks; update controller tests.
- Extend Test 5 in `bitbucket-mcp-integration-tests.md`.

**Task 1.4: Completion**

- Run `make lint` and `make test`.

**Task 1.5: Compress implementation summaries**

- Follow [compress-implementation-summaries.md](/.context/compress-implementation-summaries.md) when per-task summary files exist (optional if none).
