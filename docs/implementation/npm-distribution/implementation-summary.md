# Implementation Summary: npm Distribution for atlacp Tools

**Plan:** [plan-npm-distribution.md](./plan-npm-distribution.md)

## Overview

Implemented npm-based distribution for the atlacp tools, including platform-specific package manifests, an install script, release packaging automation, and a new GitHub release publish workflow. The work also renamed the MCP binary to `bbcp`, added macOS build targets, and consolidated the default runtime data path under `~/.atlacp`.

## Tasks

### Task 1.0: Rename `cmd/mcp` binary to `bbcp`
Renamed the MCP server entrypoint to `bbcp`, updated the Cobra root command name, and switched Docker cleanup to use dynamically generated image names instead of a hardcoded package.

### Task 1.1: Extend build platforms to include darwin
Expanded the build matrix to include `darwin/amd64` and `darwin/arm64`, verifying the dist output for all four target platforms.

### Task 1.2: Create npm package structure (package.json files)
Created the main `@atlacp/install` manifest and the four platform-specific optional package manifests with the required `os`/`cpu` gating.

### Task 1.3: Implement install script with unit tests
Added the installer module and unit tests covering platform detection, binary copying, PATH handling, and shell config selection.

### Task 1.4: Add `npm-packages` Make target
Added packaging automation to populate npm `bin/` directories, inject versions with `jq`, and update build docs for the new target.

### Task 1.5: Update default accounts file path
Changed the default accounts file location to `~/.atlacp/accounts.json` and updated the resolver tests accordingly.

### Task 1.6: Create GitHub Actions workflow to publish npm
Added a release-triggered workflow that builds binaries, populates npm packages, and publishes the platform packages followed by the main installer package with provenance.

## Deviations & notes

- Task 1.1 added a small coverage test in `internal/testing/mocks` because `make test` failed on coverage collection without at least one test in that package.
- Task 1.4 also added `.npm-packages` cleanup in the Makefile `clean` target.
- Task 1.5 adjusted a `//nolint` directive in `internal/di/dig.go` to satisfy linting during task verification.

## Completion

- Lint: ✓
- Type check: ✓
- Tests: ✓
