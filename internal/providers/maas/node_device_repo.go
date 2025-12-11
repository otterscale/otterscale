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

var (
	// allowedVendors is a set of GPU vendor IDs that are supported.
	allowedVendors = map[string]bool{
		"10de": true, // NVIDIA
		"1002": true, // AMD
		"8086": true, // Intel
	}
)

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

	filteredGPUs := make([]entity.NodeDevice, 0, len(nodeGPUs))
	for _, gpu := range nodeGPUs {
		if allowedVendors[gpu.VendorID] {
			filteredGPUs = append(filteredGPUs, gpu)
		}
	}

	return filteredGPUs, nil
}
