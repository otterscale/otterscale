package harbor

import (
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

// ProvideHarborClient returns a Harbor API client if the harbor URL
// is configured, or nil otherwise. Returning the interface type
// directly (rather than *Client) allows Wire to inject nil when
// Harbor integration is disabled.
func ProvideHarborClient(conf *config.Config) core.HarborClient {
	harborURL := conf.ServerHarborURL()
	if harborURL == "" {
		return nil
	}
	return NewClient(harborURL)
}
