package tables

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

func TestToSchemaMetadata_EmptyTableName(t *testing.T) {
	// Input
	tableName := ""

	// Expected outcome
	// expectedKey := "TABLE-NAME"
	expectedValue := ""

	// Execute function
	result := toSchemaMetadata(tableName)

	// Check that result is not nil
	if result == nil {
		t.Fatal("Expected non-nil metadata result")
	}

	// Verify metadata contains the expected key-value pair
	// actualValue, exists := result.Values()[0]
	// if !exists {
	// 	t.Errorf("Expected metadata to contain key %q", expectedKey)
	// }
	actualValue := result.Values()[0]
	if actualValue != expectedValue {
		t.Errorf("For empty table name, got value %q, want %q", actualValue, expectedValue)
	}
}

func TestToSchemaMetadata_SimpleTableName(t *testing.T) {
	// Input
	tableName := "users"

	// Expected outcome
	// expectedKey := "TABLE-NAME"
	expectedValue := "users"

	// Execute function
	result := toSchemaMetadata(tableName)

	// Check that result is not nil
	if result == nil {
		t.Fatal("Expected non-nil metadata result")
	}

	// Verify metadata contains the expected key-value pair
	// actualValue, exists := result.Values()[0]
	// if !exists {
	// 	t.Errorf("Expected metadata to contain key %q", expectedKey)
	// }
	actualValue := result.Values()[0]
	if actualValue != expectedValue {
		t.Errorf("For table name %q, got value %q, want %q", tableName, actualValue, expectedValue)
	}
}
func TestBuilderAppendJson_HandleNilInputByAppendingNull(t *testing.T) {
	// Create a mock StringBuilder
	fbuilder := array.NewStringBuilder(memory.NewGoAllocator())

	// Call the function with nil input
	builderAppendJson(fbuilder, nil)

	// Verify that AppendNull was called
	if fbuilder.Len() != 1 {
		t.Errorf("Expected 1 value in builder, got %d", fbuilder.Len())
	}

	// Check if the value is actually null
	arr := fbuilder.NewStringArray()
	if !arr.IsNull(0) {
		t.Error("Expected null value in builder, got non-null value")
	}
}
func TestBuilderAppendJson_AppendValidJsonStringForNonNullStruct(t *testing.T) {
	// Setup
	pool := memory.NewGoAllocator()
	fbuilder := array.NewStringBuilder(pool)
	defer fbuilder.Release()

	// Input
	type testStruct struct {
		Name string
	}
	v := testStruct{Name: "Test"}

	// Execute
	builderAppendJson(fbuilder, v)

	// Verify
	result := fbuilder.NewStringArray()
	defer result.Release()
	if result.Len() != 1 {
		t.Fatalf("Expected 1 value in builder, got %d", result.Len())
	}
	if result.IsNull(0) {
		t.Error("Expected non-null value in builder, got null")
	}
	got := result.Value(0)
	expected := `{"Name":"Test"}`
	if got != expected {
		t.Errorf("Expected builder to contain %q, got %q", expected, got)
	}
}
