package standalone

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/canonical/gomaasclient/entity/node"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core/configuration"
	"github.com/otterscale/otterscale/internal/core/facility"
	"github.com/otterscale/otterscale/internal/core/machine"
	"github.com/otterscale/otterscale/internal/core/machine/tag"
	"github.com/otterscale/otterscale/internal/core/network"
)

type charm struct {
	Name           string
	Channel        string
	PlacementScope string
	Subordinate    bool
}

type StandaloneUseCase struct {
	conf *config.Config

	facility     facility.FacilityRepo
	ipRange      network.IPRangeRepo
	machine      machine.MachineRepo
	orchestrator machine.OrchestratorRepo
	provisioner  configuration.ProvisionerRepo
	relation     facility.RelationRepo
	subnet       network.SubnetRepo
	tag          tag.TagRepo
}

func NewStandaloneUseCase(conf *config.Config, facility facility.FacilityRepo, ipRange network.IPRangeRepo, machine machine.MachineRepo, orchestrator machine.OrchestratorRepo, provisioner configuration.ProvisionerRepo, relation facility.RelationRepo, subnet network.SubnetRepo, tag tag.TagRepo) *StandaloneUseCase {
	return &StandaloneUseCase{
		conf:         conf,
		facility:     facility,
		ipRange:      ipRange,
		machine:      machine,
		orchestrator: orchestrator,
		provisioner:  provisioner,
		relation:     relation,
		subnet:       subnet,
		tag:          tag,
	}
}

func (uc *StandaloneUseCase) CreateNode(ctx context.Context, scope, machineID string, virtualIPs []string, calicoCIDR string, osdDevices []string) (err error) {
	// cleanup reserved IPs on error
	releaseFunc := []func() error{}
	defer func() {
		if err != nil {
			for _, fn := range releaseFunc {
				_ = fn()
			}
		}
	}()

	// check maas machine status
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return err
	}

	if machine.Status != node.StatusDeployed {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("machine is not deployed"))
	}

	// check juju agent status
	jujuID, err := uc.machine.ExtractJujuID(machine)
	if err != nil {
		return err
	}

	agentStatus, err := uc.orchestrator.AgentStatus(ctx, scope, jujuID)
	if err != nil {
		return err
	}

	if agentStatus != "started" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("machine is not started"))
	}

	// prepare arguments
	if len(osdDevices) == 0 {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("no OSD devices provided"))
	}

	nfsVIP, releaseNFSVIP, err := uc.reserveIP(ctx, machine, fmt.Sprintf("Ceph NFS IP for %q", scope))
	if err != nil {
		return err
	}

	releaseFunc = append(releaseFunc, releaseNFSVIP)

	if calicoCIDR == "" {
		calicoCIDR = defaultCalicoCIDR
	}

	if len(virtualIPs) == 0 {
		vip, releaseVIP, err := uc.reserveIP(ctx, machine, fmt.Sprintf("Kubernetes Load Balancer IP for %q", scope))
		if err != nil {
			return err
		}

		virtualIPs = append(virtualIPs, vip.String())
		releaseFunc = append(releaseFunc, releaseVIP)
	}

	ceph := newCeph(scope, osdDevices, nfsVIP.String())
	kubernetes := newKubernetes(scope, virtualIPs, calicoCIDR)
	addons := newAddons(scope)

	for _, base := range []base{kubernetes, ceph, addons} {
		if err := uc.deploy(ctx, scope, machineID, jujuID, base); err != nil {
			return err
		}
	}

	return uc.createCOS(ctx, scope)
}

func (uc *StandaloneUseCase) applyTags(ctx context.Context, machineID string, tags []string) error {
	for _, t := range tags {
		_, _ = uc.tag.Create(ctx, t, tag.BuiltIn) // Ignore error if tag already exists

		if err := uc.tag.AddMachines(ctx, t, []string{machineID}); err != nil {
			return err
		}
	}

	return nil
}

func (uc *StandaloneUseCase) deploy(ctx context.Context, scope, maasID, jujuID string, base base) error {
	series, err := uc.provisioner.Get(ctx, "default_distro_series")
	if err != nil {
		return err
	}

	if err := uc.applyTags(ctx, maasID, base.Tags()); err != nil {
		return err
	}

	eg, egctx := errgroup.WithContext(ctx)

	for _, charm := range base.Charms() {
		eg.Go(func() error {
			name := scope + "-" + strings.TrimPrefix(charm.Name, "ch:")

			configs, err := base.Configs()
			if err != nil {
				return err
			}

			return uc.facility.Create(egctx, scope, name, configs, charm.Name, charm.Channel, charm.PlacementScope, charm.Subordinate, jujuID, series)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	eg, egctx = errgroup.WithContext(ctx)

	for _, endpoints := range base.Relations() {
		eg.Go(func() error {
			return uc.relation.Create(egctx, scope, endpoints)
		})
	}

	return eg.Wait()
}

func buildConfigs(scope string, configs map[string]map[string]any) (string, error) {
	m := map[string]string{}

	for name, config := range configs {
		m := map[string]any{
			scope + "-" + name: config,
		}

		value, err := yaml.Marshal(m)
		if err != nil {
			return "", fmt.Errorf("failed to marshal config for %s: %w", name, err)
		}

		m["ch:"+name] = string(value)
	}

	data, err := yaml.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal final configs: %w", err)
	}

	return string(data), nil
}

func buildRelations(scope string, relationList [][]string) [][]string {
	relations := [][]string{}

	for _, r := range relationList {
		endpoints := []string{}

		for _, endpoint := range r {
			endpoints = append(endpoints, scope+"-"+endpoint)
		}

		relations = append(relations, endpoints)
	}

	return relations
}
