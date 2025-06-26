# Bitbucket OpenAPI Navigation Tools

This directory contains tools and resources for efficiently working with the large Bitbucket OpenAPI specification (33K+ lines). These navigation aids help make the massive specification file more accessible and easier to work with.

## Directory Structure

```
openapi/
├── navigation/             # Quick reference and navigation aids
│   ├── endpoints-index.json # Index of key API endpoints with line numbers
│   ├── models-index.json   # Index of key data models with line numbers
│   ├── quick-reference.md  # Human-readable guide to key sections
│   └── search-patterns.md  # Common search patterns for finding endpoints/models
├── subsets/                # Focused subsets of the OpenAPI spec
│   ├── pullrequest-model.json      # Just the PullRequest model
│   └── pullrequest-endpoints.json  # Just the PullRequest endpoints
├── analysis/               # Analysis and documentation
│   ├── model-differences.md        # Differences between OpenAPI and implementation
│   └── implementation-status.md    # Status of API implementation
├── verify.sh               # Script to verify navigation infrastructure
└── README.md               # This file
```

## Working with the Navigation Tools

### Using the Quick Reference

The `navigation/quick-reference.md` file contains line numbers for key sections of the OpenAPI spec, such as:

- Line numbers for important model definitions
- Line numbers for common API endpoints
- Model dependency information
- Common search patterns

### Navigating with the Indexes

The JSON index files make it easy to programmatically find parts of the spec:

```bash
# Find all pull request endpoints
jq '.pullrequests.paths' navigation/endpoints-index.json

# Get line numbers for specific endpoints
jq '.pullrequests.line_ranges.get' navigation/endpoints-index.json

# Get properties of the pull request model
jq '.pullrequest.properties' navigation/models-index.json
```

### Using the Search Patterns

The `search-patterns.md` file contains common grep patterns for finding information in the OpenAPI spec:

```bash
# Find pull request endpoints
grep -n "pullrequests" ../openapi.json

# Find pull request model definition
grep -n '"pullrequest":' ../openapi.json
```

### Working with Subsets

Instead of loading the entire 33K line file, work with focused subsets:

- `pullrequest-model.json` - Just the PullRequest model definition
- `pullrequest-endpoints.json` - Just the PullRequest endpoints

## Implementation Status

Check `analysis/implementation-status.md` for the current status of API implementation, showing which endpoints are completed and which are pending.

## Model Differences

The `analysis/model-differences.md` file tracks differences between the OpenAPI spec models and our implementation, including field name mismatches and parameter differences.

## Verification

Run the verification script to ensure all navigation tools are working:

```bash
./verify.sh
```

## Development Workflow

1. **Start with the quick reference** to find the section you need
2. **Use the indexes** to get precise line numbers for navigation
3. **Work with subsets** rather than the full file when possible
4. **Update the status tracking** as you implement new endpoints
5. **Document differences** between the spec and your implementation

Following this workflow will make working with the large OpenAPI spec significantly easier and more efficient. 