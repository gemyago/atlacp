# atlacp

[![Build](https://github.com/gemyago/atlacp/actions/workflows/build-flow.yml/badge.svg)](https://github.com/gemyago/atlacp/actions/workflows/build-flow.yml)
[![Coverage](https://raw.githubusercontent.com/gemyago/atlacp/test-artifacts/coverage/golang-coverage.svg)](https://htmlpreview.github.io/?https://raw.githubusercontent.com/gemyago/atlacp/test-artifacts/coverage/golang-coverage.html)

An MCP (Model Context Protocol) interface for Atlassian products (Jira and Bitbucket).

## Available tools

- `bitbucket_create_pr` - create a pull request
- `bitbucket_read_pr` - read a pull request
- `bitbucket_update_pr` - update a pull request
- `bitbucket_approve_pr` - approve a pull request
- `bitbucket_merge_pr` - merge a pull request
- `bitbucket_list_pr_tasks` - list tasks on a pull request
- `bitbucket_update_pr_task` - update a task on a pull request
- `bitbucket_create_pr_task` - create a task on a pull request

## Quick Start

### Configure accounts

The tool is designed to be running locally on developer's machine. In order to run the tool you need to configure your Atlassian accounts first. For bitbucket you need to create a personal access token that can read and write to the repository.

Example `accounts-config.json` file:
```json
{
  "accounts": [
    {
      "name": "user",
      "default": true,
      "bitbucket": { "type": "Basic", "value": "ATBBxxxxxxxxxxxxxxxx" }
    }
  ]
} 
```

You may optionally configure multiple accounts for different roles or different workspaces, for example you may have `user` and `bot` accounts. See `quick-start/atlassian-accounts-stub.json` for more details.

More on Atlassian tokens:
- [Personal API Tokens](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/#Create-an-API-token) 
 (keep in mind to create a basic token for API use). When using personal access tokens, all requests will be made on behalf of the user.
- [Bitbucket Access Tokens](https://support.atlassian.com/bitbucket-cloud/docs/access-tokens/) - good for bots and other automation tools.

### Start the MCP server

Start docker container pointing on the `accounts-config.json` file:

```bash
docker run -d --name atlacp-mcp \
  --restart=always \
  -p 8080:8080 \
  -v $(pwd)/accounts-config.json:/app/accounts-config.json \
  ghcr.io/gemyago/atlacp-mcp:latest
```

### Integrate with AI editors

Cursor MCP config (.cursor/mcp.json) section may look like this:

```json
{
  "mcpServers": {
    "Atlassian MCP": {
      "url": "http://localhost:8080"
    }
  }
}
```

Once configured, try to send a prompt to the editor, similar to below:
```text
Check pull request 12345 from bitbucket repo owner/repo-name
```

You should see a response with PR details.

## License

MIT