package pg

import (
	"testing"
)

func TestDropTableStatementWithBasicTableName(t *testing.T) {
	// Input
	tableName := "users"

	// Expected outcome
	expected := "drop table \"users\""

	// Call the function
	actual := dropTableStatement(tableName)

	// Assert the result
	if actual != expected {
		t.Errorf("dropTableStatement(%q) = %q, want %q", tableName, actual, expected)
	}
}