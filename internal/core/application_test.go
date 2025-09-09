// application_test.go
package core

import (
	"context"
	"testing"

	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockFacilityRepo struct {
	mock.Mock
}

func (m *mockFacilityRepo) Create(ctx context.Context, uuid, name string, configYAML string, charmName, channel string, revision, number int, base *base.Base, placements []instance.Placement, constraint *constraints.Value, trust bool) (*application.DeployInfo, error) {
	args := m.Called(ctx, uuid, name, configYAML, charmName, channel, revision, number, base, placements, constraint, trust)
	return args.Get(0).(*application.DeployInfo), args.Error(1)
}

func (m *mockFacilityRepo) Update(ctx context.Context, uuid, name string, configYAML string) error {
	args := m.Called(ctx, uuid, name, configYAML)
	return args.Error(0)
}

func (m *mockFacilityRepo) Delete(ctx context.Context, uuid, name string, destroyStorage, force bool) error {
	args := m.Called(ctx, uuid, name, destroyStorage, force)
	return args.Error(0)
}

func (m *mockFacilityRepo) Expose(ctx context.Context, uuid, name string, endpoints map[string]params.ExposedEndpoint) error {
	args := m.Called(ctx, uuid, name, endpoints)
	return args.Error(0)
}

func (m *mockFacilityRepo) AddUnits(ctx context.Context, uuid, name string, number int, placements []instance.Placement) ([]string, error) {
	args := m.Called(ctx, uuid, name, number, placements)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockFacilityRepo) ResolveUnitErrors(ctx context.Context, uuid string, units []string) error {
	args := m.Called(ctx, uuid, units)
	return args.Error(0)
}

func (m *mockFacilityRepo) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	args := m.Called(ctx, uuid, endpoints)
	return args.Get(0).(*params.AddRelationResults), args.Error(1)
}

func (m *mockFacilityRepo) DeleteRelation(ctx context.Context, uuid string, id int) error {
	args := m.Called(ctx, uuid, id)
	return args.Error(0)
}

func (m *mockFacilityRepo) GetConfig(ctx context.Context, uuid string, name string) (map[string]any, error) {
	args := m.Called(ctx, uuid, name)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *mockFacilityRepo) GetLeader(ctx context.Context, uuid, name string) (string, error) {
	args := m.Called(ctx, uuid, name)
	return args.String(0), args.Error(1)
}

func (m *mockFacilityRepo) GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error) {
	args := m.Called(ctx, uuid, name)
	return args.Get(0).(*application.UnitInfo), args.Error(1)
}

func (m *mockFacilityRepo) Consume(ctx context.Context, uuid string, args *crossmodel.ConsumeApplicationArgs) error {
	argsCalled := m.Called(ctx, uuid, args)
	return argsCalled.Error(0)
}

func TestApplicationUseCase_GetPublicAddress(t *testing.T) {
	ctx := context.Background()
	uc := &ApplicationUseCase{
		facility: &mockFacilityRepo{},
		// Mock other repos if needed
	}

	mockFacility := uc.facility.(*mockFacilityRepo)
	mockFacility.On("GetLeader", ctx, "uuid", "facility").Return("leader", nil)
	mockFacility.On("GetUnitInfo", ctx, "uuid", "leader").Return(&application.UnitInfo{PublicAddress: "1.2.3.4"}, nil)

	addr, err := uc.GetPublicAddress(ctx, "uuid", "facility")
	assert.NoError(t, err)
	assert.Equal(t, "1.2.3.4", addr)

	mockFacility.AssertExpectations(t)
}
