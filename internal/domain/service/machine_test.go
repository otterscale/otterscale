package service

import (
	"context"
	"errors"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	params "github.com/juju/juju/rpc/params"
	"github.com/openhdc/otterscale/internal/domain/model"
)

type Status struct {
	Machines map[string]model.MachineStatus
}

// Minimal mock for model.Status to satisfy test usage
// Remove this if model.Status is imported from the correct package.
var _ = Status{}

// Ensure this struct matches the fields in model.Machine if needed for tests.
type Machine struct {
	ID                  string
	WorkloadAnnotations map[string]string
}

// Mock implementation for MAASMachine interface
type fakeMAASMachine struct {
	getFunc        func(ctx context.Context, systemID string) (*entity.Machine, error)
	listFunc       func(ctx context.Context) ([]entity.Machine, error)
	releaseFunc    func(ctx context.Context, systemID string, force bool) (*entity.Machine, error)
	powerOnFunc    func(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error)
	powerOffFunc   func(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error)
	commissionFunc func(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error)
}

// NexusService is defined in nexus.go, do not redeclare here.

func (m *fakeMAASMachine) Get(ctx context.Context, systemID string) (*entity.Machine, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, systemID)
	}
	return nil, nil
}
func (m *fakeMAASMachine) List(ctx context.Context) ([]entity.Machine, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return nil, nil
}
func (m *fakeMAASMachine) Release(ctx context.Context, systemID string, force bool) (*entity.Machine, error) {
	if m.releaseFunc != nil {
		return m.releaseFunc(ctx, systemID, force)
	}
	return nil, nil
}
func (m *fakeMAASMachine) PowerOn(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error) {
	if m.powerOnFunc != nil {
		return m.powerOnFunc(ctx, systemID, params)
	}
	return nil, nil
}
func (m *fakeMAASMachine) PowerOff(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error) {
	if m.powerOffFunc != nil {
		return m.powerOffFunc(ctx, systemID, params)
	}
	return nil, nil
}
func (m *fakeMAASMachine) Commission(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error) {
	if m.commissionFunc != nil {
		return m.commissionFunc(ctx, systemID, params)
	}
	return nil, nil
}

func TestMAASMachine_Get(t *testing.T) {
	mock := &fakeMAASMachine{
		getFunc: func(ctx context.Context, systemID string) (*entity.Machine, error) {
			if systemID == "valid" {
				return &entity.Machine{SystemID: "valid"}, nil
			}
			return nil, errors.New("not found")
		},
	}
	m, err := mock.Get(context.Background(), "valid")
	if err != nil || m == nil || m.SystemID != "valid" {
		t.Errorf("expected valid machine, got %v, err %v", m, err)
	}
	_, err = mock.Get(context.Background(), "invalid")
	if err == nil {
		t.Error("expected error for invalid systemID")
	}
}

func TestMAASMachine_List(t *testing.T) {
	mock := &fakeMAASMachine{
		listFunc: func(ctx context.Context) ([]entity.Machine, error) {
			return []entity.Machine{{SystemID: "m1"}, {SystemID: "m2"}}, nil
		},
	}
	list, err := mock.List(context.Background())
	if err != nil || len(list) != 2 {
		t.Errorf("expected 2 machines, got %v, err %v", list, err)
	}
}

func TestMAASMachine_Release(t *testing.T) {
	mock := &fakeMAASMachine{
		releaseFunc: func(ctx context.Context, systemID string, force bool) (*entity.Machine, error) {
			if force {
				return &entity.Machine{SystemID: systemID}, nil
			}
			return nil, errors.New("force required")
		},
	}
	m, err := mock.Release(context.Background(), "m1", true)
	if err != nil || m == nil || m.SystemID != "m1" {
		t.Errorf("expected released machine, got %v, err %v", m, err)
	}
	_, err = mock.Release(context.Background(), "m1", false)
	if err == nil {
		t.Error("expected error when force is false")
	}
}

func TestMAASMachine_PowerOnOff(t *testing.T) {
	mock := &fakeMAASMachine{
		powerOnFunc: func(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error) {
			return &entity.Machine{SystemID: systemID}, nil
		},
		powerOffFunc: func(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error) {
			return &entity.Machine{SystemID: systemID}, nil
		},
	}
	m, err := mock.PowerOn(context.Background(), "m1", nil)
	if err != nil || m == nil || m.SystemID != "m1" {
		t.Errorf("expected powered on machine, got %v, err %v", m, err)
	}
	m, err = mock.PowerOff(context.Background(), "m1", nil)
	if err != nil || m == nil || m.SystemID != "m1" {
		t.Errorf("expected powered off machine, got %v, err %v", m, err)
	}
}

func TestMAASMachine_Commission(t *testing.T) {
	mock := &fakeMAASMachine{
		commissionFunc: func(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error) {
			return &entity.Machine{SystemID: systemID}, nil
		},
	}
	m, err := mock.Commission(context.Background(), "m1", nil)
	if err != nil || m == nil || m.SystemID != "m1" {
		t.Errorf("expected commissioned machine, got %v, err %v", m, err)
	}
}

func TestNexusService_ListMachines(t *testing.T) {
	mock := &fakeMachineRepo{
		listFunc: func(ctx context.Context) ([]model.Machine, error) {
			return []model.Machine{
				{SystemID: "machine1", WorkloadAnnotations: map[string]string{"juju-model-uuid": "scopeUUID"}},
				{SystemID: "machine2", WorkloadAnnotations: map[string]string{"juju-model-uuid": "otherUUID"}},
			}, nil
		},
	}

	nexusService := &NexusService{machine: mock}
	machines, err := nexusService.ListMachines(context.Background(), "scopeUUID")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(machines) != 1 || machines[0].SystemID != "machine1" {
		t.Errorf("expected 1 machine with SystemID 'machine1', got %v", machines)
	}
}

// Add commissionFunc to fakeMachineRepo for testability
type fakeMachineRepo struct {
	listFunc     func(ctx context.Context) ([]model.Machine, error)
	releaseFunc  func(ctx context.Context, systemID string, force bool) (*entity.Machine, error)
	powerOnFunc  func(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error)
	powerOffFunc func(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error)
	getFunc      func(ctx context.Context, systemID string) (*entity.Machine, error)
}

// Implement List method for the repository interface.
func (f *fakeMachineRepo) List(ctx context.Context) ([]model.Machine, error) {
	if f.listFunc != nil {
		return f.listFunc(ctx)
	}
	return nil, nil
}

// Implement Get method for the MAASMachine interface.
func (f *fakeMachineRepo) Get(ctx context.Context, systemID string) (*entity.Machine, error) {
	if f.getFunc != nil {
		return f.getFunc(ctx, systemID)
	}
	return nil, nil
}

// Add dummy implementations for other MAASMachine interface methods if required.
func (f *fakeMachineRepo) Release(ctx context.Context, systemID string, force bool) (*entity.Machine, error) {
	if f.releaseFunc != nil {
		return f.releaseFunc(ctx, systemID, force)
	}
	return nil, nil
}
func (f *fakeMachineRepo) PowerOn(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error) {
	if f.powerOnFunc != nil {
		return f.powerOnFunc(ctx, systemID, params)
	}
	return nil, nil
}
func (f *fakeMachineRepo) PowerOff(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error) {
	if f.powerOffFunc != nil {
		return f.powerOffFunc(ctx, systemID, params)
	}
	return nil, nil
}
func (f *fakeMachineRepo) Commission(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error) {
	return nil, nil
}

func TestGetJujuModelUUID(t *testing.T) {
	tests := []struct {
		name    string
		ann     map[string]string
		want    string
		wantErr bool
	}{
		{"found", map[string]string{"juju-model-uuid": "uuid-123"}, "uuid-123", false},
		{"not found", map[string]string{}, "", true},
	}
	for _, tt := range tests {
		got, err := getJujuModelUUID(tt.ann)
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: expected error %v, got %v", tt.name, tt.wantErr, err)
		}
		if got != tt.want {
			t.Errorf("%s: expected %q, got %q", tt.name, tt.want, got)
		}
	}
}

func TestGetJujuMachineID(t *testing.T) {
	tests := []struct {
		name    string
		ann     map[string]string
		want    string
		wantErr bool
	}{
		{"found", map[string]string{"juju-machine-id": "juju-1"}, "1", false},
		{"not found", map[string]string{}, "", true},
		{"complex", map[string]string{"juju-machine-id": "juju-foo-bar-42"}, "42", false},
	}
	for _, tt := range tests {
		got, err := getJujuMachineID(tt.ann)
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: expected error %v, got %v", tt.name, tt.wantErr, err)
		}
		if got != tt.want {
			t.Errorf("%s: expected %q, got %q", tt.name, tt.want, got)
		}
	}
}

func TestBoolToInt(t *testing.T) {
	if boolToInt(true) != 1 {
		t.Error("expected 1 for true")
	}
	if boolToInt(false) != 0 {
		t.Error("expected 0 for false")
	}
}

func TestToPlacement(t *testing.T) {
	p := &model.MachinePlacement{LXD: true}
	pla := toPlacement(p, "foo")
	if pla.Scope != "lxd" {
		t.Errorf("expected lxd, got %v", pla.Scope)
	}
	p = &model.MachinePlacement{KVM: true}
	pla = toPlacement(p, "foo")
	if pla.Scope != "kvm" {
		t.Errorf("expected kvm, got %v", pla.Scope)
	}
	p = &model.MachinePlacement{Machine: true}
	pla = toPlacement(p, "bar")
	if pla.Scope != "#" || pla.Directive != "bar" {
		t.Errorf("expected # and bar, got %v %v", pla.Scope, pla.Directive)
	}
	p = &model.MachinePlacement{}
	if toPlacement(p, "baz") != nil {
		t.Error("expected nil for empty placement")
	}
}

func TestToConstraint(t *testing.T) {
	c := &model.MachineConstraint{
		Architecture: "amd64",
		CPUCores:     4,
		MemoryMB:     8192,
		Tags:         []string{"tag1", "tag2"},
	}
	con := toConstraint(c)
	if con.Arch == nil || *con.Arch != "amd64" {
		t.Error("expected arch amd64")
	}
	if con.CpuCores == nil || *con.CpuCores != 4 {
		t.Error("expected cpu 4")
	}
	if con.Mem == nil || *con.Mem != 8192 {
		t.Error("expected mem 8192")
	}
	if con.Tags == nil || len(*con.Tags) != 2 {
		t.Error("expected 2 tags")
	}
}

func TestNexusService_DeleteMachine_JujuNotFound(t *testing.T) {
	mock := &fakeMachineRepo{
		releaseFunc: func(ctx context.Context, systemID string, force bool) (*entity.Machine, error) {
			return &entity.Machine{SystemID: systemID, WorkloadAnnotations: map[string]string{}}, nil
		},
	}
	nexusService := &NexusService{machine: mock}
	err := nexusService.DeleteMachine(context.Background(), "m1", true)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestNexusService_PowerOnOffMachine(t *testing.T) {
	mock := &fakeMachineRepo{
		powerOnFunc: func(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error) {
			return &entity.Machine{SystemID: systemID}, nil
		},
		powerOffFunc: func(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error) {
			return &entity.Machine{SystemID: systemID}, nil
		},
	}
	nexusService := &NexusService{machine: mock}
	m, err := nexusService.PowerOnMachine(context.Background(), "m1", "test")
	if err != nil || m == nil || m.SystemID != "m1" {
		t.Errorf("expected powered on machine, got %v, err %v", m, err)
	}
	m, err = nexusService.PowerOffMachine(context.Background(), "m1", "test")
	if err != nil || m == nil || m.SystemID != "m1" {
		t.Errorf("expected powered off machine, got %v, err %v", m, err)
	}
}
func TestNexusService_GetMachine(t *testing.T) {
	mock := &fakeMachineRepo{
		getFunc: func(ctx context.Context, systemID string) (*entity.Machine, error) {
			if systemID == "m1" {
				return &entity.Machine{SystemID: "m1"}, nil
			}
			return nil, errors.New("not found")
		},
	}
	nexusService := &NexusService{machine: mock}
	m, err := nexusService.GetMachine(context.Background(), "m1")
	if err != nil || m == nil || m.SystemID != "m1" {
		t.Errorf("expected machine m1, got %v, err %v", m, err)
	}
	_, err = nexusService.GetMachine(context.Background(), "notfound")
	if err == nil {
		t.Error("expected error for notfound")
	}
}

func TestNexusService_waitForMachineReady_getError(t *testing.T) {
	nexusService := &NexusService{
		machine: &fakeMachineRepo{
			getFunc: func(ctx context.Context, systemID string) (*entity.Machine, error) {
				return nil, errors.New("get error")
			},
		},
	}
	err := nexusService.waitForMachineReady(context.Background(), "m1")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetJujuMachineID_invalidFormat(t *testing.T) {
	ann := map[string]string{"juju-machine-id": ""}
	_, err := getJujuMachineID(ann)
	if err != nil {
		t.Errorf("expected no error for empty string, got %v", err)
	}
}

// Define a minimal MachineAddResult type for testing if not present in model package
type MachineAddResult struct {
	Machine string
}

func TestNexusService_DeleteMachine_JujuError(t *testing.T) {
	mock := &fakeMachineRepo{
		releaseFunc: func(ctx context.Context, systemID string, force bool) (*entity.Machine, error) {
			return &entity.Machine{
				SystemID: systemID,
				WorkloadAnnotations: map[string]string{
					"juju-model-uuid": "uuid",
					"juju-machine-id": "juju-1",
				},
			}, nil
		},
	}
	mockMachineManager := &fakeMachineManager{
		destroyMachinesFunc: func(ctx context.Context, uuid string, force bool, machines ...string) ([]params.DestroyMachineResult, error) {
			return []params.DestroyMachineResult{{Error: nil}}, nil
		},
	}
	nexusService := &NexusService{
		machine:        mock,
		machineManager: mockMachineManager,
	}
	err := nexusService.DeleteMachine(context.Background(), "m1", true)
	if err == nil {
		t.Error("expected error from juju destroy, got nil")
	}
}

type fakeMachineManager struct {
	addMachinesFunc     func(ctx context.Context, uuid string, params []params.AddMachineParams) ([]params.AddMachinesResult, error)
	destroyMachinesFunc func(ctx context.Context, uuid string, force bool, machines ...string) ([]params.DestroyMachineResult, error)
}

func (f *fakeMachineManager) AddMachines(ctx context.Context, uuid string, params []params.AddMachineParams) ([]params.AddMachinesResult, error) {
	if f.addMachinesFunc != nil {
		return f.addMachinesFunc(ctx, uuid, params)
	}
	return nil, nil
}
func (f *fakeMachineManager) DestroyMachines(ctx context.Context, uuid string, force bool, machines ...string) ([]params.DestroyMachineResult, error) {
	if f.destroyMachinesFunc != nil {
		return f.destroyMachinesFunc(ctx, uuid, force, machines...)
	}
	return nil, nil
}
