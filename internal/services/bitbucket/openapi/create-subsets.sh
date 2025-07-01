#!/bin/bash

# Create directories if they don't exist
mkdir -p "$(dirname "$0")/subsets"

# Source file
SOURCE_FILE="$(dirname "$0")/../../bitbucket/openapi.json"
OUTDIR="$(dirname "$0")/subsets"

echo "Creating subset files from $SOURCE_FILE..."

# Extract pull request related content
echo "Generating pullrequests subset..."
jq '{
  info: .info,
  paths: (.paths | with_entries(select(.key | contains("pullrequests")))),
  definitions: (.definitions | with_entries(select(.key | contains("pullrequest") or .key | contains("paginated"))))
}' "$SOURCE_FILE" > "$OUTDIR/pullrequests.json"

# Extract repository related content
echo "Generating repositories subset..."
jq '{
  info: .info,
  paths: (.paths | with_entries(select(.key | contains("repositories")))),
  definitions: (.definitions | with_entries(select(.key | contains("repository") or .key | contains("workspace"))))
}' "$SOURCE_FILE" > "$OUTDIR/repositories.json"

# Extract core models only
echo "Generating core-models subset..."
jq '.definitions | {pullrequest, account, repository, workspace, commit}' "$SOURCE_FILE" > "$OUTDIR/core-models.json"

# Create feature-specific endpoint subsets
echo "Generating feature-specific endpoint subsets..."

# PR creation subset
jq '{
  info: .info,
  paths: (.paths | with_entries(select(.key | test("pullrequests$")))),
  definitions: (.definitions | with_entries(select(.key | contains("pullrequest") or .key == "repository" or .key == "account" or .key == "commit")))
}' "$SOURCE_FILE" > "$OUTDIR/pr-creation-endpoints.json"

# PR approval subset
jq '{
  info: .info,
  paths: (.paths | with_entries(select(.key | contains("approve")))),
  definitions: (.definitions | with_entries(select(.key | contains("pullrequest") or .key == "account")))
}' "$SOURCE_FILE" > "$OUTDIR/pr-approval-endpoints.json"

# PR merge subset
jq '{
  info: .info,
  paths: (.paths | with_entries(select(.key | contains("merge")))),
  definitions: (.definitions | with_entries(select(.key | contains("pullrequest") or .key == "repository" or .key == "account" or .key == "commit")))
}' "$SOURCE_FILE" > "$OUTDIR/pr-merge-endpoints.json"

echo "All subset files created successfully in $OUTDIR." 