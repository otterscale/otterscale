package model

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Container struct {
	Image           string `json:"image"`
	ImagePullPolicy string `json:"imagePullPolicy"`
}

type CronJob struct {
	UID                        string       `json:"uid"`
	Name                       string       `json:"name"`
	Schedule                   string       `json:"schedule"`
	Generation                 int64        `json:"generation"`
	Suspend                    *bool        `json:"suspend"`
	SuccessfulJobsHistoryLimit *int32       `json:"successfulJobsHistoryLimit"`
	FailedJobsHistoryLimit     *int32       `json:"failedJobsHistoryLimit"`
	LastScheduleTime           *metav1.Time `json:"lastScheduleTime"`
	LastSuccessfulTime         *metav1.Time `json:"lastSuccessfulTime"`
	TTLSecondsAfterFinished    *int32       `json:"ttlSecondsAfterFinished"`
	Containers                 []*Container `json:"containers"`
}

type Job struct {
	UID                     string       `json:"uid"`
	Name                    string       `json:"name"`
	Succeeded               bool         `json:"succeeded"`
	StartTime               *metav1.Time `json:"startTime"`
	CompletionTime          *metav1.Time `json:"completionTime"`
	TTLSecondsAfterFinished *int32       `json:"ttlSecondsAfterFinished"`
	Containers              []*Container `json:"containers"`
}

type ControlPlaneCredential struct {
	ClientToken  string `json:"client_token"`
	KubeletToken string `json:"kubelet_token"`
	ProxyToken   string `json:"proxy_token"`
	Scope        string `json:"scope"`
}

type Applications struct {
	Deployments            []appv1.Deployment
	Services               []corev1.Service
	Pods                   []corev1.Pod
	PersistentVolumeClaims []corev1.PersistentVolumeClaim
	StorageClasses         []storagev1.StorageClass
}
