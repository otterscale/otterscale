package status

import (
	"sync"

	"github.com/otterscale/otterscale/internal/config"
)

type Status struct {
	Phase   string
	Message string
	NewURL  string
}

type StatusUseCase struct {
	conf *config.Config

	m sync.Map
}

func NewStatusUseCase(conf *config.Config) *StatusUseCase {
	return &StatusUseCase{
		conf: conf,
	}
}

func (uc *StatusUseCase) LoadStatus() *Status {
	v, ok := uc.m.Load("default")
	if ok {
		return v.(*Status)
	}
	return &Status{}
}

func (uc *StatusUseCase) StoreStatus(phase, message, newURL string) {
	uc.m.Store("default", &Status{
		Phase:   phase,
		Message: message,
		NewURL:  newURL,
	})
}
