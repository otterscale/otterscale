package service

import (
	"context"
	"time"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/subnet"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/api/client/application"
	corebase "github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"

	helmaction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/client-go/rest"

	"github.com/openhdc/otterscale/internal/domain/model"
)

type MAASServer interface {
	Get(ctx context.Context, name string) ([]byte, error)
	Update(ctx context.Context, name, value string) error
}

type MAASPackageRepository interface {
	List(ctx context.Context) ([]entity.PackageRepository, error)
	Update(ctx context.Context, id int, params *entity.PackageRepositoryParams) (*entity.PackageRepository, error)
}

type MAASBootResource interface {
	List(ctx context.Context) ([]entity.BootResource, error)
	Import(ctx context.Context) error
	IsImporting(ctx context.Context) (bool, error)
}

type MAASBootSource interface {
	List(ctx context.Context) ([]entity.BootSource, error)
}

type MAASBootSourceSelection interface {
	List(ctx context.Context, id int) ([]entity.BootSourceSelection, error)
	Create(ctx context.Context, bootSourceID int, params *entity.BootSourceSelectionParams) (*entity.BootSourceSelection, error)
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
	Delete(ctx context.Context, id int) error
}

type MAASMachine interface {
	List(ctx context.Context) ([]entity.Machine, error)
	Get(ctx context.Context, systemID string) (*entity.Machine, error)
	Release(ctx context.Context, systemID string, params *entity.MachineReleaseParams) (*entity.Machine, error)
	PowerOn(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error)
	PowerOff(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error)
	Commission(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error)
}

type MAASTag interface {
	List(ctx context.Context) ([]entity.Tag, error)
	Get(ctx context.Context, name string) (*entity.Tag, error)
	Create(ctx context.Context, name, comment string) (*entity.Tag, error)
	Delete(ctx context.Context, name string) error
	AddMachines(ctx context.Context, name string, machineIDs []string) error
	RemoveMachines(ctx context.Context, name string, machineIDs []string) error
}

type MAASSSHKey interface {
	List(ctx context.Context) ([]entity.SSHKey, error)
}

type JujuKey interface {
	Add(ctx context.Context, uuid, key string) ([]params.ErrorResult, error)
}

type JujuMachine interface {
	AddMachines(ctx context.Context, uuid string, params []params.AddMachineParams) ([]params.AddMachinesResult, error)
	DestroyMachines(_ context.Context, uuid string, force, keep, dryRun bool, maxWait *time.Duration, machines ...string) ([]params.DestroyMachineResult, error)
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
	Create(ctx context.Context, uuid, name string, configYAML string, charmName, channel string, revision, number int, base *corebase.Base, placements []instance.Placement, constraint *constraints.Value, trust bool) (*application.DeployInfo, error)
	Update(ctx context.Context, uuid, name string, configYAML string) error
	Delete(ctx context.Context, uuid, name string, destroyStorage, force bool) error
	Expose(ctx context.Context, uuid, name string, endpoints map[string]params.ExposedEndpoint) error
	AddUnits(ctx context.Context, uuid, name string, number int, placements []instance.Placement) ([]string, error)
	ResolveUnitErrors(ctx context.Context, uuid string, units []string) error
	CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error)
	DeleteRelation(ctx context.Context, uuid string, id int) error
	GetConfig(ctx context.Context, uuid string, name string) (map[string]any, error)
	GetLeader(ctx context.Context, uuid, name string) (string, error)
	GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error)
}

type JujuAction interface {
	List(ctx context.Context, uuid, appName string) (map[string]action.ActionSpec, error)
}

type JujuCharm interface {
	List(ctx context.Context) ([]model.Charm, error)
	Get(ctx context.Context, name string) (*model.Charm, error)
	ListArtifacts(ctx context.Context, name string) ([]model.CharmArtifact, error)
}

type KubeApps interface {
	ListDeployments(ctx context.Context, config *rest.Config, namespace string) ([]appsv1.Deployment, error)
	GetDeployment(ctx context.Context, config *rest.Config, namespace, name string) (*appsv1.Deployment, error)
	ListStatefulSets(ctx context.Context, config *rest.Config, namespace string) ([]appsv1.StatefulSet, error)
	GetStatefulSet(ctx context.Context, config *rest.Config, namespace, name string) (*appsv1.StatefulSet, error)
	ListDaemonSets(ctx context.Context, config *rest.Config, namespace string) ([]appsv1.DaemonSet, error)
	GetDaemonSet(ctx context.Context, config *rest.Config, namespace, name string) (*appsv1.DaemonSet, error)
}

type KubeCore interface {
	ListServices(ctx context.Context, config *rest.Config, namespace string) ([]corev1.Service, error)
	ListPods(ctx context.Context, config *rest.Config, namespace string) ([]corev1.Pod, error)
	ListPersistentVolumeClaims(ctx context.Context, config *rest.Config, namespace string) ([]corev1.PersistentVolumeClaim, error)
}

type KubeHelm interface {
	ListReleases(config *rest.Config, namespace string) ([]release.Release, error)
	InstallRelease(config *rest.Config, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error)
	UninstallRelease(config *rest.Config, namespace, name string, dryRun bool) (*release.Release, error)
	UpgradeRelease(config *rest.Config, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error)
	RollbackRelease(config *rest.Config, namespace, name string, dryRun bool) error
	GetValues(config *rest.Config, namespace, name string) (map[string]any, error)
}

type KubeHelmChart interface {
	List(ctx context.Context) ([]*repo.IndexFile, error)
	Show(chartRef string, format helmaction.ShowOutputFormat) (string, error)
}

type KubeStorage interface {
	ListStorageClasses(ctx context.Context, config *rest.Config) ([]storagev1.StorageClass, error)
	GetStorageClass(ctx context.Context, config *rest.Config, name string) (*storagev1.StorageClass, error)
}
