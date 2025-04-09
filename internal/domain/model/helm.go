package model

import "helm.sh/helm/v3/pkg/repo"

type HelmRepo struct {
	*repo.Entry
	ChartNames []string
}
