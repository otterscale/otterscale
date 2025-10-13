package core

import (
	"sync"

	"github.com/otterscale/otterscale/internal/config"
)

type BootstrapStatus struct {
	Started  bool
	Finished bool
	Phase    string
	Message  string
}

type BootstrapUseCase struct {
	conf      *config.Config
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
	return &BootstrapStatus{
		Finished: isMAASConfigured(uc.conf),
	}
}

func (uc *BootstrapUseCase) StoreStatus(phase, message string) {
	uc.statusMap.Store("", &BootstrapStatus{
		Started:  true,
		Finished: isMAASConfigured(uc.conf),
		Phase:    phase,
		Message:  message,
	})
}

func isMAASConfigured(conf *config.Config) bool {
	return conf.MAAS.Key != "::"
}
