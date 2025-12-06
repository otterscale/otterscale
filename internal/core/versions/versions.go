package versions

// OtterScale Helm Chart
const (
	GPUOperator   = "1.0.6" // HAMi 2.6.1 & GPU Operator 25.3.2
	KubeVirtInfra = "0.1.1" // KubeVirt 1.6.2 & CDI 1.63.1
	ModelArtifact = "0.1.2"
	Registry      = "0.2.0"
	SambaOperator = "0.1.0"
)

// Official Helm Chart
const (
	KubePrometheusStack = "79.9.0"
	Istio               = "1.28.0"
)

// llm-d Components
const (
	LLMDCuda               = "0.4.0"
	LLMDInferenceScheduler = "0.4.0"
	LLMDRoutingSidecar     = "0.4.0"
	LLMDModelService       = "0.3.10"
	LLMDInfra              = "1.3.4"
)

// Kubernetes Custom Resource Definition
const (
	GatewayAPI                   = "1.4.0"
	GatewayAPIInferenceExtension = "1.2.0"
)

// Kubernetes
const (
	CephCSI = "3.13.0"
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
