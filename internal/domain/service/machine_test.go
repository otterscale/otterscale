package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/canonical/gomaasclient/entity/node"
	"github.com/juju/juju/rpc/params"
	"github.com/openhdc/otterscale/internal/domain/model"
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestNexusService_AddMachines_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachineManager := mocks.NewMockJujuMachine(ctrl)
	mockClient := mocks.NewMockJujuClient(ctrl) // 建立 JujuClient 的 mock

	// 初始化 NexusService，並傳入 mock 物件
	mockService := &NexusService{
		machineManager: mockMachineManager,
		client:         mockClient, // 確保 client 被初始化
	}

	// 設置 mock 行為
	mockClient.EXPECT().Status(gomock.Any(), gomock.Any(), gomock.Any()).Return(&params.FullStatus{}, nil)

	mockMachineManager.EXPECT().
		AddMachines(gomock.Any(), "uuid", gomock.Any()).
		Return([]params.AddMachinesResult{{Machine: "juju-0"}}, nil)

	factors := []model.MachineFactor{
		{
			MachinePlacement: &model.MachinePlacement{Machine: true, MachineID: "mid"},
			MachineConstraint: &model.MachineConstraint{
				Architecture: "amd64",
				CPUCores:     2,
				MemoryMB:     4096,
			},
		},
	}

	machines, err := mockService.AddMachines(context.Background(), "uuid", factors)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(machines) != 1 || machines[0] != "juju-0" {
		t.Errorf("expected [juju-0], got %v", machines)
	}
}

func TestNexusService_AddMachines_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachineManager := mocks.NewMockJujuMachine(ctrl)
	mockClient := mocks.NewMockJujuClient(ctrl) // 建立 JujuClient 的 mock

	// 初始化 NexusService，並傳入 mock 物件
	mockService := &NexusService{
		machineManager: mockMachineManager,
		client:         mockClient, // 確保 client 被初始化
	}

	// 設置 mock 行為
	mockClient.EXPECT().Status(gomock.Any(), gomock.Any(), gomock.Any()).Return(&params.FullStatus{}, nil)

	mockMachineManager.EXPECT().
		AddMachines(gomock.Any(), "uuid", gomock.Any()).
		Return(nil, errors.New("error"))

	factors := []model.MachineFactor{
		{
			MachinePlacement: &model.MachinePlacement{Machine: true, MachineID: "mid"},
			MachineConstraint: &model.MachineConstraint{
				Architecture: "amd64",
				CPUCores:     2,
				MemoryMB:     4096,
			},
		},
	}

	machines, err := mockService.AddMachines(context.Background(), "uuid", factors)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
	if machines != nil {
		t.Errorf("expected nil machines, got %v", machines)
	}
}

func TestNexusService_DeleteMachine_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockMachineManager := mocks.NewMockJujuMachine(ctrl)

	mockService := &NexusService{
		machine:        mockMachine,
		machineManager: mockMachineManager,
	}

	id := "machine-id"
	force := true

	// 模擬 MAASMachine 的 Release 方法
	mockMachine.EXPECT().Release(gomock.Any(), id, force).Return(&model.Machine{
		WorkloadAnnotations: map[string]string{
			"juju-model-uuid": "model-uuid",
			"juju-machine-id": "juju-0",
		},
	}, nil)

	// 模擬 JujuMachine 的 DestroyMachines 方法
	mockMachineManager.EXPECT().DestroyMachines(gomock.Any(), "model-uuid", force, "0").Return([]params.DestroyMachineResult{}, nil)

	err := mockService.DeleteMachine(context.Background(), id, force)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNexusService_DeleteMachine_ReleaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)

	mockService := &NexusService{
		machine: mockMachine,
	}

	id := "machine-id"
	force := true

	// 模擬 MAASMachine 的 Release 方法返回錯誤
	mockMachine.EXPECT().Release(gomock.Any(), id, force).Return(nil, errors.New("release error"))

	err := mockService.DeleteMachine(context.Background(), id, force)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestNexusService_DeleteMachine_DestroyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockMachineManager := mocks.NewMockJujuMachine(ctrl)

	mockService := &NexusService{
		machine:        mockMachine,
		machineManager: mockMachineManager,
	}

	id := "machine-id"
	force := true

	// 模擬 MAASMachine 的 Release 方法
	mockMachine.EXPECT().Release(gomock.Any(), id, force).Return(&model.Machine{
		WorkloadAnnotations: map[string]string{
			"juju-model-uuid": "model-uuid",
			"juju-machine-id": "juju-0",
		},
	}, nil)

	// 模擬 JujuMachine 的 DestroyMachines 方法返回錯誤
	mockMachineManager.EXPECT().DestroyMachines(gomock.Any(), "model-uuid", force, "0").Return(nil, errors.New("destroy error"))

	err := mockService.DeleteMachine(context.Background(), id, force)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestCreateMachine_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockMachineManager := mocks.NewMockJujuMachine(ctrl)
	mockTag := mocks.NewMockMAASTag(ctrl)
	mockServer := mocks.NewMockMAASServer(ctrl) // Mock the server dependency

	mockService := &NexusService{
		machine:        mockMachine,
		machineManager: mockMachineManager,
		tag:            mockTag,    // 確保 client 被初始化
		server:         mockServer, // Add the server mock
	}

	id := "machine-1"
	uuid := "uuid-1"
	tags := []string{"tag1", "tag2"}

	mockTag.EXPECT().AddMachines(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
	mockMachine.EXPECT().Commission(gomock.Any(), id, gomock.Any()).Return(&model.Machine{}, nil).Times(1)
	mockMachine.EXPECT().Get(gomock.Any(), id).Return(&model.Machine{Status: node.StatusReady}, nil).AnyTimes()
	mockMachineManager.EXPECT().AddMachines(gomock.Any(), uuid, gomock.Any()).Return([]params.AddMachinesResult{{}}, nil).Times(1)

	// Mock the server's Get method to avoid nil pointer dereference
	mockServer.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte(`"focal"`), nil).AnyTimes()

	machine, err := mockService.CreateMachine(context.Background(), id, true, false, false, false, uuid, tags)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if machine == nil {
		t.Errorf("expected machine to be created, got nil")
	}
}

func TestCreateMachine_AddTagsFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTag := mocks.NewMockMAASTag(ctrl)

	// 初始化 NexusService，並傳入 mock 物件
	mockService := &NexusService{
		tag: mockTag, // 確保 client 被初始化
	}

	id := "machine-1"
	tags := []string{"tag1", "tag2"}

	mockTag.EXPECT().AddMachines(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to add tags")).Times(2)

	machine, err := mockService.CreateMachine(context.Background(), id, true, false, false, false, "uuid-1", tags)
	if err == nil {
		t.Errorf("expected error, got none")
	}
	if machine != nil {
		t.Errorf("expected no machine to be created, got %v", machine)
	}
}

func TestCreateMachine_CommissionFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockTag := mocks.NewMockMAASTag(ctrl)

	// 初始化 NexusService，並傳入 mock 物件
	mockService := &NexusService{
		machine: mockMachine,
		tag:     mockTag, // 確保 client 被初始化
	}

	id := "machine-1"
	tags := []string{"tag1", "tag2"}

	mockTag.EXPECT().AddMachines(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
	mockMachine.EXPECT().Commission(gomock.Any(), id, gomock.Any()).Return(nil, errors.New("commission failed")).Times(1)

	machine, err := mockService.CreateMachine(context.Background(), id, true, false, false, false, "uuid-1", tags)
	if err == nil {
		t.Errorf("expected error, got none")
	}
	if machine != nil {
		t.Errorf("expected no machine to be created, got %v", machine)
	}
}

func TestCreateMachine_TimeoutWaitingForMachineReady(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockTag := mocks.NewMockMAASTag(ctrl)

	// 初始化 NexusService，並傳入 mock 物件
	mockService := &NexusService{
		machine: mockMachine,
		tag:     mockTag, // 確保 client 被初始化
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	id := "machine-1"
	tags := []string{"tag1", "tag2"}

	mockTag.EXPECT().AddMachines(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
	mockMachine.EXPECT().Commission(gomock.Any(), id, gomock.Any()).Return(&model.Machine{}, nil).Times(1)
	mockMachine.EXPECT().Get(gomock.Any(), id).Return(&model.Machine{Status: node.StatusBroken}, nil).AnyTimes()

	machine, err := mockService.CreateMachine(ctx, id, true, false, false, false, "uuid-1", tags)
	if err == nil {
		t.Errorf("expected timeout error, got none")
	}
	if machine != nil {
		t.Errorf("expected no machine to be created, got %v", machine)
	}
}

func TestListMachines_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)

	// 初始化 NexusService，並傳入 mock 物件
	mockService := &NexusService{
		machine: mockMachine,
	}

	scopeUUID := "test-scope"

	// 模擬返回的機器列表
	machines := []model.Machine{
		{WorkloadAnnotations: map[string]string{"juju-model-uuid": "test-scope-1"}},
	}

	mockMachine.EXPECT().List(gomock.Any()).Return(machines, nil).Times(1)

	result, err := mockService.ListMachines(context.Background(), scopeUUID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 machine, got %d", len(result))
	}
	if result[0].WorkloadAnnotations["juju-model-uuid"] != "test-scope-1" {
		t.Errorf("expected juju-model-uuid to be 'test-scope', got %s", result[0].WorkloadAnnotations["juju-model-uuid"])
	}
}

func TestListMachines_ListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)

	// 初始化 NexusService，並傳入 mock 物件
	mockService := &NexusService{
		machine: mockMachine,
	}

	ctx := context.Background()
	scopeUUID := "test-scope"

	mockMachine.EXPECT().List(gomock.Any()).Return(nil, errors.New("list error")).Times(1)

	result, err := mockService.ListMachines(ctx, scopeUUID)
	if err == nil {
		t.Fatal("expected error, got none")
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestGetMachine_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockService := &NexusService{
		machine: mockMachine,
	}

	ctx := context.Background()
	id := "test-machine-id"
	expectedMachine := &model.Machine{SystemID: id}

	mockMachine.EXPECT().Get(gomock.Any(), id).Return(expectedMachine, nil).Times(1)

	machine, err := mockService.GetMachine(ctx, id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if machine == nil {
		t.Fatal("expected machine, got nil")
	}
	if machine.SystemID != id {
		t.Errorf("expected machine id %s, got %s", id, machine.SystemID)
	}
}

func TestGetMachine_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockService := &NexusService{
		machine: mockMachine,
	}

	ctx := context.Background()
	id := "non-existent-machine-id"

	mockMachine.EXPECT().Get(gomock.Any(), id).Return(nil, errors.New("machine not found")).Times(1)

	machine, err := mockService.GetMachine(ctx, id)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err == nil || err.Error() != "machine not found" {
		t.Errorf("expected error 'machine not found', got %v", err)
	}
	if machine != nil {
		t.Errorf("expected nil machine, got %v", machine)
	}
}

func TestGetMachine_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockService := &NexusService{
		machine: mockMachine,
	}

	ctx := context.Background()
	id := "error-machine-id"

	mockMachine.EXPECT().Get(gomock.Any(), id).Return(nil, errors.New("some error")).Times(1)

	machine, err := mockService.GetMachine(ctx, id)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if machine != nil {
		t.Errorf("expected nil machine, got %v", machine)
	}
}

func TestPowerOnMachine_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockService := &NexusService{
		machine: mockMachine,
	}

	ctx := context.Background()
	id := "test-machine-id"
	expectedMachine := &model.Machine{SystemID: id, PowerState: "on"}
	comment := "test comment"
	params := &model.MachinePowerOnParams{Comment: comment} // Use model.MachinePowerOnParams

	mockMachine.EXPECT().PowerOn(gomock.Any(), id, params).Return(expectedMachine, nil).Times(1)

	machine, err := mockService.PowerOnMachine(ctx, id, comment)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if machine == nil {
		t.Fatal("expected machine, got nil")
	}
	if machine.SystemID != id {
		t.Errorf("expected machine id %s, got %s", id, machine.SystemID)
	}
	if machine.PowerState != "on" { // Check the power state
		t.Errorf("expected power state 'on', got '%s'", machine.PowerState)
	}
}

func TestPowerOffMachine_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockService := &NexusService{
		machine: mockMachine,
	}

	ctx := context.Background()
	id := "test-machine-id"
	expectedMachine := &model.Machine{SystemID: id, PowerState: "off"}
	comment := "test comment"
	params := &model.MachinePowerOffParams{Comment: comment} // Use model.MachinePowerOffParams

	mockMachine.EXPECT().PowerOff(gomock.Any(), id, params).Return(expectedMachine, nil).Times(1)

	machine, err := mockService.PowerOffMachine(ctx, id, comment)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if machine == nil {
		t.Fatal("expected machine, got nil")
	}
	if machine.SystemID != id {
		t.Errorf("expected machine id %s, got %s", id, machine.SystemID)
	}
	if machine.PowerState != "off" { // Check the power state
		t.Errorf("expected power state 'off', got '%s'", machine.PowerState)
	}
}
