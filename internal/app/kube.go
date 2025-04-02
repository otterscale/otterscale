package app

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/duration"

	v1 "github.com/openhdc/openhdc/api/kube/v1"
	"github.com/openhdc/openhdc/api/kube/v1/v1connect"
	"github.com/openhdc/openhdc/internal/domain/model"
	"github.com/openhdc/openhdc/internal/domain/service"
)

// KubeApp implements the KubeServiceHandler interface
type KubeApp struct {
	v1connect.UnimplementedKubeServiceHandler
	svc *service.KubeService
}

// NewKubeApp creates a new KubeApp instance
func NewKubeApp(svc *service.KubeService) *KubeApp {
	return &KubeApp{svc: svc}
}

// Ensure KubeApp implements the KubeServiceHandler interface
var _ v1connect.KubeServiceHandler = (*KubeApp)(nil)

// ListApplications retrieves applications for a given model UUID and cluster name
func (a *KubeApp) ListApplications(ctx context.Context, req *connect.Request[v1.ListApplicationsRequest]) (*connect.Response[v1.ListApplicationsResponse], error) {
	apps, err := a.svc.ListApplications(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName())
	if err != nil {
		return nil, err
	}
	res := &v1.ListApplicationsResponse{}
	res.SetApplications(a.toApplications(apps))
	return connect.NewResponse(res), nil
}

// toApplications converts domain model applications to API applications
func (a *KubeApp) toApplications(apps *model.Applications) []*v1.Application {
	ret := []*v1.Application{}
	storageClassMap := toStorageClassMap(apps.StorageClasses)

	for i := range apps.Deployments {
		selector, err := metav1.LabelSelectorAsSelector(apps.Deployments[i].Spec.Selector)
		if err != nil {
			continue
		}

		pvcs := filterPersistentVolumeClaim(apps.Deployments[i].Spec.Template.Spec.Volumes, apps.PersistentVolumeClaims)
		pods := filterPods(apps.Pods, selector)
		services := filterServices(apps.Services, pods)

		app := &v1.Application{}
		app.SetName(apps.Deployments[i].Name)
		app.SetNamespace(apps.Deployments[i].Namespace)
		app.SetLabels(apps.Deployments[i].Labels)

		if replicas := apps.Deployments[i].Spec.Replicas; replicas != nil {
			app.SetReplicas(*replicas)
		}

		app.SetStrategyType(string(apps.Deployments[i].Spec.Strategy.Type))
		app.SetContainers(toApplicationContainers(apps.Deployments[i].Spec.Template.Spec.Containers))
		app.SetPersistentVolumeClaims(toApplicationPersistentVolumeClaims(pvcs, storageClassMap))
		app.SetPods(toApplicationPods(pods))
		app.SetService(toApplicationServices(services))
		ret = append(ret, app)
	}

	return ret
}

// Filtering functions

// filterPersistentVolumeClaim filters PVCs based on volumes in a deployment
func filterPersistentVolumeClaim(vs []corev1.Volume, pvcs []corev1.PersistentVolumeClaim) []corev1.PersistentVolumeClaim {
	ret := []corev1.PersistentVolumeClaim{}
	for i := range vs {
		pvc := vs[i].PersistentVolumeClaim
		if pvc == nil {
			continue
		}
		for j := range pvcs {
			if pvc.ClaimName == pvcs[j].Name {
				ret = append(ret, pvcs[j])
				break
			}
		}
	}
	return ret
}

// filterPods filters pods based on a label selector
func filterPods(pods []corev1.Pod, selector labels.Selector) []corev1.Pod {
	ret := []corev1.Pod{}
	for idx := range pods {
		if selector.Matches(labels.Set(pods[idx].Labels)) {
			ret = append(ret, pods[idx])
		}
	}
	return ret
}

// filterServices filters services that match any of the pods
func filterServices(svcs []corev1.Service, pods []corev1.Pod) []corev1.Service {
	ret := []corev1.Service{}
	for i := range svcs {
		selector := labels.Set(svcs[i].Spec.Selector).AsSelector()
		for j := range pods {
			if selector.Matches(labels.Set(pods[j].Labels)) {
				ret = append(ret, svcs[i])
				break
			}
		}
	}
	return ret
}

// Utility functions

// toStorageClassMap creates a map of storage class name to storage class
func toStorageClassMap(scs []storagev1.StorageClass) map[string]storagev1.StorageClass {
	ret := map[string]storagev1.StorageClass{}
	for idx := range scs {
		ret[scs[idx].Name] = scs[idx]
	}
	return ret
}

// accessModesToStrings converts access modes to string slice
func accessModesToStrings(modes []corev1.PersistentVolumeAccessMode) []string {
	ret := make([]string, len(modes))
	for idx := range modes {
		ret[idx] = string(modes[idx])
	}
	return ret
}

// containerStatusesReadyString returns a string representing container readiness
func containerStatusesReadyString(statuses []corev1.ContainerStatus) string {
	ready := 0
	for idx := range statuses {
		if statuses[idx].Ready {
			ready++
		}
	}
	return fmt.Sprintf("%d/%d", ready, len(statuses))
}

// containerStatusesRestartString returns a string representing container restarts
func containerStatusesRestartString(statuses []corev1.ContainerStatus) string {
	restart := int32(0)
	var lastTerminatedAt metav1.Time
	for idx := range statuses {
		restart += statuses[idx].RestartCount
		if statuses[idx].LastTerminationState.Terminated != nil {
			lastTerminatedAt = statuses[idx].LastTerminationState.Terminated.FinishedAt
		}
	}
	if lastTerminatedAt.IsZero() {
		return fmt.Sprintf("%d", restart)
	}
	return fmt.Sprintf("%d (%s ago)", restart, duration.HumanDuration(time.Since(lastTerminatedAt.Time)))
}

// Conversion functions

// toApplicationContainers converts k8s containers to application containers
func toApplicationContainers(containers []corev1.Container) []*v1.Application_Container {
	ret := make([]*v1.Application_Container, 0, len(containers))
	for idx := range containers {
		container := &v1.Application_Container{}
		container.SetImageName(containers[idx].Image)
		container.SetImagePullPolicy(string(containers[idx].ImagePullPolicy))
		ret = append(ret, container)
	}
	return ret
}

// toApplicationPods converts k8s pods to application pods
func toApplicationPods(pods []corev1.Pod) []*v1.Application_Pod {
	ret := make([]*v1.Application_Pod, 0, len(pods))
	for idx := range pods {
		pod := &v1.Application_Pod{}
		pod.SetName(pods[idx].Name)
		pod.SetStatus(string(pods[idx].Status.Phase))
		pod.SetReady(containerStatusesReadyString(pods[idx].Status.ContainerStatuses))
		pod.SetRestarts(containerStatusesRestartString(pods[idx].Status.ContainerStatuses))
		ret = append(ret, pod)
	}
	return ret
}

// toApplicationServicePorts converts k8s service ports to application service ports
func toApplicationServicePorts(ports []corev1.ServicePort) []*v1.Application_Service_Port {
	ret := make([]*v1.Application_Service_Port, 0, len(ports))
	for idx := range ports {
		port := &v1.Application_Service_Port{}
		port.SetPort(ports[idx].Port)
		port.SetNodePort(ports[idx].NodePort)
		port.SetProtocol(string(ports[idx].Protocol))
		port.SetTargetPort(ports[idx].TargetPort.String())
		ret = append(ret, port)
	}
	return ret
}

// toApplicationServices converts k8s services to application services
func toApplicationServices(services []corev1.Service) []*v1.Application_Service {
	ret := make([]*v1.Application_Service, 0, len(services))
	for idx := range services {
		service := &v1.Application_Service{}
		service.SetName(services[idx].Name)
		service.SetType(string(services[idx].Spec.Type))
		service.SetClusterIp(services[idx].Spec.ClusterIP)
		service.SetPorts(toApplicationServicePorts(services[idx].Spec.Ports))
		ret = append(ret, service)
	}
	return ret
}

// toApplicationStorageClass converts k8s storage class to application storage class
func toApplicationStorageClass(sc *storagev1.StorageClass) *v1.Application_StorageClass {
	ret := &v1.Application_StorageClass{}
	ret.SetName(sc.Name)
	ret.SetProvisioner(sc.Provisioner)

	if reclaimPolicy := sc.ReclaimPolicy; reclaimPolicy != nil {
		ret.SetReclaimPolicy(string(*reclaimPolicy))
	}

	if volumeBindingMode := sc.VolumeBindingMode; volumeBindingMode != nil {
		ret.SetVolumeBindingMode(string(*volumeBindingMode))
	}

	ret.SetParameters(sc.Parameters)
	return ret
}

// toApplicationPersistentVolumeClaims converts k8s PVCs to application PVCs
func toApplicationPersistentVolumeClaims(pvcs []corev1.PersistentVolumeClaim, storageClassMap map[string]storagev1.StorageClass) []*v1.Application_PersistentVolumeClaim {
	ret := make([]*v1.Application_PersistentVolumeClaim, 0, len(pvcs))
	for idx := range pvcs {
		pvc := &v1.Application_PersistentVolumeClaim{}
		pvc.SetName(pvcs[idx].Name)
		pvc.SetStatus(pvcs[idx].Status.String())
		pvc.SetAccessModes(accessModesToStrings(pvcs[idx].Spec.AccessModes))
		pvc.SetCapacity(pvcs[idx].Spec.Resources.String())

		if storageClassName := pvcs[idx].Spec.StorageClassName; storageClassName != nil {
			if sc, ok := storageClassMap[*storageClassName]; ok {
				pvc.SetStorageClass(toApplicationStorageClass(&sc))
			}
		}
		ret = append(ret, pvc)
	}
	return ret
}
