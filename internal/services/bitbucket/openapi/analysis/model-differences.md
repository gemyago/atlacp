# Model Differences Analysis

## PullRequest Model
| OpenAPI Field | Implemented Field | Status | Notes |
|---------------|-------------------|--------|-------|
| `id` | `ID` | Match | Integer type in both |
| `title` | `Title` | Match | String type in both |
| `description` | `Description` | Match | String type in both |
| `state` | `State` | Match | String type in both |
| `author` | `Author` | Match | Object type in both |
| `source` | `Source` | Match | Object type in both |
| `destination` | `Destination` | Match | Object type in both |
| `merge_commit` | `MergeCommit` | Match | Object type in both |
| `draft` | `Draft` | Match | Boolean type in both |

## Path Parameter Differences
| Implementation | OpenAPI Spec | Status |
|----------------|--------------|--------|
| `{username}` | `{workspace}` | Mismatch | Need to standardize on `workspace` |
| `{repo_slug}` | `{repo_slug}` | Match | |
| `{pr_id}` | `{pull_request_id}` | Mismatch | Need to standardize on `pull_request_id` |

## Important Notes
1. The Bitbucket API uses `workspace` as the parameter name for what is sometimes referred to as `username` in our code.
2. The PullRequest ID in the OpenAPI spec is named `pull_request_id` but our code may use `pr_id` in places.
3. The state field in the PullRequest model is an enumeration with valid values: "OPEN", "MERGED", "DECLINED", "SUPERSEDED". 