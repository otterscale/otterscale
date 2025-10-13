package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
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

const KubernetesResourceStorage = corev1.ResourceStorage

type (
	KubernetesTime             = metav1.Time
	DaemonSet                  = appsv1.DaemonSet
	Deployment                 = appsv1.Deployment
	StatefulSet                = appsv1.StatefulSet
	Job                        = batchv1.Job
	JobSpec                    = batchv1.JobSpec
	ConfigMap                  = corev1.ConfigMap
	Container                  = corev1.Container
	ContainerStatus            = corev1.ContainerStatus
	Namespace                  = corev1.Namespace
	Node                       = corev1.Node
	PersistentVolumeClaim      = corev1.PersistentVolumeClaim
	PersistentVolumeAccessMode = corev1.PersistentVolumeAccessMode
	Pod                        = corev1.Pod
	PodCondition               = corev1.PodCondition
	PodPhase                   = corev1.PodPhase
	ResourceList               = corev1.ResourceList
	Secret                     = corev1.Secret
	Service                    = corev1.Service
	ServicePort                = corev1.ServicePort
	ServiceProtocol            = corev1.Protocol
	StorageClass               = storagev1.StorageClass
)

type Application struct {
	*ChartFile
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

type Storage struct {
	*PersistentVolumeClaim
	*StorageClass
}

type TTYSession struct {
	id        string
	inReader  *io.PipeReader
	inWriter  *io.PipeWriter
	outReader *io.PipeReader
	outWriter *io.PipeWriter
}

type KubeAppsRepo interface {
	ListDaemonSets(ctx context.Context, config *rest.Config, namespace string) ([]DaemonSet, error)
	GetDaemonSet(ctx context.Context, config *rest.Config, namespace, name string) (*DaemonSet, error)
	UpdateDaemonSet(ctx context.Context, config *rest.Config, namespace string, daemonSet *DaemonSet) (*DaemonSet, error)

	ListDeployments(ctx context.Context, config *rest.Config, namespace string) ([]Deployment, error)
	ListDeploymentsByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]Deployment, error)
	GetDeployment(ctx context.Context, config *rest.Config, namespace, name string) (*Deployment, error)
	UpdateDeployment(ctx context.Context, config *rest.Config, namespace string, deployment *Deployment) (*Deployment, error)

	ListStatefulSets(ctx context.Context, config *rest.Config, namespace string) ([]StatefulSet, error)
	GetStatefulSet(ctx context.Context, config *rest.Config, namespace, name string) (*StatefulSet, error)
	UpdateStatefulSet(ctx context.Context, config *rest.Config, namespace string, statefulSet *StatefulSet) (*StatefulSet, error)
}

type KubeBatchRepo interface {
	ListJobs(ctx context.Context, config *rest.Config, namespace string) ([]Job, error)
	ListJobsByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]Job, error)
	CreateJob(ctx context.Context, config *rest.Config, namespace, name string, labels, annotations map[string]string, spec *JobSpec) (*Job, error)
	DeleteJob(ctx context.Context, config *rest.Config, namespace, name string) error
}

type KubeCoreRepo interface {
	GetConfigMap(ctx context.Context, config *rest.Config, namespace, name string) (*ConfigMap, error)
	CreateConfigMap(ctx context.Context, config *rest.Config, namespace, name string, data map[string]string) (*ConfigMap, error)

	ListNamespaces(ctx context.Context, config *rest.Config) ([]Namespace, error)
	GetNamespace(ctx context.Context, config *rest.Config, name string) (*Namespace, error)
	CreateNamespace(ctx context.Context, config *rest.Config, name string) (*Namespace, error)

	ListNodes(ctx context.Context, config *rest.Config) ([]Node, error)
	GetNode(ctx context.Context, config *rest.Config, name string) (*Node, error)
	UpdateNode(ctx context.Context, config *rest.Config, node *Node) (*Node, error)

	ListPersistentVolumeClaims(ctx context.Context, config *rest.Config, namespace string) ([]PersistentVolumeClaim, error)
	GetPersistentVolumeClaim(ctx context.Context, config *rest.Config, namespace, name string) (*PersistentVolumeClaim, error)
	PatchPersistentVolumeClaim(ctx context.Context, config *rest.Config, namespace, name string, data []byte) (*PersistentVolumeClaim, error)

	ListPods(ctx context.Context, config *rest.Config, namespace string) ([]Pod, error)
	ListPodsByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]Pod, error)
	GetLogs(ctx context.Context, config *rest.Config, namespace, podName, containerName string) (string, error)
	DeletePod(ctx context.Context, config *rest.Config, namespace, name string) error
	StreamLogs(ctx context.Context, config *rest.Config, namespace, podName, containerName string) (io.ReadCloser, error)
	CreateExecutor(config *rest.Config, namespace, podName, containerName string, command []string) (remotecommand.Executor, error)

	GetSecret(ctx context.Context, config *rest.Config, namespace, name string) (*Secret, error)

	ListServices(ctx context.Context, config *rest.Config, namespace string) ([]Service, error)
	ListServicesByOptions(ctx context.Context, config *rest.Config, namespace, label, field string) ([]Service, error)
	GetService(ctx context.Context, config *rest.Config, namespace, name string) (*Service, error)
	UpdateService(ctx context.Context, config *rest.Config, namespace string, service *Service) (*Service, error)
	DeleteService(ctx context.Context, config *rest.Config, namespace, name string) error
	ListVirtualMachineServices(ctx context.Context, config *rest.Config, namespace, vmName string) ([]Service, error)
	CreateVirtualMachineService(ctx context.Context, config *rest.Config, namespace, name, vmName string, ports []corev1.ServicePort) (*Service, error)
}

type KubeStorageRepo interface {
	ListStorageClasses(ctx context.Context, config *rest.Config) ([]StorageClass, error)
	ListStorageClassesByLabel(ctx context.Context, config *rest.Config, label string) ([]StorageClass, error)
	GetStorageClass(ctx context.Context, config *rest.Config, name string) (*StorageClass, error)
}

type KubernetesUseCase struct {
	action      ActionRepo
	client      ClientRepo
	facility    FacilityRepo
	kubeApps    KubeAppsRepo
	kubeCore    KubeCoreRepo
	kubeStorage KubeStorageRepo

	ttySessionMap sync.Map
}

func NewKubernetesUseCase(action ActionRepo, client ClientRepo, facility FacilityRepo, kubeApps KubeAppsRepo, kubeCore KubeCoreRepo, kubeStorage KubeStorageRepo) *KubernetesUseCase {
	return &KubernetesUseCase{
		action:      action,
		client:      client,
		facility:    facility,
		kubeApps:    kubeApps,
		kubeCore:    kubeCore,
		kubeStorage: kubeStorage,
	}
}

func (uc *KubernetesUseCase) GetPublicAddress(ctx context.Context, uuid, facility string) (string, error) {
	leader, err := uc.facility.GetLeader(ctx, uuid, facility)
	if err != nil {
		return "", err
	}
	unitInfo, err := uc.facility.GetUnitInfo(ctx, uuid, leader)
	if err != nil {
		return "", err
	}
	return unitInfo.PublicAddress, nil
}

func (uc *KubernetesUseCase) ListNamespaces(ctx context.Context, uuid, facility string) ([]Namespace, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeCore.ListNamespaces(ctx, config)
}

func (uc *KubernetesUseCase) ListApplications(ctx context.Context, uuid, facility string) ([]Application, error) {
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
		app, err := fromDeployment(&deployments[i], services, pods, persistentVolumeClaims, storageClassMap)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}
	for i := range statefulSets {
		app, err := fromStatefulSet(&statefulSets[i], services, pods, persistentVolumeClaims, storageClassMap)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}
	for i := range daemonSets {
		app, err := fromDaemonSet(&daemonSets[i], services, pods, persistentVolumeClaims, storageClassMap)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}
	return apps, nil
}

func (uc *KubernetesUseCase) GetApplication(ctx context.Context, uuid, facility, namespace, name string) (*Application, error) {
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
		return fromDeployment(deployment, services, pods, persistentVolumeClaims, storageClassMap)
	case statefulSet != nil:
		return fromStatefulSet(statefulSet, services, pods, persistentVolumeClaims, storageClassMap)
	case daemonSet != nil:
		return fromDaemonSet(daemonSet, services, pods, persistentVolumeClaims, storageClassMap)
	}
	return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("application %q in namespace %q not found", name, namespace))
}

func (uc *KubernetesUseCase) RestartApplication(ctx context.Context, uuid, facility, namespace, name, appType string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	switch appType {
	case ApplicationTypeDeployment:
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
	case ApplicationTypeStatefulSet:
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
	case ApplicationTypeDaemonSet:
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

func (uc *KubernetesUseCase) ScaleApplication(ctx context.Context, uuid, facility, namespace, name, appType string, replicas int32) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	switch appType {
	case ApplicationTypeDeployment:
		deployment, err := uc.kubeApps.GetDeployment(ctx, config, namespace, name)
		if err != nil {
			return err
		}
		deployment.Spec.Replicas = &replicas
		_, err = uc.kubeApps.UpdateDeployment(ctx, config, namespace, deployment)
		return err
	case ApplicationTypeStatefulSet:
		statefulSet, err := uc.kubeApps.GetStatefulSet(ctx, config, namespace, name)
		if err != nil {
			return err
		}
		statefulSet.Spec.Replicas = &replicas
		_, err = uc.kubeApps.UpdateStatefulSet(ctx, config, namespace, statefulSet)
		return err
	case ApplicationTypeDaemonSet:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("daemon set does not support replica scaling"))
	default:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unknown application type: %s", appType))
	}
}

func (uc *KubernetesUseCase) DeletePod(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeCore.DeletePod(ctx, config, namespace, name)
}

func (uc *KubernetesUseCase) StreamLogs(ctx context.Context, uuid, facility, namespace, podName, containerName string) (io.ReadCloser, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeCore.StreamLogs(ctx, config, namespace, podName, containerName)
}

func (uc *KubernetesUseCase) WriteToTTYSession(sessionID string, stdIn []byte) error {
	value, ok := uc.ttySessionMap.Load(sessionID)
	if !ok {
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("session %s not found", sessionID))
	}
	if _, err := value.(*TTYSession).inWriter.Write(stdIn); err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to write to session: %w", err))
	}
	return nil
}

func (uc *KubernetesUseCase) CreateTTYSession() (string, error) {
	sessionID := uuid.New().String()

	inReader, inWriter := io.Pipe()
	outReader, outWriter := io.Pipe()

	uc.ttySessionMap.Store(sessionID, &TTYSession{
		id:        sessionID,
		inReader:  inReader,
		inWriter:  inWriter,
		outReader: outReader,
		outWriter: outWriter,
	})
	return sessionID, nil
}

func (uc *KubernetesUseCase) CleanupTTYSession(sessionID string) error {
	value, ok := uc.ttySessionMap.Load(sessionID)
	if !ok {
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("session %s not found", sessionID))
	}
	ttySession := value.(*TTYSession)
	ttySession.inReader.Close()
	ttySession.inWriter.Close()
	ttySession.outReader.Close()
	ttySession.outWriter.Close()
	uc.ttySessionMap.Delete(sessionID)
	return nil
}

func (uc *KubernetesUseCase) ExecuteTTY(ctx context.Context, sessionID, uuid, facility, namespace, podName, containerName string, command []string, stdOut chan<- []byte) error {
	value, ok := uc.ttySessionMap.Load(sessionID)
	if !ok {
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("session %s not found", sessionID))
	}
	ttySession := value.(*TTYSession)

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	exec, err := uc.kubeCore.CreateExecutor(config, namespace, podName, containerName, command)
	if err != nil {
		return err
	}
	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		buf := make([]byte, 1024) //nolint:mnd // 1KB buffer
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				n, err := ttySession.outReader.Read(buf)
				if err != nil {
					if err == io.EOF {
						return nil
					}
					return err
				}
				if n > 0 {
					// write message to std out
					stdOut <- buf[:n]
				}
			}
		}
	})
	eg.Go(func() error {
		defer close(stdOut)
		return exec.StreamWithContext(egctx, remotecommand.StreamOptions{
			Stdin:  ttySession.inReader,
			Stdout: ttySession.outWriter, // raw TTY manages stdout and stderr over the stdout stream
			Tty:    true,
		})
	})
	return eg.Wait()
}

func (uc *KubernetesUseCase) ListStorageClasses(ctx context.Context, uuid, facility string) ([]storagev1.StorageClass, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeStorage.ListStorageClasses(ctx, config)
}

func fromDeployment(d *appsv1.Deployment, svcs []corev1.Service, pods []corev1.Pod, pvcs []corev1.PersistentVolumeClaim, scm map[string]storagev1.StorageClass) (*Application, error) {
	return toApplication(d.Spec.Selector, ApplicationTypeDeployment, d.Name, d.Namespace, &d.ObjectMeta, d.Labels, d.Spec.Replicas, d.Spec.Template.Spec.Containers, d.Spec.Template.Spec.Volumes, svcs, pods, pvcs, scm)
}

func fromStatefulSet(d *appsv1.StatefulSet, svcs []corev1.Service, pods []corev1.Pod, pvcs []corev1.PersistentVolumeClaim, scm map[string]storagev1.StorageClass) (*Application, error) {
	return toApplication(d.Spec.Selector, ApplicationTypeStatefulSet, d.Name, d.Namespace, &d.ObjectMeta, d.Labels, d.Spec.Replicas, d.Spec.Template.Spec.Containers, d.Spec.Template.Spec.Volumes, svcs, pods, pvcs, scm)
}

func fromDaemonSet(d *appsv1.DaemonSet, svcs []corev1.Service, pods []corev1.Pod, pvcs []corev1.PersistentVolumeClaim, scm map[string]storagev1.StorageClass) (*Application, error) {
	return toApplication(d.Spec.Selector, ApplicationTypeDaemonSet, d.Name, d.Namespace, &d.ObjectMeta, d.Labels, nil, d.Spec.Template.Spec.Containers, d.Spec.Template.Spec.Volumes, svcs, pods, pvcs, scm)
}

func toApplication(ls *metav1.LabelSelector, appType, name, namespace string, objectMeta *metav1.ObjectMeta, labels map[string]string, replicas *int32, containers []corev1.Container, vs []corev1.Volume, svcs []corev1.Service, pods []corev1.Pod, pvcs []corev1.PersistentVolumeClaim, scm map[string]storagev1.StorageClass) (*Application, error) {
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

func IsPodHealthy(pod Pod) bool {
	phases := []corev1.PodPhase{corev1.PodRunning, corev1.PodSucceeded}
	return slices.Contains(phases, pod.Status.Phase)
}
