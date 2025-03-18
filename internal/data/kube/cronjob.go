package kube

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

var (
	ttlSecondsAfterFinished    int32 = 86400
	successfulJobsHistoryLimit int32 = 5
	failedJobsHistoryLimit     int32 = 3
)

type cronJob struct {
	kubes Kubes
}

func NewCronJob(kubes Kubes) service.KubeCronJob {
	return &cronJob{
		kubes: kubes,
	}
}

var _ service.KubeCronJob = (*cronJob)(nil)

func (r *cronJob) Get(ctx context.Context, cluster, namespace, name string) (*batchv1.CronJob, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return client.BatchV1().CronJobs(namespace).Get(ctx, name, opts)
}

func (r *cronJob) Create(ctx context.Context, cluster, namespace, name, image, schedule string) (*batchv1.CronJob, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	cronJob := r.toCronJob(name, image, schedule)
	opts := metav1.CreateOptions{}
	return client.BatchV1().CronJobs(namespace).Create(ctx, cronJob, opts)
}

func (r *cronJob) Update(ctx context.Context, cluster, namespace, name, image, schedule string) (*batchv1.CronJob, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	cronJob := r.toCronJob(name, image, schedule)
	opts := metav1.UpdateOptions{}
	return client.BatchV1().CronJobs(namespace).Update(ctx, cronJob, opts)
}

func (r *cronJob) Delete(ctx context.Context, cluster, namespace, name string) error {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return client.BatchV1().CronJobs(namespace).Delete(ctx, name, opts)
}

func (r *cronJob) toCronJob(name, image, schedule string) *batchv1.CronJob {
	return &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:                   schedule,
			SuccessfulJobsHistoryLimit: &successfulJobsHistoryLimit,
			FailedJobsHistoryLimit:     &failedJobsHistoryLimit,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"cronjob-name": name,
					},
				},
				Spec: batchv1.JobSpec{
					TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
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
}
