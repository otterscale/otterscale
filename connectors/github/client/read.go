package client

import (
	"context"
	"fmt"
	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func (c *Client) Read(ctx context.Context, msgs chan<- *pb.Message, rdr *openhdc.Reader) error {
	// get ReadMe from specific github repo
	readMe, resp, err := c.githubClient.Repositories.GetReadme(ctx, c.opts.owner, c.opts.repo, nil)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("issue with getting readme struct %v\n", err)
		return err
	}

	// get content
	content, err := readMe.GetContent()
	if err != nil {
		fmt.Printf("issue with getting readme content %v\n", err)
		return err
	}

	// send schema
	readmeSchema, record := NewReadmeSchema()
	if err := rdr.Send(pb.Migrate, msgs, record.NewRecord()); err != nil {
		return err
	}

	// send data
	readmeSchema.Append(c, content)
	if err := rdr.Send(pb.Insert, msgs, record.NewRecord()); err != nil {
		return err
	}

	return err
}
