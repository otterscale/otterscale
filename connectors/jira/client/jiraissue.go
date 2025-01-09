package client

import (
	"encoding/json"
	"log"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/openhdc/openhdc/metadata"
)

type IssueSchema struct {
	builder               *array.RecordBuilder
	expandBuilder         *array.StringBuilder
	idBuilder             *array.StringBuilder
	selfBuilder           *array.StringBuilder
	keyBuilder            *array.StringBuilder
	renderedFieldsBuilder *array.StringBuilder
	namesBuilder          *array.StringBuilder
	transitionsBuilder    *array.StringBuilder
	changelogBuilder      *array.StringBuilder
	fieldsBuilder         *array.StringBuilder
	projectIdBuilder      *array.StringBuilder
	projectKeyBuilder     *array.StringBuilder
	createdBuilder        *array.TimestampBuilder
	updatedBuilder        *array.TimestampBuilder
}

func toSchemaMetadata(tableName string) *arrow.Metadata {
	m := map[string]string{}
	metadata.SetTableName(m, tableName)
	md := arrow.MetadataFrom(m)
	return &md
}

func NewIssueSchema() *IssueSchema {
	// create schema
	sch := arrow.NewSchema(
		[]arrow.Field{
			{Name: "expand", Type: arrow.BinaryTypes.String},
			{Name: "id", Type: arrow.BinaryTypes.String},
			{Name: "self", Type: arrow.BinaryTypes.String},
			{Name: "key", Type: arrow.BinaryTypes.String},
			{Name: "renderedFields", Type: arrow.BinaryTypes.String},
			{Name: "names", Type: arrow.BinaryTypes.String},
			{Name: "transitions", Type: arrow.BinaryTypes.String},
			{Name: "changelog", Type: arrow.BinaryTypes.String},
			{Name: "fields", Type: arrow.BinaryTypes.String},
			{Name: "projectId", Type: arrow.BinaryTypes.String},
			{Name: "projectKey", Type: arrow.BinaryTypes.String},
			{Name: "created", Type: arrow.FixedWidthTypes.Timestamp_us},
			{Name: "updated", Type: arrow.FixedWidthTypes.Timestamp_us},
		}, toSchemaMetadata("jiraissue"),
	)

	// record builder
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sch)

	// append builder
	expandBuilder := builder.Field(0).(*array.StringBuilder)
	idBuilder := builder.Field(1).(*array.StringBuilder)
	selfBuilder := builder.Field(2).(*array.StringBuilder)
	keyBuilder := builder.Field(3).(*array.StringBuilder)
	renderedFieldsBuilder := builder.Field(4).(*array.StringBuilder)
	namesBuilder := builder.Field(5).(*array.StringBuilder)
	transitionsBuilder := builder.Field(6).(*array.StringBuilder)
	changelogBuilder := builder.Field(7).(*array.StringBuilder)
	fieldsBuilder := builder.Field(8).(*array.StringBuilder)
	projectIdBuilder := builder.Field(9).(*array.StringBuilder)
	projectKeyBuilder := builder.Field(10).(*array.StringBuilder)
	createdBuilder := builder.Field(11).(*array.TimestampBuilder)
	updatedBuilder := builder.Field(12).(*array.TimestampBuilder)

	return &IssueSchema{
		builder:               builder,
		expandBuilder:         expandBuilder,
		idBuilder:             idBuilder,
		selfBuilder:           selfBuilder,
		keyBuilder:            keyBuilder,
		renderedFieldsBuilder: renderedFieldsBuilder,
		namesBuilder:          namesBuilder,
		transitionsBuilder:    transitionsBuilder,
		changelogBuilder:      changelogBuilder,
		fieldsBuilder:         fieldsBuilder,
		projectIdBuilder:      projectIdBuilder,
		projectKeyBuilder:     projectKeyBuilder,
		createdBuilder:        createdBuilder,
		updatedBuilder:        updatedBuilder,
	}
}

func builderAppendJson(fbuilder *array.StringBuilder, v any) {
	if v != nil {
		// convert Fields into JSON string
		js, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		fbuilder.Append(string(js))
	} else {
		fbuilder.AppendNull()
	}
}

func (i *IssueSchema) Append(issue *jira.Issue) {
	i.expandBuilder.Append(issue.Expand)
	i.idBuilder.Append(issue.ID)
	i.selfBuilder.Append(issue.Self)
	i.keyBuilder.Append(issue.Key)
	i.projectIdBuilder.Append(issue.Fields.Project.ID)
	i.projectKeyBuilder.Append(issue.Fields.Project.Key)
	i.createdBuilder.Append(arrow.Timestamp(time.Time(issue.Fields.Created).UnixMilli()))
	i.updatedBuilder.Append(arrow.Timestamp(time.Time(issue.Fields.Updated).UnixMilli()))
	builderAppendJson(i.renderedFieldsBuilder, issue.RenderedFields)
	builderAppendJson(i.namesBuilder, issue.Names)
	builderAppendJson(i.transitionsBuilder, issue.Transitions)
	builderAppendJson(i.changelogBuilder, issue.Changelog)
	builderAppendJson(i.fieldsBuilder, issue.Fields)
}

func (i *IssueSchema) Record() arrow.Record {
	return i.builder.NewRecord()
}
