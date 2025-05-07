package or

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/openhdc/openhdc"
)

// func TestSanitizeSingleValidIdentifier(t *testing.T) {
// 	// Test input
// 	input := "test"

// 	// Expected outcome
// 	expected := `"test"`

// 	// Call the function
// 	result := sanitize(input)

// 	// Verify the result
// 	if result != expected {
// 		t.Errorf("sanitize(%q) = %q, want %q", input, result, expected)
// 	}
// }

func TestSanitizeMultipleValidIdentifiers(t *testing.T) {
	// Input
	input := []string{"test1", "test2", "test3"}

	// Expected outcome
	expected := `"test1"."test2"."test3"`

	// Call the function
	result := sanitize(input...)

	// Assert the result
	if result != expected {
		t.Errorf("sanitize() = %v, want %v", result, expected)
	}
}

func TestNewHelper_WithSchemaMissingTableName_ReturnsError(t *testing.T) {
	// Create a schema without table name metadata
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		},
		nil, // no metadata
	)

	// Create a valid codec
	codec := openhdc.NewDefaultCodec()

	// Call NewHelper
	helper, err := NewHelper(schema, codec)

	// Verify the outcome
	if err == nil {
		t.Error("Expected error when schema lacks table name metadata, but got nil")
	}

	if helper != nil {
		t.Error("Expected no helper to be returned when error occurs, but got non-nil helper")
	}
}

func TestTableName_EmptyTableName(t *testing.T) {
	// Create a Helper instance with empty tableName
	h := &Helper{
		tableName: "",
	}

	// Call the TableName method
	result := h.TableName()

	// Verify the outcome
	if result != "" {
		t.Errorf("TableName() = %q, want empty string \"\"", result)
	}
}
func TestTableName_BasicTableName(t *testing.T) {
	// Create a Helper instance with a simple table name
	h := &Helper{
		tableName: "users",
	}

	// Call the TableName method
	result := h.TableName()

	// Expected outcome
	expected := "users"

	// Verify the result
	if result != expected {
		t.Errorf("TableName() = %q, want %q", result, expected)
	}
}
func TestSchema_ReturnSchemaWhenHelperHasSchema(t *testing.T) {
	// Create a test schema
	testSchema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		},
		nil,
	)

	// Create a Helper instance with the test schema
	h := &Helper{
		sch: testSchema,
	}

	// Call the Schema method
	result := h.Schema()

	// Verify the outcome
	if result != testSchema {
		t.Errorf("Schema() returned %v, want %v", result, testSchema)
	}
}