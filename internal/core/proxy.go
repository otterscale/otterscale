package core

import "strings"

// allowedPrometheusPathPrefixes lists the Prometheus HTTP API paths
// that the metrics proxy permits. Only read-only query endpoints are
// included; administrative paths (/api/v1/admin/*, /-/*) are
// explicitly excluded.
var allowedPrometheusPathPrefixes = []string{
	"/api/v1/query",
	"/api/v1/query_range",
	"/api/v1/labels",
	"/api/v1/label/",
	"/api/v1/series",
	"/api/v1/metadata",
	"/api/v1/targets",
	"/api/v1/status/",
}

// IsAllowedPrometheusPath reports whether path is a permitted
// Prometheus query endpoint. It uses a prefix match so that paths
// like "/api/v1/query_range" and "/api/v1/label/job/values" are
// accepted, while "/api/v1/admin/tsdb/delete_series" is rejected.
func IsAllowedPrometheusPath(path string) bool {
	for _, prefix := range allowedPrometheusPathPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
