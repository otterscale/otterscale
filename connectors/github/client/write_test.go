package client

import (
	"context"
	"testing"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func TestClient_Write_NotSupported(t *testing.T) {
	client := &Client{}
	ctx := context.Background()
	msgs := make(chan *pb.Message)

	err := client.Write(ctx, msgs)
	if err != openhdc.ErrNotSupported {
		t.Errorf("expected error %v, got %v", openhdc.ErrNotSupported, err)
	}
}