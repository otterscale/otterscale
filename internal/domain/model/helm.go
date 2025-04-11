package model

import "helm.sh/helm/v3/pkg/release"

type Release struct {
	ModelName   string
	ModelUUID   string
	ClusterName string
	*release.Release
}
