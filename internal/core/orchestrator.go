package core

import (
	"context"
	"maps"
	"slices"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/juju/juju/core/status"
	"github.com/juju/juju/rpc/params"

	"github.com/otterscale/otterscale/internal/config"
)

type Essential struct {
	Type  EssentialType
	Name  string
	Scope string
	Units []EssentialUnit
}

type EssentialUnit struct {
	Name      string
	Directive string
}

type EssentialCharm struct {
	Name        string
	Channel     string
	Revision    int
	LXD         bool
	Machine     bool
	Subordinate bool
}

type OrchestratorUseCase struct {
	conf *config.Config

	action         ActionRepo
	chart          ChartRepo
	client         ClientRepo
	facility       FacilityRepo
	facilityOffers FacilityOffersRepo
	ipRange        IPRangeRepo
	kubeApps       KubeAppsRepo
	kubeCore       KubeCoreRepo
	machine        MachineRepo
	release        ReleaseRepo
	scope          ScopeRepo
	server         ServerRepo
	subnet         SubnetRepo
	tag            TagRepo
}

func NewOrchestratorUseCase(conf *config.Config, action ActionRepo, chart ChartRepo, client ClientRepo, facility FacilityRepo, facilityOffers FacilityOffersRepo, ipRange IPRangeRepo, kubeApps KubeAppsRepo, kubeCore KubeCoreRepo, machine MachineRepo, release ReleaseRepo, scope ScopeRepo, server ServerRepo, subnet SubnetRepo, tag TagRepo) *OrchestratorUseCase {
	return &OrchestratorUseCase{
		conf:           conf,
		action:         action,
		chart:          chart,
		client:         client,
		facility:       facility,
		facilityOffers: facilityOffers,
		ipRange:        ipRange,
		kubeApps:       kubeApps,
		kubeCore:       kubeCore,
		machine:        machine,
		release:        release,
		scope:          scope,
		server:         server,
		subnet:         subnet,
		tag:            tag,
	}
}

func (uc *OrchestratorUseCase) ListEssentials(ctx context.Context, esType EssentialType, scope string) ([]Essential, error) {
	eg, egctx := errgroup.WithContext(ctx)
	result := make([][]Essential, 2)
	if esType == EssentialTypeUnknown || esType == EssentialTypeKubernetes {
		eg.Go(func() error {
			v, err := uc.listEssentials(egctx, KubernetesControlPlane, scope)
			if err == nil {
				result[0] = v
			}
			return err
		})
	}
	if esType == EssentialTypeUnknown || esType == EssentialTypeCeph {
		eg.Go(func() error {
			v, err := uc.listEssentials(egctx, CephMon, scope)
			if err == nil {
				result[1] = v
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return append(result[0], result[1]...), nil
}

func (uc *OrchestratorUseCase) ListKubernetesNodeLabels(ctx context.Context, scope, facility, hostname string, all bool) (map[string]string, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	node, err := uc.kubeCore.GetNode(ctx, config, hostname)
	if err != nil {
		return nil, err
	}
	if !all {
		maps.DeleteFunc(node.Labels, func(k, _ string) bool {
			parts := strings.Split(k, "/")
			return len(parts) < 2 || !strings.HasSuffix(parts[0], LabelDomain)
		})
	}
	return node.Labels, nil
}

func (uc *OrchestratorUseCase) UpdateKubernetesNodeLabels(ctx context.Context, scope, facility, hostname string, labels map[string]string) (map[string]string, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	node, err := uc.kubeCore.GetNode(ctx, config, hostname)
	if err != nil {
		return nil, err
	}
	if node.Labels == nil {
		node.Labels = map[string]string{}
	}
	for k, v := range labels {
		if v == "" {
			delete(node.Labels, k)
		} else {
			node.Labels[k] = v
		}
	}
	updatedNode, err := uc.kubeCore.UpdateNode(ctx, config, node)
	if err != nil {
		return nil, err
	}
	return updatedNode.Labels, nil
}

func (uc *OrchestratorUseCase) listEssentials(ctx context.Context, charmName, scope string) ([]Essential, error) {
	scopes, err := uc.getAvailableScopes(ctx, scope)
	if err != nil {
		return nil, err
	}

	essentials, err := uc.fetchEssentialsFromScopes(ctx, scopes, charmName)
	if err != nil {
		return nil, err
	}

	return uc.sortEssentials(essentials), nil
}

func (uc *OrchestratorUseCase) getAvailableScopes(ctx context.Context, scope string) ([]Scope, error) {
	scopes, err := uc.scope.List(ctx)
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(scopes, func(s Scope) bool {
		return !strings.Contains(s.Name, scope) || s.Status.Status != status.Available
	}), nil
}

func (uc *OrchestratorUseCase) fetchEssentialsFromScopes(ctx context.Context, scopes []Scope, charmName string) ([]Essential, error) {
	eg, egctx := errgroup.WithContext(ctx)
	result := make([][]Essential, len(scopes))

	for i := range scopes {
		eg.Go(func() error {
			essentials, err := uc.getEssentialsForScope(egctx, &scopes[i], charmName)
			if err != nil {
				return err
			}
			result[i] = essentials
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return uc.flattenEssentials(result), nil
}

func (uc *OrchestratorUseCase) getEssentialsForScope(ctx context.Context, scope *Scope, charmName string) ([]Essential, error) {
	status, err := uc.client.Status(ctx, scope.Name, []string{"application", "*"})
	if err != nil {
		return nil, err
	}

	var essentials []Essential
	for name := range status.Applications {
		if !strings.Contains(status.Applications[name].Charm, charmName) {
			continue
		}

		units := uc.extractUnits(status.Applications[name].Units)
		essentials = append(essentials, Essential{
			Type:  essentialTypeFromCharmName(charmName),
			Name:  name,
			Scope: scope.Name,
			Units: units,
		})
	}

	return essentials, nil
}

func (uc *OrchestratorUseCase) extractUnits(statusUnits map[string]params.UnitStatus) []EssentialUnit {
	var units []EssentialUnit
	for uname := range statusUnits {
		units = append(units, EssentialUnit{
			Name:      uname,
			Directive: statusUnits[uname].Machine,
		})
	}
	return units
}

func (uc *OrchestratorUseCase) flattenEssentials(result [][]Essential) []Essential {
	var flattened []Essential
	for _, essentials := range result {
		flattened = append(flattened, essentials...)
	}
	return flattened
}

func (uc *OrchestratorUseCase) sortEssentials(essentials []Essential) []Essential {
	slices.SortFunc(essentials, func(e1, e2 Essential) int {
		return strings.Compare(e1.Name, e2.Name)
	})
	return essentials
}
