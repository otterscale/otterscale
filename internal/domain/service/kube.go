package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/sync/errgroup"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"

	"github.com/juju/juju/api/client/application"
	"github.com/moby/moby/pkg/namesgenerator"

	"github.com/openhdc/openhdc/internal/domain/model"
)

// KubeClient manages Kubernetes client connections
type KubeClient interface {
	Exists(key string) bool
	Add(key string, cfg *rest.Config) error
}

// KubeApps handles deployment-related operations
type KubeApps interface {
	ListDeployments(ctx context.Context, key, namespace string) ([]appsv1.Deployment, error)
	GetDeployment(ctx context.Context, key, namespace, name string) (*appsv1.Deployment, error)
	ListStatefulSets(ctx context.Context, key, namespace string) ([]appsv1.StatefulSet, error)
	GetStatefulSet(ctx context.Context, key, namespace, name string) (*appsv1.StatefulSet, error)
	ListDaemonSets(ctx context.Context, key, namespace string) ([]appsv1.DaemonSet, error)
	GetDaemonSet(ctx context.Context, key, namespace, name string) (*appsv1.DaemonSet, error)
}

// KubeBatch handles batch job operations
type KubeBatch interface {
	GetCronJob(ctx context.Context, key, namespace, name string) (*batchv1.CronJob, error)
	CreateCronJob(ctx context.Context, key, namespace, name, image, schedule string) (*batchv1.CronJob, error)
	UpdateCronJob(ctx context.Context, key, namespace, name, image, schedule string) (*batchv1.CronJob, error)
	DeleteCronJob(ctx context.Context, key, namespace, name string) error
	ListJobsFromCronJob(ctx context.Context, key, namespace string, cronJob *batchv1.CronJob) (*batchv1.JobList, error)
	CreateJobFromCronJob(ctx context.Context, key, namespace string, cronJob *batchv1.CronJob, createdBy string) (*batchv1.Job, error)
}

// KubeCore handles core Kubernetes resource operations
type KubeCore interface {
	GetNamespace(ctx context.Context, key, name string) (*corev1.Namespace, error)
	CreateNamespace(ctx context.Context, key, name string) (*corev1.Namespace, error)
	ListServices(ctx context.Context, key, namespace string) ([]corev1.Service, error)
	ListPods(ctx context.Context, key, namespace string) ([]corev1.Pod, error)
	ListPersistentVolumeClaims(ctx context.Context, key, namespace string) ([]corev1.PersistentVolumeClaim, error)
}

// KubeStorage handles storage-related operations
type KubeStorage interface {
	ListStorageClasses(ctx context.Context, key string) ([]storagev1.StorageClass, error)
}

// KubeHelm handles helm-related operations
type KubeHelm interface {
	ListReleases(key, namespace string) ([]*release.Release, error)
	InstallRelease(key, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error)
	UninstallRelease(key, namespace, name string, dryRun bool) (*release.Release, error)
	UpgradeRelease(key, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error)
	RollbackRelease(key, namespace, name string, dryRun bool) error
	GetChartInfo(chartRef string, format action.ShowOutputFormat) (string, error)
	ListChartVersions(ctx context.Context) (map[string]repo.ChartVersions, error)
}

// KubeService orchestrates Kubernetes operations
type KubeService struct {
	client      KubeClient
	apps        KubeApps
	batch       KubeBatch
	core        KubeCore
	storage     KubeStorage
	helm        KubeHelm
	model       JujuModel
	jujuClient  JujuClient
	application JujuApplication
}

// NewKubeService creates a new KubeService instance
func NewKubeService(
	client KubeClient,
	apps KubeApps,
	batch KubeBatch,
	core KubeCore,
	storage KubeStorage,
	helm KubeHelm,
	model JujuModel,
	jujuClient JujuClient,
	application JujuApplication,
) *KubeService {
	return &KubeService{
		client:      client,
		apps:        apps,
		batch:       batch,
		core:        core,
		storage:     storage,
		helm:        helm,
		model:       model,
		jujuClient:  jujuClient,
		application: application,
	}
}

// ensureClient ensures a Kubernetes client exists for the specified cluster
func (s *KubeService) ensureClient(ctx context.Context, uuid, cluster string) (string, error) {
	key := kubeMapKey(uuid, cluster)

	// Check if client already exists
	if ok := s.client.Exists(key); ok {
		return key, nil // Client already exists
	}

	// Create new client config
	cfg, err := s.newConfig(ctx, uuid, cluster)
	if err != nil {
		return "", err
	}

	if err := s.client.Add(key, cfg); err != nil {
		return "", err
	}

	return key, nil
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

func (s *KubeService) isKeyNotFoundError(err error) bool {
	statusErr, _ := err.(*apierrors.StatusError)
	return statusErr != nil && statusErr.Status().Code == http.StatusNotFound
}

// ListApplications retrieves all applications from the Kubernetes cluster
func (s *KubeService) ListApplications(ctx context.Context, uuid, cluster, namespace, name string) (*model.Applications, error) {
	key, err := s.ensureClient(ctx, uuid, cluster)
	if err != nil {
		return nil, err
	}

	result := &model.Applications{}
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if namespace == "" {
			deployments, err := s.apps.ListDeployments(ctx, key, namespace)
			if err == nil {
				result.Deployments = deployments
			}
			return err
		}
		deployment, err := s.apps.GetDeployment(ctx, key, namespace, name)
		if err == nil {
			result.Deployments = []appsv1.Deployment{*deployment}
		} else if s.isKeyNotFoundError(err) {
			return nil
		}
		return err
	})

	eg.Go(func() error {
		if namespace == "" {
			statefulSets, err := s.apps.ListStatefulSets(ctx, key, namespace)
			if err == nil {
				result.StatefulSets = statefulSets
			}
			return err
		}
		statefulSet, err := s.apps.GetStatefulSet(ctx, key, namespace, name)
		if err == nil {
			result.StatefulSets = []appsv1.StatefulSet{*statefulSet}
		} else if s.isKeyNotFoundError(err) {
			return nil
		}
		return err
	})

	eg.Go(func() error {
		if namespace == "" {
			daemonSets, err := s.apps.ListDaemonSets(ctx, key, namespace)
			if err == nil {
				result.DaemonSets = daemonSets
			}
			return err
		}
		daemonSet, err := s.apps.GetDaemonSet(ctx, key, namespace, name)
		if err == nil {
			result.DaemonSets = []appsv1.DaemonSet{*daemonSet}
		} else if s.isKeyNotFoundError(err) {
			return nil
		}
		return err
	})

	eg.Go(func() error {
		services, err := s.core.ListServices(ctx, key, namespace)
		if err == nil {
			result.Services = services
		}
		return err
	})

	eg.Go(func() error {
		pods, err := s.core.ListPods(ctx, key, namespace)
		if err == nil {
			result.Pods = pods
		}
		return err
	})

	eg.Go(func() error {
		pvcs, err := s.core.ListPersistentVolumeClaims(ctx, key, namespace)
		if err == nil {
			result.PersistentVolumeClaims = pvcs
		}
		return err
	})

	eg.Go(func() error {
		storageClasses, err := s.storage.ListStorageClasses(ctx, key)
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

func (s *KubeService) ListReleases(ctx context.Context) ([]*model.Release, error) {
	ms, err := s.model.List(ctx)
	if err != nil {
		return nil, err
	}

	type T struct{ name, uuid, cluster string }

	eg1, egCtx := errgroup.WithContext(ctx)
	results := make([][]T, len(ms))
	for i, m := range ms {
		eg1.Go(func() error {
			status, err := s.jujuClient.Status(egCtx, m.UUID, []string{"application", "*"})
			if err != nil {
				return err
			}
			for name := range status.Applications {
				if strings.Contains(status.Applications[name].Charm, "kubernetes-worker") {
					results[i] = append(results[i], T{
						name:    m.Name,
						uuid:    m.UUID,
						cluster: name,
					})
				}
			}
			return nil
		})
	}

	if err := eg1.Wait(); err != nil {
		return nil, err
	}

	rs := []T{}
	for _, cs := range results {
		rs = append(rs, cs...)
	}

	eg2, egCtx := errgroup.WithContext(ctx)
	res := make([][]*model.Release, len(rs))
	for i, r := range rs {
		eg2.Go(func() error {
			key, err := s.ensureClient(egCtx, r.uuid, r.cluster)
			if err != nil {
				return err
			}
			rels, err := s.helm.ListReleases(key, "")
			if err != nil {
				return err
			}
			for _, rel := range rels {
				res[i] = append(res[i], &model.Release{
					ModelName:   r.name,
					ModelUUID:   r.uuid,
					ClusterName: r.cluster,
					Release:     rel,
				})
			}
			return nil
		})
	}

	if err := eg2.Wait(); err != nil {
		return nil, err
	}

	ret := []*model.Release{}
	for _, re := range res {
		ret = append(ret, re...)
	}

	return ret, nil
}

func (s *KubeService) InstallRelease(ctx context.Context, uuid, cluster, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error) {
	key, err := s.ensureClient(ctx, uuid, cluster)
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = randomName()
	}
	return s.helm.InstallRelease(key, namespace, name, dryRun, chartRef, values)
}

func (s *KubeService) UninstallRelease(ctx context.Context, uuid, cluster, namespace, name string, dryRun bool) (*release.Release, error) {
	key, err := s.ensureClient(ctx, uuid, cluster)
	if err != nil {
		return nil, err
	}
	return s.helm.UninstallRelease(key, namespace, name, dryRun)
}

func (s *KubeService) UpgradeRelease(ctx context.Context, uuid, cluster, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error) {
	key, err := s.ensureClient(ctx, uuid, cluster)
	if err != nil {
		return nil, err
	}
	return s.helm.UpgradeRelease(key, namespace, name, dryRun, chartRef, values)
}

func (s *KubeService) RollbackRelease(ctx context.Context, uuid, cluster, namespace, name string, dryRun bool) error {
	key, err := s.ensureClient(ctx, uuid, cluster)
	if err != nil {
		return err
	}
	return s.helm.RollbackRelease(key, namespace, name, dryRun)
}

func (s *KubeService) ListCharts(ctx context.Context) (map[string]repo.ChartVersions, error) {
	return s.helm.ListChartVersions(ctx)
}

func (s *KubeService) GetChart(ctx context.Context, name string) (repo.ChartVersions, error) {
	m, err := s.helm.ListChartVersions(ctx)
	if err != nil {
		return nil, err
	}
	for k, v := range m {
		if k != name {
			continue
		}
		return v, nil
	}
	return nil, fmt.Errorf("chart %q not found", name)
}

func (s *KubeService) GetChartInfo(chartRef string) (values string, readme string, err error) {
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		info, err := s.helm.GetChartInfo(chartRef, action.ShowValues)
		if err == nil {
			values = info
		}
		return err
	})
	eg.Go(func() error {
		info, err := s.helm.GetChartInfo(chartRef, action.ShowReadme)
		if err == nil {
			readme = info
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return "", "", err
	}
	return
}

func randomName() string {
	return strings.ReplaceAll(namesgenerator.GetRandomName(0), "_", "-")
}

func kubeMapKey(uuid, cluster string) string {
	sha := sha256.New()
	sha.Write([]byte(uuid))
	sha.Write([]byte(cluster))
	return string(sha.Sum(nil))
}
