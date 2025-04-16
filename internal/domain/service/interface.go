package service

import (
	"context"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/subnet"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"

	helmaction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"

	"k8s.io/client-go/rest"
)

type MAASServer interface {
	Get(ctx context.Context, name string) (string, error)
	Update(ctx context.Context, name, value string) error
}

type MAASPackageRepository interface {
	List(ctx context.Context) ([]entity.PackageRepository, error)
	Update(ctx context.Context, id int, params *entity.PackageRepositoryParams) (*entity.PackageRepository, error)
}

type MAASBootResource interface {
	List(ctx context.Context) ([]entity.BootResource, error)
	Import(ctx context.Context) error
}

type MAASFabric interface {
	List(ctx context.Context) ([]entity.Fabric, error)
	Get(ctx context.Context, id int) (*entity.Fabric, error)
	Create(ctx context.Context, params *entity.FabricParams) (*entity.Fabric, error)
	Update(ctx context.Context, id int, params *entity.FabricParams) (*entity.Fabric, error)
	Delete(ctx context.Context, id int) error
}

type MAASVLAN interface {
	Update(ctx context.Context, fabricID, vid int, params *entity.VLANParams) (*entity.VLAN, error)
}

type MAASSubnet interface {
	List(ctx context.Context) ([]entity.Subnet, error)
	Get(ctx context.Context, id int) (*entity.Subnet, error)
	Create(ctx context.Context, params *entity.SubnetParams) (*entity.Subnet, error)
	Update(ctx context.Context, id int, params *entity.SubnetParams) (*entity.Subnet, error)
	Delete(ctx context.Context, id int) error
	GetIPAddresses(ctx context.Context, id int) ([]subnet.IPAddress, error)
	GetReservedIPRanges(ctx context.Context, id int) ([]subnet.ReservedIPRange, error)
	GetUnreservedIPRanges(ctx context.Context, id int) ([]subnet.IPRange, error)
	GetStatistics(ctx context.Context, id int) (*subnet.Statistics, error)
}

type MAASIPRange interface {
	List(ctx context.Context) ([]entity.IPRange, error)
	Create(ctx context.Context, params *entity.IPRangeParams) (*entity.IPRange, error)
	Update(ctx context.Context, id int, params *entity.IPRangeParams) (*entity.IPRange, error)
}

type MAASMachine interface {
	List(ctx context.Context) ([]entity.Machine, error)
	Get(ctx context.Context, systemID string) (*entity.Machine, error)
	PowerOn(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error)
	PowerOff(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error)
	Commission(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error)
}

type JujuMachine interface {
	AddMachines(ctx context.Context, uuid string, params []params.AddMachineParams) ([]params.AddMachinesResult, error)
}

type JujuClient interface {
	Status(ctx context.Context, uuid string, patterns []string) (*params.FullStatus, error)
}

type JujuModel interface {
	List(ctx context.Context) ([]base.UserModelSummary, error)
	Create(ctx context.Context, name string) (*base.ModelInfo, error)
}

type JujuModelConfig interface {
	List(ctx context.Context, uuid string) (map[string]any, error)
	Set(ctx context.Context, uuid string, config map[string]any) error
	Unset(ctx context.Context, uuid string, keys ...string) error
}

type JujuApplication interface {
	Create(ctx context.Context, uuid, charmName, appName, channel string, revision, number int, config map[string]string, constraint constraints.Value, placements []instance.Placement, trust bool) error
	Update(ctx context.Context, uuid, name string, config map[string]string) error
	Delete(ctx context.Context, uuid, name string, destroyStorage, force bool) error
	Expose(ctx context.Context, uuid, name string, endpoints map[string]params.ExposedEndpoint) error
	AddUnits(ctx context.Context, uuid, name string, number int, placements []instance.Placement) ([]string, error)
	ResolveUnitErrors(ctx context.Context, uuid string, units []string) error
	CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error)
	DeleteRelation(ctx context.Context, uuid string, id int) error
	GetConfigs(ctx context.Context, uuid string, name ...string) (map[string]map[string]any, error)
	GetLeader(ctx context.Context, uuid, name string) (string, error)
	GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error)
}

type JujuAction interface {
	List(ctx context.Context, uuid, appName string) (map[string]action.ActionSpec, error)
}

type KubeClient interface {
	Exists(key string) bool
	Add(key string, cfg *rest.Config) error
}

type KubeApps interface {
	ListDeployments(ctx context.Context, key, namespace string) ([]appsv1.Deployment, error)
	GetDeployment(ctx context.Context, key, namespace, name string) (*appsv1.Deployment, error)
	ListStatefulSets(ctx context.Context, key, namespace string) ([]appsv1.StatefulSet, error)
	GetStatefulSet(ctx context.Context, key, namespace, name string) (*appsv1.StatefulSet, error)
	ListDaemonSets(ctx context.Context, key, namespace string) ([]appsv1.DaemonSet, error)
	GetDaemonSet(ctx context.Context, key, namespace, name string) (*appsv1.DaemonSet, error)
}

type KubeBatch interface {
	GetCronJob(ctx context.Context, key, namespace, name string) (*batchv1.CronJob, error)
	CreateCronJob(ctx context.Context, key, namespace, name, image, schedule string) (*batchv1.CronJob, error)
	UpdateCronJob(ctx context.Context, key, namespace, name, image, schedule string) (*batchv1.CronJob, error)
	DeleteCronJob(ctx context.Context, key, namespace, name string) error
	ListJobsFromCronJob(ctx context.Context, key, namespace string, cronJob *batchv1.CronJob) (*batchv1.JobList, error)
	CreateJobFromCronJob(ctx context.Context, key, namespace string, cronJob *batchv1.CronJob, createdBy string) (*batchv1.Job, error)
}

type KubeCore interface {
	GetNamespace(ctx context.Context, key, name string) (*corev1.Namespace, error)
	CreateNamespace(ctx context.Context, key, name string) (*corev1.Namespace, error)
	ListServices(ctx context.Context, key, namespace string) ([]corev1.Service, error)
	ListPods(ctx context.Context, key, namespace string) ([]corev1.Pod, error)
	ListPersistentVolumeClaims(ctx context.Context, key, namespace string) ([]corev1.PersistentVolumeClaim, error)
}

type KubeStorage interface {
	ListStorageClasses(ctx context.Context, key string) ([]storagev1.StorageClass, error)
}

type KubeHelm interface {
	ListReleases(key, namespace string) ([]release.Release, error)
	InstallRelease(key, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error)
	UninstallRelease(key, namespace, name string, dryRun bool) (*release.Release, error)
	UpgradeRelease(key, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error)
	RollbackRelease(key, namespace, name string, dryRun bool) error
	GetChartInfo(chartRef string, format helmaction.ShowOutputFormat) (string, error)
	ListChartVersions(ctx context.Context) (map[string]repo.ChartVersions, error)
}
