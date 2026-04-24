# Plan: npm Distribution for atlacp Tools

## Introduction / Overview

Currently, atlacp tools (`bbcp`, `bbmd`) are distributed exclusively as Docker images.
This creates friction for developers who need the CLI tools on their local machines — running a Docker image just to use a CLI is cumbersome.

The industry-standard solution is to publish native binaries through npm using the **optional-platform-package pattern** (used by esbuild, Biome, etc.).
The goal is to introduce an `@atlacp/install` npm package that lets developers install all atlacp tools with a single command:

```bash
npm install -g @atlacp/install
```

After install, the binaries (`bbcp`, `bbmd`) are available at `~/.atlacp/bin/` and that directory is automatically added to the user's shell PATH.

---

## Business Logic

1. The user runs `npm install -g @atlacp/install`.
2. npm evaluates `optionalDependencies` in `@atlacp/install` and downloads only the platform-specific sub-package that matches the host OS and CPU (e.g., `@atlacp/install-darwin-arm64`).
3. npm runs the `postinstall` script from `@atlacp/install`.
4. The install script:
   - Detects the current platform and locates the correct platform package inside `node_modules`.
   - Creates `~/.atlacp/` directory structure (`bin/`) if it does not exist.
   - Copies all binaries (`bbcp`, `bbmd`) from the platform package's `bin/` directory to `~/.atlacp/bin/`, setting executable permissions (`0755`).
   - Detects the user's shell (from `$SHELL`) and appends `export PATH="$HOME/.atlacp/bin:$PATH"` to the appropriate config file if `~/.atlacp/bin` is not already in `$PATH`.
   - Prints a message instructing the user to restart their shell or source the modified config.
5. The operation is **idempotent**: re-running it replaces old binaries without duplicating PATH entries.

---

## Business Logic: Configuration Path Consolidation

The `bbcp` binary uses a runtime data file (Atlassian accounts/credentials) currently defaulting to `~/.config/atlacp/accounts.json`.
To avoid scattering atlacp files across multiple home directory locations, this default path will be changed to `~/.atlacp/accounts.json`, so the entire user-facing footprint of atlacp lives under a single `~/.atlacp/` directory.

**Analysis of the config system**: The application configuration (timeouts, server ports, log levels, etc.) is fully embedded in the binary at compile time via Go's `embed` directive. There are no runtime config files to relocate. The **only** runtime file that lives in the user's home directory is the Atlassian accounts data file. The current default is resolved in `internal/services/config_path.go` using `os.UserConfigDir()` + `atlacp/accounts.json`, which maps to `~/.config/atlacp/accounts.json` on Linux/macOS.

**Decision: Change the default to `~/.atlacp/accounts.json`.**

**Justification:**
- Single directory (`~/.atlacp/`) holds everything: binaries + runtime data
- `rm -rf ~/.atlacp/` is a complete, clean uninstall
- Consistent with tools like Docker (`~/.docker/`) and kubectl (`~/.kube/`)
- Users learn one path instead of two (`~/.atlacp/bin/` and `~/.config/atlacp/`)
- The existing `--atlassian-accounts-file` flag and `APP_ATLASSIAN_ACCOUNTSFILEPATH` env var still allow users to override the path if needed

---

## High Level Architecture

**Components involved:**

| Component | Purpose |
|---|---|
| `@atlacp/install` (main npm package) | Declares platform packages as `optionalDependencies`; contains the postinstall script |
| `@atlacp/install-linux-x64` | Pre-built binaries for Linux amd64 |
| `@atlacp/install-linux-arm64` | Pre-built binaries for Linux arm64 |
| `@atlacp/install-darwin-x64` | Pre-built binaries for macOS Intel |
| `@atlacp/install-darwin-arm64` | Pre-built binaries for macOS Apple Silicon |
| `cmd/mcp/` → `cmd/bbcp/` | Rename the MCP server binary to avoid `mcp` name collision in shell PATH |
| `internal/services/config_path.go` | Change default accounts file path to `~/.atlacp/accounts.json` |
| `build/build.cfg` | Extend platform list to include darwin targets |
| `build/Makefile` | New `npm-packages` target to populate npm package `bin/` dirs and set versions |
| `.github/workflows/publish-npm.yml` | New release workflow to build binaries and publish all npm packages |

---

## Detailed Architecture

### npm Package Directory Structure

All npm-related files live under a new top-level `npm/` directory:

```
npm/
  install/                        <- @atlacp/install  (main package)
    package.json
    install.mjs                   <- postinstall script
    install.test.mjs              <- unit tests for install logic
  install-linux-x64/              <- @atlacp/install-linux-x64
    package.json
    bin/                          <- populated by Make target at release time (not in git)
  install-linux-arm64/            <- @atlacp/install-linux-arm64
    package.json
    bin/
  install-darwin-x64/             <- @atlacp/install-darwin-x64
    package.json
    bin/
  install-darwin-arm64/           <- @atlacp/install-darwin-arm64
    package.json
    bin/
```

> **Note on `bin/` directories**: They are NOT committed to git. They are created and
> populated by the `make -C build npm-packages` target during the release process.

---

### Main package — `npm/install/package.json`

```json
{
  "name": "@atlacp/install",
  "version": "0.0.0",
  "description": "Install atlacp tools (bbcp, bbmd)",
  "scripts": {
    "postinstall": "node install.mjs"
  },
  "optionalDependencies": {
    "@atlacp/install-linux-x64":    "0.0.0",
    "@atlacp/install-linux-arm64":  "0.0.0",
    "@atlacp/install-darwin-x64":   "0.0.0",
    "@atlacp/install-darwin-arm64": "0.0.0"
  },
  "publishConfig": {
    "access": "public",
    "provenance": true
  }
}
```

`version` and all `optionalDependencies` values are `0.0.0` as a placeholder; the CI workflow replaces them with the actual release version before publishing.

**Why `optionalDependencies` and not `peerDependencies`?**  
`optionalDependencies` with `os`/`cpu` fields in the platform packages is the npm-native mechanism. npm will install the matching package and silently skip the others. `peerDependencies` are for plugin relationships and do not have platform-gating.

---

### Platform package — `npm/install-linux-x64/package.json`

```json
{
  "name": "@atlacp/install-linux-x64",
  "version": "0.0.0",
  "os": ["linux"],
  "cpu": ["x64"],
  "files": ["bin/"],
  "publishConfig": {
    "access": "public",
    "provenance": true
  }
}
```

All four platform packages follow this pattern. The `os` / `cpu` pair for each:

| Package suffix | `os` | `cpu` | GOOS/GOARCH |
|---|---|---|---|
| `linux-x64` | `linux` | `x64` | `linux/amd64` |
| `linux-arm64` | `linux` | `arm64` | `linux/arm64` |
| `darwin-x64` | `darwin` | `x64` | `darwin/amd64` |
| `darwin-arm64` | `darwin` | `arm64` | `darwin/arm64` |

---

### Install script — `npm/install/install.mjs`

The script is written as an ES module (`.mjs`) targeting Node.js 18+ with no external
dependencies. It exports named functions to allow unit testing.

Key exported functions:

| Function | Description |
|---|---|
| `detectPlatformPackage(platform, arch)` | Maps `process.platform` + `process.arch` to an `@atlacp/install-*` package name, or `null` if unsupported |
| `findPlatformBinDir(packageName, baseDir)` | Resolves the `bin/` directory path inside the installed platform package under `baseDir` (defaults to searching `node_modules` up the directory tree) |
| `ensureDir(dirPath)` | Creates a directory if it does not exist |
| `copyBinaries(srcDir, destDir)` | Copies all files from `srcDir` to `destDir`, sets permissions `0o755` on each |
| `detectShellConfigFiles()` | Returns an ordered list of shell config file paths based on `process.env.SHELL` |
| `isPathAlreadyConfigured(configFile, binDir)` | Returns `true` if `configFile` already references `binDir` |
| `appendToPath(configFile, binDir)` | Appends `export PATH="$HOME/.atlacp/bin:$PATH"` to `configFile` |
| `run()` | Main entry point — orchestrates all of the above |

**Shell detection and PATH extension logic:**

| `$SHELL` ends with | Config file to modify | Line appended |
|---|---|---|
| `bash` | `~/.bashrc` | `export PATH="$HOME/.atlacp/bin:$PATH"` |
| `zsh` | `~/.zshrc` | `export PATH="$HOME/.atlacp/bin:$PATH"` |
| anything else | `~/.profile` | `export PATH="$HOME/.atlacp/bin:$PATH"` |

The script is **idempotent**: it checks if the config file already contains `/.atlacp/bin`
before appending.

---

### Binary rename — `cmd/mcp/` → `cmd/bbcp/`

The existing `mcp` binary name conflicts with the widely-used `mcp` command (Model Context Protocol CLI and other tools). All references to the binary in `cmd/mcp/` must be updated to `bbcp`:
- Rename the directory `cmd/mcp/` → `cmd/bbcp/`
- Update any internal references to the binary name within that package

---

### Configuration path change — `internal/services/config_path.go`

The `AccountsFilePathResolver.defaultBaseDir()` method currently returns:
- macOS: `~/.config`
- Linux: `os.UserConfigDir()` (typically `~/.config`)

Combined with `"atlacp/accounts.json"`, this gives `~/.config/atlacp/accounts.json`.

Change `DefaultPath()` to unconditionally return `filepath.Join(home, ".atlacp", "accounts.json")`, removing the OS-conditional logic. The result is `~/.atlacp/accounts.json` on all platforms.

---

### Build system — `build/build.cfg`

Add darwin targets so the cross-compilation produces binaries for all four platforms:

```ini
platforms = linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
```

Cross-compilation from Linux (the CI runner) to `darwin` works because CGO is already
disabled (`CGO_ENABLED=0`) — Go's standard library handles the rest.

---

### Build system — `build/Makefile` new target `npm-packages`

```makefile
# Populate npm platform packages with binaries and set versions.
# Usage: make npm-packages VERSION=1.2.3
# Depends on 'dist' having already been run.
npm-packages: dist
    $(eval NPM_DIR := $(CURDIR)/../npm)
    @$(foreach platform,$(subst $(comma), ,$(platforms)), \
        $(eval GOOS   := $(word 1,$(subst /, ,$(platform)))) \
        $(eval GOARCH := $(word 2,$(subst /, ,$(platform)))) \
        $(eval NPMARCH := $(subst amd64,x64,$(GOARCH))) \
        $(eval PKG_DIR := $(NPM_DIR)/install-$(GOOS)-$(NPMARCH)) \
        mkdir -p $(PKG_DIR)/bin && \
        cp dist/$(GOOS)/$(GOARCH)/* $(PKG_DIR)/bin/ && \
        chmod +x $(PKG_DIR)/bin/* && \
        jq ".version = \"$(VERSION)\"" $(PKG_DIR)/package.json > $(PKG_DIR)/package.json.tmp && \
        mv $(PKG_DIR)/package.json.tmp $(PKG_DIR)/package.json; \
    )
    @jq ".version = \"$(VERSION)\" | .optionalDependencies |= with_entries(.value = \"$(VERSION)\")" \
        $(NPM_DIR)/install/package.json > $(NPM_DIR)/install/package.json.tmp
    @mv $(NPM_DIR)/install/package.json.tmp $(NPM_DIR)/install/package.json
    @touch .npm-packages
```

`jq` is a standard tool available on all CI runners. The target leaves a `.npm-packages`
sentinel file so Make can track whether it has run.

---

### GitHub Actions — `.github/workflows/publish-npm.yml`

Triggered when a GitHub Release is **published** (same event as `tag-release-images.yml`).

**Authentication**: npm's trusted publishing approach — a granular access token (automation type, scoped to the `@atlacp` npm scope) is stored as `NPM_TOKEN` in GitHub repository secrets. Combined with `--provenance`, each published package is cryptographically linked to the specific GitHub Actions workflow run that produced it.

```yaml
on:
  release:
    types: [published]

permissions:
  id-token: write   # required for npm provenance attestation
  contents: read

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          ref: ${{ github.event.release.tag_name }}

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache-dependency-path: go.sum

      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          registry-url: 'https://registry.npmjs.org'

      - name: Install Go dependencies
        run: go mod download && go install tool

      - name: Build binaries (all platforms)
        working-directory: build
        run: make dist

      - name: Populate npm packages
        working-directory: build
        env:
          NPM_VERSION: ${{ github.event.release.tag_name }}
        run: |
          VERSION="${NPM_VERSION#v}"
          make npm-packages VERSION="$VERSION"

      - name: Publish platform packages
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
        run: |
          for pkg in npm/install-linux-x64 npm/install-linux-arm64 \
                     npm/install-darwin-x64 npm/install-darwin-arm64; do
            npm publish "$pkg" --provenance --access public
          done

      - name: Publish main package
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
        run: npm publish npm/install --provenance --access public
```

---

## Key Architectural Decisions

| # | Decision | Rationale |
|---|---|---|
| 1 | `optionalDependencies` (not `peerDependencies`) | npm-native platform gating via `os`/`cpu` fields; standard pattern used by esbuild, biome, etc. |
| 2 | Install to `~/.atlacp/bin/` | Puts binaries on the user's PATH system-wide, independent of any npm project context. |
| 3 | Distribute both `bbcp` and `bbmd` | `bbmd` is the user-facing CLI; `bbcp` (formerly `mcp`) is the MCP server used as a Claude Desktop integration. Both are useful on developer machines. |
| 4 | Rename `mcp` → `bbcp` | The name `mcp` conflicts with other tools in the developer PATH (e.g., the Model Context Protocol CLI). The `bb` prefix keeps it within the atlacp namespace. |
| 5 | Consolidate runtime data to `~/.atlacp/accounts.json` | Single directory for all atlacp user state; clean uninstall via `rm -rf ~/.atlacp/`; consistent with Docker/kubectl conventions. |
| 6 | `0.0.0` placeholder in committed package.json files | Avoids merge conflicts on every release; version is injected by the CI workflow using `jq`. |
| 7 | Cross-compile darwin from Linux CI runner | CGO is already disabled; Go can cross-compile to any target OS/arch with just env vars. No macOS runner needed. |
| 8 | No Windows support (initial scope) | Explicitly out of scope for this iteration. |
| 9 | npm trusted publishing with provenance | Granular automation token scoped to `@atlacp` + `--provenance` flag ties each published package to the exact GitHub Actions run. No broad-access tokens. |
| 10 | `jq` for version injection in Makefile | Correct JSON manipulation; `jq` is available on all Ubuntu CI runners. Avoids fragile `sed` patterns on JSON. |

---

## Uncertainties

1. **`npx @atlacp/install` flow**: `npx` performs a temporary install and runs the package. Since `postinstall` runs during `npm install`, using `npx` should trigger the install script. However, the PATH modification from an ephemeral `npx` context may behave differently. The primary supported flow is `npm install -g`.

2. **Binary names in platform packages**: The Makefile target copies all binaries from `build/dist/<GOOS>/<GOARCH>/` (currently `bbcp` and `bbmd`). If a future binary is not meant for user-machine distribution, the target will need a filter list.

3. **`jq` availability for local development**: `jq` is pre-installed on GitHub-hosted Ubuntu runners. Developers running `make -C build npm-packages` locally need `jq` installed. This should be documented in `build/AGENTS.md`.

---

## Related Files

### New files to create

| File | Description |
|---|---|
| `npm/install/package.json` | Main npm package manifest |
| `npm/install/install.mjs` | Postinstall script (the installer) |
| `npm/install/install.test.mjs` | Unit tests for install script logic |
| `npm/install-linux-x64/package.json` | Linux amd64 platform package |
| `npm/install-linux-arm64/package.json` | Linux arm64 platform package |
| `npm/install-darwin-x64/package.json` | macOS Intel platform package |
| `npm/install-darwin-arm64/package.json` | macOS Apple Silicon platform package |
| `.github/workflows/publish-npm.yml` | GitHub Actions workflow: publish npm on release |

### Files to modify

| File | Change |
|---|---|
| `cmd/mcp/` | Rename directory to `cmd/bbcp/`; update internal binary name references |
| `internal/services/config_path.go` | Change default accounts file path to `~/.atlacp/accounts.json` |
| `build/build.cfg` | Add `darwin/amd64` and `darwin/arm64` to `platforms` |
| `build/Makefile` | Add `npm-packages` target |
| `build/AGENTS.md` | Document `npm-packages` target and `jq` requirement |
| `.github/workflows/cleanup-docker-images.yml` | Replace hardcoded `atlacp-mcp` with dynamic image list from `docker/.remote-image-names` |

---

## Task List

**Task 1.0: Rename `cmd/mcp` binary to `bbcp`**

> **Docker image impact**: The build system derives both the binary name and the Docker image
> name from the same source — the `cmd/` subdirectory name (via `$(notdir $(wildcard ../cmd/*))`).
> The Dockerfile is fully generic (`TARGET_BINARY` build-arg). Renaming `cmd/mcp/` → `cmd/bbcp/`
> therefore changes the Docker image from `ghcr.io/<namespace>/atlacp-mcp` to
> `ghcr.io/<namespace>/atlacp-bbcp` automatically, with no Dockerfile changes needed.
> There is no easy path to keep the old image name without adding an override mechanism to the
> build system. This is an acceptable breaking change given the project is pre-1.0 and Docker
> is being superseded by npm as the primary distribution channel.

- Rename directory `cmd/mcp/` → `cmd/bbcp/`
- Update any references to the binary name within the package (e.g., command root name, usage strings)
- Verify the renamed binary builds: `go build -o /tmp/bbcp ./cmd/bbcp`
- Fix the hardcoded image name in `.github/workflows/cleanup-docker-images.yml` (line 33):
  - Currently: `./scripts/ghcr.py cleanup-versions --namespace "..." --package atlacp-mcp --really-remove`
  - The file has an existing TODO acknowledging this should be dynamic
  - Replace with a loop driven by `make -C build docker/.remote-image-names`, which already generates the image list dynamically (same mechanism used in `tag-release-images.yml`):
    ```yaml
    - name: Get image names
      working-directory: build
      run: make docker/.remote-image-names

    - name: Cleanup Docker Images
      working-directory: build
      run: |
        while IFS= read -r image; do
          package=$(basename "$image")
          ./scripts/ghcr.py cleanup-versions \
            --namespace "users/${GITHUB_REPOSITORY_OWNER}" \
            --package "$package" --really-remove
        done < docker/.remote-image-names
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    ```
- Run `make lint` and `make test`
- Note: No new tests required (rename only, no logic change)
- Write summary to `docs/implementation/npm-distribution/summary-task-1.0.md`
- **Success criteria**: Binary builds as `bbcp`; cleanup workflow references no hardcoded image names; lint and all tests pass

---

**Task 1.1: Extend build platforms to include darwin**

- Edit `build/build.cfg`: change `platforms` to `linux/amd64,linux/arm64,darwin/amd64,darwin/arm64`
- Run `make -C build dist` to confirm all four platform directories are populated:
  - `build/dist/linux/amd64/bbcp`, `build/dist/linux/amd64/bbmd`
  - `build/dist/linux/arm64/bbcp`, `build/dist/linux/arm64/bbmd`
  - `build/dist/darwin/amd64/bbcp`, `build/dist/darwin/amd64/bbmd`
  - `build/dist/darwin/arm64/bbcp`, `build/dist/darwin/arm64/bbmd`
- Run `make lint` and `make test`
- Note: No new tests required (config-only change with no logic)
- Write summary to `docs/implementation/npm-distribution/summary-task-1.1.md`
- **Success criteria**: All 8 binary files exist after `make -C build dist`; lint and tests pass

---

**Task 1.2: Create npm package structure (package.json files)**

- Create the following files with content exactly as described in Detailed Architecture:
  - `npm/install/package.json`
  - `npm/install-linux-x64/package.json`
  - `npm/install-linux-arm64/package.json`
  - `npm/install-darwin-x64/package.json`
  - `npm/install-darwin-arm64/package.json`
- Verify all `os` and `cpu` fields match the table in the Detailed Architecture section
- Run `make lint` and `make test` (no code changes, should pass trivially)
- Write summary to `docs/implementation/npm-distribution/summary-task-1.2.md`
- **Success criteria**: All 5 `package.json` files exist with correct content; lint and tests pass

---

**Task 1.3: Implement install script with unit tests (TDD)**

- Create `npm/install/install.mjs` — export all functions listed in the Detailed Architecture section as named exports; `run()` is the default execution path called at the bottom of the file when `import.meta.url === pathToFileURL(process.argv[1]).href`
- **TDD — write failing tests first** in `npm/install/install.test.mjs` using `node:test` and `node:assert` (no npm dependencies):
  - `detectPlatformPackage('linux', 'x64')` → `'@atlacp/install-linux-x64'`
  - `detectPlatformPackage('linux', 'arm64')` → `'@atlacp/install-linux-arm64'`
  - `detectPlatformPackage('darwin', 'x64')` → `'@atlacp/install-darwin-x64'`
  - `detectPlatformPackage('darwin', 'arm64')` → `'@atlacp/install-darwin-arm64'`
  - `detectPlatformPackage('win32', 'x64')` → `null`
  - `isPathAlreadyConfigured(configFileContent, '/home/user/.atlacp/bin')` → returns `true` when `/.atlacp/bin` is present in content, `false` when absent
  - `appendToPath` appends the correct `export PATH=…` line for bash/zsh
  - `appendToPath` is idempotent — calling twice on the same file does not duplicate the line
  - `detectShellConfigFiles` returns `['~/.zshrc']` when `$SHELL` ends with `zsh`
  - `detectShellConfigFiles` returns `['~/.bashrc']` when `$SHELL` ends with `bash`
  - `detectShellConfigFiles` returns `['~/.profile']` for unknown shells
- Run tests: `node --test npm/install/install.test.mjs`
  - Verify tests fail with expected reasons (missing implementations), not syntax errors
  - Missing stubs should be added before re-running
- Implement all functions in `install.mjs`
- Re-run tests: all must pass
- Run `make lint` and `make test`
- Write summary to `docs/implementation/npm-distribution/summary-task-1.3.md`
- **Success criteria**: All unit tests pass; lint and project tests pass

---

**Task 1.4: Add `npm-packages` Make target**

- Add the `npm-packages` target to `build/Makefile` as described in Detailed Architecture
- The target must:
  - Depend on `dist`
  - Accept `VERSION` as a Make variable (error if not set)
  - Iterate over all platforms in `build.cfg`, mapping `amd64`→`x64` for npm CPU naming
  - Copy binaries from `build/dist/<GOOS>/<GOARCH>/` to `npm/install-<goos>-<npmcpu>/bin/`
  - Set `chmod +x` on each copied binary
  - Update `version` in each platform `package.json` using `jq`
  - Update `version` and all `optionalDependencies` values in `npm/install/package.json` using `jq`
  - Touch `.npm-packages` sentinel file on completion
- Manual verification: run `make -C build npm-packages VERSION=0.0.1-test` (after `make -C build dist`)
  - Verify binary files exist in `npm/install-linux-x64/bin/`, `npm/install-darwin-arm64/bin/`, etc.
  - Verify all 5 `package.json` files now show `"version": "0.0.1-test"`
  - Verify main package's `optionalDependencies` values are all `"0.0.1-test"`
- Update `build/AGENTS.md`: add a line documenting the `npm-packages` target and the `jq` local requirement
- Run `make lint` and `make test`
- Write summary to `docs/implementation/npm-distribution/summary-task-1.4.md`
- **Success criteria**: Target runs without errors; files and versions are correct; lint and tests pass

---

**Task 1.5: Update default accounts file path to `~/.atlacp/accounts.json`**

- In `internal/services/config_path.go`, update `AccountsFilePathResolver.DefaultPath()` (or `defaultBaseDir()`) to unconditionally return `filepath.Join(home, ".atlacp", "accounts.json")` instead of the current OS-conditional `os.UserConfigDir()/atlacp/accounts.json`
- Remove the darwin-specific branch in `defaultBaseDir()` since the result is now the same on all platforms
- **TDD — update existing tests first** in the same file's test file:
  - Update or add test case: `DefaultPath()` returns `~/.atlacp/accounts.json` on both linux and darwin
  - Run affected tests to confirm they fail as expected
- Implement the change
- Re-run affected tests: all must pass
- Run `make lint` and `make test`
- Write summary to `docs/implementation/npm-distribution/summary-task-1.5.md`
- **Success criteria**: `DefaultPath()` returns `~/.atlacp/accounts.json`; all tests pass; lint passes

---

**Task 1.6: Create GitHub Actions workflow to publish npm on release**

- Create `.github/workflows/publish-npm.yml` with content as described in Detailed Architecture
- The workflow must:
  - Trigger on `release: types: [published]`
  - Have `id-token: write` and `contents: read` permissions
  - Strip the `v` prefix from the release tag to derive `NPM_VERSION`
  - Build all binaries with `make -C build dist`
  - Run `make -C build npm-packages VERSION="$VERSION"`
  - Publish platform packages first (in any order), then publish the main package last
  - Use `NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}` for authentication
  - Use `--provenance` and `--access public` flags on all `npm publish` calls
- Note: No unit tests for the workflow YAML itself; verify syntax with `gh workflow view publish-npm.yml` after pushing
- Run `make lint` and `make test`
- Write summary to `docs/implementation/npm-distribution/summary-task-1.6.md`
- **Success criteria**: Workflow YAML is syntactically valid; lint and tests pass

---

**Compress implementation summaries**

- Follow [compress-implementation-summaries.md](/.context/compress-implementation-summaries.md) to compress the implementation summaries.
