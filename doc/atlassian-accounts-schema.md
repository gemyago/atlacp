# Atlassian Accounts Schema Design

## Overview

This document describes the JSON schema for the Atlassian accounts configuration file. This file supports configuring multiple named Atlassian accounts with separate credentials for Bitbucket and Jira services, along with specification of a default account.

## Schema Definition

The accounts configuration file uses the following JSON schema:

```json
{
  "accounts": [
    {
      "name": "string",
      "default": boolean,
      "bitbucket": {
        "token": "string",
        "workspace": "string"
      },
      "jira": {
        "token": "string",
        "domain": "string"
      }
    }
  ]
}
```

### Field Descriptions

#### Root Object

| Field | Type | Description |
|-------|------|-------------|
| `accounts` | Array | Array of account configuration objects. At least one account must be defined. |

#### Account Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | String | Yes | Friendly name of the account. Must be unique across all accounts. |
| `default` | Boolean | No (default: false) | Specifies if this is the default account. Only one account can be marked as default. |
| `bitbucket` | Object | No | Bitbucket-specific configuration. Required only if using Bitbucket services. |
| `jira` | Object | No | Jira-specific configuration. Required only if using Jira services. |

#### Bitbucket Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `token` | String | Yes | API token for Bitbucket authentication. |
| `workspace` | String | Yes | Bitbucket workspace or username for this account. Used to construct API URLs. |

#### Jira Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `token` | String | Yes | API token for Jira authentication. |
| `domain` | String | Yes | Jira cloud instance domain (e.g., "mycompany" for mycompany.atlassian.net). Used to construct API URLs. |

## Validation Rules

The following validation rules apply to the accounts configuration:

1. At least one account must be defined.
2. Account names must be unique across all accounts.
3. Exactly one account must be marked as default.
4. At least one of Bitbucket or Jira configuration must be present per account.
5. If a service configuration is present, all its required fields must be non-empty.

## Example Configuration

```json
{
  "accounts": [
    {
      "name": "default-user",
      "default": true,
      "bitbucket": {
        "token": "ATBBxxxxxxxxxxxxxxxx",
        "workspace": "my-workspace"
      },
      "jira": {
        "token": "ATATxxxxxxxxxxxxxxxx",
        "domain": "mycompany"
      }
    },
    {
      "name": "bitbucket-only",
      "default": false,
      "bitbucket": {
        "token": "ATBBxxxxxxxxxxxxxxxx",
        "workspace": "my-workspace"
      }
    },
    {
      "name": "jira-only",
      "default": false,
      "jira": {
        "token": "ATATxxxxxxxxxxxxxxxx",
        "domain": "mycompany"
      }
    }
  ]
}
```

## File Location

The accounts configuration file will be searched for in the following locations (in order):

1. Path specified by the `config.atlassian.accountsFilePath` configuration value
2. `$HOME/.config/atlacp/accounts.json`
3. `./accounts.json` (current working directory)

## Security Considerations

The accounts configuration file contains sensitive API tokens. It is recommended to:

1. Store the file outside of the repository
2. Set appropriate file permissions (e.g., `chmod 600`)
3. Consider using environment variables or a secrets manager for production environments 