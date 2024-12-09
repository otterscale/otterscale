package client

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/internal/connector"
)

func (c *Client) Write(ctx context.Context, rec <-chan arrow.Record, opts ...connector.WriteOption) error {
	for {
		select {
		case msg, ok := <-rec:
			if !ok {
				continue // ?
			}
			fmt.Printf("%+v", msg)
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
