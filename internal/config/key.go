// Package config provides unified configuration loading from files,
// environment variables, and CLI flags using viper and pflag.
//
// Resolution order (highest wins):
//  1. CLI flags
//  2. Environment variables (prefix OTTERSCALE_)
//  3. Config file (config.yaml in . or /etc/otterscale/)
//  4. Compiled defaults
package config

// Viper keys for server-mode configuration.
const (
	keyServerAddress           = "server.address"
	keyServerAllowedOrigins    = "server.allowed_origins"
	keyServerTunnelAddress     = "server.tunnel.address"
	keyServerKeycloakRealmURL  = "server.keycloak.realm_url"
	keyServerKeycloakClientID  = "server.keycloak.client_id"
	keyServerExternalTunnelURL = "server.external_tunnel_url"
)

// Server-mode keys for the join secret, grouped separately because the
// "join token" subcommand binds them without the rest of the server flags.
const (
	keyServerJoinSecret     = "server.join_secret"      //nolint:gosec // configuration key name, not a credential
	keyServerJoinSecretFile = "server.join_secret_file" //nolint:gosec // configuration key name, not a credential
	keyServerTrustedCAFile       = "server.trusted_ca_file"       //nolint:gosec // configuration key name, not a credential
)

// Viper keys for agent-mode configuration.
const (
	keyAgentCluster            = "agent.cluster"
	keyAgentServerURL          = "agent.server_url"
	keyAgentTunnelServerURL    = "agent.tunnel.server_url"
	keyAgentProxyPrometheusURL = "agent.proxy.prometheus_url"
	keyAgentJoinToken     = "agent.join_token"      //nolint:gosec // configuration key name, not a credential
	keyAgentJoinTokenFile = "agent.join_token_file" //nolint:gosec // configuration key name, not a credential
)
