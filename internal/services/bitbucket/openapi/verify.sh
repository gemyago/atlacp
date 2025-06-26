#!/bin/bash

# Script to verify the navigation files are working correctly

# Directory setup
SCRIPT_DIR="$(dirname "$0")"
OPENAPI_FILE="$SCRIPT_DIR/../openapi.json"
NAVIGATION_DIR="$SCRIPT_DIR/navigation"
SUBSETS_DIR="$SCRIPT_DIR/subsets"

echo "Verifying OpenAPI navigation infrastructure..."

# Check that key files exist
echo "Checking file structure..."
files_to_check=(
  "$OPENAPI_FILE"
  "$NAVIGATION_DIR/endpoints-index.json"
  "$NAVIGATION_DIR/models-index.json"
  "$NAVIGATION_DIR/quick-reference.md"
  "$SUBSETS_DIR/pullrequest-model.json"
  "$SUBSETS_DIR/pullrequest-endpoints.json"
)

for file in "${files_to_check[@]}"; do
  if [ -f "$file" ]; then
    echo "✅ Found: $file"
  else
    echo "❌ Missing: $file"
  fi
done

# Verify JSON files are valid
echo -e "\nValidating JSON files..."
json_files=(
  "$NAVIGATION_DIR/endpoints-index.json"
  "$NAVIGATION_DIR/models-index.json"
  "$SUBSETS_DIR/pullrequest-model.json"
  "$SUBSETS_DIR/pullrequest-endpoints.json"
)

for file in "${json_files[@]}"; do
  if jq empty "$file" 2>/dev/null; then
    echo "✅ Valid JSON: $file"
  else
    echo "❌ Invalid JSON: $file"
  fi
done

# Test example navigation lookups
echo -e "\nTesting navigation lookups..."

# Test endpoints lookup
echo "Looking up PR endpoints..."
pr_paths=$(jq -r '.pullrequests.paths | join(", ")' "$NAVIGATION_DIR/endpoints-index.json" 2>/dev/null)
if [ $? -eq 0 ] && [ -n "$pr_paths" ]; then
  echo "✅ Found PR endpoints: $pr_paths"
else
  echo "❌ Failed to find PR endpoints in endpoints-index.json"
fi

# Test models lookup
echo "Looking up PR model..."
pr_props=$(jq -r '.pullrequest.properties | join(", ")' "$NAVIGATION_DIR/models-index.json" 2>/dev/null)
if [ $? -eq 0 ] && [ -n "$pr_props" ]; then
  echo "✅ Found PR model properties: $pr_props"
else
  echo "❌ Failed to find PR model properties in models-index.json"
fi

echo -e "\nVerification complete!" 