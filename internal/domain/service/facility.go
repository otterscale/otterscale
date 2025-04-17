package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	jujuyaml "gopkg.in/yaml.v2"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (s *NexusService) ListFacilities(ctx context.Context, uuid string) ([]model.Facility, error) {
	st, err := s.client.Status(ctx, uuid, []string{"application", "*"})
	if err != nil {
		return nil, err
	}
	fs := []model.Facility{}
	for key := range st.Applications {
		status := st.Applications[key]
		fs = append(fs, model.Facility{
			Name:   key,
			Status: &status,
		})
	}
	return fs, nil
}

func (s *NexusService) GetFacility(ctx context.Context, uuid, name string) (*model.Facility, error) {
	st, err := s.client.Status(ctx, uuid, []string{"application", name})
	if err != nil {
		return nil, err
	}
	for key := range st.Applications {
		if key != name {
			continue
		}
		status := st.Applications[key]
		return &model.Facility{
			Name:   key,
			Status: &status,
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "facility %q not found", name)
}

func (s *NexusService) GetFacilityMetadata(ctx context.Context, uuid, name string) (*model.FacilityMetadata, error) {
	cfg, err := s.facility.GetConfig(ctx, uuid, name)
	if err != nil {
		return nil, err
	}
	configYAML, _ := jujuyaml.Marshal(cfg)
	return &model.FacilityMetadata{
		ConfigYAML: string(configYAML),
	}, nil
}

func (s *NexusService) CreateFacility(ctx context.Context, uuid, name, configYAML, charmName, channel string, revision, number int, mps []model.MachinePlacement, mc *model.MachineConstraint, trust bool) (*model.Facility, error) {
	placements, err := s.toPlacements(ctx, uuid, mps)
	if err != nil {
		return nil, err
	}
	constraint := toConstraint(mc)
	if _, err := s.facility.Create(ctx, uuid, name, configYAML, charmName, channel, revision, number, placements, &constraint, trust); err != nil {
		return nil, err
	}
	return &model.Facility{}, nil
}

func (s *NexusService) UpdateFacility(ctx context.Context, uuid, name, configYAML string) (*model.Facility, error) {
	if err := s.facility.Update(ctx, uuid, name, configYAML); err != nil {
		return nil, err
	}
	return s.GetFacility(ctx, uuid, name)
}

func (s *NexusService) DeleteFacility(ctx context.Context, uuid, name string, destroyStorage, force bool) error {
	return s.facility.Delete(ctx, uuid, name, destroyStorage, force)
}

func (s *NexusService) ExposeFacility(ctx context.Context, uuid, name string) error {
	return s.facility.Expose(ctx, uuid, name, nil)
}

func (s *NexusService) AddFacilityUnits(ctx context.Context, uuid, name string, number int, mps []model.MachinePlacement) ([]string, error) {
	placements, err := s.toPlacements(ctx, uuid, mps)
	if err != nil {
		return nil, err
	}
	return s.facility.AddUnits(ctx, uuid, name, number, placements)
}

func (s *NexusService) ListActions(ctx context.Context, uuid, appName string) ([]model.Action, error) {
	asm, err := s.action.List(ctx, uuid, appName)
	if err != nil {
		return nil, err
	}
	as := []model.Action{}
	for name, spec := range asm {
		as = append(as, model.Action{
			Name: name,
			Spec: &spec,
		})
	}
	return as, nil
}

func (s *NexusService) toPlacements(ctx context.Context, uuid string, mps []model.MachinePlacement) ([]model.Placement, error) {
	ps := []model.Placement{}
	for _, mp := range mps {
		directives, err := s.maasToJujuMachineMap(ctx, uuid)
		if err != nil {
			return nil, err
		}
		p := toPlacement(&mp, directives[mp.MachineID])
		if p != nil {
			ps = append(ps, *p)
		}
	}
	return ps, nil
}
