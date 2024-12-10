package connector

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

const maxMsgSize = 100 * 1024 * 1024

type Server struct {
	opts         serverOptions
	grpcServer   *grpc.Server
	healthServer *health.Server
}

func NewServer(svc *Service, opts ...ServerOption) *Server {
	o := serverOptions{
		grpcOpts: []grpc.ServerOption{
			grpc.MaxRecvMsgSize(maxMsgSize),
			grpc.MaxSendMsgSize(maxMsgSize),
		},
	}
	for _, opt := range opts {
		opt(&o)
	}
	gs := grpc.NewServer(o.grpcOpts...)
	hs := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gs, hs)
	reflection.Register(gs)
	pb.RegisterConnectorServer(gs, svc)
	return &Server{
		opts:         o,
		grpcServer:   gs,
		healthServer: hs,
	}
}

// TODO: ctx
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen(s.opts.network, s.opts.address)
	if err != nil {
		return err
	}
	slog.Debug(fmt.Sprintf("[gRPC] server listening on: %s", lis.Addr().String()))
	s.healthServer.Resume()
	return s.grpcServer.Serve(lis)
}

// TODO: ctx
func (s *Server) Stop(ctx context.Context) error {
	s.healthServer.Shutdown()
	s.grpcServer.GracefulStop()
	slog.Debug("[gRPC] server stopping")
	return nil
}
