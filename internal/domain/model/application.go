package model

import (
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
)

type PersistentVolumeClaim struct {
	*corev1.PersistentVolumeClaim
	*storagev1.StorageClass
}

type Application struct {
	Type                   string
	Name                   string
	Namespace              string
	Labels                 map[string]string
	Replicas               *int32
	Containers             []corev1.Container
	Services               []corev1.Service
	Pods                   []corev1.Pod
	PersistentVolumeClaims []PersistentVolumeClaim
}

type Kubernetes struct {
	ScopeName    string
	ScopeUUID    string
	FacilityName string
}

type Release struct {
	ScopeName    string
	ScopeUUID    string
	FacilityName string
	*release.Release
}

type Chart struct {
	Name     string
	Versions repo.ChartVersions
}

type ChartMetadata struct {
	ReadmeMD   string
	ValuesYAML string
}

type ControlPlaneCredential struct {
	ClientToken  string `json:"client_token"`
	KubeletToken string `json:"kubelet_token"`
	ProxyToken   string `json:"proxy_token"`
	Scope        string `json:"scope"`
}
