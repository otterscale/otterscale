package model

import (
	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/subnet"
)

type NetworkSubnet struct {
	*entity.Subnet
	IPAddresses      []subnet.IPAddress
	ReservedIPRanges []entity.IPRange
	Statistics       *subnet.Statistics
}

type NetworkSetting struct {
	*entity.VLAN
	Subnet *NetworkSubnet
}

type Network struct {
	*entity.Fabric
	Settings []*NetworkSetting
}
