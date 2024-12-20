package client

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
)

func (c *Client) validateHeader(header arrow.Record) error {
	vld := make(map[string]bool)
	for ind, clm := range header.Schema().Fields() {
		if clm.Name == "" {
			return fmt.Errorf("header row has empty value in field %d", ind)
		}
		if _, ok := vld[clm.Name]; ok {
			return fmt.Errorf("header row has duplicate value in field %d", ind)
		}
		vld[clm.Name] = true
	}
	return nil
}
