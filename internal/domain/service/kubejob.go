package service

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"

	"github.com/openhdc/openhdc/internal/domain/model"
)

type KubeJob interface {
	ListFromCronJob(ctx context.Context, cluster, namespace string, cronJob *batchv1.CronJob) (*batchv1.JobList, error)
	CreateFromCronJob(ctx context.Context, cluster, namespace string, cronJob *batchv1.CronJob, createdBy string) (*batchv1.Job, error)
}

func (s *KubeService) ListJobsFromCronJob(ctx context.Context, cluster, namespace, name string) ([]*model.Job, error) {
	cj, err := s.cronJob.Get(ctx, cluster, namespace, name)
	if err != nil {
		return nil, err
	}
	js, err := s.job.ListFromCronJob(ctx, cluster, namespace, cj)
	if err != nil {
		return nil, err
	}
	return toJobs(js), nil
}

func (s *KubeService) CreateJobFromCronJob(ctx context.Context, cluster, namespace, name, createdBy string) (*model.Job, error) {
	cj, err := s.cronJob.Get(ctx, cluster, namespace, name)
	if err != nil {
		return nil, err
	}
	j, err := s.job.CreateFromCronJob(ctx, cluster, namespace, cj, createdBy)
	if err != nil {
		return nil, err
	}
	return toJob(j), nil
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
