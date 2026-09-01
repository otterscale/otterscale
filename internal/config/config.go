package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config wraps a viper instance with typed accessors for every configuration
// key. Create one via New().
type Config struct {
	v *viper.Viper
}

// New loads the config file, environment variables, and compiled defaults, in
// that priority order. CLI flags, bound later via BindFlags, outrank all three.
func New() (*Config, error) {
	v := viper.New()

	for _, o := range ServerOptions {
		v.SetDefault(o.Key, o.Default)
	}
	for _, o := range AgentOptions {
		v.SetDefault(o.Key, o.Default)
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/otterscale/")

	if err := v.ReadInConfig(); err != nil {
		var notFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &notFoundErr) && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// OTTERSCALE_-prefixed, dots replaced by underscores
	// (e.g. OTTERSCALE_SERVER_ADDRESS).
	v.SetEnvPrefix("OTTERSCALE")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &Config{v: v}, nil
}

// BindFlags registers CLI flags and binds them to their viper keys, so flag
// values override the file and environment sources.
func (c *Config) BindFlags(fs *pflag.FlagSet, options []Option) error {
	for _, o := range options {
		switch v := o.Default.(type) {
		case string:
			fs.String(o.Flag, v, o.Description)
		case int:
			fs.Int(o.Flag, v, o.Description)
		case bool:
			fs.Bool(o.Flag, v, o.Description)
		case []string:
			fs.StringSlice(o.Flag, v, o.Description)
		case time.Duration:
			fs.Duration(o.Flag, v, o.Description)
		default:
			return fmt.Errorf("unsupported flag type for key: %s", o.Key)
		}

		if err := c.v.BindPFlag(o.Key, fs.Lookup(o.Flag)); err != nil {
			return fmt.Errorf("failed to bind flag %s: %w", o.Flag, err)
		}
	}

	return nil
}

func (c *Config) ServerAddress() string {
	return c.v.GetString(keyServerAddress)
}

// ServerAllowedOrigins lists the allowed CORS origins.
func (c *Config) ServerAllowedOrigins() []string {
	return c.v.GetStringSlice(keyServerAllowedOrigins)
}

func (c *Config) ServerTunnelAddress() string {
	return c.v.GetString(keyServerTunnelAddress)
}

// ServerKeycloakRealmURL is the issuer URL for OIDC token verification.
func (c *Config) ServerKeycloakRealmURL() string {
	return c.v.GetString(keyServerKeycloakRealmURL)
}

// ServerKeycloakClientID is expected in the "aud" claim of incoming tokens.
func (c *Config) ServerKeycloakClientID() string {
	return c.v.GetString(keyServerKeycloakClientID)
}

// ServerExternalTunnelURL is the address agents dial to establish tunnels.
func (c *Config) ServerExternalTunnelURL() string {
	return c.v.GetString(keyServerExternalTunnelURL)
}

// ServerJoinSecret is the root secret for issuing and verifying agent
// join tokens, read from the configured file when one is set.
func (c *Config) ServerJoinSecret() (string, error) {
	return c.secret(keyServerJoinSecretFile, keyServerJoinSecret)
}

// ServerTrustedCAFile is the path to the CA certificate agents need to verify
// this server. Empty when the server's certificate chains to a public CA.
func (c *Config) ServerTrustedCAFile() string {
	return c.v.GetString(keyServerTrustedCAFile)
}

// AgentCluster is the name this agent registers under.
func (c *Config) AgentCluster() string {
	return c.v.GetString(keyAgentCluster)
}

func (c *Config) AgentServerURL() string {
	return c.v.GetString(keyAgentServerURL)
}

func (c *Config) AgentTunnelServerURL() string {
	return c.v.GetString(keyAgentTunnelServerURL)
}

// AgentProxyPrometheusURL is the in-cluster Prometheus the metrics proxy
// targets. Empty disables the proxy.
func (c *Config) AgentProxyPrometheusURL() string {
	return c.v.GetString(keyAgentProxyPrometheusURL)
}

// AgentJoinToken is what this agent presents when it registers, read from
// the configured file when one is set.
func (c *Config) AgentJoinToken() (string, error) {
	return c.secret(keyAgentJoinTokenFile, keyAgentJoinToken)
}

// secret reads a credential from the file named by fileKey, falling back to the
// inline value at valueKey.
//
// The file form exists because an inline secret reaches the process as a flag
// or an environment variable, visible in the container spec and in
// /proc/<pid>/environ; a mounted file is not. Surrounding whitespace is trimmed
// so a file written with a trailing newline works.
func (c *Config) secret(fileKey, valueKey string) (string, error) {
	path := c.v.GetString(fileKey)
	if path == "" {
		return c.v.GetString(valueKey), nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s from %q: %w", valueKey, path, err)
	}
	return strings.TrimSpace(string(content)), nil
}
