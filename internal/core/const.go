package core

const (
	LabelDomain                          = "otterscale.com"
	DataVolumeBootImageLabel             = "otterscale.com/data-volume.boot-image"
	VirtualMachineNameLabel              = "otterscale.com/virtual-machine.name"
	ApplicationReleaseNameLabel          = "otterscale.com/application-release.name"
	ApplicationReleaseLLMDModelNameLabel = "otterscale.com/application-release.llmd-model-name"
	ApplicationReleaseChartRefAnnotation = "otterscale.com/application-release.chart-ref"
)

const BuiltInMachineTagComment = "built-in"

const (
	Kubernetes             = "kubernetes"
	KubernetesControlPlane = "kubernetes-control-plane"
	KubernetesWorker       = "kubernetes-worker"
	Ceph                   = "ceph"
	CephMon                = "ceph-mon"
	CephOSD                = "ceph-osd"
	KubeVirt               = "kubevirt"
	LLMd                   = "llm-d"
)

const (
	kubevirtHealthOK           = 11
	kubevirtHealthNotInstalled = 21
	kubevirtHealthFailed       = 22
	kubevirtHealthPending      = 31
)

const (
	llmdInfraHealthOK           = 11
	llmdInfraHealthNotInstalled = 21
	llmdInfraHealthFailed       = 22
	llmdInfraHealthPending      = 31
)

const (
	environmentHealthOK           = 11
	environmentHealthNotInstalled = 21
)

const (
	ApplicationTypeDeployment  = "Deployment"
	ApplicationTypeStatefulSet = "StatefulSet"
	ApplicationTypeDaemonSet   = "DaemonSet"
)

const (
	annotationHAMINodeNvidiaRegister   = "hami.io/node-nvidia-register"
	annotationHAMIVGPUNode             = "hami.io/vgpu-node"
	annotationHAMIVGPUDevicesAllocated = "hami.io/vgpu-devices-allocated"
	annotationHAMIBindTime             = "hami.io/bind-time"
	annotationHAMIBindPhase            = "hami.io/bind-phase"
)

const (
	bistKindFIO             = "fio"
	bistKindWarp            = "warp"
	bistNamespace           = "bist"
	bistLabel               = "bist.otterscale.com/name=bist"
	bistAnnotationCreatedBy = "bist.otterscale.com/created-by"
	bistAnnotationKind      = "bist.otterscale.com/kind"
	bistAnnotationFIO       = "bist.otterscale.com/fio"
	bistAnnotationWarp      = "bist.otterscale.com/warp"
	bistBlockPool           = "otterscale_bist_pool"
	bistBlockImage          = "otterscale_bist_image"
)

const (
	minioLabel       = "app.kubernetes.io/name=minio"
	minioField       = "spec.type=NodePort"
	minioServiceName = "minio-api"
)

type DataVolumeSourceType int64

const (
	DataVolumeSourceTypeBlank DataVolumeSourceType = iota
	DataVolumeSourceTypeHTTP
	DataVolumeSourceTypePVC
)

type EssentialType int32

const (
	EssentialTypeUnknown EssentialType = iota
	EssentialTypeKubernetes
	EssentialTypeCeph
)
