package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/machine"
)

// Note: MAAS API do not support context.
type nodeDeviceRepo struct {
	maas *MAAS
}

func NewNodeDeviceRepo(maas *MAAS) machine.NodeDeviceRepo {
	return &nodeDeviceRepo{
		maas: maas,
	}
}

var _ machine.NodeDeviceRepo = (*nodeDeviceRepo)(nil)

func (r *nodeDeviceRepo) ListGPUs(_ context.Context, machineID string) ([]machine.GPU, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.NodeDeviceParams{
		HardwareType: "gpu",
	}

	return client.NodeDevices.Get(machineID, params)
}
