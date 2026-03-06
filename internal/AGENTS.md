<!-- Nearest AGENTS.md takes precedence. Scope: guidance for packages under internal/. Keep concise; link to canonical code. -->

## Purpose

Please review project level [AGENTS.md](../AGENTS.md). This file complements it with internal/ specific details.

## Architecture Overview

**Key architectural decisions**:
- **Consumer-Defined interface** All components should follow "accept interface and return struct" principle for dependencies:
  - Component dependencies (services, repositories, etc.) should be accepted as interfaces (for flexibility and testability)
  - Return types should be concrete structs (for clarity and avoiding unnecessary abstraction)
  - Strong justification is required to deviate from this pattern
- Define interfaces next to consumer by default (in a same file). Move to a separate if getting bigger or used by multiple consumers in the same package.

The project follows Pragmatic Layered Architecture with layers mapped as follows:
- APIs: `internal/api` (HTTP, MCP, ...) - "world" interacts with the system here
- Application layer: `internal/app` (business logic)
- Services: `internal/services` (external API clients, HTTP, Bitbucket, Jira, etc.)

Additional notes:
- Config loading: `internal/config` (embedded JSON via viper)
- Register components and app wiring. Example: [internal/app/register.go](./app/register.go)

## Application Layer

The Application layer contains all business logic and defines the core data structures of the system.

### Dependency Rules
- **Inward Imports**: The dependency flow is strictly inward. The Application layer must not import from the Services or API layers.
- **Interface-Based Interaction**: Interactions with external systems (third-party APIs, etc.) are defined via interfaces within the Application layer (`ports.go`). Services layer provides the concrete implementations.
- **Boundary Isolation:** Request/Response types from the API/CLI layers must never enter the Application layer. They must be mapped to Application types at the boundary.

### Application Components
Example components:
- Bitbucket service: [internal/app/bitbucket.go](./app/bitbucket.go)
- Bitbucket service tests: [internal/app/bitbucket_test.go](./app/bitbucket_test.go)
- Ports (outbound interfaces): [internal/app/ports.go](./app/ports.go)

## API Layer

API layer follows the dependency inversion principle by defining interfaces for the required application layer components and relying on DI to provide concrete implementations.

DI registration example: [internal/api/mcp/controllers/register.go](./api/mcp/controllers/register.go):
* Concrete implementation of application layer interface is bound using `di.ProvideAs` (e.g. `di.ProvideAs[*app.BitbucketService, bitbucketService]`)

### MCP Tools
- Example MCP tool controller: [internal/api/mcp/controllers/bitbucket.go](./api/mcp/controllers/bitbucket.go)

## Logging and Diagnostics

- Use log/slog via DI; no globals. See [internal/diag/slog.go](./diag/slog.go) and [internal/diag/testing.go](./diag/testing.go)
- Follow `.golangci.yml` slog rules; prefer context-aware logging.
- Components that need logging should accept `RootLogger` as dependency and create a child logger with component name: `logger := deps.RootLogger.WithGroup("http-server")`

## Task completion protocol

Follow project wide completion protocol. No exceptions.
