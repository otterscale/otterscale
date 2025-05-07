package pg

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	// "github.com/jackc/pgx/v5"
)

func TestUpsertStatement_BasicInsertWithoutUpdate(t *testing.T) {
	// Create an empty schema
	fields := []arrow.Field{}
	schema := arrow.NewSchema(fields, nil)

	// Test inputs
	tableName := "users"
	update := false

	// Call the function
	result := upsertStatement(tableName, schema, update)

	// Verify the output
	// expected := "insert into users () values () on conflict do nothing"
	expected := "insert into \"" + tableName + "\" () values () on conflict do nothing"
	if result != expected {
		t.Errorf("The generated SQL statement should match the expected basic insert with conflict handling. Expected: %s, got: %s", expected, result)
	}
}

func TestUpsertStatement_BasicInsertWithUpdate(t *testing.T) {
	// Create an empty schema
	metadata := arrow.NewMetadata(nil, nil)
	schema := arrow.NewSchema([]arrow.Field{}, &metadata)

	// Test inputs
	tableName := "users"
	update := true

	// Expected output
	expected := "insert into \"" + tableName + "\" () values () on conflict () do nothing"

	// Call the function
	actual := upsertStatement(tableName, schema, update)

	// Verify the result
	if actual != expected {
		t.Errorf("Upsert statement mismatch.\nExpected: %s\nActual: %s", expected, actual)
	}
}
