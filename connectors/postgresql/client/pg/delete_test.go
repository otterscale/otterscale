package pg

import (
	"testing"
)

func TestDeleteStatement_EmptyTableNameInput(t *testing.T) {
	// Input
	tableName := ""

	// Expected outcome
	expected := "delete from \"\" where _openhdc_synced_at < $1"

	// Call the function
	actual := deleteStatement(tableName)

	// Assert the result
	if actual != expected {
		t.Errorf("deleteStatement() with empty table name = %v, want %v", actual, expected)
	}
}