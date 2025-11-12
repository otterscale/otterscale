package network

import (
	"context"
)

type VLAN struct {
	ID          int
	VID         int
	Name        string
	MTU         int
	Description string
	DHCPOn      bool
}

type VLANRepo interface {
	Update(ctx context.Context, fabricID, vid int, name string, mtu int, description string, dhcpOn bool) (*VLAN, error)
}

func (uc *NetworkUseCase) UpdateVLAN(ctx context.Context, fabricID, vid int, name string, mtu int, description string, dhcpOn bool) (*VLAN, error) {
	return uc.vlan.Update(ctx, fabricID, vid, name, mtu, description, dhcpOn)
}
