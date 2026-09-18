package values

// The Go shape of the otterscale-agent-flux chart's values file. Field order is
// the file's order, which the encoder preserves.
//
// Typed structs rather than a template: the file nests six levels deep with
// conditional blocks, and a wrong indent there produces valid YAML with the
// wrong shape, which Helm ignores silently.

// values is the whole file.
type values struct {
	Repositories repositories `yaml:"repositories"`
	Agent        release      `yaml:"agent"`
	Flux         release      `yaml:"flux"`
}

// A struct rather than a map: Go sorts map keys, and a mistyped one would be
// accepted while silently dropping the chart defaults it meant to override.
type repositories struct {
	Modules   repository `yaml:"modules"`
	Operators repository `yaml:"operators"`
}

// Only the two keys the chart already defines are ever emitted, so interval
// and the rest come from its own values.yaml by deep merge.
type repository struct {
	URL    string            `yaml:"url"`
	Labels map[string]string `yaml:"labels,omitempty"`
}

// release is one HelmRelease's overrides, passed through untouched, so the
// shape is the target chart's.
type release struct {
	Values     any          `yaml:"values,omitempty"`
	ValuesFrom []valuesFrom `yaml:"valuesFrom"`
}

// valuesFrom mirrors a Flux HelmRelease spec.valuesFrom entry.
type valuesFrom struct {
	Kind      string `yaml:"kind"`
	Name      string `yaml:"name"`
	ValuesKey string `yaml:"valuesKey"`
	Optional  bool   `yaml:"optional"`
}

// agentValues is the otterscale-agent chart's own values.
type agentValues struct {
	Agent          agentSection   `yaml:"agent"`
	ClusterAdmin   clusterAdmin   `yaml:"clusterAdmin"`
	ClusterInfo    clusterInfo    `yaml:"clusterInfo"`
	TenantOperator tenantOperator `yaml:"tenantOperator"`
	// The one optional block: a publicly signed server gives an agent nothing
	// to trust, and the chart's empty default is then already right.
	TrustedCA *trustedCA `yaml:"trustedCA,omitempty"`
}

type agentSection struct {
	ServerURL       string `yaml:"serverURL"`
	TunnelServerURL string `yaml:"tunnelServerURL"`
	Cluster         string `yaml:"cluster"`
	JoinToken       string `yaml:"joinToken"`
}

type clusterAdmin struct {
	Enabled bool     `yaml:"enabled"`
	Users   []string `yaml:"users"`
}

type clusterInfo struct {
	Enabled         bool   `yaml:"enabled"`
	ExternalAddress string `yaml:"externalAddress"`
	NodePortRange   string `yaml:"nodePortRange"`
	InferenceURL    string `yaml:"inferenceURL"`
}

type tenantOperator struct {
	Enabled bool         `yaml:"enabled"`
	Harbor  harborConfig `yaml:"harbor"`
}

type harborConfig struct {
	URL   string      `yaml:"url"`
	Robot harborRobot `yaml:"robot"`
}

type harborRobot struct {
	Name   string `yaml:"name"`
	Secret string `yaml:"secret"`
}

type trustedCA struct {
	SecretName string `yaml:"secretName"`
	Key        string `yaml:"key"`
}

// fluxValues carries only the CA that lets source-controller verify a
// privately signed registry.
type fluxValues struct {
	Flux2 flux2 `yaml:"flux2"`
}

type flux2 struct {
	SourceController sourceController `yaml:"sourceController"`
}

type sourceController struct {
	ExtraEnv     []envVar      `yaml:"extraEnv"`
	VolumeMounts []volumeMount `yaml:"volumeMounts"`
	Volumes      []volume      `yaml:"volumes"`
}

type envVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type volumeMount struct {
	MountPath string `yaml:"mountPath"`
	Name      string `yaml:"name"`
	ReadOnly  bool   `yaml:"readOnly"`
}

type volume struct {
	Name   string       `yaml:"name"`
	Secret secretSource `yaml:"secret"`
}

// Deliberately not optional: the same Secret backs trustedCA on the agent, so
// a failed mount is a clearer signal than x509 errors on every pull.
type secretSource struct {
	SecretName string `yaml:"secretName"`
}
