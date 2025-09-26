package core

const (
	DataVolumeBootImageLabel    = "otterscale.com/data-volume.boot-image"
	VirtualMachineNameLabel     = "otterscale.com/virtual-machine.name"
	ApplicationReleaseNameLabel = "otterscale.com/application-release.name"
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
