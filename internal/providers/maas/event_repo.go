package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/machine"
)

type eventRepo struct {
	maas *MAAS
}

func NewEventRepo(maas *MAAS) machine.EventRepo {
	return &eventRepo{
		maas: maas,
	}
}

var _ machine.EventRepo = (*eventRepo)(nil)

func (r *eventRepo) Get(_ context.Context, machineID string) ([]machine.Event, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.EventParams{
		ID: machineID,
	}

	resp, err := client.Events.Get(params)
	if err != nil {
		return nil, err
	}

	return r.toEvents(resp.Events), nil
}

func (r *eventRepo) toEvents(es []entity.Event) []machine.Event {
	ret := []machine.Event{}

	for _, e := range es {
		ret = append(ret, machine.Event{
			Type:    e.Type,
			Created: e.Created,
		})
	}

	return ret
}
