package core

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/juju/juju/core/base"
)

var kubeConfigMap sync.Map

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func defaultBase(ctx context.Context, server ServerRepo) (base.Base, error) {
	series, err := server.Get(ctx, "default_distro_series")
	if err != nil {
		return base.Base{}, err
	}
	return base.GetBaseFromSeries(series)
}

func getJujuModelUUID(m map[string]string) (string, error) {
	v, ok := m["juju-model-uuid"]
	if !ok {
		return "", errors.New("juju model uuid not found")
	}
	return v, nil
}

func getJujuMachineID(m map[string]string) (string, error) {
	v, ok := m["juju-machine-id"]
	if !ok {
		return "", errors.New("juju machine id not found")
	}
	token := strings.Split(v, "-")
	return token[len(token)-1], nil
}

func toPlacement(p *MachinePlacement, directive string) *Placement {
	placement := &Placement{}
	if p.LXD {
		placement.Scope = "lxd"
	} else if p.KVM {
		placement.Scope = "kvm"
	} else if p.Machine {
		placement.Scope = "#"
		placement.Directive = directive
	} else {
		return nil
	}
	return placement
}

func toConstraint(c *MachineConstraint) Constraint {
	constraint := Constraint{}
	if c.Architecture != "" {
		constraint.Arch = &c.Architecture
	}
	if c.CPUCores > 0 {
		constraint.CpuCores = &c.CPUCores
	}
	if c.MemoryMB > 0 {
		constraint.Mem = &c.MemoryMB
	}
	if len(c.Tags) > 0 {
		constraint.Tags = &c.Tags
	}
	return constraint
}
