package core

import (
	"context"
	"fmt"
	"slices"
	"time"

	"connectrpc.com/connect"
	jujuyaml "gopkg.in/yaml.v2"

	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/api/client/application"
	corebase "github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
)

type (
	ActionSpec     = action.ActionSpec
	FacilityStatus = params.ApplicationStatus
	UnitStatus     = params.UnitStatus
	DetailedStatus = params.DetailedStatus
)

type FacilityMetadata struct {
	ConfigYAML string
}

type Facility struct {
	Name     string
	Metadata *FacilityMetadata
	Status   *FacilityStatus
}

type Action struct {
	Name string
	Spec *ActionSpec
}

type Charm struct {
	ID              string          `json:"id"`
	Type            string          `json:"type"`
	Name            string          `json:"name"`
	Result          CharmResult     `json:"result"`
	DefaultArtifact CharmArtifact   `json:"default-release"`
	Artifacts       []CharmArtifact `json:"channel-map"`
}

type CharmBase struct {
	Architecture string `json:"architecture"`
	Channel      string `json:"channel"`
	Name         string `json:"name"`
}

type CharmChannel struct {
	Base       CharmBase `json:"base"`
	Name       string    `json:"name"`
	ReleasedAt time.Time `json:"released-at"`
	Risk       string    `json:"risk"`
	Track      string    `json:"track"`
}

type CharmRevision struct {
	Bases     []CharmBase `json:"bases"`
	CreatedAt time.Time   `json:"created-at"`
	Revision  int         `json:"revision"`
	Version   string      `json:"version"`
}

type CharmArtifact struct {
	Channel  CharmChannel  `json:"channel"`
	Revision CharmRevision `json:"revision"`
}

type CharmResultCategory struct {
	Featured bool   `json:"featured"`
	Name     string `json:"name"`
}

type CharmResultMedia struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type CharmResultPublisher struct {
	DisplayName string `json:"display-name"`
	ID          string `json:"id"`
	Username    string `json:"username"`
	Validation  string `json:"validation"`
}

type CharmResult struct {
	BugsURL      string                `json:"bugs-url"`
	Categories   []CharmResultCategory `json:"categories"`
	DeployableOn []string              `json:"deployable-on"`
	Description  string                `json:"description"`
	License      string                `json:"license"`
	Media        []CharmResultMedia    `json:"media"`
	Publisher    *CharmResultPublisher `json:"publisher"`
	StoreURL     string                `json:"store-url"`
	StoreURLOld  string                `json:"store-url-old"`
	Summary      string                `json:"summary"`
	Title        string                `json:"title"`
	Unlisted     bool                  `json:"unlisted"`
	Website      string                `json:"website"`
}

type FacilityRepo interface {
	Create(ctx context.Context, uuid, name string, configYAML string, charmName, channel string, revision, number int, base *corebase.Base, placements []instance.Placement, constraint *constraints.Value, trust bool) (*application.DeployInfo, error)
	Update(ctx context.Context, uuid, name string, configYAML string) error
	Delete(ctx context.Context, uuid, name string, destroyStorage, force bool) error
	Expose(ctx context.Context, uuid, name string, endpoints map[string]params.ExposedEndpoint) error
	AddUnits(ctx context.Context, uuid, name string, number int, placements []instance.Placement) ([]string, error)
	ResolveUnitErrors(ctx context.Context, uuid string, units []string) error
	CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error)
	DeleteRelation(ctx context.Context, uuid string, id int) error
	GetConfig(ctx context.Context, uuid string, name string) (map[string]any, error)
	GetLeader(ctx context.Context, uuid, name string) (string, error)
	GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error)
	Consume(ctx context.Context, uuid string, args *crossmodel.ConsumeApplicationArgs) error
}

type FacilityOffersRepo interface {
	GetConsumeDetails(ctx context.Context, url string) (params.ConsumeOfferDetails, error)
}

type ActionRepo interface {
	List(ctx context.Context, uuid, appName string) (map[string]ActionSpec, error)
	RunCommand(ctx context.Context, uuid, unitName, command string) (string, error)
	RunAction(ctx context.Context, uuid, unitName, actionName string, parameters map[string]any) (string, error)
	GetResult(ctx context.Context, uuid, id string) (*action.ActionResult, error)
}

type CharmRepo interface {
	List(ctx context.Context) ([]Charm, error)
	Get(ctx context.Context, name string) (*Charm, error)
	ListArtifacts(ctx context.Context, name string) ([]CharmArtifact, error)
}

type FacilityUseCase struct {
	facility FacilityRepo
	server   ServerRepo
	client   ClientRepo
	action   ActionRepo
	charm    CharmRepo
	machine  MachineRepo
}

func NewFacilityUseCase(facility FacilityRepo, server ServerRepo, client ClientRepo, action ActionRepo, charm CharmRepo, machine MachineRepo) *FacilityUseCase {
	return &FacilityUseCase{
		facility: facility,
		server:   server,
		client:   client,
		action:   action,
		charm:    charm,
		machine:  machine,
	}
}

func (uc *FacilityUseCase) ListFacilities(ctx context.Context, uuid string) ([]Facility, error) {
	s, err := uc.client.Status(ctx, uuid, []string{"application", "*"})
	if err != nil {
		return nil, err
	}

	facilities := []Facility{}
	for appName := range s.Applications {
		appStatus := s.Applications[appName]
		facilities = append(facilities, Facility{
			Name:   appName,
			Status: &appStatus,
		})
	}
	return facilities, nil
}

func (uc *FacilityUseCase) GetFacility(ctx context.Context, uuid, name string) (*Facility, error) {
	s, err := uc.client.Status(ctx, uuid, []string{"application", name})
	if err != nil {
		return nil, err
	}

	app, ok := s.Applications[name]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("facility %q not found", name))
	}

	config, err := uc.facility.GetConfig(ctx, uuid, name)
	if err != nil {
		return nil, err
	}
	configYAML, _ := jujuyaml.Marshal(config)
	metadata := FacilityMetadata{
		ConfigYAML: string(configYAML),
	}

	return &Facility{
		Name:     name,
		Metadata: &metadata,
		Status:   &app,
	}, nil
}

func (uc *FacilityUseCase) CreateFacility(ctx context.Context, uuid, name, configYAML, charmName, channel string, revision, number int, mps []MachinePlacement, mc *MachineConstraint, trust bool) (*Facility, error) {
	base, err := defaultBase(ctx, uc.server)
	if err != nil {
		return nil, err
	}

	placements, err := uc.toPlacements(ctx, mps)
	if err != nil {
		return nil, err
	}

	constraint := toConstraint(mc)
	if _, err := uc.facility.Create(ctx, uuid, name, configYAML, charmName, channel, revision, number, &base, placements, &constraint, trust); err != nil {
		return nil, err
	}

	return &Facility{}, nil
}

func (uc *FacilityUseCase) UpdateFacility(ctx context.Context, uuid, name, configYAML string) (*Facility, error) {
	if err := uc.facility.Update(ctx, uuid, name, configYAML); err != nil {
		return nil, err
	}
	return &Facility{}, nil
}

func (uc *FacilityUseCase) DeleteFacility(ctx context.Context, uuid, name string, destroyStorage, force bool) error {
	return uc.facility.Delete(ctx, uuid, name, destroyStorage, force)
}

func (uc *FacilityUseCase) ExposeFacility(ctx context.Context, uuid, name string) error {
	return uc.facility.Expose(ctx, uuid, name, nil)
}

func (uc *FacilityUseCase) AddFacilityUnits(ctx context.Context, uuid, name string, number int, mps []MachinePlacement) ([]string, error) {
	placements, err := uc.toPlacements(ctx, mps)
	if err != nil {
		return nil, err
	}
	return uc.facility.AddUnits(ctx, uuid, name, number, placements)
}

func (uc *FacilityUseCase) ListActions(ctx context.Context, uuid, appName string) ([]Action, error) {
	actions, err := uc.action.List(ctx, uuid, appName)
	if err != nil {
		return nil, err
	}

	results := []Action{}
	for name, spec := range actions {
		results = append(results, Action{
			Name: name,
			Spec: &spec,
		})
	}
	return results, nil
}

func (uc *FacilityUseCase) ListCharms(ctx context.Context) ([]Charm, error) {
	charms, err := uc.charm.List(ctx)
	if err != nil {
		return nil, err
	}
	return uc.filterCharms(charms), nil
}

func (uc *FacilityUseCase) GetCharm(ctx context.Context, name string) (*Charm, error) {
	return uc.charm.Get(ctx, name)
}

func (uc *FacilityUseCase) ListArtifacts(ctx context.Context, name string) ([]CharmArtifact, error) {
	return uc.charm.ListArtifacts(ctx, name)
}

func (uc *FacilityUseCase) JujuToMAASMachineMap(ctx context.Context, uuid string) (map[string]string, error) {
	status, err := uc.client.Status(ctx, uuid, []string{"machine", "*"})
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	for name := range status.Machines {
		m[name] = string(status.Machines[name].InstanceId)
	}
	return m, nil
}

func (uc *FacilityUseCase) toPlacements(ctx context.Context, mps []MachinePlacement) ([]instance.Placement, error) {
	placements := []instance.Placement{}
	for _, mp := range mps {
		machine, err := uc.machine.Get(ctx, mp.MachineID)
		if err != nil {
			return nil, err
		}
		directive, err := getJujuMachineID(machine.WorkloadAnnotations)
		if err != nil {
			return nil, err
		}
		placement := toPlacement(&mp, directive)
		if placement != nil {
			placements = append(placements, *placement)
		}
	}
	return placements, nil
}

func (uc *FacilityUseCase) filterCharms(charms []Charm) []Charm {
	return slices.DeleteFunc(charms, func(charm Charm) bool {
		return slices.Contains(charm.Result.DeployableOn, "kubernetes") || charm.Type != "charm"
	})
}

func waitForActionCompleted(ctx context.Context, actionRepo ActionRepo, uuid, id string, tickInterval, timeoutDuration time.Duration) (*action.ActionResult, error) {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	timeout := time.After(timeoutDuration)
	for {
		select {
		case <-ticker.C:
			result, err := actionRepo.GetResult(ctx, uuid, id)
			if err != nil {
				return nil, err
			}
			if result.Status == "completed" { // state.ActionCompleted
				return result, nil
			}
			continue

		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for action %s to become completed", id)

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
