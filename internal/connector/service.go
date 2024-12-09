package connector

import (
	"context"
	"errors"
	"io"

	"github.com/apache/arrow-go/v18/arrow"
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pb.ConnectorServer = (*Service)(nil)

type Service struct {
	pb.UnimplementedConnectorServer

	opts      options
	connector Connector
}

func NewService(c Connector, opts ...Option) *Service {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &Service{
		opts:      o,
		connector: c,
	}
}

func (s *Service) Close(ctx context.Context, _ *pb.CloseRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.connector.Close(ctx)
}

func (s *Service) Pull(req *pb.PullRequest, stream pb.Connector_PullServer) error {
	msgs := make(chan arrow.Record)
	defer close(msgs)
	eg, ctx := errgroup.WithContext(stream.Context())
	eg.Go(func() error {
		return s.connector.Read(ctx, msgs, WithTables(req.Tables), WithSkipTables(req.SkipTables))
	})
	eg.Go(func() error {
		for msg := range msgs {
			rec, err := FromArrowRecord(msg)
			if err != nil {
				return err
			}
			res := &pb.PullResponse{
				Record: rec,
			}
			// TODO: CHECK SIZE
			// if proto.Size(res) > MaxMsgSize {
			// 	continue
			// }
			if err := stream.Send(res); err != nil {
				return err
			}
		}
		return nil
	})
	return eg.Wait()
}

func (s *Service) Push(stream pb.Connector_PushServer) error {
	msgs := make(chan arrow.Record)
	defer close(msgs)
	eg, ctx := errgroup.WithContext(stream.Context())
	eg.Go(func() error {
		return s.connector.Write(ctx, msgs)
	})
	eg.Go(func() error {
		for {
			req, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
			rec, err := ToArrowRecord(req.GetRecord())
			if err != nil {
				return err
			}
			msgs <- rec
		}
		return nil
	})
	return eg.Wait()
}
