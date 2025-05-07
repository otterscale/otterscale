package client

import (
	"context"
	"testing"

	"github.com/openhdc/openhdc"
)

func TestNewConnector_Success(t *testing.T) {
	codec := openhdc.NewDefaultCodec()
	connector, err := NewConnector(codec, WithOwner("test-owner"), WithName("test-name"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if connector == nil {
		t.Fatal("expected connector to be non-nil")
	}

	client, ok := connector.(*Client)
	if !ok {
		t.Fatal("expected connector to be of type *Client")
	}

	if client.Name() != "test-name" {
		t.Errorf("expected name to be 'test-name', got %s", client.Name())
	}
}

func TestNewConnector_MissingOwner(t *testing.T) {
	codec := openhdc.NewDefaultCodec()
	_, err := NewConnector(codec)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "owner is empty"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestClient_Close(t *testing.T) {
	codec := openhdc.NewDefaultCodec()
	connector, err := NewConnector(codec, WithOwner("test-owner"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	client, ok := connector.(*Client)
	if !ok {
		t.Fatal("expected connector to be of type *Client")
	}

	err = client.Close(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}