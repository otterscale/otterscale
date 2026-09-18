# API surface

Three ConnectRPC services, defined in [proto/](../proto/) and served over gRPC, gRPC-Web, and Connect (JSON over HTTP) on the same port.

| Service                       | RPCs                                                                                                                                                    | Purpose                                                                                                                                                                        |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `link.v1.LinkService`         | `Register`, `ListLinks`, `IssueAgentValues`                                                                                                             | Agent registration via CSR, the roster of connected clusters, and what an operator needs to install an agent on a joining one.                                                 |
| `resource.v1.ResourceService` | `Discovery`, `Schema`, `List`, `Get`, `Describe`, `Create`, `Apply`, `Update`, `Delete`, `Watch`                                                        | Generic typed and unstructured access to any Kubernetes resource, including server-side apply and streaming watches. Discovery and OpenAPI schemas are TTL-cached per cluster. |
| `runtime.v1.RuntimeService`   | `PodLog`, `ExecuteTTY`/`WriteTTY`/`ResizeTTY`, `PortForward`/`WritePortForward`, `VNC`/`WriteVNC`, `Scale`, `Restart`, `SubResourceAction`, `ShowChart` | Interactive operations: log follow, exec with a resizable TTY, port-forward, VNC consoles, scaling, rolling restarts, and Helm chart inspection.                               |

Every request carries a `cluster` field naming the target; the server routes it down that cluster's tunnel.

## Additional endpoints

- **`/proxy/{cluster}/prometheus/{path...}`** — a read-only proxy to each cluster's in-cluster Prometheus. Only query endpoints pass the allowlist; admin paths (`/api/v1/admin/*`, `/-/*`) are rejected after path normalisation, so `..` segments cannot slip past the check.

  Authorisation here is deliberately cluster-agnostic: any authenticated user may query any registered cluster's metrics. Metrics are treated as shared operational data, unlike the Kubernetes paths, which enforce per-cluster RBAC by impersonation. Namespace-level isolation is not provided and cannot be added at this layer — the allowlist gates endpoints, not PromQL. Restricting tenants to their own namespaces takes an enforced label matcher (prom-label-proxy or equivalent) in front of Prometheus.

- **`/link/values/{id}`** — the Helm override values that install an agent on a joining cluster, as raw YAML, so the file can be piped into `helm install -f -`. **Public**, because the id in the path is the whole credential: the response carries a join token and a registry secret, which is why it is served `no-store` and why `IssueAgentValues`, which mints these URLs, is not marked `NO_SIDE_EFFECTS` and so cannot be fetched by GET or cached.

  The id is 128 bits from a CSPRNG, naming an entry the server holds for an hour. It carries nothing: the cluster name, the cluster-admin identities and the addresses stay server-side rather than travelling in a URL, which also keeps the whole thing at 62 characters regardless of how many identities the values bind. Every rejection is the same `401 invalid or expired token` whether the id was never issued or has expired, and nothing about it reaches the logs.

  Serving it performs no writes and contacts no registry, so a repeated fetch — Helm retrying, an operator re-running an install — cannot disturb a cluster that is already running. The entries live in the server's memory, so a restart invalidates outstanding URLs; the procedure returns the file inline as well, so re-issuing is the whole remedy.

  There is no rate limit: for an hour after issue, any caller holding the id can make the server perform one map lookup and one render.

- **`/metrics`** — OpenTelemetry-derived Prometheus metrics, **authenticated like every other route**. These carry cluster names and per-procedure call patterns across every managed cluster, so a scrape must present a bearer token; Prometheus can obtain one with an `oauth2:` section pointed at the same Keycloak client. A deployment that needs open scraping should expose it on a separate listener rather than making this path public, which would also open it on the internet-facing API port.

- **gRPC health checking and server reflection** — both public.

## Authentication

Requests are authenticated by OIDC against the configured Keycloak realm. Four paths are reachable without a token:

- `link.v1.LinkService/Register` — agents have no credential until they have registered; the call is authorised by a [join token](operations.md#joining-a-cluster) instead.
- `/link/values/` and everything under it — authorised by the signed token in the path, as above.
- `grpc.health.v1.Health/Check` and `/Watch`
- `grpc.reflection.v1.ServerReflection/ServerReflectionInfo`

`IssueAgentValues` additionally requires membership of the admin group (`oidc:admin`): what it returns claims a cluster, and thereby cluster-admin on it.

## Generated clients

The proto definitions are the single source of truth; three client surfaces are generated from them by `make proto`:

- **Go** — [api/](../api/), importable as `github.com/otterscale/otterscale/api/...`.
- **TypeScript** — [ts/](../ts/), published as [`@otterscale/api`](https://www.npmjs.com/package/@otterscale/api).
- **OpenAPI** — [openapi.yaml](../openapi.yaml), for anything that speaks plain HTTP+JSON. RPCs marked `NO_SIDE_EFFECTS` also accept `GET`.
