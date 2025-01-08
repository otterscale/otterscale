package client

import (
	"context"
	"fmt"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/api/property/v1"
)

func (c *Client) Read(ctx context.Context, msg chan<- *pb.Message, opts openhdc.ReadOptions) error {
	// sync issues table
	err := c.readIssue(msg)

	return err
}

func (c *Client) readIssue(msg chan<- *pb.Message) error {
	// TODO: change setting
	batchSize := 100 // 每次抓取的 Issue 數量
	startAt := 0     // 起始位置
	projectName := "[System] IFAS"

	// timestamp
	syncedAt := time.Now().UTC().Truncate(time.Second)

	// migration
	new, err := pb.NewMessage(property.MessageKind_migrate, c.issueSchema.Record(), c.opts.name, syncedAt)
	if err != nil {
		return err
	}
	msg <- new

	// read jira issues
	for {
		options := &jira.SearchOptions{
			StartAt:    startAt,
			MaxResults: batchSize,
		}
		// get issues
		issues, resp, err := c.jiraClient.Issue.Search(fmt.Sprintf("project = '%s'", projectName), options)
		if err != nil {
			return err
		}

		// append issue
		for _, issue := range issues {
			c.issueSchema.Append(&issue)
		}

		// new message
		new, err := pb.NewMessage(property.MessageKind_insert, c.issueSchema.Record(), c.opts.name, syncedAt)
		if err != nil {
			return err
		}
		msg <- new

		// check is there next page
		if resp.StartAt+resp.MaxResults >= resp.Total {
			break // 已經抓取完所有 Issues
		}

		// update startAt to get next batch jira issues
		startAt = resp.StartAt + resp.MaxResults
	}

	return err
}
