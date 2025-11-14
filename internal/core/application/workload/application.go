package workload

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/otterscale/otterscale/internal/core/application/chart"
	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/application/storage"
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	ApplicationTypeDeployment  = "Deployment"
	ApplicationTypeStatefulSet = "StatefulSet"
	ApplicationTypeDaemonSet   = "DaemonSet"
)

type Application struct {
	Type       string
	Name       string
	Namespace  string
	Labels     map[string]string
	Replicas   *int32
	ObjectMeta ObjectMeta
	Pods       []Pod
	Containers []Container
	Services   []service.Service
	Storages   []storage.Storage
	ChartFile  *chart.File // return only when fetching from GetApplication
}

func (uc *WorkloadUseCase) ListApplications(ctx context.Context, scope string) (apps []Application, endpoint string, err error) {
	var (
		deployments            []Deployment
		statefulSets           []StatefulSet
		daemonSets             []DaemonSet
		pods                   []Pod
		services               []service.Service
		persistentVolumeClaims []storage.PersistentVolumeClaim
		storageClasses         []storage.StorageClass
	)

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.deployment.List(egctx, scope, "", "")
		if err == nil {
			deployments = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.statefulSet.List(egctx, scope, "", "")
		if err == nil {
			statefulSets = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.daemonSet.List(egctx, scope, "", "")
		if err == nil {
			daemonSets = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.pod.List(egctx, scope, "", "")
		if err == nil {
			pods = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.service.List(egctx, scope, "", "")
		if err == nil {
			services = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.persistentVolumeClaim.List(egctx, scope, "", "")
		if err == nil {
			persistentVolumeClaims = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.storageClass.List(egctx, scope, "")
		if err == nil {
			storageClasses = v
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, "", err
	}

	apps = []Application{}

	for i := range deployments {
		app, err := uc.fromDeployment(&deployments[i], pods, services, persistentVolumeClaims, storageClasses)
		if err != nil {
			return nil, "", err
		}

		apps = append(apps, *app)
	}

	for i := range statefulSets {
		app, err := uc.fromStatefulSet(&statefulSets[i], pods, services, persistentVolumeClaims, storageClasses)
		if err != nil {
			return nil, "", err
		}

		apps = append(apps, *app)
	}

	for i := range daemonSets {
		app, err := uc.fromDaemonSet(&daemonSets[i], pods, services, persistentVolumeClaims, storageClasses)
		if err != nil {
			return nil, "", err
		}

		apps = append(apps, *app)
	}

	return apps, uc.service.Host(scope), nil
}

func (uc *WorkloadUseCase) GetApplication(ctx context.Context, scope, namespace, name string) (*Application, error) {
	var (
		deployment             *Deployment
		statefulSet            *StatefulSet
		daemonSet              *DaemonSet
		pods                   []Pod
		services               []service.Service
		persistentVolumeClaims []storage.PersistentVolumeClaim
		storageClasses         []storage.StorageClass
	)

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.deployment.Get(egctx, scope, namespace, name)
		if err == nil {
			deployment = v
		} else if uc.isKeyNotFoundError(err) {
			return nil
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.statefulSet.Get(egctx, scope, namespace, name)
		if err == nil {
			statefulSet = v
		} else if uc.isKeyNotFoundError(err) {
			return nil
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.daemonSet.Get(egctx, scope, namespace, name)
		if err == nil {
			daemonSet = v
		} else if uc.isKeyNotFoundError(err) {
			return nil
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.pod.List(egctx, scope, namespace, "")
		if err == nil {
			pods = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.service.List(egctx, scope, namespace, "")
		if err == nil {
			services = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.persistentVolumeClaim.List(egctx, scope, namespace, "")
		if err == nil {
			persistentVolumeClaims = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.storageClass.List(egctx, scope, "")
		if err == nil {
			storageClasses = v
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	switch {
	case deployment != nil:
		return uc.fromDeployment(deployment, pods, services, persistentVolumeClaims, storageClasses)

	case statefulSet != nil:
		return uc.fromStatefulSet(statefulSet, pods, services, persistentVolumeClaims, storageClasses)

	case daemonSet != nil:
		return uc.fromDaemonSet(daemonSet, pods, services, persistentVolumeClaims, storageClasses)
	}

	return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("application %q in namespace %q not found", name, namespace))
}

func (uc *WorkloadUseCase) filterServices(selector labels.Selector, namespace string, services []service.Service) []service.Service {
	ret := []service.Service{}

	for i := range services {
		if services[i].Namespace == namespace && selector.Matches(labels.Set(services[i].Spec.Selector)) {
			ret = append(ret, services[i])
		}
	}

	return ret
}

func (uc *WorkloadUseCase) filterPods(selector labels.Selector, namespace string, pods []Pod) []Pod {
	ret := []Pod{}

	for i := range pods {
		if pods[i].Namespace == namespace && selector.Matches(labels.Set(pods[i].Labels)) {
			ret = append(ret, pods[i])
		}
	}

	return ret
}

func (uc *WorkloadUseCase) filterStorages(namespace string, volumes []storage.Volume, persistentVolumeClaims []storage.PersistentVolumeClaim, storageClasses []storage.StorageClass) []storage.Storage {
	storageClassMap := map[string]*storage.StorageClass{}
	for i := range storageClasses {
		sc := storageClasses[i]
		storageClassMap[sc.Name] = &sc
	}

	pvcMap := map[string]*storage.PersistentVolumeClaim{}
	for i := range persistentVolumeClaims {
		pvc := persistentVolumeClaims[i]
		if pvc.Namespace == namespace {
			pvcMap[pvc.Name] = &pvc
		}
	}

	ret := []storage.Storage{}

	for i := range volumes {
		vol := volumes[i]

		if vol.PersistentVolumeClaim == nil {
			continue
		}

		pvc, found := pvcMap[vol.PersistentVolumeClaim.ClaimName]
		if !found {
			continue
		}

		storage := storage.Storage{
			PersistentVolumeClaim: pvc,
		}

		scName := pvc.Spec.StorageClassName
		if scName != nil && *scName != "" {
			sc, found := storageClassMap[*scName]
			if found {
				storage.StorageClass = sc
			}
		}

		ret = append(ret, storage)
	}

	return ret
}

func (uc *WorkloadUseCase) toApplication(ls *v1.LabelSelector, appType, name, namespace string, labels map[string]string, replicas *int32, objectMeta ObjectMeta, pods []Pod, containers []Container, services []service.Service, volumes []storage.Volume, persistentVolumeClaims []storage.PersistentVolumeClaim, storageClasses []storage.StorageClass) (*Application, error) {
	selector, err := v1.LabelSelectorAsSelector(ls)
	if err != nil {
		return nil, fmt.Errorf("failed to create selector: %w", err)
	}

	return &Application{
		Type:       appType,
		Name:       name,
		Namespace:  namespace,
		Labels:     labels,
		Replicas:   replicas,
		ObjectMeta: objectMeta,
		Containers: containers,
		Services:   uc.filterServices(selector, namespace, services),
		Pods:       uc.filterPods(selector, namespace, pods),
		Storages:   uc.filterStorages(namespace, volumes, persistentVolumeClaims, storageClasses),
	}, nil
}

func (uc *WorkloadUseCase) fromDeployment(workload *Deployment, pods []Pod, services []service.Service, persistentVolumeClaims []storage.PersistentVolumeClaim, storageClasses []storage.StorageClass) (*Application, error) {
	return uc.toApplication(
		workload.Spec.Selector,
		ApplicationTypeDeployment,
		workload.Name,
		workload.Namespace,
		workload.Labels,
		workload.Spec.Replicas,
		workload.ObjectMeta,
		pods,
		workload.Spec.Template.Spec.Containers,
		services,
		workload.Spec.Template.Spec.Volumes,
		persistentVolumeClaims,
		storageClasses,
	)
}

func (uc *WorkloadUseCase) fromStatefulSet(workload *StatefulSet, pods []Pod, services []service.Service, persistentVolumeClaims []storage.PersistentVolumeClaim, storageClasses []storage.StorageClass) (*Application, error) {
	return uc.toApplication(
		workload.Spec.Selector,
		ApplicationTypeStatefulSet,
		workload.Name,
		workload.Namespace,
		workload.Labels,
		workload.Spec.Replicas,
		workload.ObjectMeta,
		pods,
		workload.Spec.Template.Spec.Containers,
		services,
		workload.Spec.Template.Spec.Volumes,
		persistentVolumeClaims,
		storageClasses,
	)
}

func (uc *WorkloadUseCase) fromDaemonSet(workload *DaemonSet, pods []Pod, services []service.Service, persistentVolumeClaims []storage.PersistentVolumeClaim, storageClasses []storage.StorageClass) (*Application, error) {
	return uc.toApplication(
		workload.Spec.Selector,
		ApplicationTypeDaemonSet,
		workload.Name,
		workload.Namespace,
		workload.Labels,
		nil,
		workload.ObjectMeta,
		pods,
		workload.Spec.Template.Spec.Containers,
		services,
		workload.Spec.Template.Spec.Volumes,
		persistentVolumeClaims,
		storageClasses,
	)
}
