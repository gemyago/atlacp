<!-- Nearest AGENTS.md takes precedence. Scope: build tooling under build/. Keep concise; link to root for globals. -->

## Purpose & Scope
- Module-specific rules for build assets, images, and packaging in `build/`.
- For global setup/CI/workflows, see [AGENTS.md](../AGENTS.md).

## Quick Setup
- Install crane tool (auto): `make install-crane`
- Optional (iterate Python scripts): `python -m venv .venv && source .venv/bin/activate && pip install -r ../requirements.txt`

## Build Artifacts
- Build multi-platform binaries (from this dir): `make dist`
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
