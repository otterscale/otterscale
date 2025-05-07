package tables

import (
	"reflect"
	"testing"

	"github.com/andygrunwald/go-jira"
	"github.com/apache/arrow-go/v18/arrow"
)

func TestNewIssueFieldSchema_BasicSchemaCreation(t *testing.T) {
	// Call the function to test
	got := NewIssueFieldSchema()

	// Verify the schema fields
	expectedFields := []arrow.Field{
		{Name: "id", Type: arrow.BinaryTypes.String},
		{Name: "key", Type: arrow.BinaryTypes.String},
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "custom", Type: arrow.FixedWidthTypes.Boolean},
		{Name: "navigable", Type: arrow.FixedWidthTypes.Boolean},
		{Name: "searchable", Type: arrow.FixedWidthTypes.Boolean},
		{Name: "clausenames", Type: arrow.BinaryTypes.String},
		{Name: "schema", Type: arrow.BinaryTypes.String},
	}

	// Check the schema in the builder
	if !reflect.DeepEqual(got.builder.Schema().Fields(), expectedFields) {
		t.Errorf("NewIssueFieldSchema() created schema fields = %v, want %v",
			got.builder.Schema().Fields(), expectedFields)
	}

	// Verify builder types are as expected
	// Using type names for validation instead of type assertions
	if got.idBuilder == nil {
		t.Error("idBuilder is nil")
	}
	if got.keyBuilder == nil {
		t.Error("keyBuilder is nil")
	}
	if got.nameBuilder == nil {
		t.Error("nameBuilder is nil")
	}
	if got.customBuilder == nil {
		t.Error("customBuilder is nil")
	}
	if got.navigableBuilder == nil {
		t.Error("navigableBuilder is nil")
	}
	if got.searchableBuilder == nil {
		t.Error("searchableBuilder is nil")
	}
	if got.clausenamesBuilder == nil {
		t.Error("clausenamesBuilder is nil")
	}
	if got.schemaBuilder == nil {
		t.Error("schemaBuilder is nil")
	}

	// Verify the record builder is initialized with default allocator
	// if got.builder.Allocator() != memory.DefaultAllocator {
	// 	t.Error("builder was not initialized with default allocator")
	// }
}

func TestIssueFieldSchema_AppendWithCompleteFieldData(t *testing.T) {
	// Setup
	schema := NewIssueFieldSchema()
	field := &jira.Field{
		ID:          "1",
		Key:         "KEY",
		Name:        "Name",
		Custom:      true,
		Navigable:   true,
		Searchable:  true,
		ClauseNames: []string{"Clause1", "Clause2"},
		Schema: jira.FieldSchema{
			Type: "string",
		},
	}

	// Execute
	schema.Append(field)

	// Verify
	if schema.idBuilder.Len() != 1 {
		t.Errorf("idBuilder should have 1 value, got %d", schema.idBuilder.Len())
	}
	if schema.idBuilder.Value(0) != "1" {
		t.Errorf("idBuilder should have correct value, got %s", schema.idBuilder.Value(0))
	}

	if schema.keyBuilder.Len() != 1 {
		t.Errorf("keyBuilder should have 1 value, got %d", schema.keyBuilder.Len())
	}
	if schema.keyBuilder.Value(0) != "KEY" {
		t.Errorf("keyBuilder should have correct value, got %s", schema.keyBuilder.Value(0))
	}

	if schema.nameBuilder.Len() != 1 {
		t.Errorf("nameBuilder should have 1 value, got %d", schema.nameBuilder.Len())
	}
	if schema.nameBuilder.Value(0) != "Name" {
		t.Errorf("nameBuilder should have correct value, got %s", schema.nameBuilder.Value(0))
	}

	if schema.customBuilder.Len() != 1 {
		t.Errorf("customBuilder should have 1 value, got %d", schema.customBuilder.Len())
	}
	if schema.customBuilder.Value(0) != true {
		t.Errorf("customBuilder should have correct value, got %v", schema.customBuilder.Value(0))
	}

	if schema.navigableBuilder.Len() != 1 {
		t.Errorf("navigableBuilder should have 1 value, got %d", schema.navigableBuilder.Len())
	}
	if schema.navigableBuilder.Value(0) != true {
		t.Errorf("navigableBuilder should have correct value, got %v", schema.navigableBuilder.Value(0))
	}

	if schema.searchableBuilder.Len() != 1 {
		t.Errorf("searchableBuilder should have 1 value, got %d", schema.searchableBuilder.Len())
	}
	if schema.searchableBuilder.Value(0) != true {
		t.Errorf("searchableBuilder should have correct value, got %v", schema.searchableBuilder.Value(0))
	}

	if schema.clausenamesBuilder.Len() != 1 {
		t.Errorf("clausenamesBuilder should have 1 value, got %d", schema.clausenamesBuilder.Len())
	}
	if schema.clausenamesBuilder.Value(0) != "Clause1,Clause2" {
		t.Errorf("clausenamesBuilder should have correct value, got %s", schema.clausenamesBuilder.Value(0))
	}

	if schema.schemaBuilder.Len() != 1 {
		t.Errorf("schemaBuilder should have 1 value, got %d", schema.schemaBuilder.Len())
	}
	if schema.schemaBuilder.Value(0) == "" {
		t.Errorf("schemaBuilder should have non-empty value")
	}

	// Verify
	if schema.idBuilder.Len() != 1 {
		t.Errorf("idBuilder should have 1 value, got %d", schema.idBuilder.Len())
	}
	if schema.idBuilder.Value(0) != "1" {
		t.Errorf("idBuilder should have correct value, got %s", schema.idBuilder.Value(0))
	}

	if schema.keyBuilder.Len() != 1 {
		t.Errorf("keyBuilder should have 1 value, got %d", schema.keyBuilder.Len())
	}
	if schema.keyBuilder.Value(0) != "KEY" {
		t.Errorf("keyBuilder should have correct value, got %s", schema.keyBuilder.Value(0))
	}

	if schema.nameBuilder.Len() != 1 {
		t.Errorf("nameBuilder should have 1 value, got %d", schema.nameBuilder.Len())
	}
	if schema.nameBuilder.Value(0) != "Name" {
		t.Errorf("nameBuilder should have correct value, got %s", schema.nameBuilder.Value(0))
	}

	if schema.customBuilder.Len() != 1 {
		t.Errorf("customBuilder should have 1 value, got %d", schema.customBuilder.Len())
	}
	if schema.customBuilder.Value(0) != true {
		t.Errorf("customBuilder should have correct value, got %v", schema.customBuilder.Value(0))
	}

	if schema.navigableBuilder.Len() != 1 {
		t.Errorf("navigableBuilder should have 1 value, got %d", schema.navigableBuilder.Len())
	}
	if schema.navigableBuilder.Value(0) != true {
		t.Errorf("navigableBuilder should have correct value, got %v", schema.navigableBuilder.Value(0))
	}

	if schema.searchableBuilder.Len() != 1 {
		t.Errorf("searchableBuilder should have 1 value, got %d", schema.searchableBuilder.Len())
	}
	if schema.searchableBuilder.Value(0) != true {
		t.Errorf("searchableBuilder should have correct value, got %v", schema.searchableBuilder.Value(0))
	}

	if schema.clausenamesBuilder.Len() != 1 {
		t.Errorf("clausenamesBuilder should have 1 value, got %d", schema.clausenamesBuilder.Len())
	}
	if schema.clausenamesBuilder.Value(0) != "Clause1,Clause2" {
		t.Errorf("clausenamesBuilder should have correct value, got %s", schema.clausenamesBuilder.Value(0))
	}

	if schema.schemaBuilder.Len() != 1 {
		t.Errorf("schemaBuilder should have 1 value, got %d", schema.schemaBuilder.Len())
	}
	if schema.schemaBuilder.Value(0) == "" {
		t.Errorf("schemaBuilder should have non-empty value")
	}
}

func TestIssueFieldSchema_Record_BasicFunctionalityWithValidBuilder(t *testing.T) {
	// Setup - create a properly initialized IssueFieldSchema
	schema := NewIssueFieldSchema()

	// Execute - call the Record method
	record := schema.Record()

	// Verify - check that a valid arrow.Record is returned
	if record == nil {
		t.Error("Record() returned nil, expected a valid arrow.Record")
	}

	// Verify the record has the expected schema fields
	expectedFields := []arrow.Field{
		{Name: "id", Type: arrow.BinaryTypes.String},
		{Name: "key", Type: arrow.BinaryTypes.String},
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "custom", Type: arrow.FixedWidthTypes.Boolean},
		{Name: "navigable", Type: arrow.FixedWidthTypes.Boolean},
		{Name: "searchable", Type: arrow.FixedWidthTypes.Boolean},
		{Name: "clausenames", Type: arrow.BinaryTypes.String},
		{Name: "schema", Type: arrow.BinaryTypes.String},
	}

	if !record.Schema().Equal(arrow.NewSchema(expectedFields, nil)) {
		t.Errorf("Record schema does not match expected schema. Got: %v, Want: %v",
			record.Schema().Fields(), expectedFields)
	}

	// Verify the record has 0 rows (since we didn't append any data)
	if record.NumRows() != 0 {
		t.Errorf("Expected empty record with 0 rows, got %d rows", record.NumRows())
	}
}
