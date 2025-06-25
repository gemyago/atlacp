# Atlassian MCP Integration Examples

This directory contains examples and test scripts for the Atlassian MCP Integration.

## Quick Start

### Prerequisites

1. **Docker** - Make sure Docker is installed and running
2. **Node.js** - Version 18 or higher for test scripts
3. **Atlassian Accounts** - Configure your accounts in `atlassian-accounts-stub.json`

### Configuration

1. **Prepare Atlassian Accounts**: Use `atlassian-accounts-stub.json` as a template and create `atlassian-accounts.json` with your actual credentials:

```json
{
  "accounts": [
    {
      "name": "user",
      "default": true,
      "bitbucket": {
        "token": "YOUR_BITBUCKET_TOKEN",
        "workspace": "YOUR_WORKSPACE"
      },
      "jira": {
        "token": "YOUR_JIRA_TOKEN",
        "domain": "YOUR_DOMAIN"
      }
    },
    {
      "name": "bot",
      "default": false,
      "bitbucket": {
        "token": "YOUR_BOT_TOKEN",
        "workspace": "YOUR_WORKSPACE"
      },
      "jira": {
        "token": "YOUR_BOT_TOKEN",
        "domain": "YOUR_DOMAIN"
      }
    }
  ]
}
```

Do basic check by sending a `list` request to the server:
```bash
curl -X POST http://localhost:8080/mcp/list -H "Content-Type: application/json" -d '{"jsonrpc": "2.0", "method": "mcp.list", "id": 1}'
```

## Use Docker Compose to run MCP Server

### HTTP Server

Run the HTTP server for testing HTTP transport:

```bash
# Start HTTP server
docker-compose up atlacp-http
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

### Cursor Integration

1. **Install MCP Server**: Add the MCP server to Cursor's configuration
2. **Configure Accounts**: Set up your Atlassian accounts
3. **Use Tools**: Access Bitbucket tools directly from Cursor

### Claude Desktop Integration

1. **Add Server**: Configure the MCP server in Claude Desktop
2. **Test Connection**: Verify tools are available
3. **Start Using**: Begin using Bitbucket tools in conversations

## Next Steps

- [ ] Test with real Bitbucket repositories
- [ ] Integrate with CI/CD pipelines
- [ ] Add Jira integration testing
- [ ] Create more comprehensive test scenarios 