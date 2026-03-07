package core

import "testing"

func TestIsAllowedPrometheusPath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		// Allowed read-only query paths.
		{"/api/v1/query", true},
		{"/api/v1/query?query=up", true},
		{"/api/v1/query_range", true},
		{"/api/v1/query_range?query=up&start=0&end=1&step=15s", true},
		{"/api/v1/labels", true},
		{"/api/v1/label/job/values", true},
		{"/api/v1/label/__name__/values", true},
		{"/api/v1/series", true},
		{"/api/v1/series?match[]=up", true},
		{"/api/v1/metadata", true},
		{"/api/v1/targets", true},
		{"/api/v1/targets/metadata", true},
		{"/api/v1/status/config", true},
		{"/api/v1/status/runtimeinfo", true},

		// Disallowed admin / mutating paths.
		{"/api/v1/admin/tsdb/delete_series", false},
		{"/api/v1/admin/tsdb/clean_tombstones", false},
		{"/api/v1/admin/tsdb/snapshot", false},
		{"/-/reload", false},
		{"/-/quit", false},
		{"/api/v1/write", false},
		{"/", false},
		{"", false},
		{"/api/v2/query", false},
		{"/random/path", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := IsAllowedPrometheusPath(tt.path); got != tt.want {
				t.Errorf("IsAllowedPrometheusPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
