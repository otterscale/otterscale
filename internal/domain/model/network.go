package model

import (
	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/subnet"
)

type (
	Fabric            = entity.Fabric
	FabricParams      = entity.FabricParams
	VLAN              = entity.VLAN
	VLANParams        = entity.VLANParams
	Subnet            = entity.Subnet
	SubnetParams      = entity.SubnetParams
	IPRange           = entity.IPRange
	IPRangeParams     = entity.IPRangeParams
	IPAddress         = subnet.IPAddress
	NetworkStatistics = subnet.Statistics
)

type Network struct {
	*entity.Fabric
	*entity.VLAN
	Subnet *NetworkSubnet
}

type NetworkSubnet struct {
	*Subnet
	*subnet.Statistics
	IPAddresses []subnet.IPAddress
	IPRanges    []entity.IPRange
}
