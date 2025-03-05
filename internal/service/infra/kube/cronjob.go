package kube

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openhdc/openhdc/internal/service/domain/service"
)

type cronJob struct {
	client *kubernetes.Clientset
}

func NewCronJob(client *kubernetes.Clientset) service.KubeCronJob {
	return &cronJob{
		client: client,
	}
}

var _ service.KubeCronJob = (*cronJob)(nil)

func (r *cronJob) Get(ctx context.Context, name string) (*batchv1.CronJob, error) {
	opts := metav1.GetOptions{}
	return r.client.BatchV1().CronJobs(ns).Get(ctx, name, opts)
}

func (r *cronJob) Create(ctx context.Context, name, image, schedule string) (*batchv1.CronJob, error) {
	cj := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:                   schedule,
			SuccessfulJobsHistoryLimit: &succLimit,
			FailedJobsHistoryLimit:     &failLimit,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"cronjob-name": name,
					},
				},
				Spec: batchv1.JobSpec{
					TTLSecondsAfterFinished: &ttl,
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            name,
									Image:           image,
									ImagePullPolicy: corev1.PullIfNotPresent,
									Command:         []string{"/bin/sh", "-c", "date; echo Hello from the Kubernetes cluster"},
								},
							},
							RestartPolicy: corev1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}
	opts := metav1.CreateOptions{}
	return r.client.BatchV1().CronJobs(ns).Create(ctx, cj, opts)
}
