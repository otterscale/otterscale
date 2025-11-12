package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/machine/metal"
)

type nodeDeviceRepo struct {
	maas *MAAS
}

func NewNodeDeviceRepo(maas *MAAS) metal.NodeDeviceRepo {
	return &nodeDeviceRepo{
		maas: maas,
	}
}

var _ metal.NodeDeviceRepo = (*nodeDeviceRepo)(nil)

func (r *nodeDeviceRepo) ListGPUs(_ context.Context, machineID string) ([]metal.GPU, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.NodeDeviceParams{
		HardwareType: "gpu",
	}

	nodeDevices, err := client.NodeDevices.Get(machineID, params)
	if err != nil {
		return nil, err
	}

	return r.toGPUs(nodeDevices), nil
}

func (r *nodeDeviceRepo) toGPUs(nds []entity.NodeDevice) []metal.GPU {
	gpus := make([]metal.GPU, 0, len(nds))

	for _, nd := range nds {
		gpu := metal.GPU{
			VendorID:    nd.VendorID,
			ProductID:   nd.ProductID,
			VendorName:  nd.VendorName,
			ProductName: nd.ProductName,
		}
		gpus = append(gpus, gpu)
	}

	return gpus
}
