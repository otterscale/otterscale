package metal

import (
	"context"
	"time"
)

type Event struct {
	Type    string
	Created string
}

type EventRepo interface {
	Get(ctx context.Context, machineID string) ([]Event, error)
}

func (uc *MetalUseCase) lastCommissionedAt(ctx context.Context, machineID string) (time.Time, error) {
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
