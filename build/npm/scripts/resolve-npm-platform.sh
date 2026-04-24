#!/usr/bin/env bash
# build/npm/scripts/resolve-npm-platform.sh
#
# Converts GOOS/GOARCH pairs to npm platform identifiers.
# Used by build/npm/Makefile to map Go build matrix to npm package names.
#
# Convention follows: golang-backend-boilerplate/build/scripts/resolve-docker-tags.sh
# (standalone, CLI flags, --self-test for local/CI testing)
#
# Usage:
#   resolve-npm-platform.sh --goos linux --goarch amd64 --format suffix
#     -> linux-x64
#   resolve-npm-platform.sh --goos linux --goarch amd64 --format os
#     -> linux
#   resolve-npm-platform.sh --goos linux --goarch amd64 --format cpu
#     -> x64
#   resolve-npm-platform.sh --goos windows --goarch amd64 --format suffix
#     -> win32-x64
#   resolve-npm-platform.sh --self-test
#     -> runs internal unit tests

set -euo pipefail

usage() {
  echo "Usage: $0 --goos <GOOS> --goarch <GOARCH> --format <suffix|os|cpu>"
  echo "       $0 --self-test"
}

npm_os() {
  case "$1" in
    linux)   echo "linux"  ;;
    darwin)  echo "darwin" ;;
    windows) echo "win32"  ;;
    *)       echo "Error: unsupported GOOS: $1" >&2; exit 1 ;;
  esac
}

npm_cpu() {
  case "$1" in
    amd64) echo "x64"   ;;
    arm64) echo "arm64" ;;
    *)     echo "Error: unsupported GOARCH: $1" >&2; exit 1 ;;
  esac
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

  echo "Running resolve-npm-platform self-tests..."

  assert "$(npm_os linux)"   "linux"  "linux os"
  assert "$(npm_os darwin)"  "darwin" "darwin os"
  assert "$(npm_os windows)" "win32"  "windows os"

  assert "$(npm_cpu amd64)"  "x64"   "amd64 cpu"
  assert "$(npm_cpu arm64)"  "arm64" "arm64 cpu"

  # suffix tests (inline: combine os + cpu)
  assert "$(npm_os linux)-$(npm_cpu amd64)"   "linux-x64"    "linux/amd64 suffix"
  assert "$(npm_os linux)-$(npm_cpu arm64)"   "linux-arm64"  "linux/arm64 suffix"
  assert "$(npm_os darwin)-$(npm_cpu arm64)"  "darwin-arm64" "darwin/arm64 suffix"
  assert "$(npm_os windows)-$(npm_cpu amd64)" "win32-x64"    "windows/amd64 suffix"

  if [[ $failures -eq 0 ]]; then
    echo "All tests passed!"
    return 0
  else
    echo "$failures test(s) failed."
    return 1
  fi
}

# === Parse args ===
GOOS_VAL=""
GOARCH_VAL=""
FORMAT="suffix"
SELF_TEST=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --goos)      GOOS_VAL="$2";  shift 2 ;;
    --goarch)    GOARCH_VAL="$2"; shift 2 ;;
    --format)    FORMAT="$2";    shift 2 ;;
    --self-test) SELF_TEST=true;  shift   ;;
    -h|--help)   usage; exit 0            ;;
    *) echo "Unknown option: $1" >&2; usage; exit 1 ;;
  esac
done

if $SELF_TEST; then
  run_tests
  exit $?
fi

if [[ -z "$GOOS_VAL" || -z "$GOARCH_VAL" ]]; then
  echo "Error: --goos and --goarch are required" >&2
  usage
  exit 1
fi

OS=$(npm_os "$GOOS_VAL")
CPU=$(npm_cpu "$GOARCH_VAL")

case "$FORMAT" in
  suffix) echo "${OS}-${CPU}" ;;
  os)     echo "$OS"          ;;
  cpu)    echo "$CPU"         ;;
  *)      echo "Error: unknown format: $FORMAT" >&2; usage; exit 1 ;;
esac
