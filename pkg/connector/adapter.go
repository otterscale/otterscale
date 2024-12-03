package connector

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

var _ pb.ConnectorServer = (*Adapter)(nil)

type Adapter struct {
	pb.UnimplementedConnectorServer

	source      Source
	destination Destination
}

func NewSourceAdapter(source Source) *Adapter {
	return &Adapter{
		source: source,
	}
}

func NewDestinationAdapter(destination Destination) *Adapter {
	return &Adapter{
		destination: destination,
	}
}

func (s *Adapter) Close(ctx context.Context, _ *pb.CloseRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.destination.Close(ctx)
}

func (s *Adapter) Pull(req *pb.PullRequest, stream pb.Connector_PullServer) error {
	ctx := stream.Context()
	return s.source.Read(ctx)
}

func (s *Adapter) Push(stream pb.Connector_PushServer) error {
	ctx := stream.Context()
	return s.destination.Write(ctx)
}
