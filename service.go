package openhdc

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

type Service struct {
	pb.UnimplementedConnectorServer

	connector Connector
}

var _ pb.ConnectorServer = (*Service)(nil)

func NewService(c Connector) pb.ConnectorServer {
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
	eg.Go(func() error {
		r := NewReader()
		sync := req.GetSync()
		if sync != nil {
			r = NewReader(
				WithBatchSize(sync.GetBatchSize()),
				WithKeys(sync.GetKeys()...),
				WithSkipKeys(sync.GetSkipKeys()...),
				WithOptions(sync.GetOptions()...),
			)
		}
		return s.read(ctx, msgs, r)
	})
	eg.Go(func() error {
		return s.send(msgs, stream)
	})
	return eg.Wait()
}

func (s *Service) Push(stream pb.Connector_PushServer) error {
	msgs := make(chan *pb.Message)
	eg, ctx := errgroup.WithContext(stream.Context())
	eg.Go(func() error {
		return s.receive(ctx, msgs, stream)
	})
	eg.Go(func() error {
		return s.write(ctx, msgs)
	})
	return eg.Wait()
}

// canceled by context
func (s *Service) read(ctx context.Context, msgs chan *pb.Message, r *Reader) error {
	defer close(msgs)
	return s.connector.Read(ctx, msgs, r)
}

// canceled by close channel
func (s *Service) send(msgs chan *pb.Message, stream pb.Connector_PullServer) error {
	for msg := range msgs {
		if proto.Size(msg) > defaultMaxMessageSize {
			slog.Warn("[Send] skip oversized message")
			continue
		}
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

// canceled by next recv
func (s *Service) receive(ctx context.Context, msgs chan *pb.Message, stream pb.Connector_PushServer) error {
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
}

// canceled by close channel
func (s *Service) write(ctx context.Context, msgs chan *pb.Message) error {
	return s.connector.Write(ctx, msgs)
}
