package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/csv"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/api/property/v1"
)

func (c *Client) getHeadRows() (arrow.Record, error) {
	rdr := csv.NewInferringReader(c.file, csv.WithChunk(1))
	defer func() {
		if _, err := c.file.Seek(0, io.SeekStart); err != nil {
			slog.Error(err.Error())
		}
	}()

	if !rdr.Next() {
		return nil, fmt.Errorf("file is empty")
	}

	return rdr.Record(), nil
}

func (c *Client) getReader() (*csv.Reader, error) {
	hdr, err := c.getHeadRows()
	if errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("file %s is empty", c.opts.filePath)
	}
	if err != nil {
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

	return csv.NewReader(c.file, sch, opts...), nil
}

func (c *Client) Read(ctx context.Context, msg chan<- *pb.Message, opts openhdc.ReadOptions) error {
	syncedAt := time.Now().UTC().Truncate(time.Second)

	rdr, err := c.getReader()
	if err != nil {
		return err
	}
	defer rdr.Release()

	cndFlg := false
	for rdr.Next() {
		rcr := rdr.Record()

		// send schema on first row
		if !cndFlg {
			new, err := pb.NewMessage(property.MessageKind_migrate, rcr, c.opts.name, syncedAt)
			if err != nil {
				return err
			}
			msg <- new

			cndFlg = true
		}

		new, err := pb.NewMessage(property.MessageKind_insert, rcr, c.opts.name, syncedAt)
		if err != nil {
			return err
		}
		msg <- new
	}

	return rdr.Err()
}
