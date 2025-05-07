package or

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/api/workload/v1"
)

func TestSelectStatement_BasicSelectGeneration(t *testing.T) {
	// Create arrow schema with fields ["id", "name", "email"]
	fields := []arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "email", Type: arrow.BinaryTypes.String},
	}
	schema := arrow.NewSchema(fields, nil)

	// Test inputs
	tableName := "users"
	mode := property.SyncMode_full_overwrite
	var curs []*workload.Sync_Option_Cursor = nil

	// Expected output
	expected := "select \"id\", \"name\", \"email\" from \"users\""

	// Call the function
	actual := selectStatement(tableName, schema, mode, curs)

	// Assert the result
	if actual != expected {
		t.Errorf("Expected: %s, Got: %s", expected, actual)
	}
}

func TestCursorsToWhere_SingleCursorCondition(t *testing.T) {
	// Setup test inputs
	mode := property.SyncMode_incremental_append
	cursors := make([]*workload.Sync_Option_Cursor, 1)
	cursors[0] = &workload.Sync_Option_Cursor{}
	cursors[0].SetField("updated_at")
	cursors[0].SetValue("2023-01-01")
	// cursors := []*workload.Sync_Option_Cursor{
	// 	{
	// 		Field: "updated_at",
	// 		Value: "2023-01-01",
	// 	},
	// }

	// Call the function
	result := cursorsToWhere(mode, cursors)

	// Verify the result
	expected := " where \"updated_at\" > '2023-01-01'"
	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
}
