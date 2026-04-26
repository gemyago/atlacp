# Atlassian MCP Quick Start

This directory contains setup examples for connecting `bbcp` to MCP-compatible clients.

## Prepare accounts file

Create `atlassian-accounts.json` using `atlassian-accounts-stub.json` as a template. Follow next steps to set your Atlassian tokens.

More on Atlassian tokens:
- [Personal API Tokens](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/#Create-an-API-token) 
 (keep in mind to create a basic token for API use). When using personal access tokens, all requests will be made on behalf of the user.
- [Bitbucket Access Tokens](https://support.atlassian.com/bitbucket-cloud/docs/access-tokens/) - good for bots and other automation tools.

## Adding Atlassian Tokens

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
* Copy the token value and use it in the `atlassian-accounts.json` file as a user account.

### Bot API Token
Skip this step if you don't plan to perform any actions on behalf of the bot.

* Go your repository or workspace settings, click on "Access tokens"
* Create a new token with below permissions:
  ```text
  pullrequest
  pipeline
  repository:write
  repository
  pullrequest:write
  ```
* Copy the token value and use it in the `atlassian-accounts.json` file as a bot account.

## Run MCP Server

If you installed atlacp locally, start the MCP server directly:

```bash
bbcp http
```

Or use STDIO transport:

```bash
bbcp stdio
```

Docker is optional. If you prefer containers, published images follow the per-binary naming pattern, so the MCP server image is `ghcr.io/gemyago/atlacp-bbcp`.

**Note**
Due to ghcr constraints, if you are logged in to ghcr, you may have to run `docker logout ghcr.io` to avoid authentication errors when pulling public images.

### Run using Docker directly

```bash
docker run -d --name atlacp-mcp \
  --restart=always \
  -p 8080:8080 \
  -v $(pwd)/atlassian-accounts.json:/app/atlassian-accounts.json \
  ghcr.io/gemyago/atlacp-bbcp:latest \
  -a /app/atlassian-accounts.json \
  http
```

### Run using Docker Compose

Use example `docker-compose.yml` file to run the MCP server:
```yaml
services:
  # MCP Server for testing HTTP transport
  atlacp-http:
    image: ghcr.io/gemyago/atlacp-bbcp:latest
    command:
      - http
      - --atlassian-accounts-file=/app/config/atlassian-accounts.json
      - --log-level=info
    volumes:
      - ./atlassian-accounts.json:/app/config/atlassian-accounts.json:ro
    ports:
      - "8080:8080"
```

Run the MCP server with HTTP transport:

```bash
docker compose up atlacp-http
```

**Note**
Due to ghcr constraints, if you are logged in to ghcr, you may have to run `docker logout ghcr.io` to avoid authentication errors when pulling public images.

## Integration with AI Code Editors

Example configuration for Cursor (.cursor/mcp.json):

```json
{
  "mcpServers": {
    "Atlassian MCP": {
      "url": "http://localhost:8080"
    }
  }
}
```

Use `http://localhost:8080/sse` - for SSE transport. Otherwise Streamable HTTP transport is used.
