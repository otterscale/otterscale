package values

import (
	"fmt"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"

	"github.com/otterscale/otterscale/internal/core"
)

// The agent chart's own defaults are example addresses, not sane fallbacks, so
// the render has to override every one: an omitted key publishes it.
const (
	placeholderAddress      = "192.0.2.1"
	placeholderInferenceURL = "https://inference.example.com"
	placeholderRobot        = "robot$devel"
)

// Fixture versions, so a chart bump leaves the goldens alone. All three differ
// and all three reach the output, so a mixed-up wiring fails here.
func testCharts() core.ChartVersions {
	return core.ChartVersions{
		AgentFlux: "0.1.2",
		Agent:     "1.2.3",
		Flux:      "2.3.4",
	}
}

// fullValues is a fully configured render. The token and secret in it are
// fixtures, not credentials.
func fullValues() *core.AgentValues {
	return &core.AgentValues{ //nolint:gosec // G101: fixtures mirroring a real values file
		Cluster:           "cluster",
		JoinToken:         "LEGbfEh0nknqtAtTbEmf8rJm95UQEzGxxJKtwDK01-I",
		ServerURL:         "https://192.168.196.222/api/",
		TunnelServerURL:   "https://192.168.196.221:30300",
		ClusterAdminUsers: []string{"2ae606b7-ff7c-4c47-be76-8371963ee3c5"},
		ClusterInfo: core.AgentClusterInfo{
			ExternalAddress: "192.168.196.222",
			NodePortRange:   "30000-32767",
			InferenceURL:    "https://192.168.196.223",
		},
		Harbor: core.HarborSettings{
			URL:              "https://192.168.196.222:8443",
			ModulesRepoURL:   "oci://192.168.196.222:8443/modules",
			OperatorsRepoURL: "oci://192.168.196.222:8443/operators",
			Robot: core.HarborRobotCredentials{ //nolint:gosec // G101: a test fixture, not a credential
				Name:   "robot$cluster",
				Secret: "OtrzmR4w2cuoBbhL1pnhxp8G9cwgOM5c0N0",
			},
		},
		TrustedCASecret: core.TrustedCASecretName,
		TrustedCAKey:    core.DefaultTrustedCAKey,
	}
}

const wantFull = `# otterscale-agent-flux chart version: 0.1.2
repositories:
  modules:
    url: oci://192.168.196.222:8443/modules
    labels:
      tenant.otterscale.io/from-harbor: "true"
  operators:
    url: oci://192.168.196.222:8443/operators
agent:
  version: 1.2.3
  values:
    agent:
      serverURL: https://192.168.196.222/api/
      tunnelServerURL: https://192.168.196.221:30300
      cluster: cluster
      joinToken: LEGbfEh0nknqtAtTbEmf8rJm95UQEzGxxJKtwDK01-I
    clusterAdmin:
      enabled: true
      users:
        - 2ae606b7-ff7c-4c47-be76-8371963ee3c5
    clusterInfo:
      enabled: true
      externalAddress: 192.168.196.222
      nodePortRange: 30000-32767
      inferenceURL: https://192.168.196.223
    tenantOperator:
      enabled: true
      harbor:
        url: https://192.168.196.222:8443
        robot:
          name: robot$cluster
          secret: OtrzmR4w2cuoBbhL1pnhxp8G9cwgOM5c0N0
    trustedCA:
      secretName: otterscale-ca
      key: ca.crt
  valuesFrom:
    - kind: ConfigMap
      name: otterscale-agent-values
      valuesKey: values.yaml
      optional: true
flux:
  version: 2.3.4
  values:
    flux2:
      sourceController:
        extraEnv:
          - name: SSL_CERT_DIR
            value: /etc/ssl/certs:/etc/otterscale/ca
        volumeMounts:
          - mountPath: /etc/otterscale/ca
            name: trusted-ca
            readOnly: true
        volumes:
          - name: trusted-ca
            secret:
              secretName: otterscale-ca
  valuesFrom:
    - kind: ConfigMap
      name: flux-values
      valuesKey: values.yaml
      optional: true
`

// A publicly signed server hands an agent nothing to trust, so the whole CA
// plumbing drops out. The override ConfigMaps stay either way.
const wantWithoutTrustedCA = `# otterscale-agent-flux chart version: 0.1.2
repositories:
  modules:
    url: oci://192.168.196.222:8443/modules
    labels:
      tenant.otterscale.io/from-harbor: "true"
  operators:
    url: oci://192.168.196.222:8443/operators
agent:
  version: 1.2.3
  values:
    agent:
      serverURL: https://192.168.196.222/api/
      tunnelServerURL: https://192.168.196.221:30300
      cluster: cluster
      joinToken: LEGbfEh0nknqtAtTbEmf8rJm95UQEzGxxJKtwDK01-I
    clusterAdmin:
      enabled: true
      users:
        - 2ae606b7-ff7c-4c47-be76-8371963ee3c5
    clusterInfo:
      enabled: true
      externalAddress: 192.168.196.222
      nodePortRange: 30000-32767
      inferenceURL: https://192.168.196.223
    tenantOperator:
      enabled: true
      harbor:
        url: https://192.168.196.222:8443
        robot:
          name: robot$cluster
          secret: OtrzmR4w2cuoBbhL1pnhxp8G9cwgOM5c0N0
  valuesFrom:
    - kind: ConfigMap
      name: otterscale-agent-values
      valuesKey: values.yaml
      optional: true
flux:
  version: 2.3.4
  valuesFrom:
    - kind: ConfigMap
      name: flux-values
      valuesKey: values.yaml
      optional: true
`

func TestRenderer_Render(t *testing.T) {
	tests := []struct {
		name   string
		values func() *core.AgentValues
		want   string
	}{
		{
			name:   "everything configured",
			values: fullValues,
			want:   wantFull,
		},
		{
			name: "no trusted CA",
			values: func() *core.AgentValues {
				v := fullValues()
				v.TrustedCASecret = ""
				v.TrustedCAKey = ""
				return v
			},
			want: wantWithoutTrustedCA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRenderer(testCharts()).Render(tt.values())
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("rendered output differs:\n%s", diffLines(tt.want, got))
			}
		})
	}
}

// The placeholder trap: Helm deep-merges, so an omitted key keeps the chart's
// example address and the dashboard advertises it.
func TestRenderer_RenderWritesEmptyInferenceURL(t *testing.T) {
	v := fullValues()
	v.ClusterInfo.InferenceURL = ""

	got, err := NewRenderer(testCharts()).Render(v)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(got, `inferenceURL: ""`) {
		t.Error("an empty inference URL must be written out, not omitted")
	}
	if strings.Contains(got, placeholderInferenceURL) {
		t.Errorf("the chart placeholder %q must not survive", placeholderInferenceURL)
	}
}

// The whole set at once, as the install-time grep does.
func TestRenderer_RenderOverridesEveryPlaceholder(t *testing.T) {
	got, err := NewRenderer(testCharts()).Render(fullValues())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	for _, placeholder := range []string{placeholderAddress, placeholderInferenceURL, placeholderRobot} {
		if strings.Contains(got, placeholder) {
			t.Errorf("the chart placeholder %q survived the render", placeholder)
		}
	}
}

// Why this renders through an encoder rather than a template:
// ValidateClusterName covers the cluster name, but the cluster-admin
// identities are free-form strings that land in a list.
func TestRenderer_RenderQuotesHostileValues(t *testing.T) {
	hostile := []string{
		`quote"inside`,
		"line\nbreak",
		"colon: space",
		"#hash",
		"-leading-dash",
		"*anchor",
		"{braces}",
		"",
	}

	v := fullValues()
	v.ClusterAdminUsers = hostile
	v.Harbor.Robot.Secret = "secret: with #both\nand a newline"

	got, err := NewRenderer(testCharts()).Render(v)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	var parsed struct {
		Agent struct {
			Values struct {
				ClusterAdmin struct {
					Users []string `yaml:"users"`
				} `yaml:"clusterAdmin"`
				TenantOperator struct {
					Harbor struct {
						Robot struct {
							Secret string `yaml:"secret"`
						} `yaml:"robot"`
					} `yaml:"harbor"`
				} `yaml:"tenantOperator"`
			} `yaml:"values"`
		} `yaml:"agent"`
	}
	if err := yaml.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("the rendered output is not valid YAML: %v\n%s", err, got)
	}

	users := parsed.Agent.Values.ClusterAdmin.Users
	if len(users) != len(hostile) {
		t.Fatalf("users = %q, want %q", users, hostile)
	}
	for i := range users {
		if users[i] != hostile[i] {
			t.Errorf("users[%d] = %q, want %q", i, users[i], hostile[i])
		}
	}
	if got := parsed.Agent.Values.TenantOperator.Harbor.Robot.Secret; got != v.Harbor.Robot.Secret {
		t.Errorf("robot secret = %q, want %q", got, v.Harbor.Robot.Secret)
	}
}

// The RPC and the URL must return identical bytes, so nothing in the output
// may vary between calls — no timestamp, no map iteration order.
func TestRenderer_RenderIsDeterministic(t *testing.T) {
	renderer := NewRenderer(testCharts())

	first, err := renderer.Render(fullValues())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for range 10 {
		again, err := renderer.Render(fullValues())
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		if again != first {
			t.Fatalf("successive renders differ:\n%s", diffLines(first, again))
		}
	}
}

// The pins, checked through a parse rather than the goldens: omitempty means a
// zero-value Renderer drops both keys silently.
func TestRenderer_RenderPinsChartVersions(t *testing.T) {
	charts := testCharts()

	got, err := NewRenderer(charts).Render(fullValues())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	var parsed struct {
		Agent struct {
			Version string `yaml:"version"`
		} `yaml:"agent"`
		Flux struct {
			Version string `yaml:"version"`
		} `yaml:"flux"`
	}
	if err := yaml.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("the rendered output is not valid YAML: %v\n%s", err, got)
	}

	if parsed.Agent.Version != charts.Agent {
		t.Errorf("agent.version = %q, want %q", parsed.Agent.Version, charts.Agent)
	}
	if parsed.Flux.Version != charts.Flux {
		t.Errorf("flux.version = %q, want %q", parsed.Flux.Version, charts.Flux)
	}
}

// The umbrella version cannot be pinned inside the file, so it heads it.
func TestRenderer_RenderHeaderNamesChartVersion(t *testing.T) {
	charts := testCharts()

	got, err := NewRenderer(charts).Render(fullValues())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	wantHeader := chartVersionComment + charts.AgentFlux
	header, _, _ := strings.Cut(got, "\n")
	if header != wantHeader {
		t.Errorf("first line = %q, want %q", header, wantHeader)
	}

	// A comment, so the body must survive it.
	var parsed struct {
		Repositories map[string]any `yaml:"repositories"`
	}
	if err := yaml.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("the header made the output invalid YAML: %v\n%s", err, got)
	}
	if len(parsed.Repositories) == 0 {
		t.Error("the header displaced the document body")
	}
}

// diffLines reports only the differing lines.
func diffLines(want, got string) string {
	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(got, "\n")

	var out strings.Builder
	for i := range max(len(wantLines), len(gotLines)) {
		var wantLine, gotLine string
		if i < len(wantLines) {
			wantLine = wantLines[i]
		}
		if i < len(gotLines) {
			gotLine = gotLines[i]
		}
		if wantLine != gotLine {
			fmt.Fprintf(&out, "line %d:\n  want: %q\n   got: %q\n", i+1, wantLine, gotLine)
		}
	}
	if out.Len() == 0 {
		return "(no line differs; check trailing whitespace)"
	}
	return out.String()
}
