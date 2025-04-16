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

func (r *server) Get(_ context.Context, name string) ([]byte, error) {
	return r.maas.MAASServer.Get(name)
}

func (r *server) Update(_ context.Context, name, value string) error {
	return r.maas.MAASServer.Post(name, value)
}
