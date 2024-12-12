package client

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func (c *Client) migrate(sch *arrow.Schema, msg chan<- *pb.Message) error {
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sch)
	new, err := pb.NewMessage(pb.Kind_KIND_MIGRATE, builder.NewRecord())
	if err != nil {
		return err
	}
	msg <- new
	return nil
}
