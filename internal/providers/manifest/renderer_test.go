package manifest

import (
	"strings"
	"testing"

	"github.com/otterscale/otterscale/internal/core"
)

func TestRenderAgentManifestRancherProjectEnvironment(t *testing.T) {
	params := &core.ManifestParams{
		Cluster:   "cluster",
		UserName:  "admin@example.com",
		Image:     "ghcr.io/otterscale/otterscale:test",
		ServerURL: "https://server.example.com",
		TunnelURL: "https://tunnel.example.com",
	}
	renderer := NewRenderer()

	withoutProject, err := renderer.RenderAgentManifest(params)
	if err != nil {
		t.Fatalf("RenderAgentManifest without Project: %v", err)
	}
	if strings.Contains(withoutProject, "OTTERSCALE_AGENT_RANCHER_PROJECT_ID") {
		t.Fatal("empty Project ID must not emit an environment variable")
	}

	params.RancherProjectID = "local:p-test"
	withProject, err := renderer.RenderAgentManifest(params)
	if err != nil {
		t.Fatalf("RenderAgentManifest with Project: %v", err)
	}
	want := "- name: OTTERSCALE_AGENT_RANCHER_PROJECT_ID\n              value: \"local:p-test\""
	if !strings.Contains(withProject, want) {
		t.Fatalf("manifest does not contain Project environment variable:\n%s", withProject)
	}
}
