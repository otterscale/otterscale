package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/machine/metal"
)

type eventRepo struct {
	maas *MAAS
}

func NewEventRepo(maas *MAAS) metal.EventRepo {
	return &eventRepo{
		maas: maas,
	}
}

var _ metal.EventRepo = (*eventRepo)(nil)

func (r *eventRepo) Get(_ context.Context, machineID string) ([]metal.Event, error) {
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

func (r *eventRepo) toEvents(es []entity.Event) []metal.Event {
	ret := make([]metal.Event, 0, len(es))

	for _, e := range es {
		ret = append(ret, metal.Event{
			Type:    e.Type,
			Created: e.Created,
		})
	}

	return ret
}
