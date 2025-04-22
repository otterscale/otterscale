package model

type Error struct {
	Code    string
	Level   ErrorLevel
	Message string
	Details string
	URL     string
}

var ErrCephNotFound = Error{
	Code:    "CEPH_NOT_FOUND",
	Level:   ErrorLevelCritical,
	Message: "Ceph cluster not found.",
	Details: "Please install and configure a Ceph cluster as your storage space.",
	URL:     "/docs/setup/ceph",
}

var ErrKubernetesNotFound = Error{
	Code:    "KUBERNETES_NOT_FOUND",
	Level:   ErrorLevelCritical,
	Message: "Kubernetes cluster not found.",
	Details: "Please install and configure a Kubernetes cluster as your orchestration platform.",
	URL:     "/docs/setup/kubernetes",
}
