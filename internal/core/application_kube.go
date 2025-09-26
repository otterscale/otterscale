package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/sync/errgroup"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type (
	Deployment  = appsv1.Deployment
	StatefulSet = appsv1.StatefulSet
	DaemonSet   = appsv1.DaemonSet
)

type (
	Job     = batchv1.Job
	JobSpec = batchv1.JobSpec
)

type (
	Node                  = corev1.Node
	Namespace             = corev1.Namespace
	ConfigMap             = corev1.ConfigMap
	Container             = corev1.Container
	Secret                = corev1.Secret
	Service               = corev1.Service
	ServicePort           = corev1.ServicePort
	Pod                   = corev1.Pod
	PersistentVolumeClaim = corev1.PersistentVolumeClaim
)

type StorageClass = storagev1.StorageClass

type Storage struct {
	*PersistentVolumeClaim
	*StorageClass
}

type Application struct {
	*ChartMetadata
	Type       string
	Name       string
	Namespace  string
	ObjectMeta *metav1.ObjectMeta
	Labels     map[string]string
	Replicas   *int32
	Containers []Container
	Services   []Service
	Pods       []Pod
	Storages   []Storage
}

type ControlPlaneCredential struct {
	ClientToken  string `json:"client_token"`
	KubeletToken string `json:"kubelet_token"`
	ProxyToken   string `json:"proxy_token"`
	Scope        string `json:"scope"`
}

type KubeAppsRepo interface {
	// Deployment
	ListDeployments(ctx context.Context, config *rest.Config, namespace string) ([]Deployment, error)
	GetDeployment(ctx context.Context, config *rest.Config, namespace, name string) (*Deployment, error)
	UpdateDeployment(ctx context.Context, config *rest.Config, namespace string, deployment *Deployment) (*Deployment, error)

	// StatefulSet
	ListStatefulSets(ctx context.Context, config *rest.Config, namespace string) ([]StatefulSet, error)
	GetStatefulSet(ctx context.Context, config *rest.Config, namespace, name string) (*StatefulSet, error)
	UpdateStatefulSet(ctx context.Context, config *rest.Config, namespace string, statefulSet *StatefulSet) (*StatefulSet, error)

	// DaemonSet
	ListDaemonSets(ctx context.Context, config *rest.Config, namespace string) ([]DaemonSet, error)
	GetDaemonSet(ctx context.Context, config *rest.Config, namespace, name string) (*DaemonSet, error)
	UpdateDaemonSet(ctx context.Context, config *rest.Config, namespace string, daemonSet *DaemonSet) (*DaemonSet, error)
}

type KubeBatchRepo interface {
	// Job
	ListJobs(ctx context.Context, config *rest.Config, namespace string) ([]Job, error)
	ListJobsByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]Job, error)
	CreateJob(ctx context.Context, config *rest.Config, namespace, name string, labels, annotations map[string]string, spec *JobSpec) (*Job, error)
	DeleteJob(ctx context.Context, config *rest.Config, namespace, name string) error
}

type KubeCoreRepo interface {
	// Node
	GetNode(ctx context.Context, config *rest.Config, name string) (*Node, error)
	UpdateNode(ctx context.Context, config *rest.Config, node *Node) (*Node, error)

	// Namespace
	ListNamespaces(ctx context.Context, config *rest.Config) ([]Namespace, error)

	// Service
	ListServices(ctx context.Context, config *rest.Config, namespace string) ([]Service, error)
	ListServicesByOptions(ctx context.Context, config *rest.Config, namespace, label, field string) ([]Service, error)
	ListVirtualMachineServices(ctx context.Context, config *rest.Config, namespace, vmName string) ([]Service, error)
	GetService(ctx context.Context, config *rest.Config, namespace, name string) (*Service, error)
	CreateVirtualMachineService(ctx context.Context, config *rest.Config, namespace, name, vmName string, ports []corev1.ServicePort) (*Service, error)
	UpdateService(ctx context.Context, config *rest.Config, namespace string, service *Service) (*Service, error)
	DeleteService(ctx context.Context, config *rest.Config, namespace, name string) error

	// Pod
	ListPods(ctx context.Context, config *rest.Config, namespace string) ([]Pod, error)
	ListPodsByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]Pod, error)
	GetLogs(ctx context.Context, config *rest.Config, namespace, podName, containerName string) (string, error)
	DeletePod(ctx context.Context, config *rest.Config, namespace, name string) error
	StreamLogs(ctx context.Context, config *rest.Config, namespace, podName, containerName string) (io.ReadCloser, error)
	CreateExecutor(config *rest.Config, namespace, podName, containerName string, command []string) (remotecommand.Executor, error)

	// PersistentVolumeClaim
	ListPersistentVolumeClaims(ctx context.Context, config *rest.Config, namespace string) ([]PersistentVolumeClaim, error)
	GetPersistentVolumeClaim(ctx context.Context, config *rest.Config, namespace, name string) (*PersistentVolumeClaim, error)
	PatchPersistentVolumeClaim(ctx context.Context, config *rest.Config, namespace, name string, data []byte) (*PersistentVolumeClaim, error)

	// Namespace
	GetNamespace(ctx context.Context, config *rest.Config, name string) (*Namespace, error)
	CreateNamespace(ctx context.Context, config *rest.Config, name string) (*Namespace, error)

	// ConfigMap
	GetConfigMap(ctx context.Context, config *rest.Config, namespace, name string) (*ConfigMap, error)
	CreateConfigMap(ctx context.Context, config *rest.Config, namespace, name string, data map[string]string) (*ConfigMap, error)

	// Secret
	GetSecret(ctx context.Context, config *rest.Config, namespace, name string) (*Secret, error)
}

type KubeStorageRepo interface {
	// StorageClass
	ListStorageClasses(ctx context.Context, config *rest.Config) ([]StorageClass, error)
	ListStorageClassesByLabel(ctx context.Context, config *rest.Config, label string) ([]StorageClass, error)
	GetStorageClass(ctx context.Context, config *rest.Config, name string) (*StorageClass, error)
}

func (uc *ApplicationUseCase) ListNamespaces(ctx context.Context, uuid, facility string) ([]Namespace, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeCore.ListNamespaces(ctx, config)
}

func (uc *ApplicationUseCase) ListApplications(ctx context.Context, uuid, facility string) ([]Application, error) {
	var (
		deployments            []appsv1.Deployment
		statefulSets           []appsv1.StatefulSet
		daemonSets             []appsv1.DaemonSet
		services               []corev1.Service
		pods                   []corev1.Pod
		persistentVolumeClaims []corev1.PersistentVolumeClaim
		storageClasses         []storagev1.StorageClass
	)

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		v, err := uc.kubeApps.ListDeployments(egctx, config, "")
		if err == nil {
			deployments = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeApps.ListStatefulSets(egctx, config, "")
		if err == nil {
			statefulSets = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeApps.ListDaemonSets(egctx, config, "")
		if err == nil {
			daemonSets = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeCore.ListServices(egctx, config, "")
		if err == nil {
			services = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeCore.ListPods(egctx, config, "")
		if err == nil {
			pods = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeCore.ListPersistentVolumeClaims(egctx, config, "")
		if err == nil {
			persistentVolumeClaims = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeStorage.ListStorageClasses(egctx, config)
		if err == nil {
			storageClasses = v
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	storageClassMap := toStorageClassMap(storageClasses)

	apps := []Application{}
	for i := range deployments {
		app, err := uc.fromDeployment(&deployments[i], services, pods, persistentVolumeClaims, storageClassMap)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}
	for i := range statefulSets {
		app, err := uc.fromStatefulSet(&statefulSets[i], services, pods, persistentVolumeClaims, storageClassMap)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}
	for i := range daemonSets {
		app, err := uc.fromDaemonSet(&daemonSets[i], services, pods, persistentVolumeClaims, storageClassMap)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}
	return apps, nil
}

func (uc *ApplicationUseCase) GetApplication(ctx context.Context, uuid, facility, namespace, name string) (*Application, error) {
	var (
		deployment             *appsv1.Deployment
		statefulSet            *appsv1.StatefulSet
		daemonSet              *appsv1.DaemonSet
		services               []corev1.Service
		pods                   []corev1.Pod
		persistentVolumeClaims []corev1.PersistentVolumeClaim
		storageClasses         []storagev1.StorageClass
	)

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		v, err := uc.kubeApps.GetDeployment(egctx, config, namespace, name)
		if err == nil {
			deployment = v
		} else if isKeyNotFoundError(err) {
			return nil
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeApps.GetStatefulSet(egctx, config, namespace, name)
		if err == nil {
			statefulSet = v
		} else if isKeyNotFoundError(err) {
			return nil
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeApps.GetDaemonSet(egctx, config, namespace, name)
		if err == nil {
			daemonSet = v
		} else if isKeyNotFoundError(err) {
			return nil
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeCore.ListServices(egctx, config, namespace)
		if err == nil {
			services = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeCore.ListPods(egctx, config, namespace)
		if err == nil {
			pods = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeCore.ListPersistentVolumeClaims(egctx, config, namespace)
		if err == nil {
			persistentVolumeClaims = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeStorage.ListStorageClasses(egctx, config)
		if err == nil {
			storageClasses = v
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	storageClassMap := toStorageClassMap(storageClasses)

	switch {
	case deployment != nil:
		return uc.fromDeployment(deployment, services, pods, persistentVolumeClaims, storageClassMap)
	case statefulSet != nil:
		return uc.fromStatefulSet(statefulSet, services, pods, persistentVolumeClaims, storageClassMap)
	case daemonSet != nil:
		return uc.fromDaemonSet(daemonSet, services, pods, persistentVolumeClaims, storageClassMap)
	}
	return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("application %q in namespace %q not found", name, namespace))
}

func (uc *ApplicationUseCase) RestartApplication(ctx context.Context, uuid, facility, namespace, name, appType string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	switch appType {
	case "deployment":
		deployment, err := uc.kubeApps.GetDeployment(ctx, config, namespace, name)
		if err != nil {
			return err
		}
		if deployment.Spec.Template.Annotations == nil {
			deployment.Spec.Template.Annotations = map[string]string{}
		}
		deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
		_, err = uc.kubeApps.UpdateDeployment(ctx, config, namespace, deployment)
		return err
	case "StatefulSet":
		statefulSet, err := uc.kubeApps.GetStatefulSet(ctx, config, namespace, name)
		if err != nil {
			return err
		}
		if statefulSet.Spec.Template.Annotations == nil {
			statefulSet.Spec.Template.Annotations = map[string]string{}
		}
		statefulSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
		_, err = uc.kubeApps.UpdateStatefulSet(ctx, config, namespace, statefulSet)
		return err
	case "DaemonSet":
		daemonSet, err := uc.kubeApps.GetDaemonSet(ctx, config, namespace, name)
		if err != nil {
			return err
		}
		if daemonSet.Spec.Template.Annotations == nil {
			daemonSet.Spec.Template.Annotations = map[string]string{}
		}
		daemonSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
		_, err = uc.kubeApps.UpdateDaemonSet(ctx, config, namespace, daemonSet)
		return err
	default:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unknown application type: %s", appType))
	}
}

func (uc *ApplicationUseCase) ScaleApplication(ctx context.Context, uuid, facility, namespace, name, appType string, replicas int32) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	switch appType {
	case "deployment":
		deployment, err := uc.kubeApps.GetDeployment(ctx, config, namespace, name)
		if err != nil {
			return err
		}
		deployment.Spec.Replicas = &replicas
		_, err = uc.kubeApps.UpdateDeployment(ctx, config, namespace, deployment)
		return err
	case "StatefulSet":
		statefulSet, err := uc.kubeApps.GetStatefulSet(ctx, config, namespace, name)
		if err != nil {
			return err
		}
		statefulSet.Spec.Replicas = &replicas
		_, err = uc.kubeApps.UpdateStatefulSet(ctx, config, namespace, statefulSet)
		return err
	case "DaemonSet":
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("daemon set does not support replica scaling"))
	default:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unknown application type: %s", appType))
	}
}

func (uc *ApplicationUseCase) DeletePod(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeCore.DeletePod(ctx, config, namespace, name)
}

func (uc *ApplicationUseCase) ListStorageClasses(ctx context.Context, uuid, facility string) ([]storagev1.StorageClass, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeStorage.ListStorageClasses(ctx, config)
}

func (uc *ApplicationUseCase) fromDeployment(d *appsv1.Deployment, svcs []corev1.Service, pods []corev1.Pod, pvcs []corev1.PersistentVolumeClaim, scm map[string]storagev1.StorageClass) (*Application, error) {
	return uc.toApplication(d.Spec.Selector, "Deployment", d.Name, d.Namespace, &d.ObjectMeta, d.Labels, d.Spec.Replicas, d.Spec.Template.Spec.Containers, d.Spec.Template.Spec.Volumes, svcs, pods, pvcs, scm)
}

func (uc *ApplicationUseCase) fromStatefulSet(d *appsv1.StatefulSet, svcs []corev1.Service, pods []corev1.Pod, pvcs []corev1.PersistentVolumeClaim, scm map[string]storagev1.StorageClass) (*Application, error) {
	return uc.toApplication(d.Spec.Selector, "StatefulSet", d.Name, d.Namespace, &d.ObjectMeta, d.Labels, d.Spec.Replicas, d.Spec.Template.Spec.Containers, d.Spec.Template.Spec.Volumes, svcs, pods, pvcs, scm)
}

func (uc *ApplicationUseCase) fromDaemonSet(d *appsv1.DaemonSet, svcs []corev1.Service, pods []corev1.Pod, pvcs []corev1.PersistentVolumeClaim, scm map[string]storagev1.StorageClass) (*Application, error) {
	return uc.toApplication(d.Spec.Selector, "DaemonSet", d.Name, d.Namespace, &d.ObjectMeta, d.Labels, nil, d.Spec.Template.Spec.Containers, d.Spec.Template.Spec.Volumes, svcs, pods, pvcs, scm)
}

func (uc *ApplicationUseCase) toApplication(ls *metav1.LabelSelector, appType, name, namespace string, objectMeta *metav1.ObjectMeta, labels map[string]string, replicas *int32, containers []corev1.Container, vs []corev1.Volume, svcs []corev1.Service, pods []corev1.Pod, pvcs []corev1.PersistentVolumeClaim, scm map[string]storagev1.StorageClass) (*Application, error) {
	selector, err := metav1.LabelSelectorAsSelector(ls)
	if err != nil {
		return nil, fmt.Errorf("failed to create selector: %w", err)
	}
	return &Application{
		Type:       appType,
		Name:       name,
		Namespace:  namespace,
		ObjectMeta: objectMeta,
		Labels:     labels,
		Replicas:   replicas,
		Containers: containers,
		Services:   filterServices(svcs, namespace, selector),
		Pods:       filterPods(pods, namespace, selector),
		Storages:   filterStorages(pvcs, vs, namespace, scm),
	}, nil
}

func (uc *ApplicationUseCase) StreamLogs(ctx context.Context, uuid, facility, namespace, podName, containerName string) (io.ReadCloser, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeCore.StreamLogs(ctx, config, namespace, podName, containerName)
}

func (uc *ApplicationUseCase) ExecuteTTY(ctx context.Context, uuid, facility, namespace, podName, containerName string, command []string) (remotecommand.Executor, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeCore.CreateExecutor(config, namespace, podName, containerName, command)
}

func filterServices(svcs []corev1.Service, namespace string, s labels.Selector) []corev1.Service {
	ret := []corev1.Service{}
	for i := range svcs {
		if svcs[i].Namespace == namespace && s.Matches(labels.Set(svcs[i].Spec.Selector)) {
			ret = append(ret, svcs[i])
		}
	}
	return ret
}

func filterPods(pods []corev1.Pod, namespace string, s labels.Selector) []corev1.Pod {
	ret := []corev1.Pod{}
	for i := range pods {
		if pods[i].Namespace == namespace && s.Matches(labels.Set(pods[i].Labels)) {
			ret = append(ret, pods[i])
		}
	}
	return ret
}

func filterStorages(pvcs []corev1.PersistentVolumeClaim, vs []corev1.Volume, namespace string, scm map[string]storagev1.StorageClass) []Storage {
	ret := []Storage{}
	for i := range vs {
		if vs[i].PersistentVolumeClaim == nil {
			continue
		}
		for j := range pvcs {
			if vs[i].PersistentVolumeClaim.ClaimName != pvcs[j].Name {
				continue
			}
			if pvcs[j].Namespace != namespace {
				continue
			}
			if name := pvcs[j].Spec.StorageClassName; name != nil {
				if sc, ok := scm[*name]; ok {
					ret = append(ret, Storage{
						PersistentVolumeClaim: &pvcs[j],
						StorageClass:          &sc,
					})
					continue
				}
			}
			ret = append(ret, Storage{
				PersistentVolumeClaim: &pvcs[j],
			})
			break
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

func isKeyNotFoundError(err error) bool {
	statusErr, _ := err.(*k8serrors.StatusError)
	return statusErr != nil && statusErr.Status().Code == http.StatusNotFound
}
