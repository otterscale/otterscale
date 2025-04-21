package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/moby/moby/pkg/namesgenerator"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"helm.sh/helm/v3/pkg/action"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"

	"github.com/openhdc/openhdc/internal/domain/model"
)

const (
	appTypeDeployment  = "Deployment"
	appTypeStatefulSet = "StatefulSet"
	appTypeDaemonSet   = "DaemonSet"
)

const kubernetesCharmName = "kubernetes-worker"

func (s *NexusService) ListApplications(ctx context.Context, uuid, facility string) ([]model.Application, error) {
	if err := s.setKubernetesClient(ctx, uuid, facility); err != nil {
		return nil, err
	}

	var (
		deployments            []appsv1.Deployment
		statefulSets           []appsv1.StatefulSet
		daemonSets             []appsv1.DaemonSet
		services               []corev1.Service
		pods                   []corev1.Pod
		persistentVolumeClaims []corev1.PersistentVolumeClaim
		storageClasses         []storagev1.StorageClass
	)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		v, err := s.apps.ListDeployments(ctx, uuid, facility, "")
		if err == nil {
			deployments = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.apps.ListStatefulSets(ctx, uuid, facility, "")
		if err == nil {
			statefulSets = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.apps.ListDaemonSets(ctx, uuid, facility, "")
		if err == nil {
			daemonSets = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.core.ListServices(ctx, uuid, facility, "")
		if err == nil {
			services = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.core.ListPods(ctx, uuid, facility, "")
		if err == nil {
			pods = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.core.ListPersistentVolumeClaims(ctx, uuid, facility, "")
		if err == nil {
			persistentVolumeClaims = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.storage.ListStorageClasses(ctx, uuid, facility)
		if err == nil {
			storageClasses = v
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	scm := toStorageClassMap(storageClasses)

	apps := []model.Application{}
	for i := range deployments {
		d := deployments[i]
		app, err := toApplication(d.Spec.Selector, appTypeDeployment, d.Name, d.Namespace, &d.ObjectMeta, d.Labels, d.Spec.Replicas, d.Spec.Template.Spec.Containers, d.Spec.Template.Spec.Volumes, services, pods, persistentVolumeClaims, scm)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}
	for i := range statefulSets {
		d := deployments[i]
		app, err := toApplication(d.Spec.Selector, appTypeStatefulSet, d.Name, d.Namespace, &d.ObjectMeta, d.Labels, d.Spec.Replicas, d.Spec.Template.Spec.Containers, d.Spec.Template.Spec.Volumes, services, pods, persistentVolumeClaims, scm)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}
	for i := range daemonSets {
		d := deployments[i]
		app, err := toApplication(d.Spec.Selector, appTypeDaemonSet, d.Name, d.Namespace, &d.ObjectMeta, d.Labels, nil, d.Spec.Template.Spec.Containers, d.Spec.Template.Spec.Volumes, services, pods, persistentVolumeClaims, scm)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}

	return apps, nil
}

func (s *NexusService) GetApplication(ctx context.Context, uuid, facility, namespace, name string) (*model.Application, error) {
	if err := s.setKubernetesClient(ctx, uuid, facility); err != nil {
		return nil, err
	}

	var (
		deployment             *appsv1.Deployment
		statefulSet            *appsv1.StatefulSet
		daemonSet              *appsv1.DaemonSet
		services               []corev1.Service
		pods                   []corev1.Pod
		persistentVolumeClaims []corev1.PersistentVolumeClaim
		storageClasses         []storagev1.StorageClass
	)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		v, err := s.apps.GetDeployment(ctx, uuid, facility, namespace, name)
		if err == nil {
			deployment = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.apps.GetStatefulSet(ctx, uuid, facility, namespace, name)
		if err == nil {
			statefulSet = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.apps.GetDaemonSet(ctx, uuid, facility, namespace, name)
		if err == nil {
			daemonSet = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.core.ListServices(ctx, uuid, facility, namespace)
		if err == nil {
			services = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.core.ListPods(ctx, uuid, facility, namespace)
		if err == nil {
			pods = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.core.ListPersistentVolumeClaims(ctx, uuid, facility, namespace)
		if err == nil {
			persistentVolumeClaims = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.storage.ListStorageClasses(ctx, uuid, facility)
		if err == nil {
			storageClasses = v
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	scm := toStorageClassMap(storageClasses)

	if deployment != nil {
		return toApplication(deployment.Spec.Selector, appTypeDeployment, deployment.Name, deployment.Namespace, &deployment.ObjectMeta, deployment.Labels, deployment.Spec.Replicas, deployment.Spec.Template.Spec.Containers, deployment.Spec.Template.Spec.Volumes, services, pods, persistentVolumeClaims, scm)
	} else if statefulSet != nil {
		return toApplication(statefulSet.Spec.Selector, appTypeDeployment, statefulSet.Name, statefulSet.Namespace, &statefulSet.ObjectMeta, statefulSet.Labels, statefulSet.Spec.Replicas, statefulSet.Spec.Template.Spec.Containers, statefulSet.Spec.Template.Spec.Volumes, services, pods, persistentVolumeClaims, scm)
	} else if daemonSet != nil {
		return toApplication(daemonSet.Spec.Selector, appTypeDeployment, daemonSet.Name, daemonSet.Namespace, &daemonSet.ObjectMeta, daemonSet.Labels, nil, daemonSet.Spec.Template.Spec.Containers, daemonSet.Spec.Template.Spec.Volumes, services, pods, persistentVolumeClaims, scm)
	}

	return nil, status.Errorf(codes.NotFound, "application %q in namespace %q not found", name, namespace)
}

func (s *NexusService) ListReleases(ctx context.Context) ([]model.Release, error) {
	ks, err := s.listKubernetes(ctx)
	if err != nil {
		return nil, err
	}
	return s.listReleases(ctx, ks)
}

func (s *NexusService) CreateRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool, chartRef, valuesYAML string) (*model.Release, error) {
	values := map[string]any{}
	if err := yaml.Unmarshal([]byte(valuesYAML), &values); err != nil {
		return nil, err
	}
	if err := s.setKubernetesClient(ctx, uuid, facility); err != nil {
		return nil, err
	}
	if name == "" {
		name = randomName()
	}
	r, err := s.helm.InstallRelease(uuid, facility, namespace, name, dryRun, chartRef, values)
	if err != nil {
		return nil, err
	}
	return &model.Release{Release: r}, nil
}

func (s *NexusService) UpdateRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool, chartRef, valuesYAML string) (*model.Release, error) {
	values := map[string]any{}
	if err := yaml.Unmarshal([]byte(valuesYAML), &values); err != nil {
		return nil, err
	}
	if err := s.setKubernetesClient(ctx, uuid, facility); err != nil {
		return nil, err
	}
	r, err := s.helm.UpgradeRelease(uuid, facility, namespace, name, dryRun, chartRef, values)
	if err != nil {
		return nil, err
	}
	return &model.Release{Release: r}, nil
}

func (s *NexusService) DeleteRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool) error {
	if err := s.setKubernetesClient(ctx, uuid, facility); err != nil {
		return err
	}
	if _, err := s.helm.UninstallRelease(uuid, facility, namespace, name, dryRun); err != nil {
		return err
	}
	return nil
}

func (s *NexusService) RollbackRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool) error {
	if err := s.setKubernetesClient(ctx, uuid, facility); err != nil {
		return err
	}
	return s.helm.RollbackRelease(uuid, facility, namespace, name, dryRun)
}

func (s *NexusService) ListCharts(ctx context.Context) ([]model.Chart, error) {
	fs, err := s.helm.ListChartVersions(ctx)
	if err != nil {
		return nil, err
	}
	cs := []model.Chart{}
	for _, f := range fs {
		for key := range f.Entries {
			cs = append(cs, model.Chart{
				Name:     key,
				Versions: f.Entries[key],
			})
		}
	}
	return cs, nil
}

func (s *NexusService) GetChart(ctx context.Context, name string) (*model.Chart, error) {
	fs, err := s.helm.ListChartVersions(ctx)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		for key := range f.Entries {
			if key != name {
				continue
			}
			return &model.Chart{
				Name:     key,
				Versions: f.Entries[key],
			}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "chart %q not found", name)
}

func (s *NexusService) GetChartMetadata(ctx context.Context, chartRef string) (*model.ChartMetadata, error) {
	md := &model.ChartMetadata{}
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		v, err := s.helm.ShowChart(chartRef, action.ShowValues)
		if err == nil {
			md.ValuesYAML = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := s.helm.ShowChart(chartRef, action.ShowReadme)
		if err == nil {
			md.ReadmeMD = v
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return md, nil
}

func (s *NexusService) setKubernetesClient(ctx context.Context, uuid, facility string) error {
	if s.kubernetes.Exists(uuid, facility) {
		return nil
	}
	unit, err := s.facility.GetLeader(ctx, uuid, facility)
	if err != nil {
		return err
	}
	unitInfo, err := s.facility.GetUnitInfo(ctx, uuid, unit)
	if err != nil {
		return err
	}
	cfg, err := newKubernetesConfig(unitInfo)
	if err != nil {
		return err
	}
	return s.kubernetes.Set(uuid, facility, cfg)
}

func (s *NexusService) listKubernetes(ctx context.Context) ([]model.Kubernetes, error) {
	scopes, err := s.scope.List(ctx)
	if err != nil {
		return nil, err
	}

	eg, ctx := errgroup.WithContext(ctx)
	result := make([][]model.Kubernetes, len(scopes))
	for i := range scopes {
		scope := scopes[i]
		eg.Go(func() error {
			fs, err := s.ListFacilities(ctx, scope.UUID)
			if err != nil {
				return err
			}
			for j := range fs {
				facility := fs[j]
				if strings.Contains(facility.Status.Charm, kubernetesCharmName) {
					result[i] = append(result[i], model.Kubernetes{
						ScopeName:    scope.Name,
						ScopeUUID:    scope.UUID,
						FacilityName: facility.Name,
					})
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	ks := []model.Kubernetes{}
	for _, k := range result {
		ks = append(ks, k...)
	}
	return ks, nil
}

func (s *NexusService) listReleases(ctx context.Context, ks []model.Kubernetes) ([]model.Release, error) {
	eg, ctx := errgroup.WithContext(ctx)
	result := make([][]model.Release, len(ks))
	for i := range ks {
		k := ks[i]
		eg.Go(func() error {
			if err := s.setKubernetesClient(ctx, k.ScopeUUID, k.FacilityName); err != nil {
				return err
			}
			rels, err := s.helm.ListReleases(k.ScopeUUID, k.FacilityName, "")
			if err != nil {
				return err
			}
			for _, rel := range rels {
				result[i] = append(result[i], model.Release{
					ScopeName:    k.ScopeName,
					ScopeUUID:    k.ScopeUUID,
					FacilityName: k.FacilityName,
					Release:      &rel,
				})
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	rs := []model.Release{}
	for _, r := range result {
		rs = append(rs, r...)
	}
	return rs, nil
}

func newKubernetesConfig(unitInfo *model.UnitInfo) (*rest.Config, error) {
	endpoint, err := extractEndpoint(unitInfo)
	if err != nil {
		return nil, err
	}

	clientToken, err := extractClientToken(unitInfo)
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

func extractEndpoint(unitInfo *model.UnitInfo) (string, error) {
	var endpoints []string
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
	return "", status.Error(codes.NotFound, "endpoint not found")
}

func extractClientToken(unitInfo *model.UnitInfo) (string, error) {
	credentials := make(map[string]model.ControlPlaneCredential)
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

func toApplication(ls *metav1.LabelSelector, appType, name, namespace string, objectMeta *metav1.ObjectMeta, labels map[string]string, replicas *int32, containers []corev1.Container, vs []corev1.Volume, svcs []corev1.Service, pods []corev1.Pod, pvcs []corev1.PersistentVolumeClaim, scm map[string]storagev1.StorageClass) (*model.Application, error) {
	selector, err := metav1.LabelSelectorAsSelector(ls)
	if err != nil {
		return nil, fmt.Errorf("failed to create selector: %w", err)
	}
	return &model.Application{
		Type:                   appType,
		Name:                   name,
		Namespace:              namespace,
		ObjectMeta:             objectMeta,
		Labels:                 labels,
		Replicas:               replicas,
		Containers:             containers,
		Services:               filterServices(svcs, selector),
		Pods:                   filterPods(pods, selector),
		PersistentVolumeClaims: filterPersistentVolumeClaim(pvcs, vs, scm),
	}, nil
}

func filterServices(svcs []corev1.Service, s labels.Selector) []corev1.Service {
	ret := []corev1.Service{}
	for i := range svcs {
		if s.Matches(labels.Set(svcs[i].Spec.Selector)) {
			ret = append(ret, svcs[i])
		}
	}
	return ret
}

func filterPods(pods []corev1.Pod, s labels.Selector) []corev1.Pod {
	ret := []corev1.Pod{}
	for i := range pods {
		if s.Matches(labels.Set(pods[i].Labels)) {
			ret = append(ret, pods[i])
		}
	}
	return ret
}

func filterPersistentVolumeClaim(pvcs []corev1.PersistentVolumeClaim, vs []corev1.Volume, scm map[string]storagev1.StorageClass) []model.PersistentVolumeClaim {
	ret := []model.PersistentVolumeClaim{}
	for i := range vs {
		if vs[i].PersistentVolumeClaim == nil {
			continue
		}
		for j := range pvcs {
			if vs[i].PersistentVolumeClaim.ClaimName == pvcs[j].Name {
				if name := pvcs[j].Spec.StorageClassName; name != nil {
					if sc, ok := scm[*name]; ok {
						ret = append(ret, model.PersistentVolumeClaim{
							PersistentVolumeClaim: &pvcs[j],
							StorageClass:          &sc,
						})
						continue
					}
				}
				ret = append(ret, model.PersistentVolumeClaim{
					PersistentVolumeClaim: &pvcs[j],
				})
				break
			}
		}
	}
	return ret
}

func toStorageClassMap(scs []storagev1.StorageClass) map[string]storagev1.StorageClass {
	ret := map[string]storagev1.StorageClass{}
	for i := range scs {
		ret[scs[i].Name] = scs[i]
	}
	return ret
}

func randomName() string {
	return strings.ReplaceAll(namesgenerator.GetRandomName(0), "_", "-")
}
