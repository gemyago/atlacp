# Decline Pull Request Implementation Plan

## Goal

Expose Bitbucket Cloud pull-request decline operations through the client, application service, `bbmd` CLI, and MCP server.

## Tasks

### Task 1.1: Add Bitbucket and application decline support

Implement the bodyless Bitbucket decline endpoint, application parameters and service method, account selection, generated client mock, and unit tests.

### Task 1.2: Add CLI and MCP interfaces

Add `bbmd pr decline` and `bitbucket_decline_pr`, forward the updated pull-request response, validate MCP IDs as positive integers without truncation, and cover discovery/parity/handler behavior.

### Task 1.3: Update documentation and verify

Document supported interfaces and endpoint status, run focused and full verification, and record live-test status.
