# Implementation Summary: PR comment resolution (list + resolve MCP)

**Plan:** [plan-pr-comment-resolution.md](./plan-pr-comment-resolution.md)

## Overview

The Bitbucket service and app layer now expose a flat `resolved` flag on PR comments, derived from the list response’s `resolution` JSON only. The MCP server lists that field in JSON and adds `bitbucket_resolve_pr_comment` to call Bitbucket’s resolve endpoint. Integration tests in the Bitbucket MCP doc were extended to cover resolve and re-list behavior.

## Tasks

### Task 1.1: Service models and resolution parsing

Confirmed `PRComment` resolution storage and `ResolvedStateFromResolutionJSON`; tests live in `pr_comment_resolution_test.go`.

### Task 1.2: HTTP client — resolve comment

`ResolvePRComment` (POST to the resolve URL with Bearer auth), `CommentResolution` in models, and `resolve_pr_comment_test.go`.

### Task 1.3: App layer mapping and ListPRComments enrichment

App types `BitbucketPRComment` / `BitbucketListPRCommentsResult`, list enrichment via resolution parsing only, and `ResolvePRComment` wired through ports and mocks.

### Task 1.4: MCP + integration doc

`bitbucketService` port extended with resolve, new `bitbucket_resolve_pr_comment` tool, list tool docs for `resolved`, regenerated mocks, controller tests, and Test 5 extended in `bitbucket-mcp-integration-tests.md`. `makeListPRCommentsHandler` was introduced to satisfy `funlen` after expanding the list tool description.

## Deviations & notes

- **List-only resolution:** `resolved` is derived only from each list item’s `resolution` field (no per-comment GET).
- **Summary filename vs. heading:** One MCP task summary file uses a `1.4-mcp` filename while the heading referenced Task 1.3 in the document body.

## Completion

- Lint: ✓
- Tests: ✓ (total coverage ~95.1% at time of completion)
