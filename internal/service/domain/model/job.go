package model

import "time"

type Job struct {
	UID                     string       `json:"uid"`
	Name                    string       `json:"name"`
	Succeeded               bool         `json:"succeeded"`
	StartTime               time.Time    `json:"startTime"`
	CompletionTime          time.Time    `json:"completionTime"`
	TTLSecondsAfterFinished *int32       `json:"ttlSecondsAfterFinished"`
	Containers              []*Container `json:"containers"`
}
