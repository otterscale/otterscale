package core

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
)

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func flatten[T any](data [][]T) []T {
	totalLen := 0
	for _, innerSlice := range data {
		totalLen += len(innerSlice)
	}

	result := make([]T, 0, totalLen)
	for _, innerSlice := range data {
		result = append(result, innerSlice...)
	}
	return result
}

func defaultBase(ctx context.Context, server ServerRepo) (base.Base, error) {
	series, err := server.Get(ctx, "default_distro_series")
	if err != nil {
		return base.Base{}, err
	}
	return base.GetBaseFromSeries(series)
}

func getJujuModelName(m map[string]string) (string, error) {
	v, ok := m["juju-machine-id"]
	if !ok {
		return "", errors.New("juju machine id not found")
	}
	token := strings.Split(v, "-")
	return token[0], nil
}

func getJujuMachineID(m map[string]string) (string, error) {
	v, ok := m["juju-machine-id"]
	if !ok {
		return "", errors.New("juju machine id not found")
	}
	token := strings.Split(v, "-")
	return token[len(token)-1], nil
}

func toPlacement(p *MachinePlacement, directive string) *instance.Placement {
	placement := &instance.Placement{
		Directive: directive,
	}
	switch {
	case p.LXD:
		placement.Scope = "lxd"
	case p.KVM:
		placement.Scope = "kvm"
	case p.Machine:
		placement.Scope = "#"
	}
	return placement
}

func toConstraint(c *MachineConstraint) constraints.Value {
	constraint := constraints.Value{}
	if c == nil {
		return constraint
	}
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

func checkDirExists(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
