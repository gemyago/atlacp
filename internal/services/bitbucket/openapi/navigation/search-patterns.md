# Bitbucket OpenAPI Search Patterns

## Endpoint Searches

### Pull Request Endpoints
```bash
# Find all pull request related endpoints
grep_search query="/repositories/.*/pullrequests" include_pattern="openapi.json"

# Find specific pull request operations
grep_search query="approve" include_pattern="openapi.json" 
grep_search query="merge" include_pattern="openapi.json"
grep_search query="decline" include_pattern="openapi.json"
```

## Model Searches

### Pull Request Models
```bash
# Find the pull request model definition
grep_search query="\"pullrequest\":" include_pattern="openapi.json"

# Find pull request properties
grep_search query="properties.*title" include_pattern="openapi.json"
grep_search query="properties.*state" include_pattern="openapi.json"
grep_search query="properties.*source" include_pattern="openapi.json"
```

## Related Models

### Repository
```bash
# Find repository model definition
grep_search query="\"repository\":" include_pattern="openapi.json"

# Find repository references in pull requests
grep_search query="repository.*pullrequest" include_pattern="openapi.json"
```

### Commit
```bash
# Find commit model definition
grep_search query="\"commit\":" include_pattern="openapi.json"

# Find commit references in pull requests
grep_search query="commit.*pullrequest" include_pattern="openapi.json"
```

## Parameter Searches
```bash
# Find required parameters for pull request operations
grep_search query="parameters.*required.*true" include_pattern="openapi.json"

# Find workspace parameter definitions
grep_search query="workspace.*parameter" include_pattern="openapi.json"

# Find pull request ID parameter definitions
grep_search query="pull_request_id" include_pattern="openapi.json"
``` 