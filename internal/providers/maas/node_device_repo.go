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
	// Pre-defined allowed GPU vendors with their PCI vendor IDs
	allowedVendors := map[string]bool{
		"10de": true, // NVIDIA
		"1002": true, // AMD
		"8086": true, // Intel
	}

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

	// Pre-allocate slice with estimated capacity to reduce memory allocations
	filteredGPUs := make([]entity.NodeDevice, 0, len(nodeGPUs))
	for _, gpu := range nodeGPUs {
		if allowedVendors[gpu.VendorID] {
			filteredGPUs = append(filteredGPUs, gpu)
		}
	}

	return filteredGPUs, nil
}
