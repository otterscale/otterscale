package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type core struct {
	kube *Kube
}

func NewCore(kube *Kube) service.KubeCore {
	return &core{
		kube: kube,
	}
}

var _ service.KubeCore = (*core)(nil)

func (r *core) ListServices(ctx context.Context, config *rest.Config, namespace string) ([]v1.Service, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.CoreV1().Services(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *core) ListPods(ctx context.Context, config *rest.Config, namespace string) ([]v1.Pod, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *core) ListPersistentVolumeClaims(ctx context.Context, config *rest.Config, namespace string) ([]v1.PersistentVolumeClaim, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
