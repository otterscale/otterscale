package client

import (
	"context"
	"errors"
	"github.com/google/wire"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v68/github"

	"github.com/openhdc/openhdc"
)

var ProviderSet = wire.NewSet(NewConnector, openhdc.NewDefaultCodec)

type Client struct {
	openhdc.Codec
	opts options

	githubClient *github.Client
}

func NewConnector(c openhdc.Codec, opt ...Option) (openhdc.Connector, error) {
	o := options{}
	for _, option := range opt {
		option.apply(&o)
	}

	if o.owner == "" {
		return nil, errors.New("owner is empty")
	}

	ctx := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tokenClient := oauth2.NewClient(ctx, tokenService)

	g := github.NewClient(tokenClient)

	return &Client{
		Codec:        c,
		opts:         o,
		githubClient: g,
	}, nil
}

func (c *Client) Name() string {
	return c.opts.name
}

func (c *Client) Close(ctx context.Context) error {
	return nil
}
