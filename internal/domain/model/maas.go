package model

import "github.com/canonical/gomaasclient/entity"

type PackageRepository = entity.PackageRepository

type PackageRepositoryParams = entity.PackageRepositoryParams

type Fabric = entity.Fabric

type FabricParams = entity.FabricParams

type VLAN = entity.VLAN

type VLANParams = entity.VLANParams

type Subnet = entity.Subnet

type SubnetParams = entity.SubnetParams

type IPRange = entity.IPRange

type IPRangeParams = entity.IPRangeParams

type NetworkSetting struct {
	*VLAN
	*Subnet
	*IPRange
}

type Network struct {
	ID          int
	Name        string
	ClassType   string
	ResourceURI string
	Settings    []*NetworkSetting
}
