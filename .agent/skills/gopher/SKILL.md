---
name: gopher
description: >
  Portable Go style and testing habits: interfaces, errors, slog, tests, TDD mindset.
  Pair with each project's AGENTS.md (or equivalent) for versions, layout, commands, and tooling.
---

# Gopher — Go coding and testing (generic)

**Use with project docs.** This skill is **not** a substitute for the repository’s own rules. Always read the nearest **`AGENTS.md`** (or your project’s equivalent: `CONTRIBUTING.md`, team wiki, Cursor rules) for:

- Go / toolchain versions and how they are pinned
- Repo layout, module boundaries, and where application code lives
- How to run **lint** and **tests**, and what “done” means before you report back
- Codegen (OpenAPI/protobuf/sqlc/etc.), mock generation, and migration commands
- Security and compliance expectations for this codebase
- Framework-specific rules (ORM, HTTP routers, DI, specialized SDKs, etc.)

Nearest project instructions **win** when they conflict with anything below.

---

## Go coding guide (language-level)

### Key considerations

- **Accept interfaces, return structs** canonical way to structure code:
  - Dependencies (services, repos, clients or any other components) **SHELL** be accepted as consumer defined **interfaces** (testability, boundaries).
  - Return **concrete types** (structs) from constructors.
- Define **interfaces next to the consumer** (same file by default); split only when size or reuse demands it.
- Use **`log/slog`** (or the project’s logging facade) as an injected dependency — avoid global loggers unless the project already standardized on them.
- Wrap errors: `fmt.Errorf("context: %w", err)` (or the project’s error helpers).
- Methods with more than 3 arguments (context does not count) is a warning sign. Use params struct instead.
- Names such as "tools", "helpers", "utils" e.t.c are banned. Use descriptive names instead.
- Prefer **functional options** for optional constructor parameters: a `type FooOption func(*Foo)` (or `*fooConfig`), `WithBar(...)` functions that set fields, and `NewFoo(opts ...FooOption)` applying them in order. Avoid a separate `NewFooWithOpts` when the zero-arg `NewFoo()` case is the default.


### Testing style and patterns (from the coding guide)

More detail is in **Testing best practices** below. Common points:

- Use random data in tests. Use faker (`github.com/jaswdr/faker/v2`) for variable test data; think twice before using fixed literals, this is a warning sign. If you think fixed literal is fine in some cases, you're wrong unless you prove it strongly.
- Tests in the **same package** as production code (or `_test` package if the project prefers that for black-box tests).
- Use **one top-level test function per unit** with **nested `t.Run`** for methods and scenarios.
- Avoid **static/shared state** across tests.
- Use a **single, consistent way** to build test dependencies (e.g. a shared `makeMockDeps` or fixture constructor **if the project uses that pattern**); avoid copy-pasted setup.
- Use `TestComponentName` and nested `t.Run` per API surface, avoid separate `TestComponentNameFunctionName` functions per method.
- Keep helpers **local**: if a helper is only used in one test, nest it inside that test.
- Compare **whole values** (`expected` vs `actual` struct)
- Use **`require.Error` / `require.ErrorIs`** (or equivalents) for error assertions.
- **Mocks:** follow the project’s documented approach (codegen tool, hand-written fakes, etc.).
- Use **`t.Context()`** over `context.Background()` / `context.TODO()` in tests when the Go version supports it.
- Avoid package-level test helpers; define them inside the relevant `TestXxx` (e.g. as a closure used by nested `t.Run` blocks) or inside a single `t.Run` when only that case needs them.

---

## Testing best practices (embedded, mindset-level)

#### Core Principles

##### Follow TDD

- Work in small steps; stub if needed, then add a failing test, then minimal code to pass; repeat.
- If the project documents a TDD workflow or other strategy, follow that document.

##### Testing philosophy

- **Focus on business logic** — what the code must guarantee.
- **Prefer** single test for single code branch
- **Avoid excessive tests** — skip scenarios that do not match real use cases.
- **Avoid splitting one behavior across many tests** — one behavior, one clear story when possible.
- **Test behavior, not implementation** — observable outcomes over internal structure.
- **Avoid** testing framework internals or standard library functions.

If project documents specific mocking strategy, follow it. If you have to create mock manually, follow these principles:

- **Pragmatic mocks** — simplest thing that isolates the unit under test.
- **Minimal setup** — only what the case needs.
- **Clear names** — names should say what behavior is under test.