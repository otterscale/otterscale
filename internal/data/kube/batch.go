package kube

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	oscore "github.com/otterscale/otterscale/internal/core"
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

func (r *batch) ListJobs(ctx context.Context, config *rest.Config, namespace string) ([]oscore.Job, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	JobList, err := clientset.BatchV1().Jobs(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return JobList.Items, nil
}

func (r *batch) CreateJob(ctx context.Context, config *rest.Config, namespace, name string, labels, annotations map[string]string, spec *oscore.JobSpec) (*oscore.Job, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
	}
	if spec != nil {
		job.Spec = *spec
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
