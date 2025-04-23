package model

import (
	"github.com/canonical/gomaasclient/entity"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/rpc/params"
)

const JobHostUnits = model.JobHostUnits

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
	MachineStatus    = params.MachineStatus
	MachineAddParams = params.AddMachineParams
	Placement        = instance.Placement
	Constraint       = constraints.Value
	MachineJob       = model.MachineJob
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
