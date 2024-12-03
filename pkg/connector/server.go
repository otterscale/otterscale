package connector

import (
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/pkg/transport"
)

func NewTransportServer(a *Adapter) *transport.Server {
	srv := transport.NewServer()
	pb.RegisterConnectorServer(srv, a)
	return srv
}
