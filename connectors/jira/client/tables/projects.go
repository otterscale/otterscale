package tables

import (
	jira "github.com/andygrunwald/go-jira"
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

type ProjectSchema struct {
	builder                *array.RecordBuilder
	idBuilder              *array.StringBuilder
	keyBuilder             *array.StringBuilder
	urlBuilder             *array.StringBuilder
	leadBuilder            *array.StringBuilder //Json
	nameBuilder            *array.StringBuilder
	selfBuilder            *array.StringBuilder
	emailBuilder           *array.StringBuilder
	rolesBuilder           *array.StringBuilder //Json
	expandBuilder          *array.StringBuilder
	versionsBuilder        *array.StringBuilder //Json
	avatarurlsBuilder      *array.StringBuilder //Json
	componentsBuilder      *array.StringBuilder //Json
	issuetypesBuilder      *array.StringBuilder //Json
	descriptionBuilder     *array.StringBuilder
	assigneetypeBuilder    *array.StringBuilder
	projecttypekeyBuilder  *array.StringBuilder
	projectcategoryBuilder *array.StringBuilder //Json
}

func NewProjectSchema() *ProjectSchema {
	// create schema
	sch := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.BinaryTypes.String},
			{Name: "key", Type: arrow.BinaryTypes.String},
			{Name: "url", Type: arrow.BinaryTypes.String},
			{Name: "lead", Type: arrow.BinaryTypes.String},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "self", Type: arrow.BinaryTypes.String},
			{Name: "email", Type: arrow.BinaryTypes.String},
			{Name: "roles", Type: arrow.BinaryTypes.String},
			{Name: "expand", Type: arrow.BinaryTypes.String},
			{Name: "versions", Type: arrow.BinaryTypes.String},
			{Name: "avatarurls", Type: arrow.BinaryTypes.String},
			{Name: "components", Type: arrow.BinaryTypes.String},
			{Name: "issuetypes", Type: arrow.BinaryTypes.String},
			{Name: "description", Type: arrow.BinaryTypes.String},
			{Name: "assigneetype", Type: arrow.BinaryTypes.String},
			{Name: "projecttypekey", Type: arrow.BinaryTypes.String},
			{Name: "projectcategory", Type: arrow.BinaryTypes.String},
		}, toSchemaMetadata("projects"),
	)

	// record builder
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sch)

	// append builder
	idBuilder := builder.Field(0).(*array.StringBuilder)
	keyBuilder := builder.Field(1).(*array.StringBuilder)
	urlBuilder := builder.Field(2).(*array.StringBuilder)
	leadBuilder := builder.Field(3).(*array.StringBuilder)
	nameBuilder := builder.Field(4).(*array.StringBuilder)
	selfBuilder := builder.Field(5).(*array.StringBuilder)
	emailBuilder := builder.Field(6).(*array.StringBuilder)
	rolesBuilder := builder.Field(7).(*array.StringBuilder)
	expandBuilder := builder.Field(8).(*array.StringBuilder)
	versionsBuilder := builder.Field(9).(*array.StringBuilder)
	avatarurlsBuilder := builder.Field(10).(*array.StringBuilder)
	componentsBuilder := builder.Field(11).(*array.StringBuilder)
	issuetypesBuilder := builder.Field(12).(*array.StringBuilder)
	descriptionBuilder := builder.Field(13).(*array.StringBuilder)
	assigneetypeBuilder := builder.Field(14).(*array.StringBuilder)
	projecttypekeyBuilder := builder.Field(15).(*array.StringBuilder)
	projectcategoryBuilder := builder.Field(16).(*array.StringBuilder)

	return &ProjectSchema{
		builder:                builder,
		idBuilder:              idBuilder,
		keyBuilder:             keyBuilder,
		urlBuilder:             urlBuilder,
		leadBuilder:            leadBuilder,
		nameBuilder:            nameBuilder,
		selfBuilder:            selfBuilder,
		emailBuilder:           emailBuilder,
		rolesBuilder:           rolesBuilder,
		expandBuilder:          expandBuilder,
		versionsBuilder:        versionsBuilder,
		avatarurlsBuilder:      avatarurlsBuilder,
		componentsBuilder:      componentsBuilder,
		issuetypesBuilder:      issuetypesBuilder,
		descriptionBuilder:     descriptionBuilder,
		assigneetypeBuilder:    assigneetypeBuilder,
		projecttypekeyBuilder:  projecttypekeyBuilder,
		projectcategoryBuilder: projectcategoryBuilder,
	}
}

func (i *ProjectSchema) Append(project *jira.Project, ProjectTypeKey string) {
	i.idBuilder.Append(project.ID)
	i.keyBuilder.Append(project.Key)
	i.urlBuilder.Append(project.URL)
	builderAppendJson(i.leadBuilder, project.Lead)
	i.nameBuilder.Append(project.Name)
	i.selfBuilder.Append(project.Self)
	i.emailBuilder.Append(project.Email)
	builderAppendJson(i.rolesBuilder, project.Roles)
	i.expandBuilder.Append(project.Expand)
	builderAppendJson(i.versionsBuilder, project.Versions)
	builderAppendJson(i.avatarurlsBuilder, project.AvatarUrls)
	builderAppendJson(i.componentsBuilder, project.Components)
	builderAppendJson(i.issuetypesBuilder, project.IssueTypes)
	i.descriptionBuilder.Append(project.Description)
	i.assigneetypeBuilder.Append(project.AssigneeType)
	i.projecttypekeyBuilder.Append(ProjectTypeKey)
	builderAppendJson(i.projectcategoryBuilder, project.ProjectCategory)
}

func (i *ProjectSchema) Record() arrow.Record {
	return i.builder.NewRecord()
}
