package maas

import (
	"context"

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

func (r *server) Get(_ context.Context, name string) ([]byte, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.MAASServer.Get(name)
}

func (r *server) Update(_ context.Context, name, value string) error {
	client, err := r.maas.client()
	if err != nil {
		return err
	}
	return client.MAASServer.Post(name, value)
}
