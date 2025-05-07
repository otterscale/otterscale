package client

import (
	"context"
	"testing"

	"github.com/openhdc/openhdc"
	"github.com/stretchr/testify/assert"
)

func TestNewConnector_Success(t *testing.T) {
	codec := openhdc.NewDefaultCodec()
	connector, err := NewConnector(codec, WithServer("test-server"), WithUsername("test-username"), WithName("test-name"), WithPassword("test-password"), WithToken("test-token"))

	assert.NoError(t, err)
	assert.NotNil(t, connector)

	client, ok := connector.(*Client)
	assert.True(t, ok)
	assert.Equal(t, "test-name", client.Name())
	assert.Equal(t, "test-username", client.opts.username)
}

func TestNewConnector_MissingOwner(t *testing.T) {
	codec := openhdc.NewDefaultCodec()
	connector, err := NewConnector(codec)

	assert.Error(t, err)
	assert.Nil(t, connector)
	assert.Equal(t, "server path is empty", err.Error())
}

func TestClient_Close(t *testing.T) {
	codec := openhdc.NewDefaultCodec()
	connector, err := NewConnector(codec, WithServer("test-server"), WithUsername("test-username"), WithName("test-name"), WithPassword("test-password"), WithToken("test-token"))

	assert.NoError(t, err)
	assert.NotNil(t, connector)

	client, ok := connector.(*Client)
	assert.True(t, ok)

	err = client.Close(context.Background())
	assert.NoError(t, err)
}
