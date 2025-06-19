#!/usr/bin/env bash

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_DIR=$(dirname "$SCRIPT_DIR")

set -euo pipefail

bash -c "cd $PROJECT_DIR && go run ./cmd/mcp stdio"