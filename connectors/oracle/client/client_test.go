package client

import (
	"testing"

	"github.com/openhdc/openhdc"
	"github.com/stretchr/testify/assert"
)

// TestNewConnector_EmptyConnectionString tests the behavior of NewConnector when provided
// with an empty connection string. It verifies that the connector is nil and an
// appropriate error message is returned.
func TestNewConnector_EmptyConnectionString(t *testing.T) {
	codec := openhdc.NewDefaultCodec()
	connector, err := NewConnector(codec, WithConnString(""))

	if connector != nil {
		t.Error("expected connector to be nil")
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "connection string is empty"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestNewConnector_InvalidConnectionString(t *testing.T) {
	codec := openhdc.NewDefaultCodec()
	connector, err := NewConnector(codec, WithConnString("oracle://system:system@192.168.197.171:1521/"))

	if connector != nil {
		t.Error("expected connector to be nil")
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() == "connection string is empty" {
		t.Error("expected parsing error, got empty connection string error")
	}
}

// func TestNewConnector_WithValidConnectionString(t *testing.T) {
// 	codec := openhdc.NewDefaultCodec()
// 	connector, err := NewConnector(codec, WithConnString("oracle://system:system@192.168.197.171:1521/oracle"))

// 	assert.NoError(t, err, "expected no error for valid connection string")
// 	assert.NotNil(t, connector, "expected connector to be non-nil")

// 	client, ok := connector.(*Client)
// 	assert.True(t, ok, "expected connector to be of type *Client")

// 	// Check that the connection string is set correctly in options
// 	// assert.Equal(t, "valid_connection_string", client.opts.connString, "connection string should be set correctly")
// 	assert.Equal(t, "oracle://system:system@192.168.197.171:1521/oracle", client.opts.connString, "connection string should be set correctly")

// 	// Check that the default name and namespace are empty
// 	assert.Equal(t, "", client.opts.name, "expected default name to be empty")
// 	assert.Equal(t, "", client.opts.namespace, "expected default namespace to be empty")

// 	// Test that Close does not return error
// 	err = client.Close(context.TODO())
// 	assert.NoError(t, err, "expected no error on Close with valid connector")
// }

func TestNewConnector_NilCodec(t *testing.T) {
	connector, err := NewConnector(nil, WithConnString("oracle://system:system@192.168.197.171:1521/oracle"))
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

	if client.Codec != nil {
		t.Error("expected codec to be nil")
	}
}

// func TestClient_TestConnection_InvalidPool(t *testing.T) {
// 	client := &Client{
// 		Codec: openhdc.NewDefaultCodec(),
// 		pool:  nil,
// 	}

// 	ctx := context.Background()
// 	err := client.TestConnection(ctx)

// 	assert.Error(t, err, "Expected error when testing connection with invalid pool")
// 	assert.Contains(t, err.Error(), "ping error", "Expected ping error message")
// }

// func TestTestConnection_NilContext(t *testing.T) {
// 	codec := openhdc.NewDefaultCodec()
// 	connector, err := NewConnector(codec, WithConnString("oracle://system:system@192.168.197.171:1521/oracle"))
// 	if err != nil {
// 		t.Fatalf("failed to create connector: %v", err)
// 	}

// 	client, ok := connector.(*Client)
// 	if !ok {
// 		t.Fatal("expected connector to be of type *Client")
// 	}

// 	err = client.TestConnection(context.Background())

// 	assert.NoError(t, err, "expected no error for valid context")
// }

func TestClient_Name_EmptyNameOption(t *testing.T) {
	codec := openhdc.NewDefaultCodec()
	connector, err := NewConnector(codec, WithConnString("oracle://system:system@192.168.197.171:1521/oracle"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	client, ok := connector.(*Client)
	if !ok {
		t.Fatal("expected connector to be of type *Client")
	}

	name := client.Name()
	assert.Empty(t, name, "expected empty string when name option was not set")
}
