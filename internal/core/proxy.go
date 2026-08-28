package core

import (
	"path"
	"strings"
)

// allowedPrometheusPathPrefixes covers read-only query endpoints only;
// administrative paths (/api/v1/admin/*, /-/*) are excluded.
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

// AllowedPrometheusPath matches by prefix, so "/api/v1/query_range" and
// "/api/v1/label/job/values" are accepted while "/api/v1/admin/tsdb/delete_series"
// is rejected.
//
// Normalizing and checking are deliberately one operation returning one value:
// prefix matching means "/api/v1/query/../../-/reload" satisfies the allowlist
// while addressing an admin endpoint, and handing back the normalized path is
// what stops a caller from checking one path and forwarding another.
func AllowedPrometheusPath(raw string) (normalized string, allowed bool) {
	normalized = normalizePath(raw)
	for _, prefix := range allowedPrometheusPathPrefixes {
		if strings.HasPrefix(normalized, prefix) {
			return normalized, true
		}
	}
	return normalized, false
}

// normalizePath resolves "." and ".." segments and collapses repeated
// separators, so that the path being checked is the one an upstream
// router will act on.
//
// A trailing slash is preserved. path.Clean drops it, and a check meant
// to reject dangerous paths should not quietly rewrite legitimate ones.
func normalizePath(raw string) string {
	cleaned := path.Clean(raw)
	if cleaned != "/" && strings.HasSuffix(raw, "/") {
		cleaned += "/"
	}
	return cleaned
}
