package core

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/node"
	apibase "github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/status"
	"github.com/juju/juju/rpc/params"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/otterscale/otterscale/internal/config"
)

// Test constants for GPU tests
const (
	testVGPUCoresPercent = 50.0  // vGPU cores usage percentage
	testVGPUMemoryMB     = 8192  // vGPU memory in MB
	testGPUMemoryMB      = 16384 // Physical GPU memory in MB
)

// Mock MachineRepo
type essMockMachineRepo struct {
	machines []Machine
	getErr   error
}

func (m *essMockMachineRepo) List(ctx context.Context) ([]Machine, error) {
	return m.machines, nil
}

func (m *essMockMachineRepo) Get(ctx context.Context, id string) (*Machine, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for _, mach := range m.machines {
		if mach.SystemID == id {
			return &mach, nil
		}
	}
	return nil, errors.New("not found")
}

// Add Commission method to satisfy MachineRepo interface
func (m *essMockMachineRepo) Commission(ctx context.Context, id string, params *entity.MachineCommissionParams) (*Machine, error) {
	for _, mach := range m.machines {
		if mach.SystemID == id {
			return &mach, nil
		}
	}
	return nil, errors.New("not found")
}

// Add PowerOff method to satisfy MachineRepo interface
func (m *essMockMachineRepo) PowerOff(ctx context.Context, id string, params *entity.MachinePowerOffParams) (*Machine, error) {
	for _, mach := range m.machines {
		if mach.SystemID == id {
			return &mach, nil
		}
	}
	return nil, errors.New("not found")
}

// Add Release method to satisfy MachineRepo interface
func (m *essMockMachineRepo) Release(ctx context.Context, id string, params *entity.MachineReleaseParams) (*Machine, error) {
	for _, mach := range m.machines {
		if mach.SystemID == id {
			return &mach, nil
		}
	}
	return nil, errors.New("not found")
}

// Mock ClientRepo
type essMockClientRepo struct {
	statusMap map[string]*params.FullStatus
}

func (m *essMockClientRepo) Status(ctx context.Context, uuid string, keys []string) (*params.FullStatus, error) {
	if s, ok := m.statusMap[uuid]; ok {
		return s, nil
	}
	return &params.FullStatus{}, nil
}

// Mock ScopeRepo
type essMockScopeRepo struct {
	scopes []Scope
}

func (m *essMockScopeRepo) List(ctx context.Context) ([]Scope, error) {
	return m.scopes, nil
}

func (m *essMockScopeRepo) Create(ctx context.Context, name string) (*Scope, error) {
	return &Scope{Name: name, UUID: "uuid"}, nil
}

// Mock FacilityRepo
type essMockFacilityRepo struct{}

// Add GetConfig mock implementation to satisfy FacilityRepo interface
func (m *essMockFacilityRepo) GetConfig(ctx context.Context, uuid, name string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (m *essMockFacilityRepo) GetLeader(ctx context.Context, uuid, name string) (string, error) {
	return "leader", nil
}

func (m *essMockFacilityRepo) GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error) {
	return &application.UnitInfo{PublicAddress: "10.0.0.1"}, nil
}

// AddUnits mock implementation to satisfy FacilityRepo interface
func (m *essMockFacilityRepo) AddUnits(ctx context.Context, uuid, app string, units int, placements []instance.Placement) ([]string, error) {
	return []string{}, nil
}

// Consume mock implementation to satisfy FacilityRepo interface
func (m *essMockFacilityRepo) Consume(ctx context.Context, uuid string, args *crossmodel.ConsumeApplicationArgs) error {
	return nil
}

// Create mock implementation to satisfy FacilityRepo interface
func (m *essMockFacilityRepo) Create(
	ctx context.Context,
	uuid, name, app, channel, series string,
	units, minUnits int,
	base *base.Base,
	placements []instance.Placement,
	constraints *constraints.Value,
	trusted bool,
) (*application.DeployInfo, error) {
	return &application.DeployInfo{}, nil
}

// CreateRelation mock implementation to satisfy FacilityRepo interface
func (m *essMockFacilityRepo) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	return &params.AddRelationResults{}, nil
}

// Delete mock implementation to satisfy FacilityRepo interface
func (m *essMockFacilityRepo) Delete(ctx context.Context, uuid, name string, force, destroyStorage bool) error {
	return nil
}

// Add DeleteRelation mock implementation to satisfy FacilityRepo interface
func (m *essMockFacilityRepo) DeleteRelation(ctx context.Context, uuid string, relationID int) error {
	return nil
}

// Add Expose mock implementation to satisfy FacilityRepo interface
func (m *essMockFacilityRepo) Expose(ctx context.Context, uuid, app string, endpoints map[string]params.ExposedEndpoint) error {
	return nil
}

// Add ResolveUnitErrors mock implementation to satisfy FacilityRepo interface
func (m *essMockFacilityRepo) ResolveUnitErrors(ctx context.Context, uuid string, units []string) error {
	return nil
}

// Add Update mock implementation to satisfy FacilityRepo interface
func (m *essMockFacilityRepo) Update(ctx context.Context, uuid, name, value string) error {
	return nil
}

// Mock ActionRepo
type essMockActionRepo struct{}

func (m *essMockActionRepo) List(ctx context.Context, uuid, appName string) (map[string]ActionSpec, error) {
	return map[string]ActionSpec{}, nil
}

func (m *essMockActionRepo) RunCommand(ctx context.Context, uuid, unitName, command string) (string, error) {
	return "", nil
}

func (m *essMockActionRepo) RunAction(ctx context.Context, uuid, unitName, actionName string, parameters map[string]any) (string, error) {
	return "", nil
}

func (m *essMockActionRepo) GetResult(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
	return &action.ActionResult{
		Status: "completed",
		Output: map[string]interface{}{
			"kubeconfig": `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: dGVzdA==
    server: https://1.2.3.4:6443
  name: test
contexts:
- context:
    cluster: test
    user: test
  name: test
current-context: test
kind: Config
users:
- name: test
  user:
    client-certificate-data: dGVzdA==
    client-key-data: dGVzdA==
`,
		},
	}, nil
}

type (
	mockFacilityOffersRepo struct{}
	mockConfig             struct {
		config.Config
	}
)

// Add GetConsumeDetails to satisfy FacilityOffersRepo interface
func (m *mockFacilityOffersRepo) GetConsumeDetails(ctx context.Context, url string) (params.ConsumeOfferDetails, error) {
	return params.ConsumeOfferDetails{
		Offer: &params.ApplicationOfferDetailsV5{
			OfferURL: url,
		},
		Macaroon: nil,
		ControllerInfo: &params.ExternalControllerInfo{
			ControllerTag: "controller-00000000-0000-4000-8000-000000000000",
			Alias:         "test-alias",
			Addrs:         []string{"test-addr"},
			CACert:        "test-cert",
		},
	}, nil
}

func TestEssentialUseCase_IsMachineDeployed(t *testing.T) {
	machines := []Machine{
		{
			Machine: &entity.Machine{
				SystemID:            "id1",
				Status:              node.StatusDeployed,
				WorkloadAnnotations: map[string]string{"juju-model-uuid": "uuid1"},
			},
			LastCommissioned: time.Now(),
		},
	}
	uc := NewEssentialUseCase(nil, &simpleMockKubeCoreRepo{}, &simpleMockKubeAppsRepo{}, &simpleMockActionRepo{}, nil, nil, nil, &essMockMachineRepo{machines: machines}, &simpleMockSubnetRepo{}, &simpleMockIPRangeRepo{}, nil, nil, &simpleMockTagRepo{})
	msg, ok, err := uc.IsMachineDeployed(context.Background(), "uuid1")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "", msg)
}

func TestEssentialUseCase_ListStatuses(t *testing.T) {
	client := &essMockClientRepo{
		statusMap: map[string]*params.FullStatus{
			"uuid": {
				Applications: map[string]params.ApplicationStatus{
					"ceph-mon": {
						Charm: "ch:amd64/ceph-mon-123",
						Status: params.DetailedStatus{
							Status: status.Blocked.String(),
							Info:   "blocked info",
						},
					},
					"irrelevant": {
						Charm: "ch:amd64/other-123",
						Status: params.DetailedStatus{
							Status: status.Active.String(),
							Info:   "active info",
						},
					},
				},
			},
		},
	}
	uc := NewEssentialUseCase(nil, &simpleMockKubeCoreRepo{}, &simpleMockKubeAppsRepo{}, &simpleMockActionRepo{}, nil, nil, nil, nil, &simpleMockSubnetRepo{}, &simpleMockIPRangeRepo{}, nil, client, &simpleMockTagRepo{})
	statuses, err := uc.ListStatuses(context.Background(), "uuid")
	assert.NoError(t, err)
	assert.NotEmpty(t, statuses)
	assert.Contains(t, statuses[0].Message, "[blocked]")
}

func TestEssentialUseCase_ListEssentials(t *testing.T) {
	scope := &essMockScopeRepo{
		scopes: []Scope{{UUID: "uuid", Name: "test", Status: apibase.Status{Status: status.Available}}},
	}
	client := &essMockClientRepo{
		statusMap: map[string]*params.FullStatus{
			"uuid": {
				Applications: map[string]params.ApplicationStatus{
					"kubernetes-control-plane": {
						Charm: "ch:amd64/kubernetes-control-plane-123",
						Units: map[string]params.UnitStatus{
							"unit/0": {Machine: "0"},
						},
					},
				},
			},
		},
	}
	uc := NewEssentialUseCase(nil, &simpleMockKubeCoreRepo{}, &simpleMockKubeAppsRepo{}, &simpleMockActionRepo{}, scope, nil, nil, nil, &simpleMockSubnetRepo{}, &simpleMockIPRangeRepo{}, nil, client, &simpleMockTagRepo{})
	essentials, err := uc.ListEssentials(context.Background(), 1, "uuid")
	assert.NoError(t, err)
	assert.NotEmpty(t, essentials)
	assert.Equal(t, "kubernetes-control-plane", essentials[0].Name)
}

func TestEssentialUseCase_CreateSingleNode(t *testing.T) {
	machines := []Machine{
		{
			Machine: &entity.Machine{
				SystemID:            "id1",
				WorkloadAnnotations: map[string]string{"juju-machine-id": "scope"},
			},
			LastCommissioned: time.Now(),
		},
	}
	machineRepo := &essMockMachineRepo{machines: machines}
	serverRepo := &essMockServerRepo{}
	facilityRepo := &essMockFacilityRepo{}
	scopeRepo := &essMockScopeRepo{}
	ipRangeRepo := &mockIPRangeRepo{}
	subnetRepo := &mockSubnetRepo{}
	facilityOffersRepo := &mockFacilityOffersRepo{}
	conf := &mockConfig{}
	clientRepo := &essMockClientRepo{
		statusMap: map[string]*params.FullStatus{
			"uuid": {
				Machines: map[string]params.MachineStatus{
					"1": {
						AgentStatus: params.DetailedStatus{Status: status.Started.String()},
					},
				},
			},
		},
	}
	uc := NewEssentialUseCase(&conf.Config, &simpleMockKubeCoreRepo{}, &simpleMockKubeAppsRepo{}, &simpleMockActionRepo{}, scopeRepo, facilityRepo, facilityOffersRepo, machineRepo, subnetRepo, ipRangeRepo, serverRepo, clientRepo, &simpleMockTagRepo{})
	err := uc.CreateSingleNode(context.Background(), "uuid", "id1", "prefix", []string{"10.0.0.2"}, "198.19.0.0/16", []string{"/dev/sda"})
	assert.Error(t, err) // because CreateCeph, CreateKubernetes, CreateCommon are not implemented
}

func TestEssentialUseCase_getMachineStatusMessage(t *testing.T) {
	machines := []Machine{
		{
			Machine: &entity.Machine{
				Hostname: "host1",
				Status:   node.StatusTesting,
			},
			LastCommissioned: time.Now(),
		},
		{
			Machine: &entity.Machine{
				Hostname: "host2",
				Status:   node.StatusTesting,
			},
			LastCommissioned: time.Now(),
		},
	}
	uc := NewEssentialUseCase(nil, &simpleMockKubeCoreRepo{}, &simpleMockKubeAppsRepo{}, &simpleMockActionRepo{}, nil, nil, nil, nil, &simpleMockSubnetRepo{}, &simpleMockIPRangeRepo{}, nil, nil, &simpleMockTagRepo{})
	msg := uc.getMachineStatusMessage(machines)
	assert.Contains(t, msg, "testing")
}

func TestEssentialUseCase_validateMachineStatus(t *testing.T) {
	machines := []Machine{
		{
			Machine: &entity.Machine{
				SystemID:            "id1",
				WorkloadAnnotations: map[string]string{"juju-machine-id": "1"},
				Status:              node.StatusDeployed,
			},
			LastCommissioned: time.Now(),
		},
	}
	machineRepo := &essMockMachineRepo{machines: machines}
	clientRepo := &essMockClientRepo{
		statusMap: map[string]*params.FullStatus{
			"uuid": {
				Machines: map[string]params.MachineStatus{
					"1": {
						AgentStatus: params.DetailedStatus{Status: status.Started.String()},
					},
				},
			},
		},
	}
	uc := NewEssentialUseCase(nil, &simpleMockKubeCoreRepo{}, &simpleMockKubeAppsRepo{}, &simpleMockActionRepo{}, nil, nil, nil, machineRepo, &simpleMockSubnetRepo{}, &simpleMockIPRangeRepo{}, nil, clientRepo, &simpleMockTagRepo{})
	err := uc.validateMachineStatus(context.Background(), "uuid", "id1")
	assert.NoError(t, err)
}

func TestNewCharmConfigs(t *testing.T) {
	configs := map[string]map[string]any{
		"ceph-mon":                 {"foo": "bar"},
		"kubernetes-control-plane": {"config": "value"},
	}
	result, err := NewCharmConfigs("prefix", configs)
	assert.NoError(t, err)
	assert.Contains(t, result, "ch:ceph-mon")
	assert.Contains(t, result, "ch:kubernetes-control-plane")

	// Check YAML content
	assert.Contains(t, result["ch:ceph-mon"], "prefix-ceph-mon")
	assert.Contains(t, result["ch:kubernetes-control-plane"], "prefix-kubernetes-control-plane")
}

func Test_formatAppCharm(t *testing.T) {
	name, ok := formatAppCharm("ch:amd64/ceph-mon-123")
	assert.True(t, ok)
	assert.Equal(t, "ceph-mon", name)

	_, ok = formatAppCharm("invalid-format")
	assert.False(t, ok)

	_, ok = formatAppCharm("ch:invalid")
	assert.False(t, ok)
}

func Test_formatEssentialCharm(t *testing.T) {
	assert.Equal(t, "ceph-mon", formatEssentialCharm("ch:ceph-mon"))
	assert.Equal(t, "kubernetes-control-plane", formatEssentialCharm("ch:kubernetes-control-plane"))
	assert.Equal(t, "simple-name", formatEssentialCharm("simple-name"))
}

func Test_toEssentialName(t *testing.T) {
	assert.Equal(t, "prefix-ceph-mon", toEssentialName("prefix", "ch:ceph-mon"))
	assert.Equal(t, "prefix-plain-name", toEssentialName("prefix", "plain-name"))
}

func Test_toEndpointList(t *testing.T) {
	relationList := [][]string{
		{"ch:app1", "ch:app2"},
		{"ch:app3"},
	}

	result := toEndpointList("prefix", relationList)
	assert.Equal(t, "prefix-app1", result[0][0])
	assert.Equal(t, "prefix-app2", result[0][1])
	assert.Equal(t, "prefix-app3", result[1][0])
}

func Test_isEssentialCharm(t *testing.T) {
	statusMap := map[string]params.ApplicationStatus{
		"ceph-mon":    {Charm: "ch:amd64/ceph-mon-123"},
		"invalid-app": {Charm: "invalid-charm-format"},
		"other-app":   {Charm: "ch:amd64/other-charm-456"},
	}
	charms := []EssentialCharm{{Name: "ch:ceph-mon"}}

	assert.True(t, isEssentialCharm(statusMap, "ceph-mon", charms))
	assert.False(t, isEssentialCharm(statusMap, "invalid-app", charms))
	assert.False(t, isEssentialCharm(statusMap, "other-app", charms))
}

func Test_createEssential(t *testing.T) {
	serverRepo := &essMockServerRepo{}
	machineRepo := &essMockMachineRepo{
		machines: []Machine{
			{
				Machine: &entity.Machine{
					SystemID:            "machine-id",
					WorkloadAnnotations: map[string]string{"juju-machine-id": "1"},
					Status:              node.StatusDeployed,
				},
			},
		},
	}
	facilityRepo := &essMockFacilityRepo{}
	tagRepo := &mockTagRepo{}

	err := createEssential(context.Background(), serverRepo, machineRepo, facilityRepo, tagRepo,
		"uuid", "machine-id", "prefix",
		[]EssentialCharm{
			{Name: "ch:test-charm", Machine: true},
		},
		map[string]string{"ch:test-charm": "config"}, []string{"tag1", "tag2"})
	assert.NoError(t, err)

	err = createEssential(context.Background(), serverRepo, machineRepo, facilityRepo, nil,
		"uuid", "machine-id", "prefix",
		[]EssentialCharm{
			{Name: "ch:test-charm", Machine: true},
		},
		map[string]string{"ch:test-charm": "config"}, nil)
	assert.NoError(t, err)

	err = createEssential(context.Background(), serverRepo, machineRepo, facilityRepo, tagRepo,
		"uuid", "", "prefix",
		[]EssentialCharm{
			{Name: "ch:test-charm", Machine: true},
		},
		map[string]string{"ch:test-charm": "config"}, nil)
	assert.NoError(t, err)
}

func Test_listEssentials(t *testing.T) {
	scope := &essMockScopeRepo{
		scopes: []Scope{{UUID: "uuid", Name: "test", Status: apibase.Status{Status: status.Available}}},
	}
	client := &essMockClientRepo{
		statusMap: map[string]*params.FullStatus{
			"uuid": {
				Applications: map[string]params.ApplicationStatus{
					"ceph-mon": {
						Charm: "ch:amd64/ceph-mon-123",
						Units: map[string]params.UnitStatus{
							"unit/0": {Machine: "0"},
						},
					},
				},
			},
		},
	}
	essentials, err := listEssentials(context.Background(), scope, client, "ceph-mon", 2, "uuid")
	assert.NoError(t, err)
	assert.NotEmpty(t, essentials)
	assert.Equal(t, "ceph-mon", essentials[0].Name)
}

func Test_createEssentialRelations(t *testing.T) {
	facilityRepo := &essMockFacilityRepo{}
	endpoints := [][]string{{"foo", "bar"}}
	err := createEssentialRelations(context.Background(), facilityRepo, "uuid", endpoints)
	assert.NoError(t, err)
}

func Test_getDirective(t *testing.T) {
	machines := []Machine{
		{
			Machine: &entity.Machine{
				SystemID:            "1",
				WorkloadAnnotations: map[string]string{"juju-machine-id": "1"},
				Status:              node.StatusDeployed,
			},
			LastCommissioned: time.Now(),
		},
	}
	machineRepo := &essMockMachineRepo{machines: machines}
	directive, err := getDirective(context.Background(), machineRepo, "1")
	assert.NoError(t, err)
	assert.Equal(t, "1", directive)
}

// Simple mocks for GPU testing

type simpleMockMachineRepo struct {
	machines []Machine
}

func (m *simpleMockMachineRepo) List(ctx context.Context) ([]Machine, error) {
	return m.machines, nil
}

func (m *simpleMockMachineRepo) Get(ctx context.Context, id string) (*Machine, error) {
	for _, machine := range m.machines {
		if machine.SystemID == id {
			return &machine, nil
		}
	}
	return nil, errors.New("machine not found")
}

func (m *simpleMockMachineRepo) Commission(ctx context.Context, id string, params *entity.MachineCommissionParams) (*Machine, error) {
	return nil, nil
}

func (m *simpleMockMachineRepo) PowerOff(ctx context.Context, id string, params *entity.MachinePowerOffParams) (*Machine, error) {
	return nil, nil
}

func (m *simpleMockMachineRepo) Release(ctx context.Context, id string, params *entity.MachineReleaseParams) (*Machine, error) {
	return nil, nil
}

type simpleMockKubeCoreRepo struct {
	nodes map[string]*Node
	pods  []Pod
}

func (m *simpleMockKubeCoreRepo) GetNode(ctx context.Context, config *rest.Config, name string) (*Node, error) {
	if node, exists := m.nodes[name]; exists {
		return node, nil
	}
	return nil, errors.New("node not found")
}

func (m *simpleMockKubeCoreRepo) ListPods(ctx context.Context, config *rest.Config, namespace string) ([]Pod, error) {
	return m.pods, nil
}

func (m *simpleMockKubeCoreRepo) ListPodsByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]Pod, error) {
	return m.pods, nil
}

// Minimal implementations for required interface methods
func (m *simpleMockKubeCoreRepo) UpdateNode(ctx context.Context, config *rest.Config, node *Node) (*Node, error) {
	return node, nil
}
func (m *simpleMockKubeCoreRepo) ListNamespaces(ctx context.Context, config *rest.Config) ([]Namespace, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) ListServices(ctx context.Context, config *rest.Config, namespace string) ([]Service, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) ListServicesByOptions(ctx context.Context, config *rest.Config, namespace, label, field string) ([]Service, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) ListVirtualMachineServices(ctx context.Context, config *rest.Config, namespace, vmName string) ([]Service, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) GetService(ctx context.Context, config *rest.Config, namespace, name string) (*Service, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) CreateVirtualMachineService(ctx context.Context, config *rest.Config, namespace, name, vmName string, ports []corev1.ServicePort) (*Service, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) UpdateService(ctx context.Context, config *rest.Config, namespace string, service *Service) (*Service, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) DeleteService(ctx context.Context, config *rest.Config, namespace, name string) error {
	return nil
}
func (m *simpleMockKubeCoreRepo) GetLogs(ctx context.Context, config *rest.Config, namespace, podName, containerName string) (string, error) {
	return "", nil
}
func (m *simpleMockKubeCoreRepo) DeletePod(ctx context.Context, config *rest.Config, namespace, name string) error {
	return nil
}
func (m *simpleMockKubeCoreRepo) StreamLogs(ctx context.Context, config *rest.Config, namespace, podName, containerName string) (io.ReadCloser, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) CreateExecutor(config *rest.Config, namespace, podName, containerName string, command []string) (remotecommand.Executor, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) ListPersistentVolumeClaims(ctx context.Context, config *rest.Config, namespace string) ([]PersistentVolumeClaim, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) GetPersistentVolumeClaim(ctx context.Context, config *rest.Config, namespace, name string) (*PersistentVolumeClaim, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) PatchPersistentVolumeClaim(ctx context.Context, config *rest.Config, namespace, name string, data []byte) (*PersistentVolumeClaim, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) GetNamespace(ctx context.Context, config *rest.Config, name string) (*Namespace, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) CreateNamespace(ctx context.Context, config *rest.Config, name string) (*Namespace, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) GetConfigMap(ctx context.Context, config *rest.Config, namespace, name string) (*ConfigMap, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) CreateConfigMap(ctx context.Context, config *rest.Config, namespace, name string, data map[string]string) (*ConfigMap, error) {
	return nil, nil
}
func (m *simpleMockKubeCoreRepo) GetSecret(ctx context.Context, config *rest.Config, namespace, name string) (*Secret, error) {
	return nil, nil
}

type simpleMockKubeAppsRepo struct {
	deployments []Deployment
}

func (m *simpleMockKubeAppsRepo) ListDeploymentsByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]Deployment, error) {
	return m.deployments, nil
}

func (m *simpleMockKubeAppsRepo) GetDeployment(ctx context.Context, config *rest.Config, namespace, name string) (*Deployment, error) {
	for _, deployment := range m.deployments {
		if deployment.Name == name && deployment.Namespace == namespace {
			return &deployment, nil
		}
	}
	return nil, errors.New("deployment not found")
}

// Minimal implementations for required interface methods
func (m *simpleMockKubeAppsRepo) ListDeployments(ctx context.Context, config *rest.Config, namespace string) ([]Deployment, error) {
	return m.deployments, nil
}
func (m *simpleMockKubeAppsRepo) UpdateDeployment(ctx context.Context, config *rest.Config, namespace string, deployment *Deployment) (*Deployment, error) {
	return deployment, nil
}
func (m *simpleMockKubeAppsRepo) ListStatefulSets(ctx context.Context, config *rest.Config, namespace string) ([]StatefulSet, error) {
	return nil, nil
}
func (m *simpleMockKubeAppsRepo) GetStatefulSet(ctx context.Context, config *rest.Config, namespace, name string) (*StatefulSet, error) {
	return nil, nil
}
func (m *simpleMockKubeAppsRepo) UpdateStatefulSet(ctx context.Context, config *rest.Config, namespace string, statefulSet *StatefulSet) (*StatefulSet, error) {
	return nil, nil
}
func (m *simpleMockKubeAppsRepo) ListDaemonSets(ctx context.Context, config *rest.Config, namespace string) ([]DaemonSet, error) {
	return nil, nil
}
func (m *simpleMockKubeAppsRepo) GetDaemonSet(ctx context.Context, config *rest.Config, namespace, name string) (*DaemonSet, error) {
	return nil, nil
}
func (m *simpleMockKubeAppsRepo) UpdateDaemonSet(ctx context.Context, config *rest.Config, namespace string, daemonSet *DaemonSet) (*DaemonSet, error) {
	return nil, nil
}

type simpleMockFacilityRepo struct{}

func (m *simpleMockFacilityRepo) GetLeader(ctx context.Context, uuid, name string) (string, error) {
	return "test-leader", nil
}

// Minimal implementations for required interface methods
func (m *simpleMockFacilityRepo) GetConfig(ctx context.Context, uuid, name string) (map[string]interface{}, error) {
	return nil, nil
}
func (m *simpleMockFacilityRepo) GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error) {
	return nil, nil
}
func (m *simpleMockFacilityRepo) AddUnits(ctx context.Context, uuid, app string, units int, placements []instance.Placement) ([]string, error) {
	return nil, nil
}
func (m *simpleMockFacilityRepo) Consume(ctx context.Context, uuid string, args *crossmodel.ConsumeApplicationArgs) error {
	return nil
}
func (m *simpleMockFacilityRepo) Create(ctx context.Context, uuid, name, app, channel, series string, units, minUnits int, base *base.Base, placements []instance.Placement, constraints *constraints.Value, trusted bool) (*application.DeployInfo, error) {
	return nil, nil
}
func (m *simpleMockFacilityRepo) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	return nil, nil
}
func (m *simpleMockFacilityRepo) Delete(ctx context.Context, uuid, name string, force, destroyStorage bool) error {
	return nil
}
func (m *simpleMockFacilityRepo) DeleteRelation(ctx context.Context, uuid string, relationID int) error {
	return nil
}
func (m *simpleMockFacilityRepo) Expose(ctx context.Context, uuid, app string, endpoints map[string]params.ExposedEndpoint) error {
	return nil
}
func (m *simpleMockFacilityRepo) ResolveUnitErrors(ctx context.Context, uuid string, units []string) error {
	return nil
}
func (m *simpleMockFacilityRepo) Update(ctx context.Context, uuid, name, value string) error {
	return nil
}

type simpleMockActionRepo struct{}

func (m *simpleMockActionRepo) List(ctx context.Context, uuid, appName string) (map[string]ActionSpec, error) {
	return nil, nil
}
func (m *simpleMockActionRepo) RunCommand(ctx context.Context, uuid, unitName, command string) (string, error) {
	return "", nil
}
func (m *simpleMockActionRepo) RunAction(ctx context.Context, uuid, unitName, actionName string, parameters map[string]any) (string, error) {
	return "", nil
}
func (m *simpleMockActionRepo) GetResult(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
	return &action.ActionResult{
		Status: "completed",
		Output: map[string]interface{}{
			"kubeconfig": `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: dGVzdA==
    server: https://1.2.3.4:6443
  name: test
contexts:
- context:
    cluster: test
    user: test
  name: test
current-context: test
kind: Config
users:
- name: test
  user:
    client-certificate-data: dGVzdA==
    client-key-data: dGVzdA==
`,
		},
	}, nil
}

type simpleMockSubnetRepo struct{}

func (m *simpleMockSubnetRepo) List(ctx context.Context) ([]Subnet, error)       { return nil, nil }
func (m *simpleMockSubnetRepo) Get(ctx context.Context, id int) (*Subnet, error) { return nil, nil }
func (m *simpleMockSubnetRepo) Create(ctx context.Context, params *entity.SubnetParams) (*Subnet, error) {
	return nil, nil
}
func (m *simpleMockSubnetRepo) Update(ctx context.Context, id int, params *entity.SubnetParams) (*Subnet, error) {
	return nil, nil
}
func (m *simpleMockSubnetRepo) Delete(ctx context.Context, id int) error { return nil }
func (m *simpleMockSubnetRepo) GetIPAddresses(ctx context.Context, id int) ([]IPAddress, error) {
	return nil, nil
}
func (m *simpleMockSubnetRepo) GetStatistics(ctx context.Context, id int) (*NetworkStatistics, error) {
	return nil, nil
}

type simpleMockIPRangeRepo struct{}

func (m *simpleMockIPRangeRepo) List(ctx context.Context) ([]IPRange, error)       { return nil, nil }
func (m *simpleMockIPRangeRepo) Get(ctx context.Context, id int) (*IPRange, error) { return nil, nil }
func (m *simpleMockIPRangeRepo) Create(ctx context.Context, params *entity.IPRangeParams) (*IPRange, error) {
	return nil, nil
}
func (m *simpleMockIPRangeRepo) Update(ctx context.Context, id int, params *entity.IPRangeParams) (*IPRange, error) {
	return nil, nil
}
func (m *simpleMockIPRangeRepo) Delete(ctx context.Context, id int) error { return nil }

type simpleMockTagRepo struct{}

func (m *simpleMockTagRepo) List(ctx context.Context) ([]Tag, error)            { return nil, nil }
func (m *simpleMockTagRepo) Get(ctx context.Context, name string) (*Tag, error) { return nil, nil }
func (m *simpleMockTagRepo) Create(ctx context.Context, name, comment string) (*Tag, error) {
	return nil, nil
}
func (m *simpleMockTagRepo) Delete(ctx context.Context, name string) error { return nil }
func (m *simpleMockTagRepo) AddMachines(ctx context.Context, name string, machineIDs []string) error {
	return nil
}
func (m *simpleMockTagRepo) RemoveMachines(ctx context.Context, name string, machineIDs []string) error {
	return nil
}

type essMockServerRepo struct{}

func (m *essMockServerRepo) Get(ctx context.Context, name string) (string, error) {
	return "focal", nil
}

func (m *essMockServerRepo) Update(ctx context.Context, name, value string) error {
	return nil
}

// Test for ListGPURelationsByMachine
func TestEssentialUseCase_ListGPURelationsByMachine(t *testing.T) {
	// Setup test data
	machineID := "test-machine-id"
	hostname := "test-node"

	// Mock machine with hostname
	machines := []Machine{
		{
			Machine: &entity.Machine{
				SystemID: machineID,
				Hostname: hostname,
				Status:   node.StatusDeployed,
			},
		},
	}

	// Mock node with GPU annotations
	nodes := map[string]*Node{
		hostname: {
			ObjectMeta: metav1.ObjectMeta{
				Name: hostname,
				Annotations: map[string]string{
					"hami.io/node-nvidia-register": "GPU-663aa370-535a-33b8-e01f-b325fb2025c7,10,24564,100,NVIDIA-NVIDIA GeForce RTX 4090,0,true,0,hami-core:GPU-c15ecdf3-444a-2d02-29e9-e978b2514335,10,24564,100,NVIDIA-NVIDIA GeForce RTX 3080,0,true,0,hami-core:",
				},
			},
		},
	}

	// Mock pods with vGPU annotations
	pods := []Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"hami.io/vgpu-devices-allocated": "GPU-663aa370-535a-33b8-e01f-b325fb2025c7,NVIDIA,4096,50:",
					"hami.io/bind-time":              "1609459200",
					"hami.io/bind-phase":             "Bound",
				},
				Labels: map[string]string{
					"model-name": "test-model",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: hostname,
			},
		},
	}

	// Create mocks
	machineRepo := &simpleMockMachineRepo{machines: machines}
	kubeCoreRepo := &simpleMockKubeCoreRepo{nodes: nodes, pods: pods}
	kubeAppsRepo := &simpleMockKubeAppsRepo{}
	facilityRepo := &simpleMockFacilityRepo{}
	actionRepo := &simpleMockActionRepo{}

	// Create use case
	uc := NewEssentialUseCase(
		&config.Config{}, // conf
		kubeCoreRepo,     // kubeCore
		kubeAppsRepo,     // kubeApps
		actionRepo,       // action
		nil,              // scope
		facilityRepo,     // facility
		nil,              // facilityOffers
		machineRepo,      // machine
		nil,              // subnet
		nil,              // ipRange
		nil,              // server
		nil,              // client
		nil,              // tag
	)

	// Test
	relations, err := uc.ListGPURelationsByMachine(context.Background(), "test-scope", "test-facility", machineID)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, relations)

	// Check machine relation
	machineFound := false
	gpuFound := false
	podFound := false
	vgpuFound := false

	for _, relation := range relations {
		switch {
		case relation.GetMachine() != nil:
			machineFound = true
			assert.Equal(t, machineID, relation.GetMachine().GetId())
			assert.Equal(t, hostname, relation.GetMachine().GetHostname())
		case relation.GetGpu() != nil:
			gpuFound = true
			assert.Equal(t, machineID, relation.GetGpu().GetMachineId())
			assert.Equal(t, "NVIDIA", relation.GetGpu().GetVendor())
		case relation.GetPod() != nil:
			podFound = true
			assert.Equal(t, "test-pod", relation.GetPod().GetName())
			assert.Equal(t, "default", relation.GetPod().GetNamespace())
		case relation.GetVgpu() != nil:
			vgpuFound = true
			assert.Equal(t, "test-pod", relation.GetVgpu().GetPodName())
			assert.Equal(t, float32(testVGPUCoresPercent), relation.GetVgpu().GetVcoresPercent())
		}
	}

	assert.True(t, machineFound, "Machine relation should be found")
	assert.True(t, gpuFound, "GPU relation should be found")
	assert.True(t, podFound, "Pod relation should be found")
	assert.True(t, vgpuFound, "vGPU relation should be found")
}

// Test for ListGPURelationsByModel
func TestEssentialUseCase_ListGPURelationsByModel(t *testing.T) {
	// Setup test data
	modelName := "test-model"
	namespace := "default"

	// Mock deployment with model label
	deployments := []Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-deployment",
				Namespace: namespace,
				Labels: map[string]string{
					"model-name": modelName,
				},
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "test-app",
					},
				},
			},
		},
	}

	// Mock pods with vGPU annotations
	pods := []Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: namespace,
				Annotations: map[string]string{
					"hami.io/vgpu-devices-allocated": "GPU-663aa370-535a-33b8-e01f-b325fb2025c7,NVIDIA,4096,50:",
					"hami.io/bind-time":              "1609459200",
					"hami.io/bind-phase":             "Bound",
				},
				Labels: map[string]string{
					"app":        "test-app",
					"model-name": modelName,
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "test-node",
			},
		},
	}

	// Mock machines for node mapping
	machines := []Machine{
		{
			Machine: &entity.Machine{
				SystemID: "test-machine-id",
				Hostname: "test-node",
			},
		},
	}

	// Mock nodes with GPU annotations
	nodes := map[string]*Node{
		"test-node": {
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-node",
				Annotations: map[string]string{
					"hami.io/node-nvidia-register": "GPU-663aa370-535a-33b8-e01f-b325fb2025c7,10,24564,100,NVIDIA-NVIDIA GeForce RTX 4090,0,true,0,hami-core:",
				},
			},
		},
	}

	// Create mocks
	machineRepo := &simpleMockMachineRepo{machines: machines}
	kubeCoreRepo := &simpleMockKubeCoreRepo{nodes: nodes, pods: pods}
	kubeAppsRepo := &simpleMockKubeAppsRepo{deployments: deployments}
	facilityRepo := &simpleMockFacilityRepo{}
	actionRepo := &simpleMockActionRepo{}

	// Create use case
	uc := NewEssentialUseCase(
		&config.Config{}, // conf
		kubeCoreRepo,     // kubeCore
		kubeAppsRepo,     // kubeApps
		actionRepo,       // action
		nil,              // scope
		facilityRepo,     // facility
		nil,              // facilityOffers
		machineRepo,      // machine
		nil,              // subnet
		nil,              // ipRange
		nil,              // server
		nil,              // client
		nil,              // tag
	)

	// Test
	relations, err := uc.ListGPURelationsByModel(context.Background(), "test-scope", "test-facility", namespace, modelName)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, relations)

	// Check that we have different types of relations
	machineFound := false
	gpuFound := false
	podFound := false
	vgpuFound := false

	for _, relation := range relations {
		switch {
		case relation.GetMachine() != nil:
			machineFound = true
			assert.Equal(t, "test-machine-id", relation.GetMachine().GetId())
		case relation.GetGpu() != nil:
			gpuFound = true
			assert.Equal(t, "test-machine-id", relation.GetGpu().GetMachineId())
			assert.Equal(t, "NVIDIA", relation.GetGpu().GetVendor())
		case relation.GetPod() != nil:
			podFound = true
			assert.Equal(t, "test-pod", relation.GetPod().GetName())
			assert.Equal(t, modelName, relation.GetPod().GetModelName())
		case relation.GetVgpu() != nil:
			vgpuFound = true
			assert.Equal(t, "test-pod", relation.GetVgpu().GetPodName())
			assert.Equal(t, float32(testVGPUCoresPercent), relation.GetVgpu().GetVcoresPercent())
		}
	}

	assert.True(t, machineFound, "Machine relation should be found")
	assert.True(t, gpuFound, "GPU relation should be found")
	assert.True(t, podFound, "Pod relation should be found")
	assert.True(t, vgpuFound, "vGPU relation should be found")
}
