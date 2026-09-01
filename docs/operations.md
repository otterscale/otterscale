# Operations

## Joining a cluster

Registration is the one endpoint agents reach before they have any credentials, so it is authorised by a **join token** instead. The server holds a single root secret (`--join-secret`, `--join-secret-file`, or `OTTERSCALE_SERVER_JOIN_SECRET`) and refuses to start without one; each cluster's token is derived from that secret and the cluster's name.

Issue a token wherever the root secret is available — most conveniently inside the server itself, so the secret never leaves the pod:

```console
$ kubectl exec deploy/otterscale-server -- /otterscale join token --cluster prod
xlbQpGep3w9ZJpaDyUzKpHXVTcw_5pO5mNgT3qnf3Ss
```

An agent verifies the server with its image's system roots, so a privately signed certificate has to travel with the token. The same command prints it, and says so when nothing is needed:

```console
$ kubectl exec deploy/otterscale-server -- /otterscale join ca > ca.crt
$ kubectl --context downstream -n otterscale-system \
    create secret generic otterscale-ca --from-file=ca.crt
```

Then install the agent on the joining cluster:

```console
$ helm install otterscale-agent otterscale/otterscale-agent \
    --set agent.cluster=prod \
    --set agent.joinToken=xlbQpGep3w9ZJpaDyUzKpHXVTcw_5pO5mNgT3qnf3Ss \
    --set agent.serverURL=https://otterscale.example.com/api/ \
    --set agent.tunnelServerURL=https://node1:30300 \
    --set trustedCA.secretName=otterscale-ca
```

What this does and does not give you:

- A token authorises **one cluster**. An agent holding `prod`'s token cannot register as `staging`, so a compromised agent cannot take over another cluster's traffic.
- A rejected token changes nothing. The check runs before any state is touched, so a bad registration cannot displace the agent currently serving that cluster.
- Tokens **do not expire** and cannot be revoked one by one. Rotating the root secret invalidates every token at once, after which each agent needs its new token.
- The token is sent in the registration request, so `--server-url` should be `https://`. The agent warns at startup when it is plain HTTP to a remote host — legitimate only when something else (a service mesh, for instance) provides the transport security.
- With `--set agent.joinToken`, the token is stored in the Helm release's values and is readable by anyone who can read Secrets in that namespace. To keep it out, create the Secret yourself and point `agent.existingSecret` at it. Either way the chart mounts it as a file rather than putting it in the agent's environment.

## Operating the server

The server keeps its tunnel state in memory, which shapes how it is deployed:

- **Run a single replica.** Cluster registrations, allocated loopback addresses, and live tunnel sessions live in the process that accepted them. A second replica would have its own registry and its own CA, so agents registered against one replica cannot be reached through the other, and requests routed to the wrong replica fail with "cluster not registered". Horizontal scaling needs shared state and tunnel affinity, which the current design does not provide.
- **Restarts re-key every agent.** The tunnel CA is generated at startup and never persisted, so agent certificates issued before a restart stop being trusted. Agents detect the dropped session and re-register automatically with exponential backoff, but their clusters are unreachable until they do — expect a short interruption after every restart or redeploy.
- **Agent certificates are short-lived** (24 hours) by design. Renewal happens through the same re-registration path, so no manual rotation is required.
- **Two ports must be reachable:** `8299` for the API and `8300` for the tunnel. Agents need both — the first to register, the second to stay connected.

## Monitoring

`/metrics` serves OpenTelemetry-derived Prometheus metrics, including per-procedure call counts and latencies from the ConnectRPC interceptor. The endpoint requires a bearer token like every other route; point a Prometheus scrape job at it with an `oauth2:` section using the same Keycloak client. See [api.md](api.md#additional-endpoints).

A background reaper detects disconnected tunnel clients and removes stale registrations, so `ListLinks` reflects the clusters that are actually reachable rather than every cluster that ever registered.

## Troubleshooting

| Symptom                                               | Likely cause                                                                                                                              |
| ----------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| Agent logs a TLS handshake failure against the tunnel | `--external-tunnel-url` does not match the name the agent dials. The certificate is issued for that name and the agent verifies it.       |
| Registration rejected                                 | The token does not match the cluster name, or the server's join secret has been rotated.                                                  |
| `cluster not registered`                              | The agent has not (re)connected yet, or requests are reaching a second server replica. Run one replica.                                   |
| Agent warns about plaintext HTTP at startup           | `--server-url` is `http://` to a remote host, which exposes the join token in transit. Legitimate only behind a mesh that terminates TLS. |
| Every cluster goes unreachable at once, then recovers | The server restarted; its CA is regenerated at startup and agents must re-register.                                                       |
