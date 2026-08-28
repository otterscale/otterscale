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

// EnrolmentOptions defines the entries that configure the enrolment
// secret. They are part of ServerOptions and are also bound on their
// own by the enrolment-token subcommand.
var EnrolmentOptions = []Option{
	{Key: keyServerEnrolmentSecret, Flag: toFlag(keyServerEnrolmentSecret), Default: "", Description: "Root secret used to issue and verify agent enrolment tokens (required)"},
	{Key: keyServerEnrolmentSecretFile, Flag: toFlag(keyServerEnrolmentSecretFile), Default: "", Description: "Path to a file holding the enrolment secret; takes precedence over --enrolment-secret"},
}

// ServerOptions defines the configuration entries available in server
// mode. Each entry is registered as a viper default and a CLI flag.
var ServerOptions = append([]Option{
	{Key: keyServerAddress, Flag: toFlag(keyServerAddress), Default: ":8299", Description: "Server listen address"},
	{Key: keyServerAllowedOrigins, Flag: toFlag(keyServerAllowedOrigins), Default: []string{}, Description: "Server allowed origins"},
	{Key: keyServerTunnelAddress, Flag: toFlag(keyServerTunnelAddress), Default: "127.0.0.1:8300", Description: "Server tunnel address"},
	{Key: keyServerKeycloakRealmURL, Flag: toFlag(keyServerKeycloakRealmURL), Default: "", Description: "Server keycloak realm url (required)"},
	{Key: keyServerKeycloakClientID, Flag: toFlag(keyServerKeycloakClientID), Default: "otterscale-server", Description: "Server keycloak client id"},
	{Key: keyServerExternalTunnelURL, Flag: toFlag(keyServerExternalTunnelURL), Default: "", Description: "Externally reachable tunnel URL advertised to agents"},
}, EnrolmentOptions...)

// AgentOptions defines the configuration entries available in agent
// mode.
var AgentOptions = []Option{
	{Key: keyAgentCluster, Flag: toFlag(keyAgentCluster), Default: "default", Description: "Agent cluster"},
	{Key: keyAgentServerURL, Flag: toFlag(keyAgentServerURL), Default: "http://127.0.0.1:8299", Description: "Agent control-plane server url"},
	{Key: keyAgentTunnelServerURL, Flag: toFlag(keyAgentTunnelServerURL), Default: "https://127.0.0.1:8300", Description: "Agent tunnel server url"},
	{Key: keyAgentProxyPrometheusURL, Flag: toFlag(keyAgentProxyPrometheusURL), Default: "http://otterscale-prometheus-kube-prometheus.monitoring.svc:9090", Description: "In-cluster Prometheus URL for the metrics proxy"},
	{Key: keyAgentEnrolmentToken, Flag: toFlag(keyAgentEnrolmentToken), Default: "", Description: "Enrolment token for this cluster, issued by `otterscale enrolment-token` (required)"},
	{Key: keyAgentEnrolmentTokenFile, Flag: toFlag(keyAgentEnrolmentTokenFile), Default: "", Description: "Path to a file holding the enrolment token; takes precedence over --enrolment-token"},
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
