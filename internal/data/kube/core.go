package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	oscore "github.com/openhdc/otterscale/internal/core"
)

type core struct {
	kube *Kube
}

func NewCore(kube *Kube) oscore.KubeCoreRepo {
	return &core{
		kube: kube,
	}
}

var _ oscore.KubeCoreRepo = (*core)(nil)

func (r *core) ListServices(ctx context.Context, config *rest.Config, namespace string) ([]oscore.Service, error) {
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

func (r *core) ListPods(ctx context.Context, config *rest.Config, namespace string) ([]oscore.Pod, error) {
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

func (r *core) ListPersistentVolumeClaims(ctx context.Context, config *rest.Config, namespace string) ([]oscore.PersistentVolumeClaim, error) {
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
