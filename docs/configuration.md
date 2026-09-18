# Configuration

Every option can be given four ways. Highest precedence wins:

1. CLI flags — `--keycloak-realm-url=…`
2. Environment variables — `OTTERSCALE_SERVER_KEYCLOAK_REALM_URL=…`
3. `config.yaml` in the working directory or `/etc/otterscale/`
4. Compiled defaults

`otterscale server --help` and `otterscale agent --help` are authoritative.

## Server

| Flag                           | Environment variable                           | Default                             |
| ------------------------------ | ---------------------------------------------- | ----------------------------------- |
| `--address`                    | `OTTERSCALE_SERVER_ADDRESS`                    | `:8299`                             |
| `--allowed-origins`            | `OTTERSCALE_SERVER_ALLOWED_ORIGINS`            | _(none)_                            |
| `--tunnel-address`             | `OTTERSCALE_SERVER_TUNNEL_ADDRESS`             | `127.0.0.1:8300`                    |
| `--external-url`               | `OTTERSCALE_SERVER_EXTERNAL_URL`               | _(none)_                            |
| `--external-tunnel-url`        | `OTTERSCALE_SERVER_EXTERNAL_TUNNEL_URL`        | _(none)_                            |
| `--harbor-url`                 | `OTTERSCALE_SERVER_HARBOR_URL`                 | _(none)_                            |
| `--keycloak-realm-url`         | `OTTERSCALE_SERVER_KEYCLOAK_REALM_URL`         | _(required)_                        |
| `--keycloak-client-id`         | `OTTERSCALE_SERVER_KEYCLOAK_CLIENT_ID`         | `otterscale-server`                 |
| `--join-secret`                | `OTTERSCALE_SERVER_JOIN_SECRET`                | _(required)_                        |
| `--join-secret-file`           | `OTTERSCALE_SERVER_JOIN_SECRET_FILE`           | _(takes precedence over the above)_ |
| `--trusted-ca-file`            | `OTTERSCALE_SERVER_TRUSTED_CA_FILE`            | _(none)_                            |
| `--harbor-admin-password`      | `OTTERSCALE_SERVER_HARBOR_ADMIN_PASSWORD`      | _(none)_                            |
| `--harbor-admin-password-file` | `OTTERSCALE_SERVER_HARBOR_ADMIN_PASSWORD_FILE` | _(takes precedence over the above)_ |

The last four are what `IssueAgentValues` renders into a joining cluster's installation, and none of them is required to start the server: the chart that supplies them ships on its own cadence, so an unconfigured one fails that procedure with a message naming it rather than stopping the server from serving tunnels, proxying and RBAC.

`--external-url` must include the path prefix the gateway strips — `https://otterscale.example.com/api/`, not `https://otterscale.example.com`. It is both what joining agents are told to register against and the base of the URLs that serve their values, and the two break differently when the prefix is missing: agents cannot reach the API at all, while the URL still resolves through a port-forward and 404s at the gateway. The chart appends it.

`--trusted-ca-file` names the CA agents must trust to reach this server, for the case where the certificate served for the browser-facing URL is privately signed. Its presence is also the signal that agents need a CA at all: set, and the rendered values mount `otterscale-ca` into the agent and into Flux's source-controller, and reference it as the OCI repositories' `certSecretRef`; unset, and they carry no CA settings, which is correct for a certificate chaining to a public CA.

`--harbor-admin-password-file` is read on each use rather than at startup, so the Secret the chart mounts it from may be created later, and rotating it does not need a restart.

`--external-tunnel-url` deserves attention: it is the name the tunnel certificate is issued for, and agents pin the CA and verify that hostname. When the flag is omitted, the host from `--tunnel-address` is used, which only works when that address is a concrete name rather than a wildcard such as `0.0.0.0`. Both are validated at startup, so a mismatch surfaces there rather than as an opaque handshake failure on every agent.

## Agent

| Flag                     | Environment variable                    | Default                                                            |
| ------------------------ | --------------------------------------- | ------------------------------------------------------------------ |
| `--cluster`              | `OTTERSCALE_AGENT_CLUSTER`              | `default`                                                          |
| `--server-url`           | `OTTERSCALE_AGENT_SERVER_URL`           | `http://127.0.0.1:8299`                                            |
| `--tunnel-server-url`    | `OTTERSCALE_AGENT_TUNNEL_SERVER_URL`    | `https://127.0.0.1:8300`                                           |
| `--proxy-prometheus-url` | `OTTERSCALE_AGENT_PROXY_PROMETHEUS_URL` | `http://prometheus-stack-kube-prom-prometheus.monitoring.svc:9090` |
| `--join-token`           | `OTTERSCALE_AGENT_JOIN_TOKEN`           | _(required)_                                                       |
| `--join-token-file`      | `OTTERSCALE_AGENT_JOIN_TOKEN_FILE`      | _(takes precedence over the above)_                                |

`--cluster` is the name the cluster is addressed by in every RPC, and the name its join token is bound to.

## Secrets

The join secret, the join token and the Harbor admin password each have a `-file` variant, which takes precedence over the inline flag. Prefer them: a mounted Secret keeps the value out of the process's argv and out of `kubectl describe` output.

## Config file

The file is plain YAML, keyed by the viper names the flags derive from:

```yaml
server:
  address: :8299
  tunnel:
    address: 0.0.0.0:8300
  external_url: https://otterscale.example.com/api/
  external_tunnel_url: https://tunnel.example.com:8300
  harbor_url: https://harbor.example.com
  keycloak:
    realm_url: https://sso.example.com/realms/otterscale
    client_id: otterscale-server
  join_secret_file: /etc/otterscale/join/secret
  harbor_admin_password_file: /etc/otterscale/harbor/admin-password
```

A missing file is not an error; a malformed one is.
