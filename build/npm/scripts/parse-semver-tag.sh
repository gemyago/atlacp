#!/usr/bin/env bash
# build/npm/scripts/parse-semver-tag.sh
#
# Parses a git semver tag and outputs variables for use in Makefile and CI.
# Used to derive npm dist-tag and prerelease flag from a git tag.
#
# Convention follows: golang-backend-boilerplate/build/scripts/resolve-docker-tags.sh
# (standalone, CLI flags, --self-test for local/CI testing)
#
# Usage:
#   parse-semver-tag.sh --tag v1.2.3
#   parse-semver-tag.sh --tag v1.2.3-alpha.1
#   parse-semver-tag.sh --self-test
#
# Output (one KEY=VALUE per line, sourceable):
#   VERSION=1.2.3             # semver without leading 'v'
#   PRERELEASE_ID=            # empty for stable, or: alpha, beta, rc, etc.
#   NPM_TAG=latest            # npm dist-tag: latest | alpha | beta | rc | next
#   IS_PRERELEASE=false       # true for pre-release tags (for GitHub Release flag)
#
# Examples:
#   v1.2.3          -> VERSION=1.2.3,   PRERELEASE_ID=,     NPM_TAG=latest, IS_PRERELEASE=false
#   v1.2.3-alpha.1  -> VERSION=1.2.3-alpha.1, PRERELEASE_ID=alpha, NPM_TAG=alpha, IS_PRERELEASE=true
#   v1.2.3-beta.2   -> VERSION=1.2.3-beta.2,  PRERELEASE_ID=beta,  NPM_TAG=beta,  IS_PRERELEASE=true
#   v1.2.3-rc.1     -> VERSION=1.2.3-rc.1,    PRERELEASE_ID=rc,    NPM_TAG=rc,    IS_PRERELEASE=true
#   v1.2.3-snap.1   -> VERSION=1.2.3-snap.1,  PRERELEASE_ID=snap,  NPM_TAG=next,  IS_PRERELEASE=true

set -euo pipefail

usage() {
  echo "Usage: $0 --tag <git-tag>"
  echo "       $0 --self-test"
}

parse_tag() {
  local raw_tag="$1"
  local version="${raw_tag#v}"  # strip leading 'v'
  local prerelease_id=""
  local npm_tag="latest"
  local is_prerelease="false"

  # Match: X.Y.Z-<identifier>.<number> (standard npm pre-release format)
  if [[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+-([a-zA-Z]+)\. ]]; then
    prerelease_id="${BASH_REMATCH[1]}"
    is_prerelease="true"
    case "$prerelease_id" in
      alpha) npm_tag="alpha" ;;
      beta)  npm_tag="beta"  ;;
      rc)    npm_tag="rc"    ;;
      *)     npm_tag="next"  ;;
    esac
  fi

  echo "VERSION=${version}"
  echo "PRERELEASE_ID=${prerelease_id}"
  echo "NPM_TAG=${npm_tag}"
  echo "IS_PRERELEASE=${is_prerelease}"
}

run_tests() {
  local failures=0

  assert() {
    local actual="$1"
    local expected="$2"
    local message="$3"
    if [[ "$actual" == "$expected" ]]; then
      echo "PASS: $message"
    else
      echo "FAIL: $message (expected '$expected', got '$actual')"
      ((failures++))
    fi
  }

  echo "Running parse-semver-tag self-tests..."

  # --- Stable release ---
  local result
  result=$(parse_tag "v1.2.3")
  assert "$(echo "$result" | grep ^VERSION=)"      "VERSION=1.2.3"       "stable VERSION"
  assert "$(echo "$result" | grep ^PRERELEASE_ID=)" "PRERELEASE_ID="      "stable no prerelease"
  assert "$(echo "$result" | grep ^NPM_TAG=)"      "NPM_TAG=latest"      "stable npm tag"
  assert "$(echo "$result" | grep ^IS_PRERELEASE=)" "IS_PRERELEASE=false" "stable is_prerelease=false"

  # --- alpha pre-release ---
  result=$(parse_tag "v1.2.3-alpha.1")
  assert "$(echo "$result" | grep ^VERSION=)"      "VERSION=1.2.3-alpha.1" "alpha VERSION"
  assert "$(echo "$result" | grep ^PRERELEASE_ID=)" "PRERELEASE_ID=alpha"  "alpha PRERELEASE_ID"
  assert "$(echo "$result" | grep ^NPM_TAG=)"      "NPM_TAG=alpha"        "alpha NPM_TAG"
  assert "$(echo "$result" | grep ^IS_PRERELEASE=)" "IS_PRERELEASE=true"   "alpha IS_PRERELEASE"

  # --- beta pre-release ---
  result=$(parse_tag "v1.2.3-beta.2")
  assert "$(echo "$result" | grep ^NPM_TAG=)" "NPM_TAG=beta" "beta NPM_TAG"

  # --- rc pre-release ---
  result=$(parse_tag "v1.2.3-rc.1")
  assert "$(echo "$result" | grep ^NPM_TAG=)" "NPM_TAG=rc" "rc NPM_TAG"

  # --- unknown pre-release -> next ---
  result=$(parse_tag "v1.2.3-snapshot.1")
  assert "$(echo "$result" | grep ^NPM_TAG=)"      "NPM_TAG=next"            "unknown prerelease -> next"
  assert "$(echo "$result" | grep ^PRERELEASE_ID=)" "PRERELEASE_ID=snapshot"  "unknown prerelease ID preserved"

  # --- no leading v ---
  result=$(parse_tag "1.2.3")
  assert "$(echo "$result" | grep ^VERSION=)" "VERSION=1.2.3" "no-v-prefix VERSION"

  # --- higher patch number ---
  result=$(parse_tag "v10.20.30-alpha.5")
  assert "$(echo "$result" | grep ^VERSION=)" "VERSION=10.20.30-alpha.5" "multi-digit VERSION"

  if [[ $failures -eq 0 ]]; then
    echo "All tests passed!"
    return 0
  else
    echo "$failures test(s) failed."
    return 1
  fi
}

# === Parse args ===
TAG=""
SELF_TEST=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --tag)       TAG="$2"; shift 2 ;;
    --self-test) SELF_TEST=true; shift ;;
    -h|--help)   usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage; exit 1 ;;
  esac
done

if $SELF_TEST; then
  run_tests
  exit $?
fi

if [[ -z "$TAG" ]]; then
  echo "Error: --tag is required" >&2
  usage
  exit 1
fi

parse_tag "$TAG"
