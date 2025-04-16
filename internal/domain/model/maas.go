package model

import (
	"github.com/canonical/gomaasclient/entity"
)

type NetworkSetting struct {
	*entity.VLAN
	Subnet *NetworkSubnet
}

type Networkx struct {
	*entity.Fabric
	Settings []*NetworkSetting
}
