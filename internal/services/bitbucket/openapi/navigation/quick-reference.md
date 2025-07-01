# Bitbucket OpenAPI Quick Reference

## Key Line Numbers
- **Paths section starts**: Line 24
- **Definitions section starts**: Line ~24531
- **PullRequest model**: Lines 30032-30500

## Common Endpoints
- List PRs: `/repositories/{workspace}/{repo_slug}/pullrequests` (Line ~13012)
- Get PR: `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}` (Line ~13231)
- Create PR: `/repositories/{workspace}/{repo_slug}/pullrequests` POST (Line ~13079)
- Update PR: `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}` PUT (Line ~13285)
- Approve PR: `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/approve` (Line ~13451)
- Merge PR: `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/merge` (Line ~13511)

## Model Dependencies
- PullRequest → Account, Repository, Commit
- Repository → Workspace, Account

## Search Patterns
- Find endpoints: `grep -n "pullrequests" openapi.json`
- Find models: `grep -n '"pullrequest":' openapi.json`
- Find properties: `grep -A5 -B5 '"title":' openapi.json`

## Navigation Commands
```bash
# Jump to PullRequest model definition
read_file target_file=openapi.json offset=30032 limit=500

# Jump to List PRs endpoint
read_file target_file=openapi.json offset=13012 limit=60

# Jump to Get PR endpoint
read_file target_file=openapi.json offset=13231 limit=50

# Jump to Approve PR endpoint
read_file target_file=openapi.json offset=13451 limit=50
``` 