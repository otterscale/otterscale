package client

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/stretchr/testify/assert"
)

// TestClient_Read tests the Read method of the Client struct.
// It uses table-driven tests to verify the behavior of the Read method
// with different CSV data inputs.
//
// The test cases include:
// - "valid csv": A CSV string with valid data, expecting no error.
// - "empty csv": An empty CSV string, expecting an error.
//
// For each test case, a temporary file is created and the CSV data is written to it.
// The Client struct is then initialized with the temporary file and the Read method is called.
// The test verifies whether the Read method returns an error as expected.
func TestClient_Read(t *testing.T) {
	tests := []struct {
		name    string
		csvData string
		wantErr bool
	}{
		{
			name:    "valid csv",
			csvData: "col1,col2\nval1,val2\nval3,val4",
			wantErr: false,
		},
		{
			name:    "empty csv",
			csvData: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "testfile")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(file.Name())

			if _, err := file.Write([]byte(tt.csvData)); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}
			if _, err := file.Seek(0, io.SeekStart); err != nil {
				t.Fatalf("failed to seek to start of temp file: %v", err)
			}

			client := &Client{
				opts: options{
					inferring: true,
					batchSize: 1,
					tableName: "test_table",
				},
				file: file,
			}

			msgs := make(chan *pb.Message, 10)
			rdr := &openhdc.Reader{}

			readErr := client.Read(context.Background(), msgs, rdr)
			if tt.wantErr {
				assert.Error(t, readErr)
			} else {
				assert.NoError(t, readErr)
			}
		})
	}
}

func TestClient_getFields(t *testing.T) {
	tests := []struct {
		name    string
		csvData string
		wantErr bool
	}{
		{
			name:    "valid csv",
			csvData: "col1,col2\nval1,val2\nval3,val4",
			wantErr: false,
		},
		{
			name:    "empty csv",
			csvData: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "testfile")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(file.Name())

			if _, err := file.Write([]byte(tt.csvData)); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}
			if _, err := file.Seek(0, io.SeekStart); err != nil {
				t.Fatalf("failed to seek to start of temp file: %v", err)
			}

			client := &Client{
				file: file,
			}

			_, err = client.getFields()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_newReader(t *testing.T) {
	tests := []struct {
		name    string
		csvData string
		wantErr bool
	}{
		{
			name:    "valid csv",
			csvData: "col1,col2\nval1,val2\nval3,val4",
			wantErr: false,
		},
		{
			name:    "empty csv",
			csvData: "",
			wantErr: true,
		},
		{
			name:    "validate field error",
			csvData: "col1,col2\nval1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "testfile")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(file.Name())

			if _, err := file.Write([]byte(tt.csvData)); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}
			if _, err := file.Seek(0, io.SeekStart); err != nil {
				t.Fatalf("failed to seek to start of temp file: %v", err)
			}

			// opts := options
			client := &Client{
				file: file,
				opts: options{
					inferring: true,
					batchSize: 1,
					tableName: "test_table",
				},
			}

			_, err = client.newReader()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
