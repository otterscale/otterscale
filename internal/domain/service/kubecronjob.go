package service

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"

	"github.com/openhdc/openhdc/internal/domain/model"
)

type KubeCronJob interface {
	Get(ctx context.Context, cluster, namespace, name string) (*batchv1.CronJob, error)
	Create(ctx context.Context, cluster, namespace, name, image, schedule string) (*batchv1.CronJob, error)
	Update(ctx context.Context, cluster, namespace, name, image, schedule string) (*batchv1.CronJob, error)
	Delete(ctx context.Context, cluster, namespace, name string) error
}

func (s *KubeService) GetCronJob(ctx context.Context, cluster, namespace, name string) (*model.CronJob, error) {
	cj, err := s.cronJob.Get(ctx, cluster, namespace, name)
	if err != nil {
		return nil, err
	}
	return toCronJob(cj), nil
}

func (s *KubeService) CreateCronJob(ctx context.Context, cluster, namespace, name, image, schedule string) (*model.CronJob, error) {
	cj, err := s.cronJob.Create(ctx, cluster, namespace, name, image, schedule)
	if err != nil {
		return nil, err
	}
	return toCronJob(cj), err
}

func (s *KubeService) UpdateCronJob(ctx context.Context, cluster, namespace, name, image, schedule string) (*model.CronJob, error) {
	cj, err := s.cronJob.Update(ctx, cluster, namespace, name, image, schedule)
	if err != nil {
		return nil, err
	}
	return toCronJob(cj), err
}

func (s *KubeService) DeleteCronJob(ctx context.Context, cluster, namespace, name string) error {
	return s.cronJob.Delete(ctx, cluster, namespace, name)
}

func toCronJob(cj *batchv1.CronJob) *model.CronJob {
	ret := &model.CronJob{
		UID:                        string(cj.UID),
		Name:                       cj.Name,
		Schedule:                   cj.Spec.Schedule,
		Generation:                 cj.Generation,
		Suspend:                    cj.Spec.Suspend,
		SuccessfulJobsHistoryLimit: cj.Spec.SuccessfulJobsHistoryLimit,
		FailedJobsHistoryLimit:     cj.Spec.FailedJobsHistoryLimit,
		LastScheduleTime:           cj.Status.LastScheduleTime,
		LastSuccessfulTime:         cj.Status.LastSuccessfulTime,
		TTLSecondsAfterFinished:    cj.Spec.JobTemplate.Spec.TTLSecondsAfterFinished,
	}
	for i := range cj.Spec.JobTemplate.Spec.Template.Spec.Containers {
		ret.Containers = append(ret.Containers, &model.Container{
			Image:           cj.Spec.JobTemplate.Spec.Template.Spec.Containers[i].Image,
			ImagePullPolicy: string(cj.Spec.JobTemplate.Spec.Template.Spec.Containers[i].ImagePullPolicy),
		})
	}
	return ret
}
