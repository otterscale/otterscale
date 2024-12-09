package transport

import (
	"google.golang.org/grpc"
)

type ServerOption func(*Server)

func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

func Address(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

func Connector() ServerOption {
	return func(s *Server) {
		s.grpcOpts = append(s.grpcOpts,
			grpc.MaxRecvMsgSize(maxMsgSize),
			grpc.MaxSendMsgSize(maxMsgSize),
		)
	}
}
