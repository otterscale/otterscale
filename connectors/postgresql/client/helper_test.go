package client

import (
	"testing"
	// "github.com/apache/arrow-go/v18/arrow"
)

func TestToSchemaMetadata(t *testing.T) {
	testCases := []struct {
		name      string
		tableName string
		want      map[string]string
	}{
		{
			name:      "Basic table name",
			tableName: "users",
			want:      map[string]string{"_ohdc_table_name": "users"},
		},
		{
			name:      "Empty table name",
			tableName: "",
			want:      map[string]string{"_ohdc_table_name": ""},
		},
		{
			name:      "Special characters in table name",
			tableName: "users_data-2023",
			want:      map[string]string{"_ohdc_table_name": "users_data-2023"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toSchemaMetadata(tc.tableName)

			// Check that result is not nil
			if result == nil {
				t.Fatal("Expected non-nil metadata result")
			}

			// Convert metadata back to map for comparison
			resultMap := make(map[string]string)
			keys := result.Keys()
			values := result.Values()
			for i := 0; i < len(keys); i++ {
				resultMap[keys[i]] = values[i]
			}

			// Compare maps
			if len(resultMap) != len(tc.want) {
				t.Errorf("Metadata length mismatch: got %d entries, want %d entries", len(resultMap), len(tc.want))
			}

			for k, expectedValue := range tc.want {
				actualValue, exists := resultMap[k]
				if !exists {
					t.Errorf("Missing key in metadata: %s", k)
				} else if actualValue != expectedValue {
					t.Errorf("Value mismatch for key %s: got %s, want %s", k, actualValue, expectedValue)
				}
			}
		})
	}
}
