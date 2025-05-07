package client

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/stretchr/testify/assert"
)

func TestNewReadmeSchema(t *testing.T) {
	readmeSchema, builder := NewReadmeSchema()
	defer builder.Release()

	assert.NotNil(t, readmeSchema)
	assert.NotNil(t, builder)
	assert.NotNil(t, readmeSchema.ownerBuilder)
	assert.NotNil(t, readmeSchema.repoBuilder)
	assert.NotNil(t, readmeSchema.readmeBuilder)
}

func TestReadmeSchema_Append(t *testing.T) {
	readmeSchema, builder := NewReadmeSchema()
	defer builder.Release()

	client := &Client{
		opts: options{
			owner: "test-owner",
			repo:  "test-repo",
		},
	}

	readmeContent := "This is a test README content."
	readmeSchema.Append(client, readmeContent)

	record := builder.NewRecord()
	defer record.Release()

	assert.Equal(t, int64(1), record.NumRows())
	assert.Equal(t, "test-owner", record.Column(0).(*array.String).Value(0))
	assert.Equal(t, "test-repo", record.Column(1).(*array.String).Value(0))
	assert.Equal(t, readmeContent, record.Column(2).(*array.String).Value(0))
}