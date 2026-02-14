package manifest

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/pki"
)

// ProvideAgentManifestConfig is a Wire provider that extracts the
// external URLs from the server configuration and derives an HMAC key
// for signing stateless manifest tokens. The HMAC key is derived from
// the CA's private key via HKDF, so it is deterministic for the same
// CA and survives restarts without separate persistence.
func ProvideAgentManifestConfig(conf *config.Config, ca *pki.CA) (core.AgentManifestConfig, error) {
	hmacKey, err := ca.DeriveHMACKey("manifest-token")
	if err != nil {
		return core.AgentManifestConfig{}, fmt.Errorf("derive HMAC key: %w", err)
	}
	return core.AgentManifestConfig{
		ServerURL: conf.ServerExternalURL(),
		TunnelURL: conf.ServerExternalTunnelURL(),
		HMACKey:   hmacKey,
	}, nil
}
