package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/csv"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func (c *Client) Read(ctx context.Context, msgs chan<- *pb.Message, rdr *openhdc.Reader) error {
	// new csv reader
	r, err := c.newReader()
	if err != nil {
		return err
	}
	defer r.Release()

	// start
	first := true
	for r.Next() {
		rec := r.Record()

		// send schema on first row
		if first {
			if err := rdr.Send(pb.Migrate, msgs, rec); err != nil {
				return err
			}
			first = false
		}

		// send data
		if err := rdr.Send(pb.Insert, msgs, rec); err != nil {
			return err
		}
	}
	return r.Err()
}

func (c *Client) getFields() ([]arrow.Field, error) {
	defer func() {
		if _, err := c.file.Seek(0, io.SeekStart); err != nil {
			slog.Error(err.Error())
		}
	}()

	r := csv.NewInferringReader(c.file, csv.WithChunk(1))
	if !r.Next() {
		return nil, fmt.Errorf("file is empty")
	}

	return r.Record().Schema().Fields(), nil
}

func (c *Client) newReader() (*csv.Reader, error) {
	fs, err := c.getFields()
	if err != nil {
		return nil, err
	}

	if !c.opts.inferring {
		fs = toStringFields(fs)
	}

	if err := validateFields(fs); err != nil {
		return nil, err
	}

	md := toSchemaMetadata(c.opts.tableName)
	sch := arrow.NewSchema(fs, md)
	opts := []csv.Option{
		csv.WithHeader(true),
		csv.WithChunk(int(c.opts.batchSize)),
	}
	return csv.NewReader(c.file, sch, opts...), nil
}
