package kube

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openhdc/openhdc/internal/service/domain/service"
)

type job struct {
	client *kubernetes.Clientset
}

func NewJob(client *kubernetes.Clientset) service.KubeJob {
	return &job{
		client: client,
	}
}

var _ service.KubeJob = (*job)(nil)

func (r *job) ListFromCronJob(ctx context.Context, cj *batchv1.CronJob) (*batchv1.JobList, error) {
	opts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("cronjob-name=%s", cj.Name),
	}
	return r.client.BatchV1().Jobs(ns).List(ctx, opts)
}

func (r *job) CreateFromCronJob(ctx context.Context, cj *batchv1.CronJob, createdBy string) (*batchv1.Job, error) {
	j := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: cj.Name + "-" + strings.Split(uuid.NewString(), "-")[0],
			Labels: map[string]string{
				"created-by":   createdBy,
				"cronjob-name": cj.Name,
			},
		},
		Spec: *cj.Spec.JobTemplate.Spec.DeepCopy(),
	}
	opts := metav1.CreateOptions{}
	return r.client.BatchV1().Jobs(ns).Create(ctx, j, opts)
}
