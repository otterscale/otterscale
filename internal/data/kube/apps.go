package kube

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type apps struct {
	kubeMap KubeMap
}

func NewApps(kubeMap KubeMap) service.KubeApps {
	return &apps{
		kubeMap: kubeMap,
	}
}

var _ service.KubeApps = (*apps)(nil)

func (r *apps) ListDeployments(ctx context.Context, uuid, facility, namespace string) ([]v1.Deployment, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := clientset.AppsV1().Deployments(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *apps) GetDeployment(ctx context.Context, uuid, facility, namespace, name string) (*v1.Deployment, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return clientset.AppsV1().Deployments(namespace).Get(ctx, name, opts)
}

func (r *apps) ListStatefulSets(ctx context.Context, uuid, facility, namespace string) ([]v1.StatefulSet, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *apps) GetStatefulSet(ctx context.Context, uuid, facility, namespace, name string) (*v1.StatefulSet, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return clientset.AppsV1().StatefulSets(namespace).Get(ctx, name, opts)
}

func (r *apps) ListDaemonSets(ctx context.Context, uuid, facility, namespace string) ([]v1.DaemonSet, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *apps) GetDaemonSet(ctx context.Context, uuid, facility, namespace, name string) (*v1.DaemonSet, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, opts)
}
