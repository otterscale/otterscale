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
	keyServerAddress          = "server.address"
	keyServerAllowedOrigins   = "server.allowed_origins"
	keyServerTunnelAddress    = "server.tunnel.address"
	keyServerTunnelCADir      = "server.tunnel.ca_dir"
	keyServerKeycloakRealmURL  = "server.keycloak.realm_url"
	keyServerKeycloakClientID  = "server.keycloak.client_id"
	keyServerExternalURL       = "server.external_url"
	keyServerExternalTunnelURL = "server.external_tunnel_url"
)

// Viper keys for agent-mode configuration.
const (
	keyAgentCluster         = "agent.cluster"
	keyAgentServerURL       = "agent.server_url"
	keyAgentTunnelServerURL = "agent.tunnel.server_url"
	keyAgentBootstrap       = "agent.bootstrap"
)
