# Working with Large OpenAPI Files

## Overview

This document provides comprehensive instructions for efficiently working with large OpenAPI specification files (30K+ lines) when generating HTTP API clients. It provides concrete strategies and tools optimized for AI model navigation and analysis.

## Problem Statement

Large OpenAPI files (like Bitbucket's 33K-line specification) present significant challenges for AI models:
- **Navigation difficulty**: Finding specific endpoints or models in massive files
- **Context limitations**: AI models have token limits that prevent processing entire files
- **Analysis complexity**: Understanding relationships between models and endpoints
- **Incremental development**: Need to work on specific features without full context

## Solution Strategy

### 1. Create Navigation Infrastructure

**Goal**: Transform large files into navigable, searchable resources with multiple entry points.

#### A. Directory Structure
```
internal/services/{service}/
├── openapi/
│   ├── full.json (original large file)
│   ├── navigation/
│   │   ├── endpoints-index.json
│   │   ├── models-index.json
│   │   ├── quick-reference.md
│   │   └── search-patterns.md
│   ├── subsets/
│   │   ├── {feature}-endpoints.json
│   │   ├── {feature}-models.json
│   │   └── core-models.json
│   └── analysis/
│       ├── model-differences.md
│       └── implementation-status.md
```

#### B. Navigation Index Files

**endpoints-index.json** - Quick reference for all endpoints:
```json
{
  "pullrequests": {
    "paths": [
      "/repositories/{workspace}/{repo_slug}/pullrequests",
      "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}",
      "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/approve",
      "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/merge"
    ],
    "line_ranges": {
      "list": [13012, 13068],
      "get": [13231, 13285],
      "approve": [13451, 13500],
      "merge": [13511, 13554]
    },
    "tags": ["Pullrequests"],
    "methods": ["GET", "POST", "PUT", "DELETE"]
  }
}
```

**models-index.json** - Quick reference for all models:
```json
{
  "pullrequest": {
    "line_start": 30032,
    "line_end": 30500,
    "properties": ["id", "title", "state", "author", "source", "destination"],
    "dependencies": ["account", "repository", "commit"],
    "required_fields": ["id", "title", "state"],
    "optional_fields": ["description", "reviewers", "draft"]
  }
}
```

#### C. Quick Reference Guide

**quick-reference.md** - Human-readable navigation guide:
```markdown
# {Service} OpenAPI Quick Reference

## Key Line Numbers
- **Paths section starts**: Line 24
- **Definitions section starts**: Line 24531
- **PullRequest model**: Lines 30032-30500
- **Account model**: Lines 24550-24600

## Common Endpoints
- List PRs: `/repositories/{workspace}/{repo_slug}/pullrequests` (Line ~13012)
- Get PR: `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}` (Line ~13231)
- Approve PR: `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/approve` (Line ~13451)

## Model Dependencies
- PullRequest → Account, Repository, Commit
- Repository → Workspace, Account
- Commit → Account

## Search Patterns
- Find endpoints: `grep -n "pullrequests" openapi.json`
- Find models: `grep -n '"pullrequest":' openapi.json`
- Find properties: `grep -A5 -B5 '"title":' openapi.json`
```

### 2. Create Focused Subset Files

**Goal**: Break down large files into manageable, feature-specific subsets.

#### A. Feature-Based Subsets
```bash
# Extract pull request related content
jq '{
  info: .info,
  paths: (.paths | with_entries(select(.key | contains("pullrequests")))),
  definitions: (.definitions | with_entries(select(.key | contains("pullrequest") or .key | contains("paginated"))))
}' openapi.json > subsets/pullrequests.json

# Extract repository related content
jq '{
  info: .info,
  paths: (.paths | with_entries(select(.key | contains("repositories")))),
  definitions: (.definitions | with_entries(select(.key | contains("repository") or .key | contains("workspace"))))
}' openapi.json > subsets/repositories.json
```

#### B. Core Models Subset
```bash
# Extract essential models only
jq '.definitions | {pullrequest, account, repository, workspace, commit}' openapi.json > subsets/core-models.json
```

### 3. AI-Optimized Navigation Commands

**Goal**: Provide specific commands and patterns for AI model navigation.

#### A. Line-Based Navigation
```bash
# Jump to specific sections using line numbers
read_file target_file=openapi.json start_line_one_indexed=30032 end_line_one_indexed_inclusive=30100

# Jump to endpoint definitions
read_file target_file=openapi.json start_line_one_indexed=13012 end_line_one_indexed_inclusive=13068
```

#### B. Pattern-Based Search
```bash
# Find all pull request related content
grep_search query="pullrequest" include_pattern=*.json

# Find specific model properties
grep_search query='"title":' include_pattern=*.json

# Find endpoint definitions
grep_search query='"paths":' include_pattern=*.json
```

#### C. Content Extraction
```bash
# Extract specific sections by line numbers
sed -n "30032,30100p" openapi.json > pullrequest-model.json

# Extract all endpoints for a specific tag
jq '.paths | to_entries | map(select(.value | to_entries[0].value.tags[0] == "Pullrequests"))' openapi.json
```

### 4. Implementation Workflow

**Goal**: Provide step-by-step process for working with large OpenAPI files.

#### Step 1: Initial Analysis
```bash
# Create navigation infrastructure
mkdir -p openapi/navigation openapi/subsets openapi/analysis

# Generate endpoints index
jq '.paths | to_entries | map({path: .key, methods: (.value | keys), tags: (.value | to_entries[0].value.tags[0])}) | group_by(.tags)' openapi.json > navigation/endpoints-index.json

# Generate models index
jq '.definitions | to_entries | map({name: .key, properties: (.value.properties | keys), line_start: "TBD"})' openapi.json > navigation/models-index.json
```

#### Step 2: Create Subsets
```bash
# Create feature-specific subsets
jq '{
  info: .info,
  paths: (.paths | with_entries(select(.key | contains("pullrequests")))),
  definitions: (.definitions | with_entries(select(.key | contains("pullrequest") or .key | contains("paginated"))))
}' openapi.json > subsets/pullrequests.json
```

#### Step 3: Generate Quick Reference
```bash
# Create human-readable navigation guide
cat > navigation/quick-reference.md << 'EOF'
# {Service} OpenAPI Quick Reference

## Key Sections
- Paths: Lines 24-24530
- Definitions: Lines 24531-33021

## Common Models
- PullRequest: Lines 30032-30500
- Account: Lines 24550-24600
- Repository: Lines 24600-24700

## Common Endpoints
- List PRs: Lines 13012-13068
- Get PR: Lines 13231-13285
- Approve PR: Lines 13451-13500
EOF
```

### 5. Analysis and Documentation

**Goal**: Maintain comprehensive documentation of findings and decisions.

#### A. Model Differences Analysis
Create `analysis/model-differences.md` to track differences between OpenAPI spec and implementation:
```markdown
# Model Differences Analysis

## PullRequest Model
| OpenAPI Field | Implemented Field | Status | Notes |
|---------------|-------------------|--------|-------|
| `id` | `ID` | Match | Integer type in both |
| `title` | `Title` | Match | String type in both |
| `description` | - | Missing | Not in OpenAPI spec |
| `draft` | `Draft` | Match | Boolean type in both |

## Path Parameter Differences
| Implementation | OpenAPI Spec | Status |
|----------------|--------------|--------|
| `{username}` | `{workspace}` | Mismatch |
| `{repo_slug}` | `{repo_slug}` | Match |
```

#### B. Implementation Status Tracking
Create `analysis/implementation-status.md`:
```markdown
# Implementation Status

## Completed Endpoints
- [x] GET /repositories/{workspace}/{repo_slug}/pullrequests
- [x] POST /repositories/{workspace}/{repo_slug}/pullrequests
- [x] GET /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}

## Pending Endpoints
- [ ] PUT /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}
- [ ] POST /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/approve
- [ ] POST /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/merge

## Models Status
- [x] PullRequest (partial - missing some fields)
- [x] Account (simplified)
- [ ] Repository (not implemented)
```

### 6. Quality Assurance

#### A. Validation Scripts
```bash
# Validate subset files are valid JSON
jq empty subsets/pullrequests.json && echo "Pullrequests subset is valid JSON"

# Validate navigation indexes
jq '.pullrequests.paths | length' navigation/endpoints-index.json

# Check for missing references
jq '.definitions.pullrequest.properties | keys' subsets/pullrequests.json
```

#### B. Consistency Checks
```bash
# Ensure all referenced models exist in subset
jq '.definitions | keys' subsets/pullrequests.json

# Verify endpoint paths match OpenAPI spec
jq '.paths | keys' subsets/pullrequests.json
```

### 7. Integration with Client Generation

#### A. Use Subsets for Code Generation
```bash
# Generate client from subset instead of full spec
openapi-generator-cli generate -i subsets/pullrequests.json -g go -o ./generated
```

#### B. Incremental Implementation
1. **Start with core models**: Implement essential models first
2. **Add endpoints incrementally**: Implement one endpoint group at a time
3. **Test each addition**: Ensure new code works before adding more
4. **Update documentation**: Keep navigation files current

### 8. Maintenance Guidelines

#### A. When to Update Navigation Files
- After adding new endpoints to the client
- When discovering new model relationships
- When updating OpenAPI specification
- Before major releases

#### B. Keeping Subsets Current
```bash
# Regenerate subsets when OpenAPI spec changes
./scripts/regenerate-subsets.sh

# Update line numbers in quick reference
./scripts/update-line-numbers.sh
```

## Benefits

1. **Efficient Navigation**: AI models can quickly find relevant sections
2. **Focused Analysis**: Work with manageable file sizes
3. **Incremental Development**: Build features without full context
4. **Better Documentation**: Clear understanding of what's implemented
5. **Easier Maintenance**: Track changes and dependencies
6. **Faster Development**: Reduced time spent searching large files

## Conclusion

This approach transforms large, unwieldy OpenAPI files into **navigable, searchable, and manageable** resources that AI models can efficiently work with. The key is creating **multiple entry points** and **focused subsets** rather than trying to process everything at once.

By following these patterns, AI models can:
- **Navigate quickly** to relevant sections using line numbers
- **Work incrementally** on specific features
- **Maintain context** across multiple analysis sessions
- **Generate accurate code** based on focused specifications
- **Track progress** and maintain documentation

This strategy ensures that large OpenAPI files become assets rather than obstacles in the client generation process. 