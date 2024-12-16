package openhdc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotImplemented = status.Errorf(codes.Unimplemented, "not implemented")
	ErrNotSupported   = status.Errorf(codes.InvalidArgument, "not supported")
)
