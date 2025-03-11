package client

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

type ReadmeSchema struct {
	ownerBuilder  *array.StringBuilder
	repoBuilder   *array.StringBuilder
	readmeBuilder *array.StringBuilder
}

func NewReadmeSchema() (*ReadmeSchema, *array.RecordBuilder) {
	// create schema
	sch := arrow.NewSchema(
		[]arrow.Field{
			{Name: "owner", Type: arrow.BinaryTypes.String},
			{Name: "repo", Type: arrow.BinaryTypes.String},
			{Name: "readme", Type: arrow.BinaryTypes.String},
		}, toSchemaMetadata("github_readme"),
	)

	// record builder
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sch)

	// append builder
	ownerBuilder := builder.Field(0).(*array.StringBuilder)
	repoBuilder := builder.Field(1).(*array.StringBuilder)
	readmeBuilder := builder.Field(2).(*array.StringBuilder)

	return &ReadmeSchema{
		ownerBuilder:  ownerBuilder,
		repoBuilder:   repoBuilder,
		readmeBuilder: readmeBuilder,
	}, builder
}

func (i *ReadmeSchema) Append(c *Client, content string) {
	i.ownerBuilder.Append(c.opts.owner)
	i.repoBuilder.Append(c.opts.repo)
	i.readmeBuilder.Append(content)
}
