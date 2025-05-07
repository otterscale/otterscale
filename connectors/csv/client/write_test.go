package client

import (
	"context"
	"testing"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/stretchr/testify/assert"
)

// TestClient_Write tests the Write method of the Client struct.
// It creates a new Client instance, a background context, and a channel for messages.
// The test verifies that calling the Write method returns the expected ErrNotSupported error.
func TestClient_Write(t *testing.T) {
	client := &Client{}
	ctx := context.Background()
	msgs := make(chan *pb.Message)

	err := client.Write(ctx, msgs)
	assert.Equal(t, openhdc.ErrNotSupported, err)
}

