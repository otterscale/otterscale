package transport

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: BETTER
const maxMsgSize = 100 * 1024 * 1024

type ClientOption func(*Client)

func WithEndpoint(endpoint string) ClientOption {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(c *Client) {
		c.grpcOpts = opts
	}
}

func WithConnector() ClientOption {
	return func(c *Client) {
		c.grpcOpts = append(c.grpcOpts,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(maxMsgSize),
				grpc.MaxCallSendMsgSize(maxMsgSize),
			),
		)
	}
}

type Client struct {
	endpoint string
	grpcOpts []grpc.DialOption
}

func NewClient(opts ...ClientOption) (*grpc.ClientConn, error) {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	return grpc.NewClient(c.endpoint, c.grpcOpts...)
}
