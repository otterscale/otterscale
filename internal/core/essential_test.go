package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/node"
	apibase "github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/base"

	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/status"
	"github.com/juju/juju/rpc/params"
	"github.com/stretchr/testify/assert"

	"github.com/otterscale/otterscale/internal/config"
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

// Mock ServerRepo
type essMockServerRepo struct{}

func (m *essMockServerRepo) Get(ctx context.Context, name string) (string, error) {
	return "focal", nil
}

func (m *essMockServerRepo) Update(ctx context.Context, name, value string) error {
	return nil
}

type (
	mockFacilityOffersRepo struct{}
	mockConfig             struct {
		config.Config
	}
)

// Add GetConsumeDetails to satisfy FacilityOffersRepo interface
func (m *mockFacilityOffersRepo) GetConsumeDetails(ctx context.Context, uuid string) (params.ConsumeOfferDetails, error) {
	return params.ConsumeOfferDetails{}, nil
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
	uc := NewEssentialUseCase(nil, nil, nil, nil, &essMockMachineRepo{machines: machines}, nil, nil, nil, nil)
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
	uc := NewEssentialUseCase(nil, nil, nil, nil, nil, nil, nil, nil, client)
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
	uc := NewEssentialUseCase(nil, scope, nil, nil, nil, nil, nil, nil, client)
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
	uc := NewEssentialUseCase(&conf.Config, scopeRepo, facilityRepo, facilityOffersRepo, machineRepo, subnetRepo, ipRangeRepo, serverRepo, clientRepo)
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
	uc := NewEssentialUseCase(nil, nil, nil, nil, nil, nil, nil, nil, nil)
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
	uc := NewEssentialUseCase(nil, nil, nil, nil, machineRepo, nil, nil, nil, clientRepo)
	err := uc.validateMachineStatus(context.Background(), "uuid", "id1")
	assert.NoError(t, err)
}

func TestNewCharmConfigs(t *testing.T) {
	configs := map[string]map[string]any{
		"ceph-mon": {"foo": "bar"},
	}
	result, err := NewCharmConfigs("prefix", configs)
	assert.NoError(t, err)
	assert.Contains(t, result, "ch:ceph-mon")
}

func Test_formatAppCharm(t *testing.T) {
	name, ok := formatAppCharm("ch:amd64/ceph-mon-123")
	assert.True(t, ok)
	assert.Equal(t, "ceph-mon", name)
}

func Test_formatEssentialCharm(t *testing.T) {
	assert.Equal(t, "ceph-mon", formatEssentialCharm("ch:ceph-mon"))
}

func Test_isEssentialCharm(t *testing.T) {
	statusMap := map[string]params.ApplicationStatus{
		"ceph-mon": {Charm: "ch:amd64/ceph-mon-123"},
	}
	charms := []EssentialCharm{{Name: "ch:ceph-mon"}}
	assert.True(t, isEssentialCharm(statusMap, "ceph-mon", charms))
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

func Test_toEssentialName(t *testing.T) {
	assert.Equal(t, "prefix-ceph-mon", toEssentialName("prefix", "ceph-mon"))
	assert.Equal(t, "prefix-ceph-mon", toEssentialName("prefix", "ch:ceph-mon"))
}

func Test_toEndpointList(t *testing.T) {
	relations := [][]string{{"ceph-mon", "ceph-osd"}}
	endpoints := toEndpointList("prefix", relations)
	assert.Equal(t, "prefix-ceph-mon", endpoints[0][0])
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
