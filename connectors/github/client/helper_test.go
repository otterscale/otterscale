package client

import (
	"testing"
)

func TestToSchemaMetadata(t *testing.T) {
	tableName := "test_table"
	expectedKey := "TABLE-NAME" //"table_name"
	expectedValue := tableName

	metadata := toSchemaMetadata(tableName)

	if metadata == nil {
		t.Fatalf("Expected metadata to be non-nil")
	}

	if metadata.Len() != 1 {
		t.Fatalf("Expected metadata to have 1 entry, got %d", metadata.Len())
	}

	key, value := metadata.Keys()[0], metadata.Values()[0]
	if key != expectedKey {
		t.Errorf("Expected key to be %s, got %s", expectedKey, key)
	}

	if value != expectedValue {
		t.Errorf("Expected value to be %s, got %s", expectedValue, value)
	}
}
