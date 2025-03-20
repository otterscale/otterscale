package model

import "github.com/canonical/gomaasclient/entity"

type NetworkSetting struct {
	*entity.VLAN
	*entity.Subnet
	*entity.IPRange
}

type Network struct {
	*entity.Fabric
	Settings []*NetworkSetting
}
