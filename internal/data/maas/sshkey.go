package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/otterscale/internal/domain/service"
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

func (r *sshKey) List(_ context.Context) ([]entity.SSHKey, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.SSHKeys.Get()
}
