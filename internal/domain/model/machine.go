package model

import (
	"github.com/canonical/gomaasclient/entity"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
)

type (
	Machine                 = entity.Machine
	MachineCommissionParams = entity.MachineCommissionParams
	MachinePowerOnParams    = entity.MachinePowerOnParams
	MachinePowerOffParams   = entity.MachinePowerOffParams
	NUMANode                = entity.NUMANode
	BlockDevice             = entity.BlockDevice
	NetworkInterface        = entity.NetworkInterface
)

type (
	MachineAddParams = params.AddMachineParams
	Placement        = instance.Placement
	Constraint       = constraints.Value
)

type MachinePlacement struct {
	LXD       bool
	KVM       bool
	Machine   bool
	MachineID string
}

type MachineConstraint struct {
	Architecture string
	CPUCores     uint64
	MemoryMB     uint64
	Tags         []string
}

type MachineFactor struct {
	*MachinePlacement
	*MachineConstraint
}
