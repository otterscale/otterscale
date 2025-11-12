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

type BootstrapUseCase struct {
	conf *config.Config

	m sync.Map
}

func NewBootstrapUseCase(conf *config.Config) *BootstrapUseCase {
	return &BootstrapUseCase{
		conf: conf,
	}
}

func (uc *BootstrapUseCase) LoadStatus() *Status {
	v, ok := uc.m.Load("default")
	if ok {
		return v.(*Status)
	}
	return &Status{}
}

func (uc *BootstrapUseCase) StoreStatus(phase, message, newURL string) {
	uc.m.Store("default", &Status{
		Phase:   phase,
		Message: message,
		NewURL:  newURL,
	})
}
