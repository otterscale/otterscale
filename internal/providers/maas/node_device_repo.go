package maas

import (
	"context"
	"slices"
	"strings"

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

	nodeGPUs, err := client.NodeDevices.Get(machineID, params)
	if err != nil {
		return nil, err
	}

	return r.filterNvidiaGPUs(nodeGPUs), nil
}

func (r *nodeDeviceRepo) filterNvidiaGPUs(gpus []machine.GPU) []machine.GPU {
	return slices.DeleteFunc(gpus, func(gpu machine.GPU) bool {
		return !strings.EqualFold(gpu.VendorID, "10de")
	})
}
