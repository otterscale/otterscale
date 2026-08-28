package core

import "testing"

func TestAllowedPrometheusPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		// Allowed read-only query paths.
		{path: "/api/v1/query", want: true},
		{path: "/api/v1/query?query=up", want: true},
		{path: "/api/v1/query_range", want: true},
		{path: "/api/v1/query_range?query=up&start=0&end=1&step=15s", want: true},
		{path: "/api/v1/labels", want: true},
		{path: "/api/v1/label/job/values", want: true},
		{path: "/api/v1/label/__name__/values", want: true},
		{path: "/api/v1/series", want: true},
		{path: "/api/v1/series?match[]=up", want: true},
		{path: "/api/v1/metadata", want: true},
		{path: "/api/v1/targets", want: true},
		{path: "/api/v1/targets/metadata", want: true},
		{path: "/api/v1/status/config", want: true},
		{path: "/api/v1/status/runtimeinfo", want: true},

		// Disallowed admin / mutating paths.
		{path: "/api/v1/admin/tsdb/delete_series", want: false},
		{path: "/api/v1/admin/tsdb/clean_tombstones", want: false},
		{path: "/api/v1/admin/tsdb/snapshot", want: false},
		{path: "/-/reload", want: false},
		{path: "/-/quit", want: false},
		{path: "/api/v1/write", want: false},
		{path: "/", want: false},
		{name: "empty", path: "", want: false},
		{path: "/api/v2/query", want: false},
		{path: "/random/path", want: false},

		// Traversal: these satisfy an allowed prefix as written, but
		// address something else entirely once resolved.
		{path: "/api/v1/query/../../-/reload", want: false},
		{path: "/api/v1/query/../admin/tsdb/delete_series", want: false},
		{path: "/api/v1/query/../../../etc/passwd", want: false},
		{path: "/api/v1/labels/../../v1/write", want: false},
		{path: "/api/v1/query/..", want: false},

		// Traversal that resolves back inside the allowlist stays
		// allowed: the check is on where the path lands, not on how it
		// was spelled.
		{path: "/api/v1/status/../query", want: true},
		{path: "/api/v1//query", want: true},
	}

	for _, tt := range tests {
		name := tt.name
		if name == "" {
			name = tt.path
		}
		t.Run(name, func(t *testing.T) {
			got, allowed := AllowedPrometheusPath(tt.path)
			if allowed != tt.want {
				t.Errorf("AllowedPrometheusPath(%q) = %q, %v; want allowed=%v", tt.path, got, allowed, tt.want)
			}
		})
	}
}

// TestAllowedPrometheusPathReturnsNormalizedPath pins the property the
// caller depends on: the path handed back is the one that was checked,
// so a caller cannot forward something the allowlist never saw.
func TestAllowedPrometheusPathReturnsNormalizedPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/api/v1/query", "/api/v1/query"},
		{"/api/v1/status/../query", "/api/v1/query"},
		{"/api/v1//query", "/api/v1/query"},
		{"/api/v1/./query", "/api/v1/query"},
		{"/api/v1/query/../../-/reload", "/api/-/reload"},
		{"/api/v1/query/../../../-/reload", "/-/reload"},

		// A trailing slash survives: path.Clean drops it, and a
		// security check must not quietly rewrite a valid request.
		{"/api/v1/label/job/values/", "/api/v1/label/job/values/"},
		{"/", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got, _ := AllowedPrometheusPath(tt.path); got != tt.want {
				t.Errorf("AllowedPrometheusPath(%q) normalized to %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
