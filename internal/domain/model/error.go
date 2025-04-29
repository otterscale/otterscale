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

var ErrNoMachinesDeployed = Error{
	Code:    "NO_MACHINES_DEPLOYED",
	Level:   ErrorLevelCritical,
	Message: "No machines have been deployed yet.",
	Details: "There are currently no deployed machines in the system. Please deploy at least one machine to continue.",
	URL:     "/docs/machines/deployment",
}
