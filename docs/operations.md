# Operations

## Joining a cluster

Registration is the one endpoint agents reach before they have any credentials, so it is authorised by a **join token** instead. The server holds a single root secret (`--join-secret`, `--join-secret-file`, or `OTTERSCALE_SERVER_JOIN_SECRET`) and refuses to start without one; each cluster's token is derived from that secret and the cluster's name.

`LinkService.IssueAgentValues` renders everything a joining cluster needs in one call: the token, the URLs it registers against, the cluster-admin binding for the caller, and the Harbor robot its tenant operator authenticates as. It also returns a short URL serving the identical bytes, so the file can be piped straight into Helm.

That procedure is **restricted to the admin group** (`oidc:admin` once the OIDC middleware has prefixed the token's group claims): what they return authorises claiming a cluster, and thereby cluster-admin on it.

### The CA

An agent verifies the server with its image's system roots, so a privately signed certificate has to be trusted out of band. Read it from the control plane's own Secret and create it on the joining cluster:

```console
$ kubectl -n otterscale-system get secret otterscale-ca \
    -o jsonpath='{.data.ca\.crt}' | base64 -d > ca.crt
$ kubectl --context downstream -n otterscale-system \
    create secret generic otterscale-ca --from-file=ca.crt
```

Nothing is needed when the certificate chains to a public CA, and the rendered values then carry no CA settings at all.

### Installing

Flux has to be installed on the joining cluster first: nothing reconciles a `HelmRelease` until its controllers and CRDs exist. Then fetch the values and install:

```console
$ curl -fsSL "$URL" | helm install otterscale-agent-flux \
    otterscale/otterscale-agent-flux -n otterscale-system \
    --version "$VERSION" -f -
```

Against a bare IP the certificate is necessarily privately signed, so add `-k`:

```console
$ curl -fsSLk "$URL" | helm install otterscale-agent-flux \
    otterscale/otterscale-agent-flux -n otterscale-system \
    --version "$VERSION" -f -
```

`helm install -f <url>` does not work here, and not for want of a flag: Helm fetches a `-f` URL with no TLS options at all, and `--ca-file` and `--insecure-skip-tls-verify` apply to pulling charts, not to reading values. Fetching with `curl` and piping into `-f -` covers both cases with one command.

### Chart versions

`$VERSION` is the `version` field `IssueAgentValues` returns alongside the values, and the same version heads the file as a comment, so it can be read off a saved copy without going back to the response. It matters because the values are written against one version of this chart's schema: install a different one and Helm's deep merge accepts the file while quietly dropping whatever moved.

The two nested releases need no flag. The rendered file pins `otterscale-agent` and `flux` itself, which is why an install reproduces instead of following whatever the chart repository released most recently. To override a pin, pass `--set agent.version=` or `--set flux.version=` on the command above: `--set` beats `-f`. That is the only way in: the `otterscale-agent-values` and `flux-values` ConfigMaps the file references merge into the target chart's own values, not into the umbrella chart's `version` keys.

What this does and does not give you:

- A token authorises **one cluster**. An agent holding `prod`'s token cannot register as `staging`, so a compromised agent cannot take over another cluster's traffic.
- A rejected token changes nothing. The check runs before any state is touched, so a bad registration cannot displace the agent currently serving that cluster.
- Join tokens **do not expire** and cannot be revoked one by one. Rotating the root secret invalidates every token at once, after which each agent needs its new token. It also changes the Harbor robot secret each cluster is issued.
- The URL **does** expire, an hour after it was issued, and the unguessable id in its path is the only thing authorising the fetch. Treat the URL as the credential it is. Re-issuing costs nothing: the values are derived, so the same cluster always gets the same file.
- The id carries no information. The cluster name, the cluster-admin identities and the cluster's addresses stay on the server, out of shell history, terminal scrollback and any proxy log the URL passes through.
- Outstanding URLs are held in the server's memory, so a restart invalidates them. Issue a new one; nothing else is affected, and the procedure returns the file inline as well, so a dashboard never depends on the URL surviving.
- Fetching the URL is free of side effects. It performs no writes and does not contact Harbor, so Helm retrying, Flux reconciling or an operator re-running the install cannot disturb a cluster that is already running.
- The rendered file carries the join token and the Harbor robot secret inline, which puts them in the Helm release **and** in the `HelmRelease` object — readable by anyone who can `get helmrelease`, a wider group than can read Secrets. To keep them out, create the Secrets yourself and point `agent.existingSecret` and the chart's `valuesFrom` at them.
- Registration sends the token, so the server URL should be `https://`. The agent warns at startup when it is plain HTTP to a remote host — legitimate only when something else (a service mesh, for instance) provides the transport security.

## Operating the server

The server keeps its tunnel state in memory, which shapes how it is deployed:

- **Run a single replica.** Cluster registrations, allocated loopback addresses, and live tunnel sessions live in the process that accepted them. A second replica would have its own registry and its own CA, so agents registered against one replica cannot be reached through the other, and requests routed to the wrong replica fail with "cluster not registered". Horizontal scaling needs shared state and tunnel affinity, which the current design does not provide.
- **Restarts re-key every agent.** The tunnel CA is generated at startup and never persisted, so agent certificates issued before a restart stop being trusted. Agents detect the dropped session and re-register automatically with exponential backoff, but their clusters are unreachable until they do — expect a short interruption after every restart or redeploy.
- **Agent certificates are short-lived** (24 hours) by design. Renewal happens through the same re-registration path, so no manual rotation is required.
- **Two ports must be reachable:** `8299` for the API and `8300` for the tunnel. Agents need both — the first to register, the second to stay connected.

## Monitoring

`/metrics` serves OpenTelemetry-derived Prometheus metrics, including per-procedure call counts and latencies from the ConnectRPC interceptor. The endpoint requires a bearer token like every other route; point a Prometheus scrape job at it with an `oauth2:` section using the same Keycloak client. See [api.md](api.md#additional-endpoints).

A background reaper detects disconnected tunnel clients and removes stale registrations, so `ListLinks` reflects the clusters that are actually reachable rather than every cluster that ever registered.

## Who sees which metrics

Metrics are no longer shared by every signed-in user (see [api.md](api.md#additional-endpoints)). A Workspace member sees their own namespaces through `/prometheus-tenancy`; the cluster-wide pages need cluster metrics access on that cluster. To give it to someone who is not cluster-admin, bind the ClusterRole the otterscale-agent chart ships, on each cluster they should see:

    kubectl create clusterrolebinding metrics-view-<name> --clusterrole=otterscale:metrics-cluster-view --user=<their OIDC subject>

`--group=oidc:<Keycloak group>` works too: the server prefixes token groups with `oidc:`. A change takes effect within 30 seconds, the server's cache of these answers.

The server and the dashboard must be upgraded together, and so must the kubevirt module on every cluster with VMs: upgrade it to chart 1.1.2 in the same window. A new server with an old dashboard gives regular users 403 on every chart (the old dashboard asks for cluster scope everywhere); an old server with a new dashboard answers `/prometheus-tenancy` with 404. The new dashboard's compute charts select VM metrics by `namespace`, which kubevirt 1.1.2 keeps as the VM's own namespace (`honorLabels: true`); under 1.1.1 it is `kubevirt` for every VM, so compute charts are empty for everyone, admins included, until the module is upgraded. Modules are installed per cluster at a pinned chart version, so no cluster, externally joined ones included, moves to 1.1.2 by itself. Agents need no change beyond the kube-rbac-proxy release.

## Upgrading the metrics path

prometheus-stack 1.2.0 puts kube-rbac-proxy in front of Prometheus and closes port 9090; an agent from this release queries the proxy on 9091. Neither works with the other's older version, so upgrade both, prometheus-stack first; charts are empty in between. To keep an older prometheus-stack serving a new agent for a while, set `--proxy-prometheus-url=http://prometheus-stack-kube-prom-prometheus.monitoring.svc:9090`, which sends no token.

### Clusters managed externally

A cluster joined before this release runs the otterscale-agent chart version the hub pinned at the time (1.1.x), which does not grant the agent access to the proxy or set its environment. The hub now pins 1.2.x for newly joined clusters only. To move an existing cluster, either:

- re-issue its agent values (the dashboard's import flow, or `LinkService.IssueAgentValues`) and re-apply them; or
- upgrade its `otterscale-agent-flux` release with `agent.version` set to `1.2.x` (the `agent` block in that chart's values, e.g. `agent: {version: "1.2.x"}`; the agent's own settings stay under `agent.values`).

Upgrade prometheus-stack on that cluster in the same window, and use an agent image from an otterscale release that carries this change.

## Troubleshooting

| Symptom                                                                                    | Likely cause                                                                                                                                                       |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Agent logs a TLS handshake failure against the tunnel                                      | `--external-tunnel-url` does not match the name the agent dials. The certificate is issued for that name and the agent verifies it.                                |
| Registration rejected                                                                      | The token does not match the cluster name, or the server's join secret has been rotated.                                                                           |
| `cluster not registered`                                                                   | The agent has not (re)connected yet, or requests are reaching a second server replica. Run one replica.                                                            |
| Agent warns about plaintext HTTP at startup                                                | `--server-url` is `http://` to a remote host, which exposes the join token in transit. Legitimate only behind a mesh that terminates TLS.                          |
| Every cluster goes unreachable at once, then recovers                                      | The server restarted; its CA is regenerated at startup and agents must re-register.                                                                                |
| Dashboard charts empty; agent logs `failed to read the prometheus CA`                      | prometheus-stack is older than 1.2.0 or not installed, or the agent lacks `get` on the CA Secret.                                                                  |
| Dashboard charts empty; agent logs `prometheus proxy rejected the agent` with `status=403` | kube-rbac-proxy refused the agent: its ServiceAccount is not bound to `prometheuses/api` for that Prometheus (otterscale-agent chart older than 1.2.0).            |
| Dashboard charts empty; agent logs `prometheus proxy rejected the agent` with `status=401` | kube-rbac-proxy rejected the token (TokenReview failed), e.g. the token's audience does not match or the proxy cannot reach kube-apiserver.                        |
| Agent warns `prometheus proxy disabled` at startup                                         | The agent's kube config carries no ServiceAccount token (a local kubeconfig); kube-rbac-proxy needs one, so metrics are off. The kube-apiserver proxy still works. |
| Regular users see empty charts, browser shows `403` on `/prometheus/`                      | The dashboard is older than the server; upgrade both together.                                                                                                     |
| Charts fail with `503` "cannot check metrics access"                                       | The server could not create a SelfSubjectAccessReview through the agent: the tunnel or kube-apiserver is down.                                                     |
