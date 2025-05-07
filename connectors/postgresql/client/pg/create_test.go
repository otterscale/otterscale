package pg

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
)

func TestCreateTableStatement_BasicTableWithoutPrimaryKeys(t *testing.T) {
	// Create a schema with fields but no primary key metadata
	fields := []arrow.Field{
		{Name: "column1", Type: arrow.PrimitiveTypes.Int32},
		{Name: "column2", Type: arrow.PrimitiveTypes.Float64},
	}
	schema := arrow.NewSchema(fields, nil)

	tableName := "users"
	expected := "create table if not exists \"users\" (\"column1\" int4  not null, \"column2\" float8  not null)"

	// Call the function
	actual := createTableStatement(tableName, schema)

	// Verify the result
	if actual != expected {
		t.Errorf("createTableStatement() = %v, want %v", actual, expected)
	}
}
