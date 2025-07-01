package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	oscore "github.com/openhdc/otterscale/internal/core"
)

type batch struct {
	kube *Kube
}

func NewBatch(kube *Kube) oscore.KubeBatchRepo {
	return &batch{
		kube: kube,
	}
}

var _ oscore.KubeBatchRepo = (*batch)(nil)

func (r *batch) ListJobsByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]oscore.Job, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: label,
	}
	list, err := clientset.BatchV1().Jobs(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *batch) CreateJob(ctx context.Context, config *rest.Config, job *oscore.Job) (*oscore.Job, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}
	return clientset.BatchV1().Jobs(job.GetNamespace()).Create(ctx, job, opts)
}

func (r *batch) DeleteJob(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return err
	}

	propagation := metav1.DeletePropagationBackground
	opts := metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	}
	return clientset.BatchV1().Jobs(namespace).Delete(ctx, name, opts)
}
