package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/juju/juju/api/client/application"
	"golang.org/x/sync/errgroup"

	appv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/openhdc/openhdc/internal/domain/model"
)

type KubeClient interface {
	Get(cluster string) (*kubernetes.Clientset, error)
	Add(cluster string, cfg *rest.Config) error
}

type KubeApps interface {
	ListDeployments(ctx context.Context, cluster, namespace string) ([]appv1.Deployment, error)
}

type KubeBatch interface {
	GetCronJob(ctx context.Context, cluster, namespace, name string) (*batchv1.CronJob, error)
	CreateCronJob(ctx context.Context, cluster, namespace, name, image, schedule string) (*batchv1.CronJob, error)
	UpdateCronJob(ctx context.Context, cluster, namespace, name, image, schedule string) (*batchv1.CronJob, error)
	DeleteCronJob(ctx context.Context, cluster, namespace, name string) error
	ListJobsFromCronJob(ctx context.Context, cluster, namespace string, cronJob *batchv1.CronJob) (*batchv1.JobList, error)
	CreateJobFromCronJob(ctx context.Context, cluster, namespace string, cronJob *batchv1.CronJob, createdBy string) (*batchv1.Job, error)
}

type KubeCore interface {
	GetNamespace(ctx context.Context, cluster, name string) (*corev1.Namespace, error)
	CreateNamespace(ctx context.Context, cluster, name string) (*corev1.Namespace, error)
	ListServices(ctx context.Context, cluster, namespace string) ([]corev1.Service, error)
	ListPods(ctx context.Context, cluster, namespace string) ([]corev1.Pod, error)
	ListPersistentVolumeClaims(ctx context.Context, cluster, namespace string) ([]corev1.PersistentVolumeClaim, error)
}

type KubeStorage interface {
	ListStorageClasses(ctx context.Context, cluster string) ([]storagev1.StorageClass, error)
}

type KubeService struct {
	client      KubeClient
	apps        KubeApps
	batch       KubeBatch
	core        KubeCore
	storage     KubeStorage
	application JujuApplication
}

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

func extractEndpointFromUnitInfo(unitInfo *application.UnitInfo) (string, error) {
	endpoints := []string{}

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

func (s *KubeService) ListApplications(ctx context.Context, uuid, cluster string) (*model.Applications, error) {
	if err := s.ensureClient(ctx, uuid, cluster); err != nil {
		return nil, err
	}
	var (
		deployments            []appv1.Deployment
		services               []corev1.Service
		pods                   []corev1.Pod
		persistentVolumeClaims []corev1.PersistentVolumeClaim
		storageClasses         []storagev1.StorageClass
	)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		dpl, err := s.apps.ListDeployments(ctx, cluster, "")
		if err == nil {
			deployments = dpl
		}
		return err
	})
	eg.Go(func() error {
		svc, err := s.core.ListServices(ctx, cluster, "")
		if err == nil {
			services = svc
		}
		return err
	})
	eg.Go(func() error {
		pod, err := s.core.ListPods(ctx, cluster, "")
		if err == nil {
			pods = pod
		}
		return err
	})
	eg.Go(func() error {
		pvc, err := s.core.ListPersistentVolumeClaims(ctx, cluster, "")
		if err == nil {
			persistentVolumeClaims = pvc
		}
		return err
	})
	eg.Go(func() error {
		sc, err := s.storage.ListStorageClasses(ctx, cluster)
		if err == nil {
			storageClasses = sc
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return &model.Applications{
		Deployments:            deployments,
		Services:               services,
		Pods:                   pods,
		PersistentVolumeClaims: persistentVolumeClaims,
		StorageClasses:         storageClasses,
	}, nil
}
