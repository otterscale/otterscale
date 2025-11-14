package workload

import (
	"net/http"
	"sync"

	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/service"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ObjectMeta represents Kubernetes ObjectMeta.
type ObjectMeta = metav1.ObjectMeta

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
