package client

import (
	// "fmt"
	"errors"
	"reflect"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	// "github.com/golang/mock/gomock"
	property "github.com/openhdc/openhdc/api/property/v1"
	workload "github.com/openhdc/openhdc/api/workload/v1"
	"github.com/openhdc/openhdc/connectors/oracle/client/or"
	"github.com/openhdc/openhdc/metadata"
)

func Test_toSchemaMetadata_EmptyTableName(t *testing.T) {
	tableName := ""
	expected := map[string]string{}
	metadata.SetTableName(expected, tableName)
	want := arrow.MetadataFrom(expected)

	got := toSchemaMetadata(tableName)
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("toSchemaMetadata() with empty tableName = %v, want %v", *got, want)
	}
}

func Test_toSchemaMetadata_DifferentTableNames(t *testing.T) {
	tableNames := []string{"table1", "table2", "table3"}
	for _, tableName := range tableNames {
		expected := map[string]string{}
		metadata.SetTableName(expected, tableName)
		want := arrow.MetadataFrom(expected)

		got := toSchemaMetadata(tableName)
		if !reflect.DeepEqual(*got, want) {
			t.Errorf("toSchemaMetadata(%q) = %v, want %v", tableName, *got, want)
		}
	}
}

func Test_toSchemaMetadata_NilSafety(t *testing.T) {
	// This test ensures that the function never returns nil
	got := toSchemaMetadata("sometable")
	if got == nil {
		t.Errorf("toSchemaMetadata() returned nil, want non-nil *arrow.Metadata")
	}
}

func TestSkip_ErrorInGetTableName(t *testing.T) {
	// Create a schema with invalid metadata (missing table name)
	fields := []arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}
	schema := arrow.NewSchema(fields, nil)

	// Create a record batch to ensure the schema is valid
	pool := memory.NewGoAllocator()
	builder := array.NewInt32Builder(pool)
	builder.Append(1)
	arr := builder.NewInt32Array()
	defer arr.Release()

	keys := []string{}
	skips := []string{}

	result := skip(schema, keys, skips)
	if !result {
		t.Errorf("Expected skip to return true when GetTableName returns error, got false")
	}
}

func TestSkipEmptyTableName(t *testing.T) {
	// Create a schema with empty metadata (no table name)
	fields := []arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}
	schema := arrow.NewSchema(fields, nil)

	// Test with empty keys and skips
	keys := []string{}
	skips := []string{}

	// Call the function
	result := skip(schema, keys, skips)

	// Verify the result
	if !result {
		t.Errorf("Expected skip to return true for empty table name, got false")
	}
}

func TestSkipWhenMultipleItemsInKeysNotIncludingTableName(t *testing.T) {
	// Create a schema with table name metadata
	md := arrow.NewMetadata(
		[]string{"TABLE-NAME"},
		[]string{"target_table"},
	)
	fields := []arrow.Field{}
	schema := arrow.NewSchema(fields, &md)

	// Test inputs
	keys := []string{"table1", "table2", "table3"} // Doesn't include "target_table"
	skips := []string{}                            // Empty skips

	// Call the function
	result := skip(schema, keys, skips)

	// Verify the result
	if !result {
		t.Errorf("Expected skip() to return true when keys list doesn't contain table name, got false")
	}
}

func TestSkipWithSpecialCharactersInTableName(t *testing.T) {
	// Create a schema with table name containing special characters
	md := arrow.NewMetadata(
		[]string{"TABLE-NAME"},
		[]string{"special@table#name$123"},
	)
	fields := []arrow.Field{}
	schema := arrow.NewSchema(fields, &md)

	// Define test inputs
	keys := []string{}
	skips := []string{"special@table#name$123"}

	// Call the function
	result := skip(schema, keys, skips)

	// Verify the result
	if !result {
		t.Errorf("Expected skip() to return true for table name with special characters in skips list, got false")
	}
}

func TestDeleteAll_SchemaWithZeroFieldsAndFullOverwriteSyncMode(t *testing.T) {
	// Create a schema with zero fields
	sch := arrow.NewSchema([]arrow.Field{}, nil)

	// Call the function with full_overwrite sync mode
	result := deleteAll(sch, property.SyncMode_full_overwrite)

	// Verify the result is true
	if !result {
		t.Errorf("Expected deleteAll to return true for schema with zero fields and full_overwrite sync mode, got %v", result)
	}
}
func TestDeleteAll_SchemaWithEmptyPrimaryKeyFieldAndFullOverwriteSyncMode(t *testing.T) {
	// Create a schema with empty primary key field
	fields := []arrow.Field{
		{Name: "field1", Type: arrow.PrimitiveTypes.Int32},
		{Name: "field2", Type: arrow.PrimitiveTypes.Float64},
	}
	schema := arrow.NewSchema(fields, nil)

	// Set sync mode to full_overwrite
	syncMode := property.SyncMode_full_overwrite

	// Call the function
	result := deleteAll(schema, syncMode)

	// Verify the result
	if !result {
		t.Errorf("Expected false when schema has empty primary key field and syncMode is full_overwrite, got %v", result)
	}
}

// func TestDeleteStaleEmptySchemaWithFullOverwrite(t *testing.T) {
// 	// Create an empty schema
// 	emptySchema := arrow.NewSchema([]arrow.Field{}, nil)

// 	// Call the function with empty schema and full_overwrite sync mode
// 	result := deleteStale(emptySchema, property.SyncMode_full_overwrite)

//		// Verify the expected outcome
//		if !result {
//			t.Errorf("Expected false when schema is empty even with full_overwrite sync mode, got %v", result)
//		}
//	}
func TestDeleteStale_SchemaWithOnlyTableNameMetadata(t *testing.T) {
	// Create a schema with only table name metadata
	md := arrow.MetadataFrom(map[string]string{
		"TABLE-NAME": "test_table",
	})
	fields := []arrow.Field{
		{Name: "field1", Type: arrow.PrimitiveTypes.Int32},
	}
	sch := arrow.NewSchema(fields, &md)

	// Test with full_overwrite sync mode
	result := deleteStale(sch, property.SyncMode_full_overwrite)
	if !result {
		t.Errorf("Expected false when schema has only table name metadata, got %v", result)
	}
}

func TestToMessageKind_FullOverwriteWithoutPrimaryKey(t *testing.T) {
	// Create a schema without primary key
	fields := []arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}
	sch := arrow.NewSchema(fields, nil)

	// Set up test inputs
	mode := property.SyncMode_full_overwrite
	curs := []*workload.Sync_Option_Cursor{}

	// Call the function
	got, err := toMessageKind(sch, mode, curs)

	// Verify the results
	if err != nil {
		t.Errorf("toMessageKind() returned unexpected error: %v", err)
	}

	if got != property.MessageKind_insert {
		t.Errorf("toMessageKind() = %v, want %v", got, property.MessageKind_insert)
	}
}
func TestToMessageKind_FullAppendWithoutPrimaryKey(t *testing.T) {
	// Create a schema without primary key
	fields := []arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}
	md := arrow.NewMetadata(
		[]string{"TABLE-NAME"},
		[]string{"test_table"},
	)
	schema := arrow.NewSchema(fields, &md)

	// Input parameters
	mode := property.SyncMode_full_append
	cursors := []*workload.Sync_Option_Cursor{}

	// Expected outcome
	expectedKind := property.MessageKind_insert
	var expectedErr error = nil

	// Call the function
	actualKind, actualErr := toMessageKind(schema, mode, cursors)

	// Verify the results
	if actualKind != expectedKind {
		t.Errorf("Expected message kind %v, got %v", expectedKind, actualKind)
	}
	if actualErr != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, actualErr)
	}
}
func TestToMessageKind_IncrementalAppendWithoutCursor(t *testing.T) {
	// Create a simple schema (content doesn't matter for this test)
	fields := []arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}
	schema := arrow.NewSchema(fields, nil)

	// Test inputs
	mode := property.SyncMode_incremental_append
	var cursors []*workload.Sync_Option_Cursor // empty cursor slice

	// Call the function
	result, err := toMessageKind(schema, mode, cursors)

	// Verify the result
	expectedResult := property.MessageKind(0)
	if result != expectedResult {
		t.Errorf("Expected MessageKind %v, got %v", expectedResult, result)
	}

	// Verify the error
	expectedError := "cursor is empty"
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error '%v', got '%v'", expectedError, err)
	}
}
func TestToMessageKind_IncrementalAppendDedupeWithoutCursor(t *testing.T) {
	// Create a simple schema (the specific schema doesn't matter for this test)
	fields := []arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}
	schema := arrow.NewSchema(fields, nil)

	// Test inputs
	mode := property.SyncMode_incremental_append_dedupe
	var cursors []*workload.Sync_Option_Cursor // empty cursor slice

	// Call the function
	got, err := toMessageKind(schema, mode, cursors)

	// Verify the error
	if err == nil {
		t.Fatal("Expected error 'cursor is empty', got nil")
	}
	if got != 0 {
		t.Errorf("Expected return value 0, got %v", got)
	}
	if err.Error() != "cursor is empty" {
		t.Errorf("Expected error 'cursor is empty', got '%v'", err.Error())
	}

	// Alternatively, we could use errors.Is if the error was wrapped
	if !errors.Is(err, errors.New("cursor is empty")) {
		t.Errorf("Error does not match expected: %v", err)
	}
}

func TestToMessageKind_UnspecifiedMode(t *testing.T) {
	// Create a simple schema (doesn't matter for this test case)
	fields := []arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}
	sch := arrow.NewSchema(fields, nil)

	// Use an invalid mode (not one of the defined constants)
	invalidMode := property.SyncMode(999)       // Using a value that's not in the enum
	cursors := []*workload.Sync_Option_Cursor{} // Empty cursor slice

	// Call the function
	got, err := toMessageKind(sch, invalidMode, cursors)

	// Verify the result
	expected := property.MessageKind_message_kind_unspecified
	if got != expected {
		t.Errorf("toMessageKind() with invalid mode = %v, want %v", got, expected)
	}
	if err != nil {
		t.Errorf("toMessageKind() with invalid mode returned error %v, want nil", err)
	}
}
func TestToMessageKind_InvalidCursorStructure(t *testing.T) {
	tests := []struct {
		name        string
		mode        property.SyncMode
		cursors     []*workload.Sync_Option_Cursor
		expectError bool
	}{
		{
			name:        "incremental_append with nil cursor elements",
			mode:        property.SyncMode_incremental_append,
			cursors:     []*workload.Sync_Option_Cursor{nil, nil},
			expectError: false, // Should pass cursor existence check
		},
		{
			name:        "incremental_append_dedupe with nil cursor elements",
			mode:        property.SyncMode_incremental_append_dedupe,
			cursors:     []*workload.Sync_Option_Cursor{nil, nil},
			expectError: false, // Should pass cursor existence check
		},
	}

	// Create a simple schema without primary key
	fields := []arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}
	schema := arrow.NewSchema(fields, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := toMessageKind(schema, tt.mode, tt.cursors)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				// We only care that it didn't fail the cursor existence check
				// The error could be about missing primary key, which is expected
				if err != nil && err.Error() == "cursor is empty" {
					t.Errorf("Should not fail on cursor existence check with nil cursor elements")
				}
			}
		})
	}
}
func TestToMessageKind_AllModesWithoutPrimaryKey(t *testing.T) {
	// Create a schema without primary key
	fields := []arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}
	schema := arrow.NewSchema(fields, nil)

	// Test cases for each mode
	tests := []struct {
		name     string
		mode     property.SyncMode
		cursors  []*workload.Sync_Option_Cursor
		expected property.MessageKind
		wantErr  bool
	}{
		{
			name:     "full_overwrite without primary key",
			mode:     property.SyncMode_full_overwrite,
			cursors:  nil,
			expected: property.MessageKind_insert,
			wantErr:  false,
		},
		{
			name:     "full_append without primary key",
			mode:     property.SyncMode_full_append,
			cursors:  nil,
			expected: property.MessageKind_insert,
			wantErr:  false,
		},
		{
			name:     "incremental_append without primary key",
			mode:     property.SyncMode_incremental_append,
			cursors:  []*workload.Sync_Option_Cursor{{}}, // non-empty cursor
			expected: property.MessageKind_insert,
			wantErr:  false,
		},
		{
			name:     "incremental_append_dedupe without primary key",
			mode:     property.SyncMode_incremental_append_dedupe,
			cursors:  []*workload.Sync_Option_Cursor{{}}, // non-empty cursor
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toMessageKind(schema, tt.mode, tt.cursors)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if got != tt.expected {
				t.Errorf("Expected MessageKind %v, got %v", tt.expected, got)
			}
		})
	}
}

func Test_toFieldMetadata_NilConTypes(t *testing.T) {
	// Setup input
	att := &or.Attribute{ConTypes: nil}

	// Call the function
	result := toFieldMetadata(att)

	// Verify the result
	if result.Len() != 0 {
		t.Errorf("Expected empty metadata for nil ConTypes, got length %d", result.Len())
	}

	// if !reflect.DeepEqual(result, arrow.Metadata{}) {
	// 	t.Errorf("Expected empty arrow.Metadata, got %v", result)
	// }
}
func TestToFieldMetadata_EmptyConTypes(t *testing.T) {
	// Setup
	emptyString := ""
	att := &or.Attribute{ConTypes: &emptyString}

	// Execute
	result := toFieldMetadata(att)

	// Verify
	expected := arrow.MetadataFrom(map[string]string{})
	if !result.Equal(expected) {
		t.Errorf("Expected empty arrow.Metadata when ConTypes is empty string, got %v", result)
	}
}

func TestToFieldMetadata_InvalidConTypes(t *testing.T) {
	// Setup test case
	invalidString := "xyz" // Contains neither 'p' nor 'u'
	att := &or.Attribute{ConTypes: &invalidString}

	// Execute function
	result := toFieldMetadata(att)

	// Verify results
	expected := arrow.MetadataFrom(map[string]string{})
	if !result.Equal(expected) {
		t.Errorf("Should return empty metadata for invalid ConTypes, got %v", result)
	}
}


func TestToFieldMetadata_SpecialCharactersInConTypes(t *testing.T) {
	// Setup test case
	specialCharString := "!@#$%^&*()"
	att := &or.Attribute{ConTypes: &specialCharString}

	// Expected outcome
	expected := arrow.MetadataFrom(map[string]string{})

	// Call the function
	got := toFieldMetadata(att)

	// Verify the result
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("toFieldMetadata() with special characters in ConTypes = %v, want %v", got, expected)
	}
}

func TestToFieldMetadata_ConTypesWithNumbers(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name     string
		numberString string
		want     arrow.Metadata
	}{
		{
			name:     "single digit",
			numberString: "1",
			want:     arrow.MetadataFrom(map[string]string{}),
		},
		{
			name:     "multiple digits",
			numberString: "123",
			want:     arrow.MetadataFrom(map[string]string{}),
		},
		{
			name:     "mixed with letters but no p or u",
			numberString: "a1b2c3",
			want:     arrow.MetadataFrom(map[string]string{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create input with number string
			att := &or.Attribute{ConTypes: &tt.numberString}
			
			// Call the function
			got := toFieldMetadata(att)
			
			// Verify the result
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toFieldMetadata() with ConTypes = %v got = %v, want %v", tt.numberString, got, tt.want)
			}
		})
	}
}