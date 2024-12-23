package openhdc

import (
	"context"
	"errors"
	"fmt"
	"io"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

var _ pb.ConnectorServer = (*Service)(nil)

type Service struct {
	pb.UnimplementedConnectorServer

	connector Connector
}

func NewService(c Connector) *Service {
	return &Service{
		connector: c,
	}
}

func (s *Service) Close(ctx context.Context, _ *pb.CloseRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.connector.Close(ctx)
}

func (s *Service) Pull(req *pb.PullRequest, stream pb.Connector_PullServer) error {
	msgs := make(chan *pb.Message)
	eg, ctx := errgroup.WithContext(stream.Context())
	// canceled by context
	eg.Go(func() error {
		defer close(msgs)
		o := ReadOptions{}
		opts := []ReadOption{
			WithTables(req.GetTables()),
			WithSkipTables(req.GetSkipTables()),
		}
		for _, opt := range opts {
			opt(&o)
		}
		return s.connector.Read(ctx, msgs, o)
	})
	// canceled by close channel
	eg.Go(func() error {
		for msg := range msgs {
			if proto.Size(msg) > maxMsgSize {
				fmt.Println("[pull] skip oversized message")
				continue
			}
			if err := stream.Send(msg); err != nil {
				return err
			}
		}
		return nil
	})
	return eg.Wait()
}

func (s *Service) Push(stream pb.Connector_PushServer) error {
	msgs := make(chan *pb.Message)
	eg, ctx := errgroup.WithContext(stream.Context())
	// canceled by next recv
	eg.Go(func() error {
		defer close(msgs)
		for {
			msg, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return stream.SendAndClose(&emptypb.Empty{})
			}
			if err != nil {
				return err
			}
			select {
			case msgs <- msg:
			case <-ctx.Done():
				return nil
			}
		}
	})
	// canceled by close channel
	eg.Go(func() error {
		return s.connector.Write(ctx, msgs)
	})
	return eg.Wait()
}
