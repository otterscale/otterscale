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

// Viper keys for the enrolment secret. They are server-mode keys, but
// live in their own group because the enrolment-token subcommand binds
// them without the rest of the server flags.
const (
	keyServerEnrolmentSecret     = "server.enrolment_secret"      //nolint:gosec // configuration key name, not a credential
	keyServerEnrolmentSecretFile = "server.enrolment_secret_file" //nolint:gosec // configuration key name, not a credential
)

// Viper keys for agent-mode configuration.
const (
	keyAgentCluster            = "agent.cluster"
	keyAgentServerURL          = "agent.server_url"
	keyAgentTunnelServerURL    = "agent.tunnel.server_url"
	keyAgentProxyPrometheusURL = "agent.proxy.prometheus_url"
	keyAgentEnrolmentToken     = "agent.enrolment_token"      //nolint:gosec // configuration key name, not a credential
	keyAgentEnrolmentTokenFile = "agent.enrolment_token_file" //nolint:gosec // configuration key name, not a credential
)
