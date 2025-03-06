package service

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/openhdc/openhdc/internal/service/domain/model"
)

type KubeNamespace interface {
	Get(ctx context.Context) (*corev1.Namespace, error)
	Create(ctx context.Context) (*corev1.Namespace, error)
}

type KubeCronJob interface {
	Get(ctx context.Context, name string) (*batchv1.CronJob, error)
	Create(ctx context.Context, name, image, schedule string) (*batchv1.CronJob, error)
	Update(ctx context.Context, name, image, schedule string) (*batchv1.CronJob, error)
	Delete(ctx context.Context, name string) error
}

type KubeJob interface {
	ListFromCronJob(ctx context.Context, cj *batchv1.CronJob) (*batchv1.JobList, error)
	CreateFromCronJob(ctx context.Context, cronJob *batchv1.CronJob, createdBy string) (*batchv1.Job, error)
}

type KubeService struct {
	namespace KubeNamespace
	cronJob   KubeCronJob
	job       KubeJob
}

func NewKubeService(namespace KubeNamespace, cronJob KubeCronJob, job KubeJob) *KubeService {
	return &KubeService{
		namespace: namespace,
		cronJob:   cronJob,
		job:       job,
	}
}

func (s *KubeService) GetCronJob(ctx context.Context, name string) (*model.CronJob, error) {
	cj, err := s.cronJob.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return toCronJob(cj), nil
}

func (s *KubeService) CreateCronJob(ctx context.Context, name, image, schedule string) (*model.CronJob, error) {
	if err := s.createNamespaceIfNotExists(ctx); err != nil {
		return nil, err
	}
	cj, err := s.cronJob.Create(ctx, name, image, schedule)
	if err != nil {
		return nil, err
	}
	return toCronJob(cj), err
}

func (s *KubeService) UpdateCronJob(ctx context.Context, name, image, schedule string) (*model.CronJob, error) {
	cj, err := s.cronJob.Update(ctx, name, image, schedule)
	if err != nil {
		return nil, err
	}
	return toCronJob(cj), err
}

func (s *KubeService) DeleteCronJob(ctx context.Context, name string) error {
	return s.cronJob.Delete(ctx, name)
}

func (s *KubeService) ListJobsFromCronJob(ctx context.Context, name string) ([]*model.Job, error) {
	cj, err := s.cronJob.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	js, err := s.job.ListFromCronJob(ctx, cj)
	if err != nil {
		return nil, err
	}
	return toJobs(js), nil
}

func (s *KubeService) CreateJobFromCronJob(ctx context.Context, name, createdBy string) (*model.Job, error) {
	if err := s.createNamespaceIfNotExists(ctx); err != nil {
		return nil, err
	}
	cj, err := s.cronJob.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	j, err := s.job.CreateFromCronJob(ctx, cj, createdBy)
	if err != nil {
		return nil, err
	}
	return toJob(j), nil
}

func (s *KubeService) createNamespaceIfNotExists(ctx context.Context) error {
	_, err := s.namespace.Get(ctx)
	if err == nil {
		return nil
	}
	if !apierrors.IsNotFound(err) {
		return err
	}
	_, err = s.namespace.Create(ctx)
	return err
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

func toJob(j *batchv1.Job) *model.Job {
	ret := &model.Job{
		UID:                     string(j.UID),
		Name:                    j.Name,
		Succeeded:               j.Status.Succeeded > 0,
		StartTime:               j.Status.StartTime,
		CompletionTime:          j.Status.CompletionTime,
		TTLSecondsAfterFinished: j.Spec.TTLSecondsAfterFinished,
	}
	for i := range j.Spec.Template.Spec.Containers {
		ret.Containers = append(ret.Containers, &model.Container{
			Image:           j.Spec.Template.Spec.Containers[i].Image,
			ImagePullPolicy: string(j.Spec.Template.Spec.Containers[i].ImagePullPolicy),
		})
	}
	return ret
}

func toJobs(l *batchv1.JobList) []*model.Job {
	ret := []*model.Job{}
	for i := range l.Items {
		ret = append(ret, toJob(&l.Items[i]))
	}
	return ret
}
