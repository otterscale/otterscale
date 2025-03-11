package client

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	"github.com/google/wire"
	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/connectors/jira/client/tables"
)

var ProviderSet = wire.NewSet(NewConnector, openhdc.NewDefaultCodec)

// var ProviderSet = wire.NewSet(NewConnector, NewIssueCodec)

type Client struct {
	openhdc.Codec
	opts options

	jiraClient       *jira.Client
	issueSchema      *tables.IssueSchema
	issueFieldSchema *tables.IssueFieldSchema
	projectSchema    *tables.ProjectSchema
}

func NewConnector(c openhdc.Codec, opt ...Option) (openhdc.Connector, error) {
	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	if opts.server == "" {
		return nil, errors.New("server path is empty")
	}

	if opts.token == "" && (opts.username == "" || opts.password == "") {
		return nil, errors.New("username, password or token is empty")
	}

	// create jira client
	var jc *jira.Client
	var err error
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} // For Phison jira test
	if opts.token != "" {
		tp := jira.BearerAuthTransport{Token: opts.token, Transport: tr}
		jc, err = jira.NewClient(tp.Client(), opts.server)

	} else {
		tp := jira.BasicAuthTransport{Username: opts.username, Password: opts.password, Transport: tr}
		jc, err = jira.NewClient(tp.Client(), opts.server)
	}
	if err != nil {
		return nil, err
	}

	// create issue Schema
	is := tables.NewIssueSchema()
	ifs := tables.NewIssueFieldSchema()
	ps := tables.NewProjectSchema()

	return &Client{
		Codec:            c,
		opts:             opts,
		jiraClient:       jc,
		issueSchema:      is,
		issueFieldSchema: ifs,
		projectSchema:    ps,
	}, nil
}

func (c *Client) Name() string {
	return c.opts.name
}

func (c *Client) Close(ctx context.Context) error {
	return nil
}
