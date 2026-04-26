# atlacp

[![Build](https://github.com/gemyago/atlacp/actions/workflows/build-flow.yml/badge.svg)](https://github.com/gemyago/atlacp/actions/workflows/build-flow.yml)
[![Coverage](https://raw.githubusercontent.com/gemyago/atlacp/test-artifacts/coverage/golang-coverage.svg)](https://htmlpreview.github.io/?https://raw.githubusercontent.com/gemyago/atlacp/test-artifacts/coverage/golang-coverage.html)

Bitbucket tooling for developers and AI agents.

The project provides two entrypoints:

- `bbmd` - direct Bitbucket CLI for humans, scripts, and agent workflows.
- `bbcp` - MCP (Model Context Protocol) server for MCP-compatible editors and clients.

The `bbmd` cli should be the superior to MCP when used by AI agents.

## Quick Install

Install the binaries:

```bash
npx --yes @atlacp/install
```

The installer places `bbmd` and `bbcp` in `~/.atlacp/bin` and updates your shell profile. Restart the shell, or source the updated profile before continuing. Follow [accounts setup](#account-setup) to configure Bitbucket account(s).

Use `bbmd skill` to teach your agent to use the cli. You can either write it to the agent relevant skill folder or add it somewhere in AGENTS.md.

```markdown
### Bitbucket integration

Run `bbmd skill` to learn how to work with bitbucket.
```

Or write the skill:
```bash
# Opencode
mkdir -p .opencode/skills/bitbucket && bbmd skill > .opencode/skills/bitbucket/SKILL.md

# Claude
mkdir -p .claude/skills/bitbucket && bbmd skill > .claude/skills/bitbucket/SKILL.md

# Codex
mkdir -p .codex/skills/bitbucket && bbmd skill > .codex/skills/bitbucket/SKILL.md

# Cursor
mkdir -p .cursor/skills/bitbucket && bbmd skill > .cursor/skills/bitbucket/SKILL.md
```
__I'm so tired to have 20 folders with AI stuff...__

## Account Setup

You can have multiple accounts configured. This may be useful if you want AI to use different account for some tasks, like posting PR reviews or similar. The account `name` is just a label.

Configure a Bitbucket account:

```bash
bbmd auth add \
  --name user \
  --default \
  --token-type Bearer|Basic \
  --token-value "<token>"
```

See [Bitbucket Access Tokens](#bitbucket-access-tokens) for more details on creating tokens.

Check the configured accounts:
```bash
bbmd auth status
```

## Supported Commands

### bbmd (CLI)

- `bbmd auth status` - list configured Atlassian accounts with redacted tokens.
- `bbmd auth add` - add or replace an account in the local accounts file.
- `bbmd auth remove` - remove an account by name.
- `bbmd auth set-default` - choose the default account.
- `bbmd pr create` - create a pull request.
- `bbmd pr read` - get pull request details.
- `bbmd pr update` - update pull request title, description, or draft state.
- `bbmd pr approve` - approve a pull request.
- `bbmd pr request-changes` - remove approval / request changes.
- `bbmd pr merge` - merge a pull request.
- `bbmd pr diffstat` - list changed files summary for a pull request.
- `bbmd pr diff` - get raw diff text.
- `bbmd pr add-comment` - add a general or inline pull request comment.
- `bbmd pr list-comments` - list pull request comments.
- `bbmd pr resolve-comment` - resolve a pull request comment thread.
- `bbmd pr list-tasks` - list pull request tasks.
- `bbmd pr create-task` - create a pull request task.
- `bbmd pr update-task` - update a pull request task.
- `bbmd file content` - get file content at a commit.

### bbcp (MCP tools)

- `bitbucket_add_pr_comment` - add a comment to a pull request.
- `bitbucket_approve_pr` - approve a pull request.
- `bitbucket_create_pr` - create a pull request.
- `bitbucket_create_pr_task` - create a task on a pull request.
- `bitbucket_get_file_content` - get file content from a repository at a specific commit.
- `bitbucket_get_pr_diff` - get the diff of a pull request.
- `bitbucket_get_pr_diffstat` - get the diffstat of a pull request.
- `bitbucket_list_pr_comments` - list pull request comments.
- `bitbucket_list_pr_tasks` - list tasks on a pull request.
- `bitbucket_merge_pr` - merge a pull request.
- `bitbucket_read_pr` - read a pull request.
- `bitbucket_resolve_pr_comment` - resolve a pull request comment thread.
- `bitbucket_request_pr_changes` - request changes on a pull request.
- `bitbucket_update_pr` - update a pull request.
- `bitbucket_update_pr_task` - update a task on a pull request.

## MCP Server

Supported transports:
- HTTP at `http://localhost:8080`.
- SSE at `http://localhost:8080/sse`.
- STDIO via `bbcp stdio`.

To Run HTTP/SSE server:
```bash
bbcp http
```

Or just use STDIO:
```bash
bbcp stdio
```

### Integrate AI Tools

Cursor MCP config (`.cursor/mcp.json`) may look like this for STDIO transport:
```json
{
  "mcpServers": {
    "Bitbucket": {
      "command": "bbcp",
      "args": ["stdio"]
    }
  }
}
```

Or if you prefer HTTP transport:
```json
{
  "mcpServers": {
    "Bitbucket": {
      "url": "http://localhost:8080"
    }
  }
}
```

Once configured, send a prompt similar to:

```text
Check titbucket pull request https://bitbucket.org/workspace/repo-slug/pull-requests/123
```

You should see a response with PR details.

## Bitbucket Access Tokens

Please review the official documentation:
- [Personal API Tokens](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/#Create-an-API-token) 
 (keep in mind to create a basic token for API use). When using personal access tokens, all requests will be made on behalf of the user.
- [Bitbucket Access Tokens](https://support.atlassian.com/bitbucket-cloud/docs/access-tokens/) - good for bots and other automation tools.

### User API Token

* Go to https://id.atlassian.com/manage-profile/security/api-tokens
* Create a new token, with at least below scopes:
  ```text
  read:account
  read:issue:bitbucket
  read:me
  read:pipeline:bitbucket
  read:project:bitbucket
  read:pullrequest:bitbucket
  read:repository:bitbucket
  read:runner:bitbucket
  read:snippet:bitbucket
  read:user:bitbucket
  write:issue:bitbucket
  write:pullrequest:bitbucket
  ```
* Create a Basic token from it using shell command below:
  ```bash
  echo -n "<your-email>:<your-api-token>" | base64
  ```
* Add user account to the local accounts file:
  ```bash
  bbmd auth add \
    --name user \
    --default \
    --token-type Basic \
    --token-value "<base64-email-colon-api-token>"
  ```

### Bot API Token

If you plan your AI agent to operate as a bot, best option is to create "Access token" instead of "User API token".

* Go your repository or workspace settings, click on "Access tokens"
* Create a new token with below permissions:
  ```text
  pullrequest
  pipeline
  repository:write
  repository
  pullrequest:write
  ```
* Add bot account to the local accounts file:
  ```bash
  bbmd auth add \
    --name bot \
    --token-type Bearer \
    --token-value "<your-bot-api-token>"
  ```
  Note: You may use any name in place of `bot` in above command. It's just a label.

## Contribute

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for more details.

## License

MIT
