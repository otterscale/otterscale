package connector

import (
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/pkg/transport/grpc"
)

func NewGrpcServer(a *Adapter) *grpc.Server {
	srv := grpc.NewServer()
	pb.RegisterConnectorServer(srv, a)
	return srv
}
