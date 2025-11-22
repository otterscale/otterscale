package versions

// OtterScale Helm Chart
const (
	GPUOperator   = "1.0.5" // HAMi 2.6.1 & GPU Operator 25.3.2
	KubeVirtInfra = "0.1.0" // KubeVirt 1.6.2 & CDI 1.63.1
	ModelArtifact = "0.1.0"
	SambaOperator = "0.1.0"
)

// Official Helm Chart
const (
	KubePrometheusStack = "79.7.1"
	Istio               = "1.28.0"
	InferencePool       = "0.1.0" // ?
	LLMDInfra           = "1.3.3"
	LLMDModelService    = "0.3.6"
)

// Kubernetes Custom Resource Definition
const (
	GatewayAPI                   = "v1.3.0"
	GatewayAPIInferenceExtension = "v1.1.0"
)

// Juju Charm
const (
	Kubernetes = "1.33/stable"
	Ceph       = "squid/stable"
	HACluster  = "2.8/stable"
)

// Controller Versions
const (
	Juju = "3.9"
	MAAS = "3.6"
)
