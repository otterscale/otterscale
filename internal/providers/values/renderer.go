// Package values renders the Helm override values that install the otterscale
// agent on a joining cluster, keeping the otterscale-agent-flux chart's schema
// out of the domain layer.
package values

import (
	"fmt"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/otterscale/otterscale/internal/core"
)

// What marks a repository backed by an OtterScale-managed Harbor.
const (
	fromHarborLabelKey   = "tenant.otterscale.io/from-harbor"
	fromHarborLabelValue = "true"
)

// ConfigMaps an operator may drop beside a release to override what this
// renders. Optional: absent, the release is what the control plane issued.
const (
	agentOverrideConfigMap = "otterscale-agent-values"
	fluxOverrideConfigMap  = "flux-values"
	overrideValuesKey      = "values.yaml"
)

// Where both charts mount a trusted CA, and the search path that makes Go
// binaries use it. The system directory stays first so public CAs keep working.
const (
	trustedCADir      = "/etc/otterscale/ca"
	trustedCAVolume   = "trusted-ca"
	sslCertDirEnv     = "SSL_CERT_DIR"
	sslCertDirDefault = "/etc/ssl/certs"
)

// yamlIndent matches the chart sources, so rendered and hand-written files
// diff cleanly.
const yamlIndent = 2

// Renderer turns resolved values into YAML.
type Renderer struct{}

var _ core.AgentValuesRenderer = (*Renderer)(nil)

func NewRenderer() *Renderer {
	return &Renderer{}
}

// Render writes the override values for one joining cluster. Nothing in the
// output varies between calls, which is what lets the RPC and the URL return
// byte-identical files.
func (r *Renderer) Render(v *core.AgentValues) (string, error) {
	file := values{
		Repositories: buildRepositories(v),
		Agent: release{
			Values:     buildAgentValues(v),
			ValuesFrom: overrideFrom(agentOverrideConfigMap),
		},
		Flux: release{
			Values:     buildFluxValues(v),
			ValuesFrom: overrideFrom(fluxOverrideConfigMap),
		},
	}

	var out strings.Builder

	enc := yaml.NewEncoder(&out)
	enc.SetIndent(yamlIndent)
	if err := enc.Encode(file); err != nil {
		return "", fmt.Errorf("encode agent values: %w", err)
	}
	if err := enc.Close(); err != nil {
		return "", fmt.Errorf("encode agent values: %w", err)
	}

	return out.String(), nil
}

func buildRepositories(v *core.AgentValues) repositories {
	return repositories{
		Modules: repository{
			URL:    v.Harbor.ModulesRepoURL,
			Labels: map[string]string{fromHarborLabelKey: fromHarborLabelValue},
		},
		Operators: repository{URL: v.Harbor.OperatorsRepoURL},
	}
}

func buildAgentValues(v *core.AgentValues) agentValues {
	out := agentValues{
		Agent: agentSection{
			ServerURL:       v.ServerURL,
			TunnelServerURL: v.TunnelServerURL,
			Cluster:         v.Cluster,
			JoinToken:       v.JoinToken,
		},
		ClusterAdmin: clusterAdmin{
			Enabled: true,
			Users:   v.ClusterAdminUsers,
		},
		ClusterInfo: clusterInfo{
			Enabled:         true,
			ExternalAddress: v.ClusterInfo.ExternalAddress,
			NodePortRange:   v.ClusterInfo.NodePortRange,
			InferenceURL:    v.ClusterInfo.InferenceURL,
		},
		TenantOperator: tenantOperator{
			Enabled: true,
			Harbor: harborConfig{
				URL: v.Harbor.URL,
				Robot: harborRobot{
					Name:   v.Harbor.Robot.Name,
					Secret: v.Harbor.Robot.Secret,
				},
			},
		},
	}

	if v.TrustedCASecret != "" {
		out.TrustedCA = &trustedCA{SecretName: v.TrustedCASecret, Key: v.TrustedCAKey}
	}

	return out
}

// buildFluxValues gives source-controller the CA it needs for a privately
// signed registry. SSL_CERT_DIR covers every pull it makes, which is why the
// HelmRepositories carry no per-repository certificate reference.
func buildFluxValues(v *core.AgentValues) any {
	if v.TrustedCASecret == "" {
		return nil
	}

	return fluxValues{
		Flux2: flux2{
			SourceController: sourceController{
				ExtraEnv: []envVar{{
					Name:  sslCertDirEnv,
					Value: sslCertDirDefault + ":" + trustedCADir,
				}},
				VolumeMounts: []volumeMount{{
					MountPath: trustedCADir,
					Name:      trustedCAVolume,
					ReadOnly:  true,
				}},
				Volumes: []volume{{
					Name:   trustedCAVolume,
					Secret: secretSource{SecretName: v.TrustedCASecret},
				}},
			},
		},
	}
}

func overrideFrom(name string) []valuesFrom {
	return []valuesFrom{{
		Kind:      "ConfigMap",
		Name:      name,
		ValuesKey: overrideValuesKey,
		Optional:  true,
	}}
}
