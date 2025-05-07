package pg

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/api/workload/v1"
)

func TestSelectStatement_BasicFunctionalityWithMultipleFields(t *testing.T) {
	// Create arrow schema with multiple fields
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

	// Verify the result
	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestCursorsToWhere_EmptyCursorListWithSupportedSyncMode(t *testing.T) {
	// Test inputs
	mode := property.SyncMode_incremental_append
	var curs []*workload.Sync_Option_Cursor = nil

	// Expected output
	expected := ""

	// Call the function
	actual := cursorsToWhere(mode, curs)

	// Verify the result
	if actual != expected {
		t.Errorf("Expected: %q, got: %q", expected, actual)
	}
}
