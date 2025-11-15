package bootstrap

import (
	"sync"

	"github.com/otterscale/otterscale/internal/config"
)

type Status struct {
	Phase   string
	Message string
	NewURL  string
}

type UseCase struct {
	conf *config.Config

	m sync.Map
}

func NewUseCase(conf *config.Config) *UseCase {
	return &UseCase{
		conf: conf,
	}
}

func (uc *UseCase) LoadStatus() *Status {
	v, ok := uc.m.Load("default")
	if ok {
		return v.(*Status)
	}
	return &Status{}
}

func (uc *UseCase) StoreStatus(phase, message, newURL string) {
	uc.m.Store("default", &Status{
		Phase:   phase,
		Message: message,
		NewURL:  newURL,
	})
}
