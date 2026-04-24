GOLANGCI_LINT_MK_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOLANGCI_LINT_REPO_ROOT := $(abspath $(GOLANGCI_LINT_MK_DIR)/../..)
ifeq ($(abspath $(CURDIR)),$(GOLANGCI_LINT_REPO_ROOT))
GOLANGCI_LINT_MODULE_PATH ?= root
else
GOLANGCI_LINT_MODULE_PATH ?= $(patsubst $(GOLANGCI_LINT_REPO_ROOT)/%,%,$(abspath $(CURDIR)))
endif

GOLANGCI_LINT_VERSION_FILE ?= $(GOLANGCI_LINT_REPO_ROOT)/.golangci-version
GOLANGCI_LINT_BIN_DIR ?= $(GOLANGCI_LINT_REPO_ROOT)/bin
GOLANGCI_LINT_BIN ?= $(GOLANGCI_LINT_BIN_DIR)/golangci-lint
GOLANGCI_LINT_CACHE_ROOT ?= $(GOLANGCI_LINT_REPO_ROOT)/.cache/golangci-lint
GOLANGCI_LINT_CACHE_DIR ?= $(GOLANGCI_LINT_CACHE_ROOT)/linter-$(GOLANGCI_LINT_MODULE_PATH)
GOLANGCI_LINT_GO_CACHE_DIR ?= $(GOLANGCI_LINT_CACHE_ROOT)/go-$(GOLANGCI_LINT_MODULE_PATH)
GOLANGCI_LINT_RUN_ENV ?= GOLANGCI_LINT_CACHE="$(GOLANGCI_LINT_CACHE_DIR)" GOCACHE="$(GOLANGCI_LINT_GO_CACHE_DIR)"

# Not caching for CI for now, haven't seen issues there yet
# But caching per module for local since it's happening pretty often

ifeq ($(CI),true)
GOLANGCI_LINT_EXEC ?= golangci-lint
GOLANGCI_LINT_PREREQS ?=
GOLANGCI_LINT_RUN ?= $(GOLANGCI_LINT_EXEC) run
else
GOLANGCI_LINT_EXEC ?= $(GOLANGCI_LINT_BIN)
GOLANGCI_LINT_PREREQS ?= $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_CACHE_DIR)
GOLANGCI_LINT_RUN ?= $(GOLANGCI_LINT_RUN_ENV) $(GOLANGCI_LINT_EXEC) run
endif


# Installing from sources is not recommended; the published binary is preferred.
# More info is here https://golangci-lint.run/docs/welcome/install/#install-from-sources
$(GOLANGCI_LINT_BIN): $(GOLANGCI_LINT_VERSION_FILE)
	@set -eu; \
	version="$$(cat "$(GOLANGCI_LINT_VERSION_FILE)")"; \
	mkdir -p "$(GOLANGCI_LINT_BIN_DIR)"; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(GOLANGCI_LINT_BIN_DIR)" "$$version"

$(GOLANGCI_LINT_CACHE_DIR):
	mkdir -p "$@"
