package openhdc

import (
	"context"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

const defaultMaxMessageSize = 1024 * 1024 * 100

type Server struct {
	opts serverOptions

	gs *grpc.Server
	hs *health.Server
}

type serverOptions struct {
	network     string
	address     string
	grpcOptions []grpc.ServerOption
}

var defaultServerOptions = serverOptions{
	grpcOptions: []grpc.ServerOption{
		grpc.MaxRecvMsgSize(defaultMaxMessageSize),
		grpc.MaxSendMsgSize(defaultMaxMessageSize),
	},
}

type ServerOption interface {
	apply(*serverOptions)
}

type funcServerOption struct {
	f func(*serverOptions)
}

var _ ServerOption = (*funcServerOption)(nil)

func (fro *funcServerOption) apply(ro *serverOptions) {
	fro.f(ro)
}

func newFuncServerOption(f func(*serverOptions)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}

func WithNetwork(n string) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.network = n
	})
}

func WithAddress(n string) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.network = n
	})
}

func WithGrpcOptions(opts ...grpc.ServerOption) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.grpcOptions = opts
	})
}

func NewServer(svc *Service, opt ...ServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	s := &Server{
		opts: opts,
		gs:   grpc.NewServer(opts.grpcOptions...),
		hs:   health.NewServer(),
	}
	grpc_health_v1.RegisterHealthServer(s.gs, s.hs)
	reflection.Register(s.gs)
	pb.RegisterConnectorServer(s.gs, svc)
	return s
}

func (s *Server) Start(_ context.Context) error {
	lis, err := net.Listen(s.opts.network, s.opts.address)
	if err != nil {
		return err
	}
	slog.Info("[gRPC] server started", "network", lis.Addr().Network(), "address", lis.Addr().String())
	s.hs.Resume()
	return s.gs.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("[gRPC] initiating server shutdown")
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		s.gs.GracefulStop()
		cancel()
	}()
	<-ctx.Done()
	s.gs.Stop()
	slog.Info("[gRPC] server stopped")
	return nil
}
