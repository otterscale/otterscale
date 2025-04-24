package maas

import (
	"context"
	"errors"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type sshKey struct {
	maas *MAAS
}

func NewSSHKey(maas *MAAS) service.MAASSSHKey {
	return &sshKey{
		maas: maas,
	}
}

var _ service.MAASSSHKey = (*sshKey)(nil)

func (r *sshKey) Default(_ context.Context) (*entity.SSHKey, error) {
	keys, err := r.maas.SSHKeys.Get()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		return &key, nil
	}
	return nil, errors.New("default ssh key not found")
}
