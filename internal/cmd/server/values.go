package server

import (
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

// ProvideAgentValuesConfig gathers the deployment-wide facts a rendered agent
// values file needs.
//
// It never fails. This runs while the server is being constructed, and the
// chart supplying these settings ships on its own cadence, so rejecting an
// unconfigured value here would take out tunnels, proxying and RBAC for the
// sake of one procedure. The use case reports what is missing instead.
func ProvideAgentValuesConfig(conf *config.Config) *core.AgentValuesConfig {
	cfg := &core.AgentValuesConfig{
		ExternalURL:     conf.ServerExternalURL(),
		TunnelServerURL: conf.ServerExternalTunnelURL(),
		HarborURL:       conf.ServerHarborURL(),
	}

	// A configured CA is the signal that agents need one at all; without it
	// the rendered values carry no CA settings. The key is not derived from
	// that path: it is the one the operator gave `kubectl create secret
	// --from-file` on the joining cluster, which is a separate fact.
	if conf.ServerTrustedCAFile() != "" {
		cfg.TrustedCASecret = core.TrustedCASecretName
		cfg.TrustedCAKey = core.DefaultTrustedCAKey
	}

	return cfg
}
