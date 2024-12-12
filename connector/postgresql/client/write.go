package client

import (
	"context"
	"fmt"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/internal/connector"
)

func (c *Client) Write(ctx context.Context, msg <-chan *pb.Message, opts connector.WriteOptions) error {
	for {
		select {
		case msg, ok := <-msg:
			fmt.Printf("[%v] %+v", ok, msg)
			if !ok {
				continue // ?
			}
		}
	}
	// fmt.Println(rec)
	// kind := connector.WriteKindInsert
	// switch kind {
	// case connector.WriteKindInsert:
	// case connector.WriteKindUpdate:
	// case connector.WriteKindUpsert:
	// case connector.WriteKindDelete:
	// case connector.WriteKindMigrate:
	// default:
	// 	return fmt.Errorf("not supported kind %v", kind)
	// }
	// return nil
}
