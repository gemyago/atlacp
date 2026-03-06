<!-- Nearest AGENTS.md takes precedence. Scope: deployment tooling under deploy/. Keep concise; link to root for globals. -->

## Purpose & Scope
- Module-specific guidance for rendering and installing Helm charts in `deploy/`.
- For global setup/CI/workflows, see [../AGENTS.md](../AGENTS.md).
- Living document: update in the same PR when release/deploy steps or values change.

## Setup Tools
- Install Helm toolchain: `make tools`

## Kubernetes Prerequisites
- Select kube context/namespace explicitly (avoid default):
  - Create namespace (example): `kubectl create namespace atlacp`
- Configure GHCR pull secret (example):
  - `kubectl create secret docker-registry ghcr-registry --docker-server=https://ghcr.io --docker-username="$(gh auth status | grep -o "account [^ ]*" | cut -d ' ' -f 2)" --docker-password="$(gh auth token)" --namespace atlacp`

## Security Cautions
- Do NOT commit secrets to VCS. Keep sensitive values in external secret stores or CI-provisioned secrets.
- Prefer image pull via `imagePullSecrets` referencing the secret created above.
- Verify current kube context before applying: `kubectl config current-context`
- Use `--dry-run` and `helm template` to validate manifests before any real install/upgrade.

## Helm Commands (run from repo root or deploy/)
- Render templates:
  - `helm template deploy/helm/api-service --debug --name-template api-service -f deploy/helm/api-service/values.yaml`
- Install/upgrade (dry-run first):
  - `helm upgrade api-service deploy/helm/api-service --install --namespace atlacp -f deploy/helm/api-service/values.yaml --create-namespace --dry-run`
- Uninstall (dry-run):
  - `helm uninstall api-service --namespace atlacp --dry-run`

## Definition of Done (deploy changes)
- `helm template` renders successfully using provided values.
- `helm upgrade --dry-run` completes without errors against the target namespace.
- If publishing charts/images, CI workflows pass and reference the updated steps.

## References
- Deploy details: [deploy/README.md](README.md)
- Global guidance: [../AGENTS.md](../AGENTS.md)
