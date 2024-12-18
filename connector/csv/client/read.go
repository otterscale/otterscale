package client

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/csv"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func (c *Client) getHeadRows() (arrow.Record, error) {
	rdr := csv.NewInferringReader(c.file, csv.WithChunk(1))
	defer c.file.Seek(0, 0)

	if !rdr.Next() {
		return nil, fmt.Errorf("file is empty")
	}
	hdr := rdr.Record()

	return hdr, nil
}

func (c *Client) getReader() (*csv.Reader, error) {
	hdr, hRsErr := c.getHeadRows()

	if err := c.validateFile(hRsErr); err != nil {
		return nil, err
	}
	if err := c.validateHeader(hdr); err != nil {
		return nil, err
	}

	opts := []csv.Option{
		csv.WithHeader(true),
		csv.WithChunk(int(c.opts.batchSize)),
	}
	flds := c.toSchemaFields(hdr)
	mtd := c.toSchemaMetadata(c.opts.tableName)
	sch := arrow.NewSchema(flds, mtd)
	rdr := csv.NewReader(c.file, sch, opts...)

	return rdr, rdr.Err()
}

func (c *Client) Read(ctx context.Context, msg chan<- *pb.Message, opts openhdc.ReadOptions) error {
	rdr, err := c.getReader()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer rdr.Release()

	cndFlg := false
	for rdr.Next() {
		rcr := rdr.Record()

		if !cndFlg {
			new, err := pb.NewMessage(pb.Kind_KIND_MIGRATE, rcr)

			if err != nil {
				slog.Error(err.Error())
				return err
			}
			msg <- new

			cndFlg = true
		}

		new, err := pb.NewMessage(pb.Kind_KIND_INSERT, rcr)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		msg <- new
	}

	if err := rdr.Err(); err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
