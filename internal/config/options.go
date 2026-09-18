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
	{Key: keyServerKeycloakRealmURL, Flag: toFlag(keyServerKeycloakRealmURL), Default: "", Description: "Server keycloak realm url (required)"},
	{Key: keyServerKeycloakClientID, Flag: toFlag(keyServerKeycloakClientID), Default: "otterscale-server", Description: "Server keycloak client id"},
	{Key: keyServerExternalURL, Flag: toFlag(keyServerExternalURL), Default: "", Description: "Externally reachable URL of this API, including the gateway's path prefix (e.g. https://otterscale.example.com/api/)"},
	{Key: keyServerExternalTunnelURL, Flag: toFlag(keyServerExternalTunnelURL), Default: "", Description: "Externally reachable tunnel URL advertised to agents"},
	{Key: keyServerHarborURL, Flag: toFlag(keyServerHarborURL), Default: "", Description: "Externally reachable Harbor URL; the module repositories are derived from it"},
	{Key: keyServerJoinSecret, Flag: toFlag(keyServerJoinSecret), Default: "", Description: "Root secret used to issue and verify agent join tokens (required)"},
	{Key: keyServerJoinSecretFile, Flag: toFlag(keyServerJoinSecretFile), Default: "", Description: "Path to a file holding the join secret; takes precedence over --join-secret"},
	{Key: keyServerTrustedCAFile, Flag: toFlag(keyServerTrustedCAFile), Default: "", Description: "Path to the CA certificate agents must trust to reach this server"},
	{Key: keyServerHarborAdminPassword, Flag: toFlag(keyServerHarborAdminPassword), Default: "", Description: "Harbor admin password used to provision per-cluster robot accounts"},
	{Key: keyServerHarborAdminPasswordFile, Flag: toFlag(keyServerHarborAdminPasswordFile), Default: "", Description: "Path to a file holding the Harbor admin password; takes precedence over --harbor-admin-password"},
}

// AgentOptions defines the configuration entries available in agent
// mode.
var AgentOptions = []Option{
	{Key: keyAgentCluster, Flag: toFlag(keyAgentCluster), Default: "default", Description: "Agent cluster"},
	{Key: keyAgentServerURL, Flag: toFlag(keyAgentServerURL), Default: "http://127.0.0.1:8299", Description: "Agent control-plane server url"},
	{Key: keyAgentTunnelServerURL, Flag: toFlag(keyAgentTunnelServerURL), Default: "https://127.0.0.1:8300", Description: "Agent tunnel server url"},
	{Key: keyAgentProxyPrometheusURL, Flag: toFlag(keyAgentProxyPrometheusURL), Default: "http://prometheus-stack-kube-prom-prometheus.monitoring.svc:9090", Description: "In-cluster Prometheus URL for the metrics proxy"},
	{Key: keyAgentJoinToken, Flag: toFlag(keyAgentJoinToken), Default: "", Description: "Join token for this cluster, issued by the control plane (required)"},
	{Key: keyAgentJoinTokenFile, Flag: toFlag(keyAgentJoinTokenFile), Default: "", Description: "Path to a file holding the join token; takes precedence over --join-token"},
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
