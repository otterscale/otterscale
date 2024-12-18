package client

import (
	"errors"
	"fmt"
	"io"

	"github.com/apache/arrow-go/v18/arrow"
)

func (c *Client) validateFile(err error) error {
	switch {
	case errors.Is(err, io.EOF):
		err := fmt.Errorf("file %s is empty", c.opts.path)
		return err
	case err != nil:
		return err
	default:
		return nil
	}
}

func (c *Client) validateHeader(header arrow.Record) error {
	vld := make(map[string]bool)
	for ind, clm := range header.Schema().Fields() {
		if len(clm.Name) == 0 {
			err := fmt.Errorf("header row has empty value in field %d", ind)
			return err
		}

		_, ok := vld[clm.Name]
		if ok {
			err := fmt.Errorf("header row has duplicate value in field %d", ind)
			return err
		}
		vld[clm.Name] = true
	}

	return nil
}
