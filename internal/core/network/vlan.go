package network

import (
	"context"

	"github.com/canonical/gomaasclient/entity"
)

// VLAN represents a MAAS VLAN resource.
type VLAN = entity.VLAN

type VLANRepo interface {
	Update(ctx context.Context, fabricID, vid int, name string, mtu int, description string, dhcpOn bool) (*VLAN, error)
}

func (uc *NetworkUseCase) UpdateVLAN(ctx context.Context, fabricID, vid int, name string, mtu int, description string, dhcpOn bool) (*VLAN, error) {
	return uc.vlan.Update(ctx, fabricID, vid, name, mtu, description, dhcpOn)
}
