package service

type KubeService struct {
	cronJob   KubeCronJob
	job       KubeJob
	namespace KubeNamespace
}

func NewKubeService(cronJob KubeCronJob, job KubeJob, namespace KubeNamespace) *KubeService {
	return &KubeService{
		cronJob:   cronJob,
		job:       job,
		namespace: namespace,
	}
}
