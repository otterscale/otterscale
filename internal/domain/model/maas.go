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

type Machine = entity.Machine

type MachinesParams = entity.MachinesParams

type MachinePowerOnParams = entity.MachinePowerOnParams

type MachinePowerOffParams = entity.MachinePowerOffParams

type MachineCommissionParams = entity.MachineCommissionParams

type NetworkSetting struct {
	*VLAN
	*Subnet
	*IPRange
}

type Network struct {
	*Fabric
	Settings []*NetworkSetting
}
