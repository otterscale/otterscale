package kube

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type job struct {
	kubes Kubes
}

func NewJob(kubes Kubes) service.KubeJob {
	return &job{
		kubes: kubes,
	}
}

var _ service.KubeJob = (*job)(nil)

func (r *job) ListFromCronJob(ctx context.Context, cluster, namespace string, cronJob *batchv1.CronJob) (*batchv1.JobList, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("cronjob-name=%s", cronJob.Name),
	}
	return client.BatchV1().Jobs(namespace).List(ctx, opts)
}

func (r *job) CreateFromCronJob(ctx context.Context, cluster, namespace string, cronJob *batchv1.CronJob, createdBy string) (*batchv1.Job, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	job := r.toJob(cronJob, createdBy)
	opts := metav1.CreateOptions{}
	return client.BatchV1().Jobs(namespace).Create(ctx, job, opts)
}

func (r *job) toJob(cronJob *batchv1.CronJob, createdBy string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: cronJob.Name + "-" + strings.Split(uuid.NewString(), "-")[0],
			Labels: map[string]string{
				"created-by":   createdBy,
				"cronjob-name": cronJob.Name,
			},
		},
		Spec: *cronJob.Spec.JobTemplate.Spec.DeepCopy(),
	}
}
