package pg

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
)

func TestSchema_ReturnsInitializedSchema(t *testing.T) {
	// Setup
	h := &Helper{
		sch: &arrow.Schema{},
	}

	// Execute
	result := h.Schema()

	// Verify
	if result == nil {
		t.Errorf("Schema() should return non-nil when sch is initialized")
	}
	if result != h.sch {
		t.Errorf("Schema() should return the initialized schema pointer, expected %v, got %v", h.sch, result)
	}
}
func TestTableName_EmptyTableName(t *testing.T) {
	// Setup
	h := &Helper{
		tableName: "",
	}

	// Execute
	result := h.TableName()

	// Verify
	if result != "" {
		t.Errorf("TableName() should return empty string when tableName is empty, got %q", result)
	}
}
func TestTableName_ReturnsBasicTableName(t *testing.T) {
	// Setup
	h := &Helper{
		tableName: "users",
	}

	// Execute
	result := h.TableName()

	// Verify
	if result != "users" {
		t.Errorf("TableName() should return the initialized table name, expected 'users', got '%s'", result)
	}
}
