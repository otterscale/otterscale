package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core"
)

type nodeDevice struct {
	maas *MAAS
}

func NewNodeDevice(maas *MAAS) core.NodeDeviceRepo {
	return &nodeDevice{
		maas: maas,
	}
}

var _ core.NodeDeviceRepo = (*nodeDevice)(nil)

func (r *nodeDevice) List(_ context.Context, systemID, hardwareType string) ([]core.NodeDevice, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	params := &entity.NodeDeviceParams{
		HardwareType: hardwareType,
	}
	return client.NodeDevices.Get(systemID, params)
}
