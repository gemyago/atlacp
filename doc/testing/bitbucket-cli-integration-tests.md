# Bitbucket CLI (`bbmd`) Integration Testing

This document provides step-by-step instructions for testing the Bitbucket CLI integration (`bbmd`).

It is expected that the prompt to start the test will have the following form:
```markdown
Given the following Bitbucket repository `<file system path to the repository>`, run the tests as per the instructions in the `bitbucket-mcp-integration-tests.md` file.
```

**Important**: If anything fails - AI must not attempt to fix the issue, instead it must report and continue with next step if possible.

## Prerequisites done by the user

Assume the following is already prepared:

1. Sub-agents run the CLI with **`go run ./cmd/bbmd …`** from this (**atlacp**) repository root (Go toolchain required). If that fails, report to the user and ask to fix.
2. Proper Atlassian account configuration is set up with at least two accounts:
   - A default account - assume it is named "user" if not otherwise mentioned
   - A secondary account named "bot"
  Check this by running `go run ./cmd/bbmd auth status`
3. AI must lern how to work with bbcp using skill sub-comment: `go run ./cmd/bbmd skill`. If documentation is incomplete or confusing, AI should report this as well. This is part of the test.

For Git/SSH issues point the user on the [SSH troubleshooting](./README.md#ssh-troubleshooting). Don't read or do anything about it yourself, report and halt.

## Running scenarios

Please follow [scenarios](./integration-test-scenarios.md) and run them one by one or as requested by user.