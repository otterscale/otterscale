package client

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/stretchr/testify/assert"
)

// TestToSchemaMetadata tests the toSchemaMetadata function to ensure it correctly
// converts a table name to schema metadata and verifies that the metadata is not nil
// and contains the correct table name.
func TestToSchemaMetadata(t *testing.T) {
	tableName := "test_table"
	md := toSchemaMetadata(tableName)
	assert.NotNil(t, md)
	assert.Equal(t, tableName, md.Values()[0])
}

// TestToStringFields tests the toStringFields function to ensure it correctly
// converts a slice of arrow fields to string fields and verifies that the result
// matches the expected fields.
func TestToStringFields(t *testing.T) {
	fields := []arrow.Field{
		{Name: "field1", Type: arrow.PrimitiveTypes.Int32},
		{Name: "field2", Type: arrow.PrimitiveTypes.Float64},
	}
	expectedFields := []arrow.Field{
		{Name: "field1", Type: arrow.BinaryTypes.String},
		{Name: "field2", Type: arrow.BinaryTypes.String},
	}
	result := toStringFields(fields)
	assert.Equal(t, expectedFields, result)
}

// TestValidateFields tests the validateFields function with various test cases to ensure
// it correctly validates the fields. It checks for valid fields, empty field names, and
// duplicate field names, and verifies that the function returns the expected error status.
func TestValidateFields(t *testing.T) {
	tests := []struct {
		name    string
		fields  []arrow.Field
		wantErr bool
	}{
		{
			name: "valid fields",
			fields: []arrow.Field{
				{Name: "field1", Type: arrow.PrimitiveTypes.Int32},
				{Name: "field2", Type: arrow.PrimitiveTypes.Float64},
			},
			wantErr: false,
		},
		{
			name: "empty field name",
			fields: []arrow.Field{
				{Name: "", Type: arrow.PrimitiveTypes.Int32},
				{Name: "field2", Type: arrow.PrimitiveTypes.Float64},
			},
			wantErr: true,
		},
		{
			name: "duplicate field names",
			fields: []arrow.Field{
				{Name: "field1", Type: arrow.PrimitiveTypes.Int32},
				{Name: "field1", Type: arrow.PrimitiveTypes.Float64},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFields(tt.fields)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
