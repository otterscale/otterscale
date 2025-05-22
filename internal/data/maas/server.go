package maas

import (
	"context"
	"strings"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type server struct {
	maas *MAAS
}

func NewServer(maas *MAAS) service.MAASServer {
	return &server{
		maas: maas,
	}
}

var _ service.MAASServer = (*server)(nil)

func (r *server) Get(_ context.Context, name string) (string, error) {
	client, err := r.maas.client()
	if err != nil {
		return "", err
	}
	data, err := client.MAASServer.Get(name)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(data), `"`, ``), nil
}

func (r *server) Update(_ context.Context, name, value string) error {
	client, err := r.maas.client()
	if err != nil {
		return err
	}
	return client.MAASServer.Post(name, value)
}
