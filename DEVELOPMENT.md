# Local Development with Kind

This guide explains how to use [Kind](https://kind.sigs.k8s.io/) (Kubernetes in Docker) to run a local Gravitee APIM environment for development and testing.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Helm](https://helm.sh/docs/intro/install/)
- [Node.js](https://nodejs.org/) (for the `npx zx` orchestration script)
- [Go](https://go.dev/doc/install) (to build and test the provider)
- [Terraform](https://developer.hashicorp.com/terraform/install) (to run acceptance and example tests)
- [OpenTofu](https://opentofu.org/docs/intro/install/) (alternative to Terraform, used by example tests)
- Access to the `graviteeio.azurecr.io` container registry (for pulling APIM images)

## Quick Start

Start a local Kind cluster running the default APIM version:

```bash
make start-cluster
```

This creates a Kind cluster named `gravitee`, pulls all required Docker images, deploys APIM via Helm, and configures it with default settings.

Once the cluster is ready, the following endpoints are available:

| Service              | URL                            |
|----------------------|--------------------------------|
| Console UI           | http://localhost:30080          |
| Gateway (HTTP)       | http://localhost:30082          |
| Management API       | http://localhost:30083/management |
| Automation API       | http://localhost:30083/automation |
| Gateway (HTTPS/mTLS) | https://localhost:30084         |

Default credentials: `admin` / `admin`.

To tear down the cluster:

```bash
make stop-cluster
```

## Deployment Modes

### Full Mode (default)

Deploys the complete APIM stack: Management API, Gateway, Console UI, and MongoDB backend services for testing. Uses `hack/kind/apim/values.yaml`.

```bash
make start-cluster
```

### Minimal Mode

Deploys only the Management API and MongoDB. This is faster and uses fewer resources — useful when you only need the Automation API for provider testing. Uses `hack/kind/apim/values-minimal.yaml`.

```bash
APIM_MINIMAL=true make start-cluster
```

In minimal mode, only the Management/Automation API endpoint is available at `http://localhost:30083`.

## Configuring the APIM Version

### Default Version

The default APIM image tag and Helm chart version are defined in `hack/apim.yaml`:

```yaml
apim:
  image:
    registry: graviteeio.azurecr.io
    version: master-latest
  chart:
    registry: oci://graviteeio.azurecr.io/helm/apim3
    version: 4.11.*
```

To change the default version for all local runs, edit this file directly.

### Overriding via Environment Variables

You can override the version without modifying any files by setting environment variables:

| Variable              | Description                          | Default (from `hack/apim.yaml`)          |
|-----------------------|--------------------------------------|------------------------------------------|
| `APIM_IMAGE_TAG`      | Docker image tag for APIM components | `4.10.x-latest`                          |
| `APIM_IMAGE_REGISTRY` | Docker registry for APIM images      | `graviteeio.azurecr.io`                  |
| `APIM_CHART_VERSION`  | Helm chart version                   | `4.10.*`                                 |
| `APIM_CHART_REGISTRY` | Helm chart OCI registry              | `oci://graviteeio.azurecr.io/helm/apim3` |
| `APIM_MINIMAL`        | Use minimal deployment (`true`/`false`) | not set (full mode)                      |
| `APIM_VALUES`         | Custom Helm values file name (in `hack/kind/apim/`) | `values.yaml`                            |
| `APIM_GRAVITEE_LICENSE` | Gravitee license key (for enterprise features) | not set                                  |

### Examples

Run a specific released version:

```bash
APIM_IMAGE_REGISTRY=graviteeio APIM_IMAGE_TAG=4.10.5 APIM_CHART_REGISTRY="graviteeio/apim3" APIM_CHART_VERSION=4.10.5 make stop-cluster start-cluster
```

Run an older version (using the latest commit) in minimal mode:

```bash
APIM_IMAGE_TAG=4.9.x-latest APIM_CHART_VERSION="4.9.*" APIM_MINIMAL=true make start-cluster
```

Run the latest development build:

```bash
APIM_IMAGE_TAG=master-latest APIM_CHART_VERSION="4.11.*" make start-cluster
```

## Running Tests

All test commands assume a running APIM instance (via Kind or otherwise) at `http://localhost:30083/automation`.

### Unit Tests (no APIM required)

```bash
make unit-tests
```

### Acceptance Tests

```bash
make acceptance-tests
```

To run a single test:

```bash
APIM_USERNAME=admin APIM_PASSWORD=admin APIM_SERVER_URL=http://localhost:30083/automation \
  TF_ACC=1 go test -count=1 -v -run TestAccApiv4Resource ./tests/acceptance
```

### Example Tests

These run actual Terraform configurations from the `examples/` directory. Requires a local provider build:

```bash
make examples-tests
```

### Testing Against a Custom APIM URL

If your APIM instance is not running on the default Kind ports, override the connection settings:

```bash
APIM_SERVER_URL=http://your-host:port/automation \
APIM_USERNAME=admin \
APIM_PASSWORD=admin \
  make acceptance-tests
```

## Port Mappings

The Kind cluster exposes the following NodePorts to the host:

| Port  | Service                     |
|-------|-----------------------------|
| 30080 | APIM Console UI             |
| 30081 | APIM Portal (not used currently) |
| 30082 | APIM Gateway (HTTP)         |
| 30083 | APIM Management API         |
| 30084 | APIM Gateway (HTTPS/mTLS)   |
| 32767 | GKO Controller debug port   |

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -n default
```

### View APIM Logs

```bash
kubectl logs -l app.kubernetes.io/name=apim3 --tail=100
```

### Cluster Won't Start

If the cluster fails to start, delete it and try again:

```bash
make stop-cluster
make start-cluster
```

### Images Fail to Pull

Ensure you are authenticated to the `graviteeio.azurecr.io` registry:

```bash
docker login graviteeio.azurecr.io
```

### Pods Stuck in Pending/CrashLoopBackOff

Check events and resource usage:

```bash
kubectl describe pod <pod-name>
kubectl top nodes
```

Kind clusters share resources with your Docker daemon. Ensure Docker has enough memory allocated (at least 4 GB recommended).
