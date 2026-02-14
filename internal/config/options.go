package config

import (
	"strings"
)

// Option describes a single configuration entry: its viper key, the
// corresponding CLI flag name, the compiled default, and a
// human-readable description shown in --help output.
type Option struct {
	Key         string
	Flag        string
	Default     any
	Description string
}

// ServerOptions defines the configuration entries available in server
// mode. Each entry is registered as a viper default and a CLI flag.
var ServerOptions = []Option{
	{Key: keyServerAddress, Flag: toFlag(keyServerAddress), Default: ":8299", Description: "Server listen address"},
	{Key: keyServerAllowedOrigins, Flag: toFlag(keyServerAllowedOrigins), Default: []string{}, Description: "Server allowed origins"},
	{Key: keyServerTunnelAddress, Flag: toFlag(keyServerTunnelAddress), Default: "127.0.0.1:8300", Description: "Server tunnel address"},
	{Key: keyServerTunnelCADir, Flag: toFlag(keyServerTunnelCADir), Default: "/var/lib/otterscale/ca", Description: "Directory for persistent CA certificate and key"},
	{Key: keyServerKeycloakRealmURL, Flag: toFlag(keyServerKeycloakRealmURL), Default: "", Description: "Server keycloak realm url (required)"},
	{Key: keyServerKeycloakClientID, Flag: toFlag(keyServerKeycloakClientID), Default: "otterscale-server", Description: "Server keycloak client id"},
	{Key: keyServerExternalURL, Flag: toFlag(keyServerExternalURL), Default: "", Description: "Externally reachable server URL for agent connections (required for manifest generation)"},
	{Key: keyServerExternalTunnelURL, Flag: toFlag(keyServerExternalTunnelURL), Default: "", Description: "Externally reachable tunnel URL for agent tunnel connections (required for manifest generation)"},
}

// AgentOptions defines the configuration entries available in agent
// mode.
var AgentOptions = []Option{
	{Key: keyAgentCluster, Flag: toFlag(keyAgentCluster), Default: "default", Description: "Agent cluster"},
	{Key: keyAgentServerURL, Flag: toFlag(keyAgentServerURL), Default: "http://127.0.0.1:8299", Description: "Agent control-plane server url"},
	{Key: keyAgentTunnelServerURL, Flag: toFlag(keyAgentTunnelServerURL), Default: "https://127.0.0.1:8300", Description: "Agent tunnel server url"},
	{Key: keyAgentBootstrap, Flag: toFlag(keyAgentBootstrap), Default: true, Description: "Run Layer 0 bootstrap on startup (install FluxCD + Module CRD)"},
}

// toFlag converts a viper key like "server.tunnel.key_seed" into a
// CLI flag like "tunnel-key-seed" by lower-casing, replacing dots and
// underscores with hyphens, and stripping the "server-" or "agent-"
// prefix.
func toFlag(key string) string {
	flag := strings.ToLower(key)
	flag = strings.ReplaceAll(flag, ".", "-")
	flag = strings.ReplaceAll(flag, "_", "-")
	flag = strings.TrimPrefix(flag, "server-")
	flag = strings.TrimPrefix(flag, "agent-")
	return flag
}
