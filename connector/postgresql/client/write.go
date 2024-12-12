package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/internal/connector"
	"github.com/openhdc/openhdc/internal/metadata"
)

func (c *Client) Write(ctx context.Context, msg <-chan *pb.Message, opts connector.WriteOptions) error {
	// TODO: FAKE
	namespace := "public"

	tables, err := c.GetTables(ctx, namespace)
	if err != nil {
		return err
	}

	for {
		select {
		case msg, ok := <-msg:
			if !ok {
				continue // ?
			}
			rec, err := pb.ToArrowRecord(msg.Record)
			if err != nil {
				return err
			}
			switch msg.Kind {
			case pb.Kind_KIND_MIGRATE:
				if err := c.migrate(ctx, tables, rec.Schema()); err != nil {
					return err
				}

			case pb.Kind_KIND_INSERT, pb.Kind_KIND_DELETE:

			default:
				return fmt.Errorf("not supported kind %v", msg.Kind)
			}
		}
	}
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

func isTableExists(tabs []arrow.Table, new string) bool {
	for _, tab := range tabs {
		current, _ := tab.Schema().Metadata().GetValue(metadata.KeySchemaTableName)
		if current == new {
			return true
		}
	}
	return false
}

func (c *Client) migrate(ctx context.Context, tabs []arrow.Table, sch *arrow.Schema) error {
	tableName, ok := sch.Metadata().GetValue(metadata.KeySchemaTableName)
	if !ok {
		return errors.New("table name not found")
	}

	if !isTableExists(tabs, tableName) {
		return createTableIfNotExists(ctx, c.pool, sch)
	}
	// TODO: check migration
	panic(errors.New("check migration"))
	if true {
		return nil
	}
	if err := dropTable(ctx, c.pool, sch); err != nil {
		return err
	}
	return createTableIfNotExists(ctx, c.pool, sch)
}
