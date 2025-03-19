package maas

import (
	"context"

	"github.com/openhdc/openhdc/internal/domain/service"
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
	value, err := r.maas.MAASServer.Get(name)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (r *server) Update(_ context.Context, name, value string) error {
	return r.maas.MAASServer.Post(name, value)
}
