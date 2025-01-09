package client

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	"github.com/google/wire"
	"github.com/openhdc/openhdc"
)

var ProviderSet = wire.NewSet(NewConnector, openhdc.NewDefaultCodec)

type Client struct {
	openhdc.Codec
	opts options

	jiraClient  *jira.Client
	issueSchema *IssueSchema
}

func NewConnector(c openhdc.Codec, opts ...Option) (openhdc.Connector, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	if o.server == "" {
		return nil, errors.New("server path is empty")
	}

	if o.token == "" && (o.username == "" || o.password == "") {
		return nil, errors.New("username, password or token is empty")
	}

	// create jira client
	var jc *jira.Client
	var err error
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} // For Phison jira test
	if o.token != "" {
		tp := jira.BearerAuthTransport{Token: o.token, Transport: tr}
		jc, err = jira.NewClient(tp.Client(), o.server)

	} else {
		tp := jira.BasicAuthTransport{Username: o.username, Password: o.password, Transport: tr}
		jc, err = jira.NewClient(tp.Client(), o.server)
	}
	if err != nil {
		return nil, err
	}

	// create issue Schema
	is := NewIssueSchema()

	return &Client{
		Codec:       c,
		opts:        o,
		jiraClient:  jc,
		issueSchema: is,
	}, nil
}

func (c *Client) Name() string {
	return c.opts.name
}

func (c *Client) Close(ctx context.Context) error {
	return nil
}
