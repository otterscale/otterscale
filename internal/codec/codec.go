package codec

import (
	"github.com/apache/arrow-go/v18/arrow/array"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotImplemented = status.Errorf(codes.Unimplemented, "not implemented")
	ErrNotSupported   = status.Errorf(codes.InvalidArgument, "not supported")
)

type Codec interface {
	Append(builder array.Builder, val any) error
}

type DefaultCodec struct{}

func (DefaultCodec) Append(builder array.Builder, val any) error {
	if val == nil {
		builder.AppendNull()
		return nil
	}
	return builder.AppendValueFromString(format(val))
}
