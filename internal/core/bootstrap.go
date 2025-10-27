package core

import (
	"sync"

	"github.com/otterscale/otterscale/internal/config"
)

type BootstrapStatus struct {
	Phase   string
	Message string
	NewURL  string
}

type BootstrapUseCase struct {
	conf *config.Config

	statusMap sync.Map
}

func NewBootstrapUseCase(conf *config.Config) *BootstrapUseCase {
	return &BootstrapUseCase{
		conf: conf,
	}
}

func (uc *BootstrapUseCase) LoadStatus() *BootstrapStatus {
	v, ok := uc.statusMap.Load("")
	if ok {
		return v.(*BootstrapStatus)
	}
	return &BootstrapStatus{}
}

func (uc *BootstrapUseCase) StoreStatus(phase, message, newURL string) {
	uc.statusMap.Store("", &BootstrapStatus{
		Phase:   phase,
		Message: message,
		NewURL:  newURL,
	})
}
