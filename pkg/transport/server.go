package transport

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const maxMsgSize = 100 * 1024 * 1024

type Server struct {
	*grpc.Server
	network      string
	address      string
	healthServer *health.Server
	grpcOpts     []grpc.ServerOption
}

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		network:      "tcp",
		address:      ":0",
		healthServer: health.NewServer(),
	}
	for _, opt := range opts {
		opt(s)
	}
	s.Server = grpc.NewServer(s.grpcOpts...)
	grpc_health_v1.RegisterHealthServer(s.Server, s.healthServer)
	reflection.Register(s.Server)
	return s
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	s.healthServer.Resume()
	return s.Serve(listener)
}

func (s *Server) Stop(ctx context.Context) error {
	s.healthServer.Shutdown()
	s.GracefulStop()
	return nil
}
