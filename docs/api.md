# API surface

Three ConnectRPC services, defined in [proto/](../proto/) and served over gRPC, gRPC-Web, and Connect (JSON over HTTP) on the same port.

| Service                       | RPCs                                                                                                                                                    | Purpose                                                                                                                                                                        |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `link.v1.LinkService`         | `Register`, `ListLinks`                                                                                                                                 | Agent registration via CSR, and the roster of connected clusters.                                                                                                              |
| `resource.v1.ResourceService` | `Discovery`, `Schema`, `List`, `Get`, `Describe`, `Create`, `Apply`, `Update`, `Delete`, `Watch`                                                        | Generic typed and unstructured access to any Kubernetes resource, including server-side apply and streaming watches. Discovery and OpenAPI schemas are TTL-cached per cluster. |
| `runtime.v1.RuntimeService`   | `PodLog`, `ExecuteTTY`/`WriteTTY`/`ResizeTTY`, `PortForward`/`WritePortForward`, `VNC`/`WriteVNC`, `Scale`, `Restart`, `SubResourceAction`, `ShowChart` | Interactive operations: log follow, exec with a resizable TTY, port-forward, VNC consoles, scaling, rolling restarts, and Helm chart inspection.                               |

Every request carries a `cluster` field naming the target; the server routes it down that cluster's tunnel.

## Additional endpoints

- **`/proxy/{cluster}/prometheus/{path...}`** — a read-only proxy to each cluster's in-cluster Prometheus. Only query endpoints pass the allowlist; admin paths (`/api/v1/admin/*`, `/-/*`) are rejected after path normalisation, so `..` segments cannot slip past the check.

  Authorisation here is deliberately cluster-agnostic: any authenticated user may query any registered cluster's metrics. Metrics are treated as shared operational data, unlike the Kubernetes paths, which enforce per-cluster RBAC by impersonation. Namespace-level isolation is not provided and cannot be added at this layer — the allowlist gates endpoints, not PromQL. Restricting tenants to their own namespaces takes an enforced label matcher (prom-label-proxy or equivalent) in front of Prometheus.

- **`/metrics`** — OpenTelemetry-derived Prometheus metrics, **authenticated like every other route**. These carry cluster names and per-procedure call patterns across every managed cluster, so a scrape must present a bearer token; Prometheus can obtain one with an `oauth2:` section pointed at the same Keycloak client. A deployment that needs open scraping should expose it on a separate listener rather than making this path public, which would also open it on the internet-facing API port.

- **gRPC health checking and server reflection** — both public.

## Authentication

Requests are authenticated by OIDC against the configured Keycloak realm. Exactly three paths are reachable without a token:

- `link.v1.LinkService/Register` — agents have no credential until they have registered; the call is authorised by a [join token](operations.md#joining-a-cluster) instead.
- `grpc.health.v1.Health/Check` and `/Watch`
- `grpc.reflection.v1.ServerReflection/ServerReflectionInfo`

## Generated clients

The proto definitions are the single source of truth; three client surfaces are generated from them by `make proto`:

- **Go** — [api/](../api/), importable as `github.com/otterscale/otterscale/api/...`.
- **TypeScript** — [ts/](../ts/), published as [`@otterscale/api`](https://www.npmjs.com/package/@otterscale/api).
- **OpenAPI** — [openapi.yaml](../openapi.yaml), for anything that speaks plain HTTP+JSON. RPCs marked `NO_SIDE_EFFECTS` also accept `GET`.
