package workload

import (
	"net/http"
	"sync"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/service"
)

const ResourceStorage = v1.ResourceStorage

type (
	// ObjectMeta represents Kubernetes ObjectMeta.
	ObjectMeta = metav1.ObjectMeta

	// Time represents Kubernetes Time.
	Time = metav1.Time

	// ResourceList represents Kubernetes ResourceList.
	ResourceList = v1.ResourceList
)

type WorkloadUseCase struct {
	daemonSet   DaemonSetRepo
	deployment  DeploymentRepo
	job         JobRepo
	pod         PodRepo
	statefulSet StatefulSetRepo

	service               service.ServiceRepo
	persistentVolumeClaim persistent.PersistentVolumeClaimRepo
	storageClass          persistent.StorageClassRepo

	ttySessions sync.Map
}

func NewWorkloadUseCase(daemonSet DaemonSetRepo, deployment DeploymentRepo, job JobRepo, pod PodRepo, statefulSet StatefulSetRepo) *WorkloadUseCase {
	return &WorkloadUseCase{
		daemonSet:   daemonSet,
		deployment:  deployment,
		job:         job,
		pod:         pod,
		statefulSet: statefulSet,
	}
}

func (uc *WorkloadUseCase) isKeyNotFoundError(err error) bool {
	statusErr, _ := err.(*k8serrors.StatusError)
	return statusErr != nil && statusErr.Status().Code == http.StatusNotFound
}
