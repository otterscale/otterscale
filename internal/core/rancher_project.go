package core

import (
	"context"
	"errors"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"
)

// ErrRancherProjectCacheNotReady indicates that the Rancher Project
// informer has not completed its initial list.
var ErrRancherProjectCacheNotReady = errors.New("rancher project cache is not ready")

// RancherProject is the subset of Rancher Project metadata exposed to
// callers that create managed clusters.
type RancherProject struct {
	ID          string
	DisplayName string
}

// RancherProjectStore provides last-known Rancher Projects from the
// management-cluster informer cache.
type RancherProjectStore interface {
	ListProjects(context.Context) ([]RancherProject, error)
	HasProject(context.Context, string) (bool, error)
}

// ParseRancherProjectID validates and splits a full Rancher Project ID.
func ParseRancherProjectID(id string) (clusterID, projectName string, err error) {
	clusterID, projectName, ok := strings.Cut(id, ":")
	if !ok || clusterID == "" || projectName == "" || strings.Contains(projectName, ":") {
		return "", "", &ErrInvalidInput{
			Field:   "rancher_project_id",
			Message: "must be in <cluster-id>:<project-name> form",
		}
	}
	if problems := validation.IsDNS1123Label(clusterID); len(problems) > 0 {
		return "", "", &ErrInvalidInput{
			Field:   "rancher_project_id",
			Message: "cluster ID " + strings.Join(problems, ", "),
		}
	}
	if problems := validation.IsDNS1123Subdomain(projectName); len(problems) > 0 {
		return "", "", &ErrInvalidInput{
			Field:   "rancher_project_id",
			Message: "project name " + strings.Join(problems, ", "),
		}
	}
	return clusterID, projectName, nil
}
