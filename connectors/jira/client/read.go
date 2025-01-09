package client

import (
	"context"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func (c *Client) Read(ctx context.Context, msgs chan<- *pb.Message, rdr *openhdc.Reader) error {
	var err error
	for _, pj := range c.opts.projects {
		// sync issues table
		err = c.readIssue(pj, c.opts.startDate, msgs, rdr)
	}

	return err
}

func (c *Client) readIssue(pj string, sd string, msgs chan<- *pb.Message, rdr *openhdc.Reader) error {
	// migration
	if err := rdr.Send(pb.Migrate, msgs, c.issueSchema.Record()); err != nil {
		return err
	}

	// jql
	var jql strings.Builder
	jql.WriteString("project = '")
	jql.WriteString(pj)
	jql.WriteString("'")
	if sd != "" {
		jql.WriteString(" AND created >= '")
		jql.WriteString(sd)
		jql.WriteString("'")
	}

	// read jira issues
	var (
		count     int64
		startAt   = 0
		batchSize = 50
	)
	for {
		options := &jira.SearchOptions{
			StartAt:    startAt,
			MaxResults: batchSize,
		}
		// get issues
		issues, resp, err := c.jiraClient.Issue.Search(jql.String(), options)
		if err != nil {
			return err
		}

		// append issue
		for _, issue := range issues {
			c.issueSchema.Append(&issue)
			count++
		}

		// new message
		if count >= rdr.BatchSize() {
			if err := rdr.Send(pb.Insert, msgs, c.issueSchema.Record()); err != nil {
				return err
			}
			count = 0
		}

		// check is there next page
		if resp.StartAt+resp.MaxResults >= resp.Total {
			break // 已經抓取完所有 Issues
		}

		// update startAt to get next batch jira issues
		startAt = resp.StartAt + resp.MaxResults
	}

	// remain
	if count > 0 {
		if err := rdr.Send(pb.Insert, msgs, c.issueSchema.Record()); err != nil {
			return err
		}
	}

	return nil
}
