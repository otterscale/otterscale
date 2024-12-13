package client

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
)

func (c *Client) insert(ctx context.Context, rec arrow.Record) error {
	fmt.Println(rec)
	return nil
}
