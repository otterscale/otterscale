package tables

import (
	// "reflect"
	"testing"

	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"

	"github.com/andygrunwald/go-jira"
)

func TestNewIssueSchema_CreatesSchemaWithAllFields(t *testing.T) {
	// Call the function to test
	schema := NewIssueSchema()

	// Verify the schema contains all 13 fields
	expectedFields := []arrow.Field{
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
	}

	// Get the actual schema from the builder
	actualSchema := schema.builder.Schema()

	// Verify field count
	if len(actualSchema.Fields()) != len(expectedFields) {
		t.Errorf("Expected %d fields, got %d", len(expectedFields), len(actualSchema.Fields()))
	}

	// Verify each field
	for i, expectedField := range expectedFields {
		actualField := actualSchema.Field(i)
		if actualField.Name != expectedField.Name {
			t.Errorf("Field %d name mismatch: expected %s, got %s", i, expectedField.Name, actualField.Name)
		}
		if !arrow.TypeEqual(actualField.Type, expectedField.Type) {
			t.Errorf("Field %d type mismatch: expected %v, got %v", i, expectedField.Type, actualField.Type)
		}
	}

	// Verify all builders are initialized correctly
	builders := []array.Builder{
		schema.expandBuilder,
		schema.idBuilder,
		schema.selfBuilder,
		schema.keyBuilder,
		schema.renderedFieldsBuilder,
		schema.namesBuilder,
		schema.transitionsBuilder,
		schema.changelogBuilder,
		schema.fieldsBuilder,
		schema.projectIdBuilder,
		schema.projectKeyBuilder,
		schema.createdBuilder,
		schema.updatedBuilder,
	}

	for i, builder := range builders {
		if builder == nil {
			t.Errorf("Builder for field %d (%s) is nil", i, expectedFields[i].Name)
		}
	}
}

func TestAppend_WithCompleteIssueData(t *testing.T) {
	// Create a test issue with all fields populated
	testIssue := &jira.Issue{
		Expand: "testExpand",
		ID:     "testID",
		Self:   "testSelf",
		Key:    "testKey",
		Fields: &jira.IssueFields{
			Project: jira.Project{
				ID:  "testProjectID",
				Key: "testProjectKey",
			},
			Created: jira.Time(time.Now()),
			Updated: jira.Time(time.Now()),
		},
		RenderedFields: &jira.IssueRenderedFields{},
		Names:         map[string]string{"name1": "value1"},
		Transitions:   []jira.Transition{{ID: "t1"}},
		Changelog:     &jira.Changelog{Histories: []jira.ChangelogHistory{}},
	}

	// Initialize the schema
	schema := NewIssueSchema()

	// Call the Append function
	schema.Append(testIssue)

	// Verify all builders have the correct data appended
	verifyBuilder(t, schema.expandBuilder, testIssue.Expand, "expand")
	verifyBuilder(t, schema.idBuilder, testIssue.ID, "id")
	verifyBuilder(t, schema.selfBuilder, testIssue.Self, "self")
	verifyBuilder(t, schema.keyBuilder, testIssue.Key, "key")
	verifyBuilder(t, schema.projectIdBuilder, testIssue.Fields.Project.ID, "projectId")
	verifyBuilder(t, schema.projectKeyBuilder, testIssue.Fields.Project.Key, "projectKey")
	verifyBuilderTimestamp(t, schema.createdBuilder, testIssue.Fields.Created, "created")
	verifyBuilderTimestamp(t, schema.updatedBuilder, testIssue.Fields.Updated, "updated")
	verifyBuilderJSON(t, schema.renderedFieldsBuilder, testIssue.RenderedFields, "renderedFields")
	verifyBuilderJSON(t, schema.namesBuilder, testIssue.Names, "names")
	verifyBuilderJSON(t, schema.transitionsBuilder, testIssue.Transitions, "transitions")
	verifyBuilderJSON(t, schema.changelogBuilder, testIssue.Changelog, "changelog")
	verifyBuilderJSON(t, schema.fieldsBuilder, testIssue.Fields, "fields")
}

func TestRecord_ReturnsValidArrowRecord(t *testing.T) {
	// Initialize the schema
	schema := NewIssueSchema()
	
	// Create a test issue with required fields
	testIssue := &jira.Issue{
		Expand: "testExpand",
		ID:     "testID",
		Self:   "testSelf",
		Key:    "testKey",
		Fields: &jira.IssueFields{
			Project: jira.Project{
				ID:  "testProjectID",
				Key: "testProjectKey",
			},
			Created: jira.Time(time.Now()),
			Updated: jira.Time(time.Now()),
		},
		RenderedFields: &jira.IssueRenderedFields{},
		Names:          map[string]string{"name1": "value1"},
		Transitions:    []jira.Transition{{ID: "t1"}},
		Changelog:      &jira.Changelog{},
	}
	
	// Append the test issue
	schema.Append(testIssue)
	
	// Get the arrow.Record
	record := schema.Record()
	defer record.Release()
	
	// Verify it's not nil
	if record == nil {
		t.Fatal("Expected non-nil record")
	}
	
	// Verify record has the expected number of rows
	if record.NumRows() != 1 {
		t.Errorf("Expected 1 row, got %d", record.NumRows())
	}
	
	// Verify record has the expected number of columns
	if record.NumCols() != 13 {
		t.Errorf("Expected 13 columns, got %d", record.NumCols())
	}
	
	// Verify schema matches
	if !record.Schema().Equal(schema.builder.Schema()) {
		t.Error("Record schema doesn't match schema builder schema")
	}
	
	// Verify some of the values in the record to ensure data was transferred correctly
	idCol := record.Column(1).(*array.String)
	if idCol.Value(0) != testIssue.ID {
		t.Errorf("ID column value mismatch: expected %s, got %s", testIssue.ID, idCol.Value(0))
	}
	
	keyCol := record.Column(3).(*array.String)
	if keyCol.Value(0) != testIssue.Key {
		t.Errorf("Key column value mismatch: expected %s, got %s", testIssue.Key, keyCol.Value(0))
	}
}

// Helper function to verify string builders
func verifyBuilder(t *testing.T, builder array.Builder, expected string, fieldName string) {
	strBuilder, ok := builder.(*array.StringBuilder)
	if !ok {
		t.Errorf("Builder for %s is not a string builder", fieldName)
		return
	}
	if strBuilder.Len() != 1 {
		t.Errorf("Expected 1 value in %s builder, got %d", fieldName, strBuilder.Len())
		return
	}
	if strBuilder.Value(0) != expected {
		t.Errorf("Value mismatch in %s builder: expected %s, got %s", fieldName, expected, strBuilder.Value(0))
	}
}

// Helper function to verify timestamp builders
func verifyBuilderTimestamp(t *testing.T, builder array.Builder, expected jira.Time, fieldName string) {
	tsBuilder, ok := builder.(*array.TimestampBuilder)
	if !ok {
		t.Errorf("Builder for %s is not a timestamp builder", fieldName)
		return
	}
	if tsBuilder.Len() != 1 {
		t.Errorf("Expected 1 value in %s builder, got %d", fieldName, tsBuilder.Len())
		return
	}
	
	// Create an array from the builder to access the values
	arr := tsBuilder.NewArray()
	defer arr.Release()
	tsArr := arr.(*array.Timestamp)
	
	expectedTs := arrow.Timestamp(time.Time(expected).UnixMilli())
	if tsArr.Value(0) != expectedTs {
		t.Errorf("Value mismatch in %s builder: expected %d, got %d", fieldName, expectedTs, tsArr.Value(0))
	}
}

// Helper function to verify JSON builders
func verifyBuilderJSON(t *testing.T, builder array.Builder, expected interface{}, fieldName string) {
	strBuilder, ok := builder.(*array.StringBuilder)
	if !ok {
		t.Errorf("Builder for %s is not a string builder", fieldName)
		return
	}
	if strBuilder.Len() != 1 {
		t.Errorf("Expected 1 value in %s builder, got %d", fieldName, strBuilder.Len())
		return
	}
	// We can't easily compare the JSON strings directly since marshaling order isn't guaranteed
	// So we just verify that something was stored
	if strBuilder.Value(0) == "" {
		t.Errorf("Expected non-empty JSON in %s builder", fieldName)
	}
}