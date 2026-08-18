package manifest

import (
	"io"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"

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

func TestRancherProjectName(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		id      string
		want    string
		wantErr bool
	}{
		{name: "empty", id: "", want: ""},
		{name: "selected", id: "c-m-abcde:p-vwxyz", want: "p-vwxyz"},
		{name: "invalid", id: "c-m-abcde:p-one:p-two", wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rancherProjectName(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("rancherProjectName(%q) error = %v, wantErr=%v", tt.id, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("rancherProjectName(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

//nolint:gocyclo,funlen // This table-driven test intentionally verifies the complete rendered security contract.
func TestRenderAgentManifestRancherSecurity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		projectID     string
		projectName   string
		wantManageNS  bool
		wantProjectID string
	}{
		{
			name:          "selected Project is scoped",
			projectID:     "c-m-abcde:p-vwxyz",
			projectName:   "p-vwxyz",
			wantManageNS:  true,
			wantProjectID: "c-m-abcde:p-vwxyz",
		},
		{name: "empty Project keeps only guarded updatepsa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := renderAgentManifestForTest(t, tt.projectID)
			docs := decodeManifestForTest(t, manifest)

			wantOrder := []string{
				"ValidatingAdmissionPolicy/otterscale-workspace-rancher-project",
				"ValidatingAdmissionPolicyBinding/otterscale-workspace-rancher-project",
				"ValidatingAdmissionPolicy/otterscale-workspace-namespace-security",
				"ValidatingAdmissionPolicyBinding/otterscale-workspace-namespace-security",
				"ClusterRole/otterscale-tenant-operator-rancher-webhook",
				"ClusterRoleBinding/otterscale-tenant-operator-rancher-webhook",
			}
			if len(docs) < len(wantOrder) {
				t.Fatalf("manifest has %d documents, want at least %d", len(docs), len(wantOrder))
			}
			counts := make(map[string]int)
			for _, doc := range docs {
				counts[doc.GetKind()+"/"+doc.GetName()]++
			}
			for i, want := range wantOrder {
				got := docs[i].GetKind() + "/" + docs[i].GetName()
				if got != want {
					t.Fatalf("document %d = %q, want %q", i, got, want)
				}
				if counts[want] != 1 {
					t.Fatalf("manifest contains %d copies of %q, want 1", counts[want], want)
				}
			}

			role := docs[4]
			rules, found, err := unstructured.NestedSlice(role.Object, "rules")
			if err != nil || !found || len(rules) != 1 {
				t.Fatalf("ClusterRole rules = %#v, found=%v, err=%v", rules, found, err)
			}
			rule := rules[0].(map[string]any)
			if got := strings.Join(stringSliceForTest(t, rule["apiGroups"]), ","); got != "management.cattle.io" {
				t.Fatalf("apiGroups = %q, want management.cattle.io", got)
			}
			if got := strings.Join(stringSliceForTest(t, rule["resources"]), ","); got != "projects" {
				t.Fatalf("resources = %q, want projects", got)
			}
			gotVerbs := stringSliceForTest(t, rule["verbs"])
			wantVerbs := "updatepsa"
			if tt.wantManageNS {
				wantVerbs += ",manage-namespaces"
			}
			if got := strings.Join(gotVerbs, ","); got != wantVerbs {
				t.Fatalf("verbs = %q, want %q", got, wantVerbs)
			}
			resourceNames, hasResourceNames := rule["resourceNames"]
			if tt.projectName == "" {
				if hasResourceNames {
					t.Fatalf("empty Project emitted resourceNames: %#v", resourceNames)
				}
			} else if got := stringSliceForTest(t, resourceNames); len(got) != 1 || got[0] != tt.projectName {
				t.Fatalf("resourceNames = %v, want [%s]", got, tt.projectName)
			}

			binding := docs[5]
			subjects, found, err := unstructured.NestedSlice(binding.Object, "subjects")
			if err != nil || !found || len(subjects) != 1 {
				t.Fatalf("ClusterRoleBinding subjects = %#v, found=%v, err=%v", subjects, found, err)
			}
			subject := subjects[0].(map[string]any)
			if subject["kind"] != "ServiceAccount" || subject["name"] != "tenant-operator-controller-manager" || subject["namespace"] != "otterscale-system" {
				t.Fatalf("ClusterRoleBinding subject = %#v", subject)
			}
			roleName, found, err := unstructured.NestedString(binding.Object, "roleRef", "name")
			if err != nil || !found || roleName != "otterscale-tenant-operator-rancher-webhook" {
				t.Fatalf("ClusterRoleBinding roleRef.name = %q, found=%v, err=%v", roleName, found, err)
			}
			workspaceRules, found, err := unstructured.NestedSlice(docs[0].Object, "spec", "matchConstraints", "resourceRules")
			if err != nil || !found || len(workspaceRules) != 1 {
				t.Fatalf("Workspace policy resourceRules = %#v, found=%v, err=%v", workspaceRules, found, err)
			}
			workspaceRule := workspaceRules[0].(map[string]any)
			if got := strings.Join(stringSliceForTest(t, workspaceRule["resources"]), ","); got != "workspaces,workspaces/status" {
				t.Fatalf("Workspace policy resources = %q", got)
			}

			securityPrefix, _, found := strings.Cut(manifest, "kind: Namespace")
			if !found {
				t.Fatal("manifest does not contain the agent Namespace")
			}
			if strings.Count(securityPrefix, "failurePolicy: Fail") != 2 || strings.Count(securityPrefix, "validationActions: [Deny]") != 2 {
				t.Fatal("both admission policies must fail closed with Deny bindings")
			}
			if tt.wantProjectID != "" && strings.Count(securityPrefix, tt.wantProjectID) < 2 {
				t.Fatal("both policies must contain the exact Rancher Project ID")
			}
			for _, value := range []string{
				"pod-security.kubernetes.io/enforce'] == 'baseline'",
				"pod-security.kubernetes.io/warn'] == 'restricted'",
				"pod-security.kubernetes.io/audit'] == 'restricted'",
				"pod-security.kubernetes.io/enforce-version",
				"The Workspace controller owner is immutable",
				"system:serviceaccount:otterscale-system:tenant-operator-controller-manager",
			} {
				if !strings.Contains(securityPrefix, value) {
					t.Fatalf("security prerequisites missing %q", value)
				}
			}
		})
	}

	invalid := renderAgentManifestParamsForTest("local:bad:project")
	if _, err := NewRenderer().RenderAgentManifest(invalid); err == nil {
		t.Fatal("invalid Rancher Project ID must be rejected")
	}
}

func renderAgentManifestForTest(t *testing.T, projectID string) string {
	t.Helper()
	manifest, err := NewRenderer().RenderAgentManifest(renderAgentManifestParamsForTest(projectID))
	if err != nil {
		t.Fatalf("RenderAgentManifest: %v", err)
	}
	return manifest
}

func renderAgentManifestParamsForTest(projectID string) *core.ManifestParams {
	return &core.ManifestParams{
		Cluster:          "cluster",
		UserName:         "admin@example.com",
		Image:            "ghcr.io/otterscale/otterscale:test",
		ServerURL:        "https://server.example.com",
		TunnelURL:        "https://tunnel.example.com",
		RancherProjectID: projectID,
	}
}

func decodeManifestForTest(t *testing.T, manifest string) []*unstructured.Unstructured {
	t.Helper()
	decoder := utilyaml.NewYAMLOrJSONDecoder(strings.NewReader(manifest), 4096)
	var docs []*unstructured.Unstructured
	for {
		object := &unstructured.Unstructured{}
		err := decoder.Decode(object)
		if err == io.EOF {
			return docs
		}
		if err != nil {
			t.Fatalf("decode manifest: %v", err)
		}
		if len(object.Object) != 0 {
			docs = append(docs, object)
		}
	}
}

func stringSliceForTest(t *testing.T, value any) []string {
	t.Helper()
	items, ok := value.([]any)
	if !ok {
		t.Fatalf("value %#v is not a slice", value)
	}
	result := make([]string, len(items))
	for i, item := range items {
		result[i], ok = item.(string)
		if !ok {
			t.Fatalf("value %#v contains a non-string", value)
		}
	}
	return result
}
