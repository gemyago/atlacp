# Testing documentation

This folder holds **manual / AI-driven integration test playbooks** for Bitbucket-related flows.

| Document | Purpose |
|----------|---------|
| [bitbucket-mcp-integration-tests.md](./bitbucket-mcp-integration-tests.md) | End-to-end checks using the Atlassian MCP tools (`mcp.bitbucket_*`). |
| [bitbucket-cli-integration-tests.md](./bitbucket-cli-integration-tests.md) | Same scenarios via `go run ./cmd/bbmd` from the atlacp repo root. |

Example TypeScript fixtures used by the review test live here: `example1.ts`, `example2.ts`.

## Checking Atlassian auth (`bbmd auth`)

Before running CLI integration tests, confirm the accounts file is loaded and accounts look right. From the **atlacp project root**:

```bash
go run ./cmd/bbmd auth status
```

Use **`auth status`** to list configured accounts (tokens appear partially redacted) and see which account is the default. If you rely on a non-default accounts file path, pass the same **`--atlassian-accounts-file`** you will use for `pr` / `file` commands.

The Bitbucket playbooks assume at least two accounts (typically **`user`** and **`bot`**). If **`auth status`** does not show the expected names, fix the file or use **`bbmd auth add`**, **`bbmd auth remove`**, and **`bbmd auth set-default`** before running tests.

## SSH troubleshooting

Use this when `git push` or `git pull` against Bitbucket fails with permission or wrong-key errors.

```bash
# See which ssh key is used
ssh -T git@bitbucket.org
```

If the wrong key is used, define a section for Bitbucket explicitly in `~/.ssh/config`, for example: `Host bitbucket.org` with `IdentityFile` and related settings.

If you use a wildcard SSH `Host *` entry, exclude Bitbucket so a dedicated host block applies, for example: `Host * !bitbucket.org` before or alongside your Bitbucket-specific `Host bitbucket.org` block.
