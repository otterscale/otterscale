# Development

Go 1.26 or later. `make help` lists every target.

```console
$ make build          # build ./bin/otterscale with FIPS 140-3
$ make test           # go test with coverage
$ make vet            # go vet
$ make lint           # golangci-lint
$ make proto          # regenerate Go, TypeScript, and OpenAPI from proto/
$ make proto-lint     # buf lint + format check
$ make proto-breaking # buf breaking-change check against main
```

`make proto` installs a pinned `buf` into `./bin` on first use, so no global install is needed.

## Layout

| Path                                          | What lives there                                                       |
| --------------------------------------------- | ---------------------------------------------------------------------- |
| [cmd/otterscale/](../cmd/otterscale/)         | Entry point and Wire injectors.                                        |
| [internal/cmd/](../internal/cmd/)             | Cobra commands and the server/agent runtimes.                          |
| [internal/core/](../internal/core/)           | Domain logic: links, resources, runtime, sessions, enrolment, caching. |
| [internal/handler/](../internal/handler/)     | ConnectRPC handlers translating proto to core.                         |
| [internal/transport/](../internal/transport/) | HTTP server, chisel tunnel, in-memory pipe listener.                   |
| [internal/providers/](../internal/providers/) | Kubernetes, Helm, chisel, and cache wiring.                            |
| [internal/pki/](../internal/pki/)             | The tunnel CA and certificate issuance.                                |
| [proto/](../proto/)                           | Service definitions — the source of truth for all three clients.       |

## Dependency injection

Dependencies are assembled with [Wire](https://github.com/google/wire). After changing a provider set, regenerate the injectors:

```console
$ go tool wire ./...
```

## Changing the API

1. Edit the `.proto` files under [proto/](../proto/).
2. `make proto-lint` — buf's lint and format rules are enforced in CI.
3. `make proto-breaking` — checks the change against `main`. Breaking a released RPC needs a new version package, not an edit in place.
4. `make proto` — regenerates [api/](../api/), [ts/](../ts/), and [openapi.yaml](../openapi.yaml). Commit the generated output; CI verifies it is up to date.

A new RPC that returns a long-lived stream must also be added to `Handler.LongRunningPaths()`, or the transport's request timeouts will cut its sessions off mid-flight.

## Running locally

The server needs a Keycloak realm and an enrolment secret; the agent needs a token derived from that secret. Against a kind or minikube cluster the agent picks up your ambient kubeconfig, so no in-cluster deployment is required:

```console
$ otterscale server \
    --keycloak-realm-url=https://sso.example.com/realms/otterscale \
    --enrolment-secret=dev-only-secret &

$ otterscale agent \
    --cluster=dev \
    --enrolment-token="$(otterscale enrolment-token --cluster dev --enrolment-secret=dev-only-secret)"
```

The defaults already point the agent at `127.0.0.1:8299` and `127.0.0.1:8300`.
