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
	keyServerExternalURL       = "server.external_url"
	keyServerExternalTunnelURL = "server.external_tunnel_url"
	keyServerHarborURL         = "server.harbor_url"
)

// Server-mode keys naming a credential, or the file holding one.
const (
	keyServerJoinSecret              = "server.join_secret"      //nolint:gosec // configuration key name, not a credential
	keyServerJoinSecretFile          = "server.join_secret_file" //nolint:gosec // configuration key name, not a credential
	keyServerTrustedCAFile           = "server.trusted_ca_file"  //nolint:gosec // configuration key name, not a credential
	keyServerHarborAdminPassword     = "server.harbor_admin_password"
	keyServerHarborAdminPasswordFile = "server.harbor_admin_password_file"
)

// Viper keys for agent-mode configuration.
const (
	keyAgentCluster            = "agent.cluster"
	keyAgentServerURL          = "agent.server_url"
	keyAgentTunnelServerURL    = "agent.tunnel.server_url"
	keyAgentProxyPrometheusURL = "agent.proxy.prometheus_url"
	keyAgentJoinToken          = "agent.join_token"
	keyAgentJoinTokenFile      = "agent.join_token_file"
)
