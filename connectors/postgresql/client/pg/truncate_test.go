package pg

import (
	// "strings"
	"testing"
)

func TestTruncateTableStatement_BasicTableName(t *testing.T) {
	// Input
	tableName := "users"

	// Expected outcome
	expected := "truncate table \"users\""

	// Call the function
	actual := truncateTableStatement(tableName)

	// Verify the result
	if actual != expected {
		t.Errorf("truncateTableStatement(%q) = %q, want %q", tableName, actual, expected)
	}
}
