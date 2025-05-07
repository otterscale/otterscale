package client

import (
	"context"
	"os"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/stretchr/testify/assert"
)

type mockCodec struct{}

func (m *mockCodec) Encode(b array.Builder, v any) error {
	return nil
}

func (m *mockCodec) Decode(data arrow.Array, i int) (any, error) {
	return nil, nil
}

func TestNewConnector(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		expectErr bool
	}{
		{"Empty file path", "", true},
		{"Invalid file path", "invalid/path.csv", true},
		// {"Valid file path", "testdata/testfile.csv", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{
				WithFilePath(tt.filePath),
			}
			// opts := []Option.Options{
			// 	func(o *options) {
			// 		o.filePath = tt.filePath
			// 	},
			// }
			connector, err := NewConnector(&mockCodec{}, opts...)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, connector)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, connector)
			}
		})
	}
}

func TestClient_Name(t *testing.T) {
	client := &Client{
		opts: options{name: "test-client"},
	}
	assert.Equal(t, "test-client", client.Name())
}

func TestClient_Close(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	client := &Client{
		file: file,
	}

	err = client.Close(context.Background())
	assert.NoError(t, err)
}
