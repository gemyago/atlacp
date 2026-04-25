<!-- Nearest AGENTS.md takes precedence. Scope: build tooling under build/. Keep concise; link to root for globals. -->

## Purpose & Scope
- Module-specific rules for build assets, images, and packaging in `build/`.
- For global setup/CI/workflows, see [AGENTS.md](../AGENTS.md).

## Quick Setup
- Install crane tool (auto): `make install-crane`
- Optional (iterate Python scripts): `python -m venv .venv && source .venv/bin/activate && pip install -r ../requirements.txt`
- To test installer script against local package workspaces: `ATLACP_PACKAGES_DIR=build/npm/packages node build/npm/install/install.mjs`
- If running via non-interactive shells, direnv is not auto-sourced; use:
  - `direnv allow .`
  - `direnv exec . <command>` (example: `direnv exec . make test/npm-install`)

## Build Artifacts
- Build multi-platform binaries (from this dir): `make dist`
- Populate npm packages with built binaries and versions: `make npm/packages VERSION=1.2.3` (compat alias: `make npm-packages VERSION=1.2.3`; requires `jq` and `npm`)
- Generated npm package workspaces are now under `build/npm/packages/@atlacp/`:
  - `@atlacp/install-*` for platform packages
  - `@atlacp/install` for installer package
- Publish npm packages: `make npm/publish VERSION=1.2.3` (requires npm auth)
- Package artifacts tarball: `make build-artifacts.tar.bz2`
- Clean outputs: `make clean`

## Docker Images (multi-platform)
- Build local images for each binary in cmd/: `make docker/.local-images`
- Build & push remote images (registries from build.cfg): `make docker/.remote-images`
- List remote image base names (no build): `make docker/.remote-image-names`

## Configuration (build/build.cfg)
- Platforms: `platforms` (e.g., `linux/amd64,linux/arm64`)
- Runtime base image: `docker_runtime_image`
- Registries:
  - Local tagging: `docker_local_registry`
  - Push registries: `docker_push_registries` (comma-separated)
  - Registry host by name: `docker_<name>_registry` (e.g., `docker_ghcr_registry=ghcr.io`)

## CI Parity & Living Doc
- CI Build Artifacts runs from repo root: `go mod download && go install tool` then `make -C build build-artifacts.tar.bz2`
- Keep this file in sync with CI and Make targets. Update in the same PR when commands or packaging change.

## Definition of Done (build changes)
- Artifacts tarball produced: `build/build-artifacts.tar.*`
- Local image metadata files created: `build/docker/.local-*-image`
- Remote images built/pushed and names recorded in: `build/docker/.remote-images`

## References
- Build details: [build/README.md](README.md)
- Global guidance: [../AGENTS.md](../AGENTS.md)

## Task Completion Protocol

If changing any script, always run this command:
```bash
# From project root
make -C build test
```

Then always report: "Tests (build/tests) passed."

### Npm contents or npm build process

If changing anything related to npm building and packaging, always test like this before reporting as done:
```bash
make clean
make npm/publish VERSION=0.0.0-dev0
make npm/unpublish VERSION=0.0.0-dev0
make npm/clean-published
```

This command defaults to no-op publish and still verifies generated artifacts.
To perform real operations, pass:
- `NPM_PUBLISH_NOOP=false`
- `NPM_UNPUBLISH_NOOP=false`

Review output of the above, make sure no errors. Verify these outputs:
- `build/npm/packs.txt` exists and contains all expected tgz paths
- `build/npm/packages/packs/*.published` files exist for each generated pack
- Command output shows no accidental publish/unpublish when in noop mode
- `NPM_PUBLISH_NOOP=false` and/or `NPM_UNPUBLISH_NOOP=false` are required for real actions
- `npm/clean-published` removes all `.published` markers
Report:
- Run make: no errors
