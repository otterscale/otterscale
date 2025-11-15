package maas

import (
	"context"
	"strings"

	"github.com/otterscale/otterscale/internal/core/configuration"
)

// Note: MAAS API do not support context.
type provisionerRepo struct {
	maas *MAAS
}

func NewProvisionerRepo(maas *MAAS) configuration.ProvisionerRepo {
	return &provisionerRepo{
		maas: maas,
	}
}

var _ configuration.ProvisionerRepo = (*provisionerRepo)(nil)

func (r *provisionerRepo) Get(_ context.Context, name string) (string, error) {
	client, err := r.maas.Client()
	if err != nil {
		return "", err
	}

	data, err := client.MAASServer.Get(name)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(string(data), `"`, ``), nil
}

func (r *provisionerRepo) Update(_ context.Context, name, value string) error {
	client, err := r.maas.Client()
	if err != nil {
		return err
	}

	return client.MAASServer.Post(name, value)
}
