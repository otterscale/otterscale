# OtterScale

[![Release](https://img.shields.io/github/v/release/otterscale/otterscale?logo=github)](https://github.com/otterscale/otterscale/releases)
[![License](https://img.shields.io/github/license/otterscale/otterscale?logo=github&color=blue)](https://opensource.org/license/apache-2-0)
[![Workflow](https://img.shields.io/github/actions/workflow/status/otterscale/otterscale/ci.yml?logo=github&label=workflow)](https://github.com/otterscale/otterscale/actions/workflows/ci.yml)
[![Codecov](https://codecov.io/gh/otterscale/otterscale/graph/badge.svg?token=I7R0YEMXER)](https://codecov.io/gh/otterscale/otterscale)
[![Go Report Card](https://goreportcard.com/badge/github.com/otterscale/otterscale)](https://goreportcard.com/report/github.com/otterscale/otterscale)
[![FIPS 140-3](https://img.shields.io/badge/FIPS%20140--3-enabled-green)](https://go.dev/doc/security/fips140)

**Multi-cluster Kubernetes API gateway — a unified ConnectRPC endpoint over Chisel reverse tunnels, secured with OIDC and mTLS.**

OtterScale provides a single, authenticated entry point to many Kubernetes clusters — including clusters behind NAT, firewalls, or in air-gapped environments. A central **server (hub)** accepts ConnectRPC requests, while lightweight **agents (spokes)** running inside each cluster dial home over an mTLS reverse tunnel and forward requests to their local `kube-apiserver` with the caller's identity preserved through impersonation. The result is consistent RBAC, discovery, and runtime operations across every connected cluster.

> The original OtterScale repository now lives at [legacy](https://github.com/otterscale/legacy). This repository houses the core application; the user interface has moved to [dashboard](https://github.com/otterscale/dashboard).

## Architecture

```mermaid
  sequenceDiagram
    participant User
    participant Server as Server (Hub)
    participant Tunnel as Chisel Tunnel
    participant Agent as Agent (Spoke)
    participant K8s as kube-apiserver

    Note over Agent, Tunnel: Agent startup
    Agent->>Server: CSR registration (Link.Register)
    Server-->>Agent: mTLS certificate
    Agent->>Tunnel: Establish reverse tunnel (mTLS)
    Tunnel-->>Agent: Assigned 127.x.x.x loopback

    Note over User, K8s: User request
    User->>Server: ConnectRPC + OIDC token
    Server->>Server: Verify OIDC (Keycloak)
    Server->>Tunnel: Route to cluster loopback
    Tunnel->>Agent: Forward request
    Agent->>K8s: Impersonation (user identity)
    K8s-->>Agent: Response
    Agent-->>Tunnel: Response
    Tunnel-->>Server: Response
    Server-->>User: ConnectRPC response
```

## Features

- **Link** — Agent registration with auto-provisioned mTLS certificates via a CSR flow.
- **Resources** — Generic Kubernetes CRUD, watch, and server-side apply across clusters.
- **Runtime** — Exec/TTY, log streaming, port-forward, scaling, and rolling restarts.
- **Discovery** — API resource discovery and OpenAPI schema resolution with a TTL cache.
- **Security** — FIPS 140-3, OIDC (Keycloak), per-tunnel mTLS, and user impersonation for RBAC.

## Documentation

Installation, configuration, and operational guides will be published in the project documentation. In the meantime, `otterscale server --help` and `otterscale agent --help` describe the available options.

## Ecosystem

OtterScale's open-source components live across these repositories:

| Repository                                                       | Description                                      |
| ---------------------------------------------------------------- | ------------------------------------------------ |
| [otterscale](https://github.com/otterscale/otterscale)           | Multi-cluster Kubernetes API gateway (this repo) |
| [dashboard](https://github.com/otterscale/dashboard)             | Web management UI                                |
| [api](https://github.com/otterscale/api)                         | Shared API contract — CRDs + ConnectRPC services |
| [types](https://github.com/otterscale/types)                     | Generated TypeScript type definitions            |
| [tenant-operator](https://github.com/otterscale/tenant-operator) | Workspace / multi-tenancy operator               |

## Contributing

Contributions are welcome. A contribution guide (`CONTRIBUTING.md`) will follow; until then, please open an issue or a pull request to get involved.

## License

This project is licensed under the [Apache License 2.0](LICENSE).

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fotterscale%2Fotterscale.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fotterscale%2Fotterscale?ref=badge_large&issueType=license)
