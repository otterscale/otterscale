package machine

import (
	"context"
	"time"

	"github.com/canonical/gomaasclient/entity"
)

// Event represents a MAAS Event resource.
type Event = entity.Event

type EventRepo interface {
	Get(ctx context.Context, machineID string) ([]Event, error)
}

func (uc *MachineUseCase) lastCommissionedAt(ctx context.Context, machineID string) (time.Time, error) {
	events, err := uc.event.Get(ctx, machineID)
	if err != nil {
		return time.Time{}, err
	}
	for i := range events { // desc
		if events[i].Type == "Commissioning" {
			return time.Parse("Mon, 02 Jan. 2006 15:04:05", events[i].Created)
		}
	}
	return time.Time{}, nil
}
