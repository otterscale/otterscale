package service

import (
	"context"
	"encoding/json"
	"errors"

	"golang.org/x/sync/errgroup"

	appv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/juju/juju/api/client/application"

	"github.com/openhdc/openhdc/internal/domain/model"
)

// KubeClient manages Kubernetes client connections
type KubeClient interface {
	Get(cluster string) (*kubernetes.Clientset, error)
	Add(cluster string, cfg *rest.Config) error
}

// KubeApps handles deployment-related operations
type KubeApps interface {
	ListDeployments(ctx context.Context, cluster, namespace string) ([]appv1.Deployment, error)
	GetDeployment(ctx context.Context, cluster, namespace, name string) (*appv1.Deployment, error)
}

// KubeBatch handles batch job operations
type KubeBatch interface {
	GetCronJob(ctx context.Context, cluster, namespace, name string) (*batchv1.CronJob, error)
	CreateCronJob(ctx context.Context, cluster, namespace, name, image, schedule string) (*batchv1.CronJob, error)
	UpdateCronJob(ctx context.Context, cluster, namespace, name, image, schedule string) (*batchv1.CronJob, error)
	DeleteCronJob(ctx context.Context, cluster, namespace, name string) error
	ListJobsFromCronJob(ctx context.Context, cluster, namespace string, cronJob *batchv1.CronJob) (*batchv1.JobList, error)
	CreateJobFromCronJob(ctx context.Context, cluster, namespace string, cronJob *batchv1.CronJob, createdBy string) (*batchv1.Job, error)
}

// KubeCore handles core Kubernetes resource operations
type KubeCore interface {
	GetNamespace(ctx context.Context, cluster, name string) (*corev1.Namespace, error)
	CreateNamespace(ctx context.Context, cluster, name string) (*corev1.Namespace, error)
	ListServices(ctx context.Context, cluster, namespace string) ([]corev1.Service, error)
	ListPods(ctx context.Context, cluster, namespace string) ([]corev1.Pod, error)
	ListPersistentVolumeClaims(ctx context.Context, cluster, namespace string) ([]corev1.PersistentVolumeClaim, error)
}

// KubeStorage handles storage-related operations
type KubeStorage interface {
	ListStorageClasses(ctx context.Context, cluster string) ([]storagev1.StorageClass, error)
}

// KubeService orchestrates Kubernetes operations
type KubeService struct {
	client      KubeClient
	apps        KubeApps
	batch       KubeBatch
	core        KubeCore
	storage     KubeStorage
	application JujuApplication
}

// NewKubeService creates a new KubeService instance
func NewKubeService(
	client KubeClient,
	apps KubeApps,
	batch KubeBatch,
	core KubeCore,
	storage KubeStorage,
	application JujuApplication,
) *KubeService {
	return &KubeService{
		client:      client,
		apps:        apps,
		batch:       batch,
		core:        core,
		storage:     storage,
		application: application,
	}
}

// ensureClient ensures a Kubernetes client exists for the specified cluster
func (s *KubeService) ensureClient(ctx context.Context, uuid, cluster string) error {
	// Check if client already exists
	if _, err := s.client.Get(cluster); err == nil {
		return nil // Client already exists
	}

	// Create new client config
	cfg, err := s.newConfig(ctx, uuid, cluster)
	if err != nil {
		return err
	}

	return s.client.Add(cluster, cfg)
}

// newConfig creates a new Kubernetes client configuration
func (s *KubeService) newConfig(ctx context.Context, uuid, name string) (*rest.Config, error) {
	unit, err := s.application.GetLeader(ctx, uuid, name)
	if err != nil {
		return nil, err
	}

	unitInfo, err := s.application.GetUnitInfo(ctx, uuid, unit)
	if err != nil {
		return nil, err
	}

	endpoint, err := extractEndpointFromUnitInfo(unitInfo)
	if err != nil {
		return nil, err
	}

	clientToken, err := extractClientTokenFromUnitInfo(unitInfo)
	if err != nil {
		return nil, err
	}

	return &rest.Config{
		Host: endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
		BearerToken: clientToken,
	}, nil
}

// extractEndpointFromUnitInfo extracts the API endpoint from unit info
func extractEndpointFromUnitInfo(unitInfo *application.UnitInfo) (string, error) {
	var endpoints []string

	// Extract endpoints from relation data
	for _, erd := range unitInfo.RelationData {
		for _, rd := range erd.UnitRelationData {
			if endpointData, ok := rd.UnitData["api-endpoints"]; ok && endpointData != nil {
				if endpointStr, ok := endpointData.(string); ok {
					if err := json.Unmarshal([]byte(endpointStr), &endpoints); err != nil {
						return "", err
					}
				}
			}
		}
	}

	if len(endpoints) > 0 {
		return endpoints[0], nil
	}

	return "", errors.New("endpoint not found")
}

// extractClientTokenFromUnitInfo extracts the client token from unit info
func extractClientTokenFromUnitInfo(unitInfo *application.UnitInfo) (string, error) {
	credentials := make(map[string]model.ControlPlaneCredential)

	// Extract credentials from relation data
	for _, erd := range unitInfo.RelationData {
		for _, rd := range erd.UnitRelationData {
			if credsData, ok := rd.UnitData["creds"]; ok && credsData != nil {
				if credsStr, ok := credsData.(string); ok {
					if err := json.Unmarshal([]byte(credsStr), &credentials); err != nil {
						return "", err
					}
				}
			}
		}
	}

	for _, cred := range credentials {
		return cred.ClientToken, nil
	}

	return "", errors.New("token not found")
}

// ListApplications retrieves all applications from the Kubernetes cluster
func (s *KubeService) ListApplications(ctx context.Context, uuid, cluster, namespace, name string) (*model.Applications, error) {
	if err := s.ensureClient(ctx, uuid, cluster); err != nil {
		return nil, err
	}

	result := &model.Applications{}
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if namespace == "" {
			deployments, err := s.apps.ListDeployments(ctx, cluster, namespace)
			if err == nil {
				result.Deployments = deployments
			}
			return err
		}
		deployment, err := s.apps.GetDeployment(ctx, cluster, namespace, name)
		if err == nil {
			result.Deployments = []appv1.Deployment{*deployment}
		}
		return err
	})

	eg.Go(func() error {
		services, err := s.core.ListServices(ctx, cluster, namespace)
		if err == nil {
			result.Services = services
		}
		return err
	})

	eg.Go(func() error {
		pods, err := s.core.ListPods(ctx, cluster, namespace)
		if err == nil {
			result.Pods = pods
		}
		return err
	})

	eg.Go(func() error {
		pvcs, err := s.core.ListPersistentVolumeClaims(ctx, cluster, namespace)
		if err == nil {
			result.PersistentVolumeClaims = pvcs
		}
		return err
	})

	eg.Go(func() error {
		storageClasses, err := s.storage.ListStorageClasses(ctx, cluster)
		if err == nil {
			result.StorageClasses = storageClasses
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}
