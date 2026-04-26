# Bitbucket MCP Integration Testing

This document provides step-by-step instructions for testing the Bitbucket MCP integration. These tests are designed to be executed by AI assistants to verify the correct functioning of the Bitbucket MCP tools.

It is expected that the prompt to start the test will have the following form:
```markdown
Given the following Bitbucket repository `<file system path to the repository>`, run the tests as per the instructions in the `bitbucket-mcp-integration-tests.md` file.
```

**Important**: If anything fails - AI must not attempt to fix the issue, instead it must report and continue with next step if possible.

## Prerequisites done by the user

It should be assumed that below is already prepared by the user:
1. The ATLACP service must be running (or stdio configured) and MCP tools are registered and available to AI assistant.
2. Proper Atlassian account configuration is set up with at least two accounts:
   - A default account - assume it is named "user" if not otherwise mentioned
   - A secondary account named "bot"

For Git/SSH issues point the user on the [SSH troubleshooting](./README.md#ssh-troubleshooting). Don't read or do anything about it yourself, report and halt.

## Running scenarios

Please follow [scenarios](./integration-test-scenarios.md) and run them one by one or as requested by user.