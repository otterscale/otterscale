package tables

import (
	"testing"

	"github.com/andygrunwald/go-jira"
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
)

func TestNewProjectSchema_SchemaCreationWithCorrectFields(t *testing.T) {
	// Call the function to get the actual schema
	projectSchema := NewProjectSchema()
	actualSchema := projectSchema.builder.Schema()

	// Define expected fields
	expectedFields := []arrow.Field{
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
	}

	// Verify the number of fields
	if len(actualSchema.Fields()) != len(expectedFields) {
		t.Errorf("Expected %d fields, got %d", len(expectedFields), len(actualSchema.Fields()))
	}

	// Verify each field
	for i, expectedField := range expectedFields {
		actualField := actualSchema.Field(i)
		if actualField.Name != expectedField.Name {
			t.Errorf("Field %d: expected name %q, got %q", i, expectedField.Name, actualField.Name)
		}
		if !arrow.TypeEqual(actualField.Type, expectedField.Type) {
			t.Errorf("Field %d (%s): expected type %v, got %v",
				i, expectedField.Name, expectedField.Type, actualField.Type)
		}
	}

	// Verify metadata (if needed)
	expectedMetadataKey := "table_name"
	expectedMetadataValue := "projects"
	if actualSchema.Metadata().FindKey(expectedMetadataKey) == -1 {
		t.Errorf("Expected metadata key %q not found", expectedMetadataKey)
	} else {
		metadataValue, ok := actualSchema.Metadata().GetValue(expectedMetadataKey)
		if !ok {
			t.Errorf("Expected to find value for metadata key %q", expectedMetadataKey)
		} else if metadataValue != expectedMetadataValue {
			t.Errorf("Expected metadata value %q for key %q, got %q",
				expectedMetadataValue, expectedMetadataKey, metadataValue)
		}
	}
}

func TestProjectSchema_Record_ReturnsValidArrowRecord(t *testing.T) {
	// Initialize ProjectSchema
	projectSchema := NewProjectSchema()

	// Setup test data
	project := &jira.Project{
		ID:          "123",
		Key:         "TEST",
		URL:         "https://example.com/projects/TEST",
		Name:        "Test Project",
		Description: "Test Description",
	}
	projectTypeKey := "software"

	// Append data to the schema
	projectSchema.Append(project, projectTypeKey)

	// Call the Record function
	record := projectSchema.Record()

	// Verify the record is not nil
	if record == nil {
		t.Fatal("Expected non-nil record, got nil")
	}

	// Verify the record has the correct number of columns
	expectedCols := 17 // Number of fields in the schema
	actualCols := int(record.NumCols())
	if actualCols != expectedCols {
		t.Errorf("Expected record to have %d columns, got %d", expectedCols, actualCols)
	}

	// Verify the record has 1 row
	expectedRows := 1
	actualRows := int(record.NumRows())
	if actualRows != expectedRows {
		t.Errorf("Expected record to have %d rows, got %d", expectedRows, actualRows)
	}

	// Verify the schema of the record matches our ProjectSchema
	if !record.Schema().Equal(projectSchema.builder.Schema()) {
		t.Error("Record schema does not match the ProjectSchema")
	}

	// Verify some specific values in the record
	idCol := record.Column(0).(*array.String)
	if idCol.Value(0) != project.ID {
		t.Errorf("Expected ID field to be %q, got %q", project.ID, idCol.Value(0))
	}

	keyCol := record.Column(1).(*array.String)
	if keyCol.Value(0) != project.Key {
		t.Errorf("Expected Key field to be %q, got %q", project.Key, keyCol.Value(0))
	}
}
func TestProjectSchema_Append_WithAllFieldsPopulated(t *testing.T) {
	// Setup test data
	project := &jira.Project{
		ID:       "123",
		Key:      "TEST",
		URL:      "https://example.com/projects/TEST",
		Lead:     jira.User{Name: "John Doe"},
		Name:     "Test Project",
		Self:     "https://example.com/rest/api/2/project/TEST",
		Email:    "test@example.com",
		Roles:    map[string]string{"admin": "https://example.com/rest/api/2/project/TEST/role/10000"},
		Expand:   "description,lead,url,projectKeys",
		Versions: []jira.Version{{Name: "1.0"}},
		// AvatarUrls:      jira.AvatarUrls{Small: "https://example.com/avatar.png"},
		Components:      []jira.ProjectComponent{{Name: "Backend"}},
		IssueTypes:      []jira.IssueType{{Name: "Bug"}},
		Description:     "This is a test project",
		AssigneeType:    "PROJECT_LEAD",
		ProjectCategory: jira.ProjectCategory{Name: "Software"},
	}
	projectTypeKey := "software"

	// Initialize ProjectSchema
	projectSchema := NewProjectSchema()

	// Execute the function
	projectSchema.Append(project, projectTypeKey)

	// Verify the results
	if projectSchema.idBuilder.Len() != 1 {
		t.Errorf("Expected idBuilder to have 1 value, got %d", projectSchema.idBuilder.Len())
	} else if projectSchema.idBuilder.Value(0) != project.ID {
		t.Errorf("Expected id to be %q, got %q", project.ID, projectSchema.idBuilder.Value(0))
	}

	if projectSchema.keyBuilder.Len() != 1 {
		t.Errorf("Expected keyBuilder to have 1 value, got %d", projectSchema.keyBuilder.Len())
	} else if projectSchema.keyBuilder.Value(0) != project.Key {
		t.Errorf("Expected key to be %q, got %q", project.Key, projectSchema.keyBuilder.Value(0))
	}

	if projectSchema.urlBuilder.Len() != 1 {
		t.Errorf("Expected urlBuilder to have 1 value, got %d", projectSchema.urlBuilder.Len())
	} else if projectSchema.urlBuilder.Value(0) != project.URL {
		t.Errorf("Expected url to be %q, got %q", project.URL, projectSchema.urlBuilder.Value(0))
	}

	if projectSchema.nameBuilder.Len() != 1 {
		t.Errorf("Expected nameBuilder to have 1 value, got %d", projectSchema.nameBuilder.Len())
	} else if projectSchema.nameBuilder.Value(0) != project.Name {
		t.Errorf("Expected name to be %q, got %q", project.Name, projectSchema.nameBuilder.Value(0))
	}

	if projectSchema.selfBuilder.Len() != 1 {
		t.Errorf("Expected selfBuilder to have 1 value, got %d", projectSchema.selfBuilder.Len())
	} else if projectSchema.selfBuilder.Value(0) != project.Self {
		t.Errorf("Expected self to be %q, got %q", project.Self, projectSchema.selfBuilder.Value(0))
	}

	if projectSchema.emailBuilder.Len() != 1 {
		t.Errorf("Expected emailBuilder to have 1 value, got %d", projectSchema.emailBuilder.Len())
	} else if projectSchema.emailBuilder.Value(0) != project.Email {
		t.Errorf("Expected email to be %q, got %q", project.Email, projectSchema.emailBuilder.Value(0))
	}

	if projectSchema.expandBuilder.Len() != 1 {
		t.Errorf("Expected expandBuilder to have 1 value, got %d", projectSchema.expandBuilder.Len())
	} else if projectSchema.expandBuilder.Value(0) != project.Expand {
		t.Errorf("Expected expand to be %q, got %q", project.Expand, projectSchema.expandBuilder.Value(0))
	}

	if projectSchema.descriptionBuilder.Len() != 1 {
		t.Errorf("Expected descriptionBuilder to have 1 value, got %d", projectSchema.descriptionBuilder.Len())
	} else if projectSchema.descriptionBuilder.Value(0) != project.Description {
		t.Errorf("Expected description to be %q, got %q", project.Description, projectSchema.descriptionBuilder.Value(0))
	}

	if projectSchema.assigneetypeBuilder.Len() != 1 {
		t.Errorf("Expected assigneetypeBuilder to have 1 value, got %d", projectSchema.assigneetypeBuilder.Len())
	} else if projectSchema.assigneetypeBuilder.Value(0) != project.AssigneeType {
		t.Errorf("Expected assigneetype to be %q, got %q", project.AssigneeType, projectSchema.assigneetypeBuilder.Value(0))
	}

	if projectSchema.projecttypekeyBuilder.Len() != 1 {
		t.Errorf("Expected projecttypekeyBuilder to have 1 value, got %d", projectSchema.projecttypekeyBuilder.Len())
	} else if projectSchema.projecttypekeyBuilder.Value(0) != projectTypeKey {
		t.Errorf("Expected projecttypekey to be %q, got %q", projectTypeKey, projectSchema.projecttypekeyBuilder.Value(0))
	}

	jsonFields := []struct {
		builder *array.StringBuilder
		name    string
	}{
		{projectSchema.leadBuilder, "lead"},
		{projectSchema.rolesBuilder, "roles"},
		{projectSchema.versionsBuilder, "versions"},
		{projectSchema.avatarurlsBuilder, "avatarurls"},
		{projectSchema.componentsBuilder, "components"},
		{projectSchema.issuetypesBuilder, "issuetypes"},
		{projectSchema.projectcategoryBuilder, "projectcategory"},
	}

	for _, field := range jsonFields {
		if field.builder.Len() != 1 {
			t.Errorf("Expected %s builder to have 1 value, got %d", field.name, field.builder.Len())
		} else if field.builder.Value(0) == "" {
			t.Errorf("Expected %s to be non-empty JSON string", field.name)
		}
	}
}
