package pg

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	// "github.com/apache/arrow-go/v18/arrow/array"
	// "github.com/stretchr/testify/assert"
)

func TestAddColumnStatement_BasicFunctionalityWithPrefix(t *testing.T) {
	// Create a field with the specified properties
	field := arrow.Field{
		Name:     "test",
		Type:     arrow.BinaryTypes.String,
		Nullable: true,
		Metadata: arrow.NewMetadata(nil, nil), // No metadata means IsUnique will be false
	}

	// Call the function with prefix=true
	result := addColumnStatement(&field, true)

	// Verify the output matches expected
	expected := "add column \"test\" text  "
	if result != expected {
		t.Errorf("Expected result to be %s, got %s", expected, result)
	}
}
func TestDropColumnStatement_BasicColumnNameSanitization(t *testing.T) {
	// Create a field with a simple name that doesn't need sanitization
	field := arrow.Field{
		Name:     "username",
		Type:     arrow.BinaryTypes.String,
		Nullable: true,
		Metadata: arrow.NewMetadata(nil, nil),
	}

	// Call the function
	result := dropColumnStatement(&field)

	// Verify the output matches expected
	expected := "drop column \"username\""
	if result != expected {
		t.Errorf("Expected result to be %s, got %s", expected, result)
	}
}
func TestAlterTableStatement_EmptyTableNameWithNoAddsOrDels(t *testing.T) {
	// Input parameters
	tableName := ""
	adds := []arrow.Field{}
	dels := []arrow.Field{}

	// Call the function
	result := alterTableStatement(tableName, adds, dels)

	// Verify the output matches expected
	expected := "alter table \"\" "
	if result != expected {
		t.Errorf("Expected result to be '%s', got '%s'", expected, result)
	}
}

func TestAlterTableStatement_SingleAddColumnWithValidTableName(t *testing.T) {
	// Input parameters
	tableName := "users"
	adds := []arrow.Field{
		{
			Name: "email",
			Type: arrow.BinaryTypes.String,
		},
	}
	dels := []arrow.Field{}

	// Call the function
	result := alterTableStatement(tableName, adds, dels)

	// Verify the output matches expected
	expected := "alter table \"users\" add column \"email\" text  not null"
	if result != expected {
		t.Errorf("Expected result to be '%s', got '%s'", expected, result)
	}
}
