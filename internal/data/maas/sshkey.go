package maas

import (
	"context"

	"github.com/openhdc/otterscale/internal/core"
)

type sshKey struct {
	maas *MAAS
}

func NewSSHKey(maas *MAAS) core.SSHKeyRepo {
	return &sshKey{
		maas: maas,
	}
}

var _ core.SSHKeyRepo = (*sshKey)(nil)

func (r *sshKey) List(_ context.Context) ([]core.SSHKey, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.SSHKeys.Get()
}
