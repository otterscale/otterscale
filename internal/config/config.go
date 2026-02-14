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

// Config wraps a viper instance and provides typed accessors for every
// configuration key. Create one via New().
type Config struct {
	v *viper.Viper
}

// New initialises a Config by loading values from the config file,
// environment variables, and compiled defaults (in that priority
// order; CLI flags, bound later via BindFlags, take highest priority).
func New() (*Config, error) {
	v := viper.New()

	// Register compiled defaults for all known options.
	for _, o := range ServerOptions {
		v.SetDefault(o.Key, o.Default)
	}
	for _, o := range AgentOptions {
		v.SetDefault(o.Key, o.Default)
	}

	// Attempt to load a config file from the current directory or
	// the system-wide location.
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/otterscale/")

	if err := v.ReadInConfig(); err != nil {
		var notFoundErr viper.ConfigFileNotFoundError
		if !(errors.As(err, &notFoundErr) || errors.Is(err, os.ErrNotExist)) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Environment variables are prefixed with OTTERSCALE_ and use
	// underscores in place of dots (e.g. OTTERSCALE_SERVER_ADDRESS).
	v.SetEnvPrefix("OTTERSCALE")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &Config{v: v}, nil
}

// BindFlags registers CLI flags for the given option slice and binds
// them to the underlying viper keys so that flag values override file
// and environment sources.
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

// ---------------------------------------------------------------------------
// Server-mode accessors
// ---------------------------------------------------------------------------

// ServerAddress returns the HTTP listen address for the server.
func (c *Config) ServerAddress() string {
	return c.v.GetString(keyServerAddress)
}

// ServerAllowedOrigins returns the list of allowed CORS origins.
func (c *Config) ServerAllowedOrigins() []string {
	return c.v.GetStringSlice(keyServerAllowedOrigins)
}

// ServerTunnelAddress returns the listen address for the chisel tunnel
// server.
func (c *Config) ServerTunnelAddress() string {
	return c.v.GetString(keyServerTunnelAddress)
}

// ServerTunnelCADir returns the directory path where the CA
// certificate and private key are persisted. On first startup the
// server generates a new CA and writes the material to this
// directory; subsequent restarts load the existing CA.
func (c *Config) ServerTunnelCADir() string {
	return c.v.GetString(keyServerTunnelCADir)
}

// ServerKeycloakRealmURL returns the Keycloak realm issuer URL used
// for OIDC token verification.
func (c *Config) ServerKeycloakRealmURL() string {
	return c.v.GetString(keyServerKeycloakRealmURL)
}

// ServerKeycloakClientID returns the Keycloak client ID expected in
// the "aud" claim of incoming tokens.
func (c *Config) ServerKeycloakClientID() string {
	return c.v.GetString(keyServerKeycloakClientID)
}

// ServerExternalURL returns the externally reachable server URL that
// agents use to connect to the control plane.
func (c *Config) ServerExternalURL() string {
	return c.v.GetString(keyServerExternalURL)
}

// ServerExternalTunnelURL returns the externally reachable tunnel URL
// that agents use to establish reverse tunnels.
func (c *Config) ServerExternalTunnelURL() string {
	return c.v.GetString(keyServerExternalTunnelURL)
}

// ---------------------------------------------------------------------------
// Agent-mode accessors
// ---------------------------------------------------------------------------

// AgentCluster returns the cluster name this agent registers under.
func (c *Config) AgentCluster() string {
	return c.v.GetString(keyAgentCluster)
}

// AgentServerURL returns the fleet server URL the agent registers
// against.
func (c *Config) AgentServerURL() string {
	return c.v.GetString(keyAgentServerURL)
}

// AgentTunnelServerURL returns the chisel tunnel server URL the agent
// connects to.
func (c *Config) AgentTunnelServerURL() string {
	return c.v.GetString(keyAgentTunnelServerURL)
}

// AgentBootstrap returns whether the agent should run the Layer 0
// bootstrap process on startup, installing FluxCD and the Module CRD.
func (c *Config) AgentBootstrap() bool {
	return c.v.GetBool(keyAgentBootstrap)
}
