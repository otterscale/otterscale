package main

import (
	"context"

	"github.com/openhdc/openhdc/pkg/connector"
)

var _ connector.Destination = (*Connector)(nil)

type Connector struct{}

func newConnector() (connector.Destination, func()) {
	c := &Connector{}
	return c, func() {}
}

func (d *Connector) Close(context.Context) error {
	return nil
}

func (d *Connector) Write(context.Context) error {
	return nil
}
