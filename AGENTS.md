<!-- AGENTS.md — README for machines. Nearest file in the tree wins (hierarchical precedence). Keep this concise, concrete, and executable. -->

## Overview

This is a golang backend project with OpenTelemetry integration for observability. Go version is defined in [go.mod](./go.mod) file and this is the primary source of truth for the version.

## Purpose & Precedence
- This file gives AI coding agents the exact commands and conventions to follow in this repo.
- Closest AGENTS.md to the edited file applies; walk up directories to root if none found.
- Treat this as living documentation: update it in the same PR as any build/test/arch changes.
- **ALWAYS** follow "Task Completion Protocol" prior to reporting task completion

## Quick Setup
- direnv is assumed to be already configured
- gobrew is used to manage Go versions
- Install deps/tools: `go mod download && go install tool`
- Lint: `make lint`
- Run tests: `make test`

## Test and Lint

- Run all tests: `make test`
- Lint entire codebase: `make lint`
- Run a specific test: `go test -v ./internal/... -run "^TestName$"`
- Attempt auto fixing linting issues: `bin/golangci-lint run --fix`

## Run (local)

AI must **always** use `--noop` flag to dry-run startup checks without external deps. Without this flag, processes will start in foreground and block.

- MCP server (stdio): `go run ./cmd/mcp stdio --env local --noop`
- MCP server (HTTP): `go run ./cmd/mcp http --env local --noop`

## Docker Images (multi-platform)
- Build local images (load): `make -C build docker/.local-images`
- Build & push remote images: `make -C build docker/.remote-images`
- Configure platforms/registries in: `build/build.cfg`

## Deploy (iteration)
- Install Helm toolchain: `make -C deploy tools`
- Render chart: `helm template deploy/helm/api-service --debug --name-template api-service -f deploy/helm/api-service/values.yaml`
- Install/upgrade (dry-run): `helm upgrade api-service deploy/helm/api-service --install --namespace atlacp -f deploy/helm/api-service/values.yaml --create-namespace --dry-run`

## Configuration & Environment
- Embedded configs: `internal/config/default.json`, `<env>.json`, optional `<env>-user.json`
- Common flags on all binaries: `--env`, `--log-level`, `--json-logs`, `--logs-file`
- Env vars prefix `APP_` (dots/dashes -> underscores). Examples: `APP_ENV=local`, `APP_DEFAULT_LOG_LEVEL=info`, `APP_JSON_LOGS=true`

## Codebase structure

Codebase is split on multiple parts:
- Primary application code is in `internal` folder, see [AGENTS.md](./internal/AGENTS.md) for more details
- Binaries are defined in `cmd` folder
- Build related stuff (docker) is in `build`, see [AGENTS.md](./build/AGENTS.md)
- Deployment related stuff is in `deploy`, see [AGENTS.md](./deploy/AGENTS.md)

## Code Style & Patterns
- Lint strictly: `make lint` (see `.golangci.yml`). Use `//nolint:<rule>` only with justification.
- Many linting issues are auto fixable with `bin/golangci-lint run --fix`, try running it to apply fixes prior to direct updates
- Wrap errors with `fmt.Errorf("<something>: %w", err)`

### Testing Style and Patterns

More detailed testing best practices are in [doc/testing-best-practices.md](./doc/testing-best-practices.md). Common principles:
- Define tests in same package
- Prefer a single top-level test function per component, with nested `t.Run` blocks organizing tests by method and their scenarios.
- Avoid static variables shared across tests
- Use makeMockDeps to initialize dependencies, no inline or repeated setup
- Use random data when possible, use faker (github.com/jaswdr/faker/v2)
- Don't pollute testing namespace - if helper functions are only used within one test, nest them inside that test function
- Compare entire structs when possible instead of individual fields (e.g `assert.Equal(t, expectedUser, actualUser)`)
- Use require.Error or require.ErrorIs when asserting errors
- Use `t.Context()` instead `context.Background()` OR `context.TODO()` in tests
- Use factory functions to create reusable random data
- Follow [mockery](.context/mockery.md) for defining and generating mocks

## Security
- NEVER hardcode secrets. Use env vars/secret stores. Authenticate to GHCR before push/pull when required.
- Validate/sanitize all external inputs. Do not disable security linters without explicit justification.

## Most common AI instructions

When asked to perform common tasks, AI must follow these instructions:
- [Create plan](.context/create-plan.md)
- [Create pull request](.context/create-pull-request.md)
- [Commit code](.context/commit.md)

## Living Doc Policy (update with code)
- Keep this file short (<150 lines) and actionable. Prefer linking to canonical code over long prose.
- Update AGENTS.md in the same PR when:
  - Build/test commands change
  - New architectural patterns are introduced
  - CI or deploy workflows change
- Avoid duplication with `README.md`, `build/README.md`, `deploy/README.md`. Link instead.

## References (human docs)
- Overview: `README.md`
- Build details: `build/README.md`
- Deploy details: `deploy/README.md`

## Task Completion Protocol

AI must always follow this protocol when completing tasks. The protocol varies by task type.

### Coding Task Completion Protocol

Apply this when any code files were changed (Go, YAML, config files, etc.).

**Always** perform these steps before reporting completion:
1. Run `make lint` and confirm no errors
2. Run `make test` and confirm all tests pass
3. Verify AGENTS.md is updated if commands, workflows, or architecture changed

Report task completion status:
- Lint: ✓ no errors
- Tests: ✓ all passing, coverage XX.XX%
- AGENTS.md: ✓ updated / no changes needed

**Note:** Failing tests or lint errors mean the task is NOT complete. All failures must be resolved before completion.

### Non-Coding Task Completion Protocol

For tasks not involving code changes (investigation, documentation review, committing, etc.):
- Summarize findings or actions taken
- Confirm any deliverables were produced
- No lint/test protocol required
