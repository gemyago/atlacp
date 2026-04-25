# atlacp

[![Build](https://github.com/gemyago/atlacp/actions/workflows/build-flow.yml/badge.svg)](https://github.com/gemyago/atlacp/actions/workflows/build-flow.yml)
[![Coverage](https://raw.githubusercontent.com/gemyago/atlacp/test-artifacts/coverage/golang-coverage.svg)](https://htmlpreview.github.io/?https://raw.githubusercontent.com/gemyago/atlacp/test-artifacts/coverage/golang-coverage.html)

Bitbucket tooling for developers and AI agents.

The project provides two entrypoints:

- `bbmd` - direct Bitbucket CLI for humans, scripts, and agent workflows.
- `bbcp` - MCP (Model Context Protocol) server for MCP-compatible editors and clients.

Most users should start with `bbmd`. It does not require a long-running service and prints JSON for easy piping into scripts or agents.

## Quick Install

Install the binaries:

```bash
npm install -g @atlacp/install
```

The installer places `bbmd` and `bbcp` in `~/.atlacp/bin` and updates your shell profile. Restart the shell, or source the updated profile before continuing.

## Account Setup

Configure a Bitbucket account:

```bash
bbmd auth add \
  --name user \
  --default \
  --token-type Basic \
  --token-value "<base64-email-colon-api-token>"
```

Check the configured accounts:

```bash
bbmd auth status
```

Read a pull request:

```bash
bbmd pr read \
  --repo-owner <workspace> \
  --repo-name <repo-slug> \
  --pr-id 123
```

Get the pull request diff:

```bash
bbmd pr diff \
  --repo-owner <workspace> \
  --repo-name <repo-slug> \
  --pr-id 123
```

Post a pull request comment:

```bash
bbmd pr add-comment \
  --repo-owner <workspace> \
  --repo-name <repo-slug> \
  --pr-id 123 \
  --content "Looks good to me"
```

Get file content at a commit:

```bash
bbmd file content \
  --repo-owner <workspace> \
  --repo-name <repo-slug> \
  --commit <commit-sha> \
  --path path/to/file.go
```

Use `--account <name>` on `pr` and `file` commands when you need a non-default account.

## Supported Commands

### bbmd (CLI)

- `bbmd auth status` - list configured Atlassian accounts with redacted tokens.
- `bbmd auth add` - add or replace an account in the local accounts file.
- `bbmd auth remove` - remove an account by name.
- `bbmd auth set-default` - choose the default account.
- `bbmd pr create` - create a pull request.
- `bbmd pr read` - read pull request details.
- `bbmd pr update` - update pull request title, description, or draft state.
- `bbmd pr approve` - approve a pull request.
- `bbmd pr request-changes` - remove approval / request changes.
- `bbmd pr merge` - merge a pull request.
- `bbmd pr diffstat` - list changed files summary.
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
- `bitbucket_get_file_content` - get the content of a file in a pull request.
- `bitbucket_get_pr_diff` - get the diff of a pull request.
- `bitbucket_get_pr_diffstat` - get the diffstat of a pull request.
- `bitbucket_list_pr_tasks` - list tasks on a pull request.
- `bitbucket_merge_pr` - merge a pull request.
- `bitbucket_read_pr` - read a pull request.
- `bitbucket_request_pr_changes` - request changes on a pull request.
- `bitbucket_update_pr` - update a pull request.
- `bitbucket_update_pr_task` - update a task on a pull request.

### bbcp transport endpoints

`bbmd auth add` writes the default accounts file to `~/.atlacp/accounts.json`. `bbmd` and `bbcp` both use this file by default.

You can also provide a specific accounts file with `--atlassian-accounts-file` or `-a`:

```bash
bbmd -a ./atlassian-accounts.json auth status
bbcp -a ./atlassian-accounts.json http
```

Example `atlassian-accounts.json` file:

```json
{
  "accounts": [
    {
      "name": "user",
      "default": true,
      "bitbucket": { "type": "Basic", "value": "<base64-email-colon-api-token>" }
    }
  ]
}
```

You may configure multiple accounts for different roles or workspaces, for example `user` and `bot` accounts. See `quick-start/atlassian-accounts-stub.json` for a fuller template.

More on Atlassian tokens:

- [Personal API Tokens](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/#Create-an-API-token) - use a `Basic` token value created from `email:api-token`, for example `printf '%s' '<email>:<api-token>' | base64`.
- [Bitbucket Access Tokens](https://support.atlassian.com/bitbucket-cloud/docs/access-tokens/) - useful for bots and automation; commonly used with `--token-type Bearer`.

Streamable HTTP at `http://localhost:8080`.

SSE at `http://localhost:8080/sse`.

STDIO via `bbcp stdio`.

## MCP Server

Use the MCP server when you need to connect atlacp to Cursor or another MCP-compatible client. The detailed Docker-based walkthrough is in [quick-start](./quick-start).

### Run With Docker

Start a Docker container pointing to an accounts file:

```bash
docker run -d --name atlacp-mcp \
  --restart=always \
  -p 8080:8080 \
  -v $(pwd)/atlassian-accounts.json:/app/atlassian-accounts.json \
  ghcr.io/gemyago/atlacp-mcp:latest \
  -a /app/atlassian-accounts.json \
  http
```

Root URL serves Streamable HTTP. Append `/sse` for SSE:

- `http://localhost:8080` - Streamable HTTP transport.
- `http://localhost:8080/sse` - SSE transport.

STDIO is supported as well. Use the `stdio` subcommand instead of `http`.

### Run With Installed Binary

If you installed the binaries locally, you can run the HTTP MCP server without Docker:

```bash
bbcp http
```

Or use STDIO transport:

```bash
bbcp stdio
```

### Integrate AI Tools

Cursor MCP config (`.cursor/mcp.json`) may look like this:

```json
{
  "mcpServers": {
    "Atlassian MCP": {
      "url": "http://localhost:8080"
    }
  }
}
```

Once configured, send a prompt similar to:

```text
Check pull request 123 from Bitbucket repo workspace/repo-slug
```

You should see a response with PR details.

## Local Development

Run commands from source:

```bash
go run ./cmd/bbmd --help
go run ./cmd/bbcp http --env local --noop
```

Build local binaries:

```bash
make dist/bin
```

## License

MIT
