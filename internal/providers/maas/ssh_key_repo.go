package maas

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/otterscale/otterscale/internal/core/scope/scope"
)

type sshKeyRepo struct {
	maas *MAAS
}

func NewSSHKeyRepo(maas *MAAS) scope.SSHKeyRepo {
	return &sshKeyRepo{
		maas: maas,
	}
}

var _ scope.SSHKeyRepo = (*sshKeyRepo)(nil)

func (r *sshKeyRepo) First(_ context.Context) (string, error) {
	client, err := r.maas.Client()
	if err != nil {
		return "", err
	}

	sshKeys, err := client.SSHKeys.Get()
	if err != nil {
		return "", err
	}

	if len(sshKeys) == 0 {
		return "", connect.NewError(connect.CodeNotFound, errors.New("ssh key not found"))
	}

	return sshKeys[0].Key, nil
}
