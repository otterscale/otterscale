package app

import (
	"context"

	"connectrpc.com/connect"

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
	ret := make([]*v1.Application, len(apps.Deployments.Items))

	if apps.Deployments != nil {
		for i := range apps.Deployments.Items {
			name := apps.Deployments.Items[i].Name
			namespace := apps.Deployments.Items[i].Namespace

			ret[i] = &v1.Application{}
			ret[i].SetName(name)
			ret[i].SetNamespace(namespace)
			ret[i].SetLabels(apps.Deployments.Items[i].Labels)
			ret[i].SetReplicas(*apps.Deployments.Items[i].Spec.Replicas)

			if apps.Services != nil {
				for j := range apps.Services.Items {
					if !a.equal(apps.Services.Items[j].Annotations, name, namespace) {
						continue
					}
				}
			}

		}
	}

	return ret
}

func (a *KubeApp) equal(annotations map[string]string, name, namespace string) bool {
	releaseName, hasName := annotations["meta.helm.sh/release-name"]
	releaseNamespace, hasNamespace := annotations["meta.helm.sh/release-namespace"]
	return hasName && hasNamespace && releaseName == name && releaseNamespace == namespace
}
