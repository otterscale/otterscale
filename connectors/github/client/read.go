package client

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func (c *Client) Read(ctx context.Context, msgs chan<- *pb.Message, rdr *openhdc.Reader) error {
	// get ReadMe from specific github repo
	readMe, _, err := c.githubClient.Repositories.GetReadme(ctx, c.opts.owner, c.opts.repo, nil)
	if err != nil {
		fmt.Printf("Problem in getting readme struct %v\n", err)
		return err
	}

	// get content
	content, err := readMe.GetContent()
	if err != nil {
		fmt.Printf("Problem in getting readme content %v\n", err)
		return err
	}

	// create Arrow memory pool
	pool := memory.NewGoAllocator()

	// define schema
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "content", Type: arrow.BinaryTypes.String},
		},
		nil,
	)

	// create a builder for a Arrow array
	builder := array.NewStringBuilder(pool)
	defer builder.Release()

	// append content to builder
	builder.Append(string(content))

	// create a  Arrow array
	arr := builder.NewArray()
	defer arr.Release()

	// create a Arrow Record
	record := array.NewRecord(schema, []arrow.Array{arr}, 1)
	defer record.Release()

	if err := rdr.Send(pb.Migrate, msgs, record); err != nil {
		return err
	}

	return err
}
