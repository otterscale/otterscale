package model

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Job struct {
	UID                     string       `json:"uid"`
	Name                    string       `json:"name"`
	Succeeded               bool         `json:"succeeded"`
	StartTime               *v1.Time     `json:"startTime"`
	CompletionTime          *v1.Time     `json:"completionTime"`
	TTLSecondsAfterFinished *int32       `json:"ttlSecondsAfterFinished"`
	Containers              []*Container `json:"containers"`
}
