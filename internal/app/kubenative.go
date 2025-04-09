package app

import (
	"context"
	"fmt"
	"slices"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/duration"

	v1 "github.com/openhdc/openhdc/api/kube/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

const (
	appTypeDeployment  = "Deployment"
	appTypeStatefulSet = "StatefulSet"
	appTypeDaemonSet   = "DaemonSet"
)

// ListApplications retrieves applications for a given model UUID and cluster name
func (a *KubeApp) ListApplications(ctx context.Context, req *connect.Request[v1.ListApplicationsRequest]) (*connect.Response[v1.ListApplicationsResponse], error) {
	apps, err := a.svc.ListApplications(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName(), "", "")
	if err != nil {
		return nil, err
	}

	res := &v1.ListApplicationsResponse{}
	res.SetApplications(a.toApplications(apps))
	return connect.NewResponse(res), nil
}

// GetApplication retrieves a specific application by name and namespace
func (a *KubeApp) GetApplication(ctx context.Context, req *connect.Request[v1.GetApplicationRequest]) (*connect.Response[v1.Application], error) {
	apps, err := a.svc.ListApplications(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}

	// Find the requested application
	if len(apps.Deployments) > 0 {
		app, err := a.toApplicationFromDeployment(&apps.Deployments[0], apps)
		if err != nil {
			return nil, err
		}
		return connect.NewResponse(app), nil
	} else if len(apps.StatefulSets) > 0 {
		app, err := a.toApplicationFromStatefulSet(&apps.StatefulSets[0], apps)
		if err != nil {
			return nil, err
		}
		return connect.NewResponse(app), nil
	} else if len(apps.DaemonSets) > 0 {
		app, err := a.toApplicationFromDaemonSet(&apps.DaemonSets[0], apps)
		if err != nil {
			return nil, err
		}
		return connect.NewResponse(app), nil
	}

	return nil, fmt.Errorf("application %s in namespace %s not found", req.Msg.GetName(), req.Msg.GetNamespace())
}

// Helper functions for creating applications from different resource types
func (a *KubeApp) toApplicationFromDeployment(dpl *appsv1.Deployment, apps *model.Applications) (*v1.Application, error) {
	return a.toApplication(dpl.Spec.Selector, dpl.Spec.Template.Spec.Volumes, apps,
		appTypeDeployment, dpl.Name, dpl.Namespace, dpl.Labels,
		dpl.Spec.Replicas, dpl.Spec.Template.Spec.Containers)
}

func (a *KubeApp) toApplicationFromStatefulSet(sfs *appsv1.StatefulSet, apps *model.Applications) (*v1.Application, error) {
	return a.toApplication(sfs.Spec.Selector, sfs.Spec.Template.Spec.Volumes, apps,
		appTypeStatefulSet, sfs.Name, sfs.Namespace, sfs.Labels,
		sfs.Spec.Replicas, sfs.Spec.Template.Spec.Containers)
}

func (a *KubeApp) toApplicationFromDaemonSet(dms *appsv1.DaemonSet, apps *model.Applications) (*v1.Application, error) {
	replica := int32(1)
	return a.toApplication(dms.Spec.Selector, dms.Spec.Template.Spec.Volumes, apps,
		appTypeDaemonSet, dms.Name, dms.Namespace, dms.Labels,
		&replica, dms.Spec.Template.Spec.Containers)
}

// toApplication creates an application from kubernetes resources
func (a *KubeApp) toApplication(s *metav1.LabelSelector, vs []corev1.Volume, apps *model.Applications, appType, name, namespace string, labels map[string]string, replicas *int32, cs []corev1.Container) (*v1.Application, error) {
	selector, err := metav1.LabelSelectorAsSelector(s)
	if err != nil {
		return nil, fmt.Errorf("failed to create selector: %w", err)
	}

	pvcs := filterPersistentVolumeClaim(vs, apps.PersistentVolumeClaims)
	pods := filterPods(apps.Pods, selector)
	svcs := filterServices(apps.Services, pods)
	scm := toStorageClassMap(apps.StorageClasses)

	ret := &v1.Application{}
	ret.SetType(appType)
	ret.SetName(name)
	ret.SetNamespace(namespace)
	ret.SetLabels(labels)

	if replicas != nil {
		ret.SetReplicas(*replicas)
	}

	ret.SetHealthies(countHealthies(pods))

	ret.SetContainers(toApplicationContainers(cs))
	ret.SetPersistentVolumeClaims(toApplicationPersistentVolumeClaims(pvcs, scm))
	ret.SetPods(toApplicationPods(pods))
	ret.SetService(toApplicationServices(svcs))

	return ret, nil
}

// toApplications converts domain model applications to API applications
func (a *KubeApp) toApplications(apps *model.Applications) []*v1.Application {
	ret := []*v1.Application{}

	// Process deployments
	for i := range apps.Deployments {
		app, err := a.toApplicationFromDeployment(&apps.Deployments[i], apps)
		if err != nil {
			continue
		}
		ret = append(ret, app)
	}

	// Process StatefulSets
	for i := range apps.StatefulSets {
		app, err := a.toApplicationFromStatefulSet(&apps.StatefulSets[i], apps)
		if err != nil {
			continue
		}
		ret = append(ret, app)
	}

	// Process DaemonSets
	for i := range apps.DaemonSets {
		app, err := a.toApplicationFromDaemonSet(&apps.DaemonSets[i], apps)
		if err != nil {
			continue
		}
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

func countHealthies(pods []corev1.Pod) int32 {
	phases := []corev1.PodPhase{corev1.PodRunning, corev1.PodSucceeded}
	count := int32(0)
	for idx := range pods {
		if slices.Contains(phases, pods[idx].Status.Phase) {
			count++
		}
	}
	return count
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

func toApplicationPodLastCondition(conditions []corev1.PodCondition) *v1.Application_Pod_Condition {
	idx := len(conditions) - 1
	ret := &v1.Application_Pod_Condition{}
	ret.SetType(string(conditions[idx].Type))
	ret.SetStatus(string(conditions[idx].Status))
	ret.SetProbedAt(timestamppb.New(conditions[idx].LastProbeTime.Time))
	ret.SetTransitionedAt(timestamppb.New(conditions[idx].LastTransitionTime.Time))
	ret.SetReason((conditions[idx].Reason))
	ret.SetMessage((conditions[idx].Message))
	return ret
}

// toApplicationPods converts k8s pods to application pods
func toApplicationPods(pods []corev1.Pod) []*v1.Application_Pod {
	ret := make([]*v1.Application_Pod, 0, len(pods))
	for idx := range pods {
		pod := &v1.Application_Pod{}
		pod.SetName(pods[idx].Name)
		pod.SetPhase(string(pods[idx].Status.Phase))
		pod.SetReady(containerStatusesReadyString(pods[idx].Status.ContainerStatuses))
		pod.SetRestarts(containerStatusesRestartString(pods[idx].Status.ContainerStatuses))
		pod.SetLastCondition(toApplicationPodLastCondition(pods[idx].Status.Conditions))
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
