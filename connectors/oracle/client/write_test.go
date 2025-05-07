package client

import (
	"context"
	"testing"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func TestClient_Write(t *testing.T) {
	type args struct {
		ctx  context.Context
		msgs <-chan *pb.Message
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Write(tt.args.ctx, tt.args.msgs); (err != nil) != tt.wantErr {
				t.Errorf("Client.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
