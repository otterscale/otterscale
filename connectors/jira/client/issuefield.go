package client

import (
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

type IssueFieldSchema struct {
	builder            *array.RecordBuilder
	idBuilder          *array.StringBuilder
	keyBuilder         *array.StringBuilder
	nameBuilder        *array.StringBuilder
	customBuilder      *array.BooleanBuilder
	navigableBuilder   *array.BooleanBuilder
	searchableBuilder  *array.BooleanBuilder
	clausenamesBuilder *array.StringBuilder
	schemaBuilder      *array.StringBuilder
}

func NewIssueFieldSchema() *IssueFieldSchema {
	// create schema
	sch := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.BinaryTypes.String},
			{Name: "key", Type: arrow.BinaryTypes.String},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "custom", Type: arrow.FixedWidthTypes.Boolean},
			{Name: "navigable", Type: arrow.FixedWidthTypes.Boolean},
			{Name: "searchable", Type: arrow.FixedWidthTypes.Boolean},
			{Name: "clausenames", Type: arrow.BinaryTypes.String},
			{Name: "schema", Type: arrow.BinaryTypes.String},
		}, toSchemaMetadata("issue_fields"),
	)

	// record builder
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sch)

	// append builder
	idBuilder := builder.Field(0).(*array.StringBuilder)
	keyBuilder := builder.Field(1).(*array.StringBuilder)
	nameBuilder := builder.Field(2).(*array.StringBuilder)
	customBuilder := builder.Field(3).(*array.BooleanBuilder)
	navigableBuilder := builder.Field(4).(*array.BooleanBuilder)
	searchableBuilder := builder.Field(5).(*array.BooleanBuilder)
	clausenamesBuilder := builder.Field(6).(*array.StringBuilder)
	schemaBuilder := builder.Field(7).(*array.StringBuilder)

	return &IssueFieldSchema{
		builder:            builder,
		idBuilder:          idBuilder,
		keyBuilder:         keyBuilder,
		nameBuilder:        nameBuilder,
		customBuilder:      customBuilder,
		navigableBuilder:   navigableBuilder,
		searchableBuilder:  searchableBuilder,
		clausenamesBuilder: clausenamesBuilder,
		schemaBuilder:      schemaBuilder,
	}
}

func (i *IssueFieldSchema) Append(issueField *jira.Field) {
	i.idBuilder.Append(issueField.ID)
	i.keyBuilder.Append(issueField.Key)
	i.nameBuilder.Append(issueField.Name)
	i.customBuilder.Append(issueField.Custom)
	i.navigableBuilder.Append(issueField.Navigable)
	i.searchableBuilder.Append(issueField.Searchable)
	i.clausenamesBuilder.Append(strings.Join(issueField.ClauseNames, ","))
	builderAppendJson(i.schemaBuilder, issueField.Schema)
}

func (i *IssueFieldSchema) Record() arrow.Record {
	return i.builder.NewRecord()
}
