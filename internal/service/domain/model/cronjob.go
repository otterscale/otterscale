package model

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type CronJob struct {
	UID                        string       `json:"uid"`
	Name                       string       `json:"name"`
	Schedule                   string       `json:"schedule"`
	Generation                 int64        `json:"generation"`
	Suspend                    *bool        `json:"suspend"`
	SuccessfulJobsHistoryLimit *int32       `json:"successfulJobsHistoryLimit"`
	FailedJobsHistoryLimit     *int32       `json:"failedJobsHistoryLimit"`
	LastScheduleTime           *v1.Time     `json:"lastScheduleTime"`
	LastSuccessfulTime         *v1.Time     `json:"lastSuccessfulTime"`
	TTLSecondsAfterFinished    *int32       `json:"ttlSecondsAfterFinished"`
	Containers                 []*Container `json:"containers"`
}
