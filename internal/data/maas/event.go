package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core"
)

type event struct {
	maas *MAAS
}

func NewEvent(maas *MAAS) core.EventRepo {
	return &event{
		maas: maas,
	}
}

var _ core.EventRepo = (*event)(nil)

func (r *event) Get(_ context.Context, systemID string) ([]core.Event, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	params := &entity.EventParams{ID: systemID}
	resp, err := client.Events.Get(params)
	if err != nil {
		return nil, err
	}
	return resp.Events, nil
}
