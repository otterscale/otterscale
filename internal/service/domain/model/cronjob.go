package model

import "time"

type CronJob struct {
	UID                        string       `json:"uid"`
	Name                       string       `json:"name"`
	Schedule                   string       `json:"schedule"`
	Generation                 int64        `json:"generation"`
	Suspend                    *bool        `json:"suspend"`
	SuccessfulJobsHistoryLimit *int32       `json:"successfulJobsHistoryLimit"`
	FailedJobsHistoryLimit     *int32       `json:"failedJobsHistoryLimit"`
	LastScheduleTime           time.Time    `json:"lastScheduleTime"`
	LastSuccessfulTime         time.Time    `json:"lastSuccessfulTime"`
	TTLSecondsAfterFinished    *int32       `json:"ttlSecondsAfterFinished"`
	Containers                 []*Container `json:"containers"`
}
