package agent

import (
	"context"
	"strings"
	"testing"

	"k8s.io/client-go/rest"
)

func TestPatch_InvalidVersion(t *testing.T) {
	u := NewUpdater(&rest.Config{})

	tests := []struct {
		name    string
		version string
	}{
		{"empty string", ""},
		{"arbitrary string", "latest"},
		{"image injection", "latest@sha256:abc123"},
		{"path traversal", "../evil"},
		{"special chars", "v1.0.0; rm -rf /"},
		{"incomplete semver", "1.0"},
		{"plain text", "not-a-version"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := u.Patch(context.Background(), tt.version)
			if err == nil {
				t.Errorf("expected error for version %q, got nil", tt.version)
			}
		})
	}
}

func TestPatch_ValidVersion(t *testing.T) {
	// These versions should pass the semver validation but fail
	// later at the Kubernetes client step (no in-cluster config).
	u := NewUpdater(&rest.Config{})

	tests := []struct {
		name    string
		version string
	}{
		{"basic semver", "1.2.3"},
		{"with v prefix", "v1.2.3"},
		{"with prerelease", "v1.2.3-rc.1"},
		{"with build metadata", "v1.2.3+build.42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := u.Patch(context.Background(), tt.version)
			if err == nil {
				// If it succeeds we're probably in-cluster, which is fine.
				return
			}
			// The error should be about creating the kube client,
			// NOT about invalid version.
			if strings.Contains(err.Error(), "invalid server version") {
				t.Errorf("version %q should pass semver validation, got: %v", tt.version, err)
			}
		})
	}
}

func TestImageRef(t *testing.T) {
	got := imageRef("v1.2.3")
	want := "ghcr.io/otterscale/otterscale:v1.2.3"
	if got != want {
		t.Errorf("imageRef(v1.2.3) = %q, want %q", got, want)
	}
}

func TestDetectNamespace_ErrorOutsideCluster(t *testing.T) {
	// Outside a Kubernetes cluster, detectNamespace should return an
	// error instead of silently falling back to "default".
	_, err := detectNamespace()
	if err == nil {
		t.Error("expected error for detectNamespace outside cluster, got nil")
	}
}

