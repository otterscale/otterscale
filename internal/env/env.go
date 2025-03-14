package env

import "os"

// Environment variables
const (
	OPENHDC_IN_CLUSTER = "OPENHDC_IN_CLUSTER" //nolint:stylecheck
)

// Environment variables for MAAS configuration
const (
	OPENHDC_MAAS_API_URL     = "OPENHDC_MAAS_API_URL"     //nolint:stylecheck
	OPENHDC_MAAS_API_KEY     = "OPENHDC_MAAS_API_KEY"     //nolint:stylecheck
	OPENHDC_MAAS_API_VERSION = "OPENHDC_MAAS_API_VERSION" //nolint:stylecheck
)

// Environment variables for JUJU configuration
const (
	OPENHDC_JUJU_CONTROLLER_ADDRESSES = "OPENHDC_JUJU_CONTROLLER_ADDRESSES" //nolint:stylecheck
	OPENHDC_JUJU_USERNAME             = "OPENHDC_JUJU_USERNAME"             //nolint:stylecheck
	OPENHDC_JUJU_PASSWORD             = "OPENHDC_JUJU_PASSWORD"             //nolint:stylecheck,gosec
	OPENHDC_JUJU_CACERT_PATH          = "OPENHDC_JUJU_CACERT_PATH"          //nolint:stylecheck
)

// GetOrDefault returns the value of the environment variable if set, otherwise returns the default value.
func GetOrDefault(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}
