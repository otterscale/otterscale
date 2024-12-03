package grpc

import (
	"context"

	"github.com/openhdc/openhdc/pkg/transport"
)

var _ transport.Server = (*Server)(nil)

type Server struct{}

func NewServer(opts ...Option) *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context) error {
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}
