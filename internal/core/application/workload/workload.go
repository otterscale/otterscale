package workload

import (
	"net/http"
	"sync"

	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/application/storage"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ObjectMeta represents Kubernetes ObjectMeta.
type ObjectMeta = v1.ObjectMeta

type WorkloadUseCase struct {
	daemonSet   DaemonSetRepo
	deployment  DeploymentRepo
	job         JobRepo
	pod         PodRepo
	statefulSet StatefulSetRepo

	service               service.ServiceRepo
	persistentVolumeClaim storage.PersistentVolumeClaimRepo
	storageClass          storage.StorageClassRepo

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
