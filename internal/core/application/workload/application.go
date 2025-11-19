package workload

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/sync/errgroup"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/otterscale/otterscale/internal/core/application/chart"
	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/service"
)

const (
	ApplicationTypeDeployment  = "Deployment"
	ApplicationTypeStatefulSet = "StatefulSet"
	ApplicationTypeDaemonSet   = "DaemonSet"
)

type Application struct {
	Type        string
	Name        string
	Namespace   string
	Labels      map[string]string
	Replicas    *int32
	ObjectMeta  *ObjectMeta
	Pods        []Pod
	Containers  []Container
	Services    []service.Service
	Persistents []persistent.Persistent
	ChartFile   *chart.File // return only when fetching from GetApplication
}

func (uc *UseCase) ListApplications(ctx context.Context, scope string) (apps []Application, hostname string, err error) {
	var (
		deployments            []Deployment
		statefulSets           []StatefulSet
		daemonSets             []DaemonSet
		pods                   []Pod
		services               []service.Service
		persistentVolumeClaims []persistent.PersistentVolumeClaim
		storageClasses         []persistent.StorageClass
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

	apps, err = uc.combineApplications(deployments, statefulSets, daemonSets, pods, services, persistentVolumeClaims, storageClasses)
	if err != nil {
		return nil, "", err
	}

	url, err := uc.service.URL(scope)
	if err != nil {
		return nil, "", err
	}

	return apps, url.Hostname(), nil
}

func (uc *UseCase) RestartApplication(ctx context.Context, scope, namespace, name, appType string) error {
	switch appType {
	case ApplicationTypeDeployment:
		deployment, err := uc.deployment.Get(ctx, scope, namespace, name)
		if err != nil {
			return err
		}

		if deployment.Spec.Template.Annotations == nil {
			deployment.Spec.Template.Annotations = map[string]string{}
		}

		deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

		_, err = uc.deployment.Update(ctx, scope, namespace, deployment)
		return err

	case ApplicationTypeStatefulSet:
		statefulSet, err := uc.statefulSet.Get(ctx, scope, namespace, name)
		if err != nil {
			return err
		}

		if statefulSet.Spec.Template.Annotations == nil {
			statefulSet.Spec.Template.Annotations = map[string]string{}
		}

		statefulSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

		_, err = uc.statefulSet.Update(ctx, scope, namespace, statefulSet)
		return err

	case ApplicationTypeDaemonSet:
		daemonSet, err := uc.daemonSet.Get(ctx, scope, namespace, name)
		if err != nil {
			return err
		}

		if daemonSet.Spec.Template.Annotations == nil {
			daemonSet.Spec.Template.Annotations = map[string]string{}
		}

		daemonSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

		_, err = uc.daemonSet.Update(ctx, scope, namespace, daemonSet)
		return err
	}

	return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unknown application type: %s", appType))
}

func (uc *UseCase) ScaleApplication(ctx context.Context, scope, namespace, name, appType string, replicas int32) error {
	switch appType {
	case ApplicationTypeDeployment:
		deployment, err := uc.deployment.Get(ctx, scope, namespace, name)
		if err != nil {
			return err
		}

		deployment.Spec.Replicas = &replicas

		_, err = uc.deployment.Update(ctx, scope, namespace, deployment)
		return err

	case ApplicationTypeStatefulSet:
		statefulSet, err := uc.statefulSet.Get(ctx, scope, namespace, name)
		if err != nil {
			return err
		}

		statefulSet.Spec.Replicas = &replicas

		_, err = uc.statefulSet.Update(ctx, scope, namespace, statefulSet)
		return err

	case ApplicationTypeDaemonSet:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("daemon set does not support replica scaling"))
	}

	return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unknown application type: %s", appType))
}

func (uc *UseCase) GetApplication(ctx context.Context, scope, namespace, name string) (*Application, error) {
	var (
		deployment             *Deployment
		statefulSet            *StatefulSet
		daemonSet              *DaemonSet
		pods                   []Pod
		services               []service.Service
		persistentVolumeClaims []persistent.PersistentVolumeClaim
		storageClasses         []persistent.StorageClass
	)

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.deployment.Get(egctx, scope, namespace, name)
		if err == nil {
			deployment = v
		} else if k8serrors.IsNotFound(err) {
			return nil
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.statefulSet.Get(egctx, scope, namespace, name)
		if err == nil {
			statefulSet = v
		} else if k8serrors.IsNotFound(err) {
			return nil
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.daemonSet.Get(egctx, scope, namespace, name)
		if err == nil {
			daemonSet = v
		} else if k8serrors.IsNotFound(err) {
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

func (uc *UseCase) filterServices(podLabels map[string]string, namespace string, services []service.Service) []service.Service {
	ret := []service.Service{}

	for i := range services {
		selector := labels.SelectorFromSet(services[i].Spec.Selector)

		if services[i].Namespace == namespace && selector.Matches(labels.Set(podLabels)) {
			ret = append(ret, services[i])
		}
	}

	return ret
}

func (uc *UseCase) filterPods(selector labels.Selector, namespace string, pods []Pod) []Pod {
	ret := []Pod{}

	for i := range pods {
		if pods[i].Namespace == namespace && selector.Matches(labels.Set(pods[i].Labels)) {
			ret = append(ret, pods[i])
		}
	}

	return ret
}

func (uc *UseCase) filterPersistents(namespace string, volumes []persistent.Volume, persistentVolumeClaims []persistent.PersistentVolumeClaim, storageClasses []persistent.StorageClass) []persistent.Persistent {
	storageClassMap := map[string]*persistent.StorageClass{}
	for i := range storageClasses {
		sc := storageClasses[i]
		storageClassMap[sc.Name] = &sc
	}

	pvcMap := map[string]*persistent.PersistentVolumeClaim{}
	for i := range persistentVolumeClaims {
		pvc := persistentVolumeClaims[i]
		if pvc.Namespace == namespace {
			pvcMap[pvc.Name] = &pvc
		}
	}

	ret := []persistent.Persistent{}

	for i := range volumes {
		vol := volumes[i]

		if vol.PersistentVolumeClaim == nil {
			continue
		}

		pvc, found := pvcMap[vol.PersistentVolumeClaim.ClaimName]
		if !found {
			continue
		}

		persistent := persistent.Persistent{
			PersistentVolumeClaim: pvc,
		}

		scName := pvc.Spec.StorageClassName
		if scName != nil && *scName != "" {
			sc, found := storageClassMap[*scName]
			if found {
				persistent.StorageClass = sc
			}
		}

		ret = append(ret, persistent)
	}

	return ret
}

func (uc *UseCase) toApplication(labelSelector *v1.LabelSelector, podLabels map[string]string, appType, name, namespace string, labels map[string]string, replicas *int32, objectMeta *ObjectMeta, pods []Pod, containers []Container, services []service.Service, volumes []persistent.Volume, persistentVolumeClaims []persistent.PersistentVolumeClaim, storageClasses []persistent.StorageClass) (*Application, error) {
	selector, err := v1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to create selector: %w", err)
	}

	return &Application{
		Type:        appType,
		Name:        name,
		Namespace:   namespace,
		Labels:      labels,
		Replicas:    replicas,
		ObjectMeta:  objectMeta,
		Containers:  containers,
		Services:    uc.filterServices(podLabels, namespace, services),
		Pods:        uc.filterPods(selector, namespace, pods),
		Persistents: uc.filterPersistents(namespace, volumes, persistentVolumeClaims, storageClasses),
	}, nil
}

func (uc *UseCase) fromDeployment(workload *Deployment, pods []Pod, services []service.Service, persistentVolumeClaims []persistent.PersistentVolumeClaim, storageClasses []persistent.StorageClass) (*Application, error) {
	return uc.toApplication(
		workload.Spec.Selector,
		workload.Spec.Template.Labels,
		ApplicationTypeDeployment,
		workload.Name,
		workload.Namespace,
		workload.Labels,
		workload.Spec.Replicas,
		&workload.ObjectMeta,
		pods,
		workload.Spec.Template.Spec.Containers,
		services,
		workload.Spec.Template.Spec.Volumes,
		persistentVolumeClaims,
		storageClasses,
	)
}

func (uc *UseCase) fromStatefulSet(workload *StatefulSet, pods []Pod, services []service.Service, persistentVolumeClaims []persistent.PersistentVolumeClaim, storageClasses []persistent.StorageClass) (*Application, error) {
	return uc.toApplication(
		workload.Spec.Selector,
		workload.Spec.Template.Labels,
		ApplicationTypeStatefulSet,
		workload.Name,
		workload.Namespace,
		workload.Labels,
		workload.Spec.Replicas,
		&workload.ObjectMeta,
		pods,
		workload.Spec.Template.Spec.Containers,
		services,
		workload.Spec.Template.Spec.Volumes,
		persistentVolumeClaims,
		storageClasses,
	)
}

func (uc *UseCase) fromDaemonSet(workload *DaemonSet, pods []Pod, services []service.Service, persistentVolumeClaims []persistent.PersistentVolumeClaim, storageClasses []persistent.StorageClass) (*Application, error) {
	return uc.toApplication(
		workload.Spec.Selector,
		workload.Spec.Template.Labels,
		ApplicationTypeDaemonSet,
		workload.Name,
		workload.Namespace,
		workload.Labels,
		nil,
		&workload.ObjectMeta,
		pods,
		workload.Spec.Template.Spec.Containers,
		services,
		workload.Spec.Template.Spec.Volumes,
		persistentVolumeClaims,
		storageClasses,
	)
}

func (uc *UseCase) combineApplications(deployments []Deployment, statefulSets []StatefulSet, daemonSets []DaemonSet, pods []Pod, services []service.Service, persistentVolumeClaims []persistent.PersistentVolumeClaim, storageClasses []persistent.StorageClass) ([]Application, error) {
	ret := []Application{}

	for i := range deployments {
		app, err := uc.fromDeployment(&deployments[i], pods, services, persistentVolumeClaims, storageClasses)
		if err != nil {
			return nil, err
		}

		ret = append(ret, *app)
	}

	for i := range statefulSets {
		app, err := uc.fromStatefulSet(&statefulSets[i], pods, services, persistentVolumeClaims, storageClasses)
		if err != nil {
			return nil, err
		}

		ret = append(ret, *app)
	}

	for i := range daemonSets {
		app, err := uc.fromDaemonSet(&daemonSets[i], pods, services, persistentVolumeClaims, storageClasses)
		if err != nil {
			return nil, err
		}

		ret = append(ret, *app)
	}

	return ret, nil
}
