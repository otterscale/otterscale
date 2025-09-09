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
	"k8s.io/client-go/rest"

	"github.com/juju/juju/api/client/action"
)

type kubeMockFacilityRepo struct{}

func (m *kubeMockFacilityRepo) GetLeader(ctx context.Context, uuid, name string) (string, error) {
	return "leader", nil
}

// Consume is a mock implementation to satisfy the FacilityRepo interface.
func (m *kubeMockFacilityRepo) Consume(ctx context.Context, uuid string, args *crossmodel.ConsumeApplicationArgs) error {
	return nil
}

func (m *kubeMockFacilityRepo) GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error) {
	return &application.UnitInfo{PublicAddress: "10.0.0.1"}, nil
}

// Create is a mock implementation to satisfy the FacilityRepo interface.
func (m *kubeMockFacilityRepo) Create(
	ctx context.Context,
	uuid, app, series, channel, config string,
	numUnits, expose int,
	base *base.Base,
	placements []instance.Placement,
	constraints *constraints.Value,
	trusted bool,
) (*application.DeployInfo, error) {
	return &application.DeployInfo{}, nil
}

// AddUnits is a mock implementation to satisfy the FacilityRepo interface.
func (m *kubeMockFacilityRepo) AddUnits(ctx context.Context, uuid, app string, count int, placements []instance.Placement) ([]string, error) {
	return []string{}, nil
}

// CreateRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *kubeMockFacilityRepo) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	return &params.AddRelationResults{}, nil
}

// Delete is a mock implementation to satisfy the FacilityRepo interface.
func (m *kubeMockFacilityRepo) Delete(ctx context.Context, uuid, app string, force, destroyStorage bool) error {
	return nil
}

// DeleteRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *kubeMockFacilityRepo) DeleteRelation(ctx context.Context, uuid string, relationID int) error {
	return nil
}

// Expose is a mock implementation to satisfy the FacilityRepo interface.
func (m *kubeMockFacilityRepo) Expose(ctx context.Context, uuid, app string, endpoints map[string]params.ExposedEndpoint) error {
	return nil
}

// GetConfig is a mock implementation to satisfy the FacilityRepo interface.
func (m *kubeMockFacilityRepo) GetConfig(ctx context.Context, uuid, app string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// ResolveUnitErrors is a mock implementation to satisfy the FacilityRepo interface.
func (m *kubeMockFacilityRepo) ResolveUnitErrors(ctx context.Context, uuid string, units []string) error {
	return nil
}

// Update is a mock implementation to satisfy the FacilityRepo interface.
func (m *kubeMockFacilityRepo) Update(ctx context.Context, uuid, app string, config string) error {
	return nil
}

type kubeMockActionRepo struct{}

func (m *kubeMockActionRepo) RunCommand(ctx context.Context, uuid, leader, command string) (string, error) {
	return "action-id", nil
}

func (m *kubeMockActionRepo) RunAction(ctx context.Context, uuid, leader, action string, params map[string]any) (string, error) {
	return "action-id", nil
}

// Implement the List method to satisfy the ActionRepo interface
func (m *kubeMockActionRepo) List(ctx context.Context, uuid, leader string) (map[string]ActionSpec, error) {
	return map[string]ActionSpec{}, nil
}

// Remove the custom ActionResult type

func (m *kubeMockActionRepo) GetResult(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
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
preferences: {}
users:
- name: test
  user:
    token: test
`,
		},
	}, nil
}

func TestKubeConfig(t *testing.T) {
	fac := &kubeMockFacilityRepo{}
	act := &kubeMockActionRepo{}
	cfg, err := kubeConfig(context.Background(), fac, act, "uuid", "k8s")
	assert.NoError(t, err)
	assert.IsType(t, &rest.Config{}, cfg)
}
