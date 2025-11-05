package core

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/canonical/gomaasclient/entity/node"
	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/status"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

const (
	defaultCalicoCIDR = "198.19.0.0/16"
)

func (uc *OrchestratorUseCase) validateOSDDevices(devices []string) (string, error) {
	osdDevices := strings.Join(devices, " ")
	if osdDevices == "" {
		return "", connect.NewError(connect.CodeInvalidArgument, errors.New("no OSD devices provided"))
	}
	return osdDevices, nil
}

func (uc *OrchestratorUseCase) resolveKubeVIPs(ctx context.Context, machineID, prefix string, userVirtualIPs []string) (string, error) {
	if len(userVirtualIPs) > 0 {
		return strings.Join(userVirtualIPs, " "), nil
	}

	ip, err := uc.reserveIP(ctx, machineID, fmt.Sprintf("Kubernetes Load Balancer IP for %s", prefix))
	if err != nil {
		return "", err
	}
	return ip.String(), nil
}

func (uc *OrchestratorUseCase) resolveCalicoCIDR(userCalicoCIDR string) string {
	if userCalicoCIDR != "" {
		return userCalicoCIDR
	}
	return defaultCalicoCIDR
}

func (uc *OrchestratorUseCase) createEssential(ctx context.Context, scope, machineID, prefix string, charms []EssentialCharm, configs map[string]string, tags []string) error {
	directive, err := uc.getMachineDirective(ctx, machineID)
	if err != nil {
		return err
	}

	if err := uc.applyTags(ctx, machineID, tags); err != nil {
		return err
	}

	return uc.deployCharms(ctx, scope, prefix, directive, charms, configs)
}

func (uc *OrchestratorUseCase) applyTags(ctx context.Context, machineID string, tags []string) error {
	for _, tag := range tags {
		_, _ = uc.tag.Create(ctx, tag, BuiltInMachineTagComment)
		if err := uc.tag.AddMachines(ctx, tag, []string{machineID}); err != nil {
			return err
		}
	}
	return nil
}

func (uc *OrchestratorUseCase) deployCharms(ctx context.Context, scope, prefix, directive string, charms []EssentialCharm, configs map[string]string) error {
	base, err := defaultBase(ctx, uc.server)
	if err != nil {
		return err
	}

	eg, egctx := errgroup.WithContext(ctx)
	for _, charm := range charms {
		eg.Go(func() error {
			return uc.deployCharm(egctx, scope, prefix, directive, charm, configs[charm.Name], &base)
		})
	}
	return eg.Wait()
}

func (uc *OrchestratorUseCase) deployCharm(ctx context.Context, scope, prefix, directive string, charm EssentialCharm, config string, base *base.Base) error {
	name := toEssentialName(prefix, charm.Name)
	placements := uc.buildPlacements(directive, charm)

	_, err := uc.facility.Create(ctx, scope, name, config, charm.Name, charm.Channel, charm.Revision, 1, base, placements, nil, true)
	return err
}

func (uc *OrchestratorUseCase) buildPlacements(directive string, charm EssentialCharm) []instance.Placement {
	if directive == "" || charm.Subordinate {
		return nil
	}

	placement := toPlacement(&MachinePlacement{LXD: charm.LXD, Machine: charm.Machine}, directive)
	return []instance.Placement{*placement}
}

func (uc *OrchestratorUseCase) createEssentialRelations(ctx context.Context, scope string, endpointList [][]string) error {
	eg, egctx := errgroup.WithContext(ctx)

	for _, endpoints := range endpointList {
		eg.Go(func() error {
			_, err := uc.facility.CreateRelation(egctx, scope, endpoints)
			return err
		})
	}

	return eg.Wait()
}

func (uc *OrchestratorUseCase) getMachineDirective(ctx context.Context, machineID string) (string, error) {
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return "", err
	}

	if machine.Status != node.StatusDeployed {
		return "", connect.NewError(connect.CodeInvalidArgument, errors.New("machine status is not deployed"))
	}

	return getJujuMachineID(machine.WorkloadAnnotations)
}

func (uc *OrchestratorUseCase) validateMachineStatus(ctx context.Context, scope, machineID string) error {
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return err
	}

	if machine.Status != node.StatusDeployed {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("machine is not deployed"))
	}

	return uc.validateJujuMachineStatus(ctx, scope, machine)
}

func (uc *OrchestratorUseCase) validateJujuMachineStatus(ctx context.Context, scope string, machine *Machine) error {
	id, err := getJujuMachineID(machine.WorkloadAnnotations)
	if err != nil {
		return err
	}

	fullStatus, err := uc.client.Status(ctx, scope, []string{"machine", id})
	if err != nil {
		return err
	}

	m, ok := fullStatus.Machines[id]
	if !ok {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("machine is not found"))
	}

	if m.AgentStatus.Status != status.Started.String() {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("machine is not started"))
	}

	return nil
}

func essentialTypeFromCharmName(charmName string) EssentialType {
	switch {
	case strings.Contains(charmName, "kubernetes"):
		return EssentialTypeKubernetes
	case strings.Contains(charmName, "ceph"):
		return EssentialTypeCeph
	default:
		return EssentialTypeUnknown
	}
}

func toEssentialName(prefix, charm string) string {
	charm = strings.TrimPrefix(charm, "ch:")
	if idx := strings.Index(charm, ":"); idx != -1 {
		charm = charm[idx+1:]
	}
	return prefix + "-" + charm
}

func toEssentialConfigs(prefix string, configs map[string]map[string]any) (map[string]string, error) {
	result := make(map[string]string, len(configs))

	for name, config := range configs {
		key := toEssentialName(prefix, name)
		value, err := yaml.Marshal(map[string]any{key: config})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config for %s: %w", name, err)
		}
		result["ch:"+name] = string(value)
	}

	return result, nil
}

func toEssentialEndpointList(prefix string, relationList [][]string) [][]string {
	endpointList := make([][]string, 0, len(relationList))

	for _, relations := range relationList {
		endpoints := make([]string, 0, len(relations))
		for _, relation := range relations {
			endpoints = append(endpoints, prefix+"-"+relation)
		}
		endpointList = append(endpointList, endpoints)
	}

	return endpointList
}
