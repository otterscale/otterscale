package kube

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

var (
	ttlSecondsAfterFinished    int32 = 86400
	successfulJobsHistoryLimit int32 = 5
	failedJobsHistoryLimit     int32 = 3
)

type batch struct {
	kubeMap KubeMap
}

func NewBatch(kubeMap KubeMap) service.KubeBatch {
	return &batch{
		kubeMap: kubeMap,
	}
}

var _ service.KubeBatch = (*batch)(nil)

func (r *batch) GetCronJob(ctx context.Context, uuid, facility, namespace, name string) (*v1.CronJob, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return clientset.BatchV1().CronJobs(namespace).Get(ctx, name, opts)
}

func (r *batch) CreateCronJob(ctx context.Context, uuid, facility, namespace, name, image, schedule string) (*v1.CronJob, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	cronJob := toCronJob(name, image, schedule)
	opts := metav1.CreateOptions{}
	return clientset.BatchV1().CronJobs(namespace).Create(ctx, cronJob, opts)
}

func (r *batch) UpdateCronJob(ctx context.Context, uuid, facility, namespace, name, image, schedule string) (*v1.CronJob, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	cronJob := toCronJob(name, image, schedule)
	opts := metav1.UpdateOptions{}
	return clientset.BatchV1().CronJobs(namespace).Update(ctx, cronJob, opts)
}

func (r *batch) DeleteCronJob(ctx context.Context, uuid, facility, namespace, name string) error {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return clientset.BatchV1().CronJobs(namespace).Delete(ctx, name, opts)
}

func (r *batch) ListJobsFromCronJob(ctx context.Context, uuid, facility, namespace string, cronJob *v1.CronJob) (*v1.JobList, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("cronjob-name=%s", cronJob.Name),
	}
	return clientset.BatchV1().Jobs(namespace).List(ctx, opts)
}

func (r *batch) CreateJobFromCronJob(ctx context.Context, uuid, facility, namespace string, cronJob *v1.CronJob, createdBy string) (*v1.Job, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
	if err != nil {
		return nil, err
	}
	job := toJob(cronJob, createdBy)
	opts := metav1.CreateOptions{}
	return clientset.BatchV1().Jobs(namespace).Create(ctx, job, opts)
}

// TODO: Container
func toCronJob(name, image, schedule string) *v1.CronJob {
	return &v1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.CronJobSpec{
			Schedule:                   schedule,
			SuccessfulJobsHistoryLimit: &successfulJobsHistoryLimit,
			FailedJobsHistoryLimit:     &failedJobsHistoryLimit,
			JobTemplate: v1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"cronjob-name": name,
					},
				},
				Spec: v1.JobSpec{
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

func toJob(cronJob *v1.CronJob, createdBy string) *v1.Job {
	return &v1.Job{
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
