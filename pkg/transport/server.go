package transport

import (
	"context"

	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
}

func NewServer(opts ...Option) *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context) error {
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}
