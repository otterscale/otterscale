package harbor

import (
	"github.com/otterscale/otterscale/internal/config"
)

// ProvideHarborClient builds the client from configuration. An unset URL is
// rejected when agent values are issued, not here.
func ProvideHarborClient(conf *config.Config) *Client {
	return NewClient(conf.ServerHarborURL(), conf.ServerHarborAdminPassword)
}
