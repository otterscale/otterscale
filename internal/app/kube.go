package app

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	corev1 "k8s.io/api/core/v1"

	v1 "github.com/openhdc/openhdc/api/kube/v1"
	"github.com/openhdc/openhdc/api/kube/v1/v1connect"
	"github.com/openhdc/openhdc/internal/domain/model"
	"github.com/openhdc/openhdc/internal/domain/service"
)

// KubeApp implements the StackServiceServer interface
type KubeApp struct {
	v1connect.UnimplementedKubeServiceHandler
	svc *service.KubeService
}

// NewKubeApp creates a new KubeApp instance
func NewKubeApp(svc *service.KubeService) *KubeApp {
	return &KubeApp{svc: svc}
}

// Ensure KubeApp implements the StackServiceServer interface
var _ v1connect.KubeServiceHandler = (*KubeApp)(nil)

func (a *KubeApp) ListApplications(ctx context.Context, req *connect.Request[v1.ListApplicationsRequest]) (*connect.Response[v1.ListApplicationsResponse], error) {
	apps, err := a.svc.ListApplications(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName())
	if err != nil {
		return nil, err
	}
	res := &v1.ListApplicationsResponse{}
	res.SetApplications(a.toApplications(apps))
	return connect.NewResponse(res), nil
}

func (a *KubeApp) toApplications(apps *model.Applications) []*v1.Application {
	ret := []*v1.Application{}

	for i := range apps.Deployments {
		name, namespace, err := a.helmName(apps.Deployments[i].Annotations)
		if err != nil {
			continue
		}

		app := &v1.Application{}
		app.SetName(apps.Deployments[i].Name)
		app.SetNamespace(apps.Deployments[i].Namespace)
		app.SetLabels(apps.Deployments[i].Labels)
		app.SetReplicas(*apps.Deployments[i].Spec.Replicas)

		for j := range apps.Services {
			if !a.helmEqual(apps.Services[j].Annotations, name, namespace) {
				continue
			}

			service := &v1.Application_Service{}
			service.SetType(string(apps.Services[j].Spec.Type))
			service.SetClusterIp(apps.Services[j].Spec.ClusterIP)

			ports := make([]*v1.Application_Service_Port, len(apps.Services[j].Spec.Ports))
			for k := range apps.Services[j].Spec.Ports {
				port := &v1.Application_Service_Port{}
				port.SetPort(apps.Services[j].Spec.Ports[k].Port)
				port.SetNodePort(apps.Services[j].Spec.Ports[k].NodePort)
				port.SetProtocol(string(apps.Services[j].Spec.Ports[k].Protocol))
				port.SetTargetPort(apps.Services[j].Spec.Ports[k].TargetPort.String())
				ports[k] = port
			}
			service.SetPorts(ports)

			app.SetService(service)
		}

		pods := []*v1.Application_Pod{}
		for j := range apps.Pods {
			if !a.helmEqual(apps.Pods[j].Annotations, name, namespace) {
				continue
			}

			pod := &v1.Application_Pod{}
			pod.SetName(apps.Pods[j].Name)
			pod.SetStatus(apps.Pods[j].Status.String())

			pods = append(pods, pod)
		}
		app.SetPods(pods)

		persistentVolumeClaims := []*v1.Application_PersistentVolumeClaim{}
		for j := range apps.PersistentVolumeClaims {
			if !a.helmEqual(apps.PersistentVolumeClaims[j].Annotations, name, namespace) {
				continue
			}

			persistentVolumeClaim := &v1.Application_PersistentVolumeClaim{}
			persistentVolumeClaim.SetStatus(apps.PersistentVolumeClaims[j].Status.String())
			persistentVolumeClaim.SetAccessModes(accessModesToStrings(apps.PersistentVolumeClaims[j].Spec.AccessModes))
			persistentVolumeClaim.SetCapacity(apps.PersistentVolumeClaims[j].Spec.Resources.String())

			storageClassName := apps.PersistentVolumeClaims[j].Spec.StorageClassName
			if storageClassName != nil {
				persistentVolumeClaim.SetStorageClassName(*storageClassName)
			}

			persistentVolumeClaims = append(persistentVolumeClaims, persistentVolumeClaim)
		}
		app.SetPersistentVolumeClaims(persistentVolumeClaims)

		ret = append(ret, app)
	}

	return ret
}

func (a *KubeApp) helmName(annotations map[string]string) (name, namespace string, err error) {
	releaseName, hasName := annotations["meta.helm.sh/release-name"]
	releaseNamespace, hasNamespace := annotations["meta.helm.sh/release-namespace"]
	if !hasName || !hasNamespace {
		return "", "", fmt.Errorf("helm name not found")
	}
	return releaseName, releaseNamespace, nil
}

func (a *KubeApp) helmEqual(annotations map[string]string, name, namespace string) bool {
	releaseName, hasName := annotations["meta.helm.sh/release-name"]
	releaseNamespace, hasNamespace := annotations["meta.helm.sh/release-namespace"]
	return hasName && hasNamespace && releaseName == name && releaseNamespace == namespace
}

func accessModesToStrings(modes []corev1.PersistentVolumeAccessMode) []string {
	ret := make([]string, len(modes))
	for idx := range modes {
		ret[idx] = string(modes[idx])
	}
	return ret
}
