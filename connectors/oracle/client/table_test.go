package client

import (
	// "context"
	// "errors"
	"reflect"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	// go_ora "github.com/sijms/go-ora/v2"
	// "github.com/stretchr/testify/assert"
)

func TestTables_Get_Empty(t *testing.T) {
	tables := Tables{}
	got, ok := tables.Get("any")
	if ok {
		t.Errorf("expected table not to be found, but it was found")
	}
	if got != nil {
		t.Errorf("expected got to be nil, but got %v", got)
	}
}

func TestTables_Get_Found(t *testing.T) {
	// Create a dummy schema with metadata containing the table name "TABLE-NAME"
	md := arrow.NewMetadata([]string{"TABLE-NAME"}, []string{"TABLE-NAME"})
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}, &md)
	tables := Tables{schema}

	got, ok := tables.Get("TABLE-NAME")
	if !ok {
		t.Errorf("expected table to be found, but it was not")
	}
	if got == nil {
		t.Errorf("expected got to be non-nil, but it was nil")
	}
	if !reflect.DeepEqual(schema, got) {
		t.Errorf("expected schema to be %v, but got %v", schema, got)
	}
}

func TestTables_Get_NotFound(t *testing.T) {
	md := arrow.NewMetadata([]string{"table1"}, []string{"value"})
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "col1", Type: arrow.PrimitiveTypes.Int32},
	}, &md)
	tables := Tables{schema}

	got, ok := tables.Get("other_table")
	if ok {
		t.Errorf("expected table not to be found, but it was found")
	}
	if got != nil {
		t.Errorf("expected got to be nil, but got %v", got)
	}
}

func TestTables_Get_MultipleTables(t *testing.T) {
	md1 := arrow.NewMetadata([]string{"table1"}, []string{"v1"})
	md2 := arrow.NewMetadata([]string{"table2"}, []string{"v2"})
	md3 := arrow.NewMetadata([]string{"TABLE-NAME"}, []string{"TABLE-NAME"})
	schema1 := arrow.NewSchema([]arrow.Field{{Name: "a", Type: arrow.PrimitiveTypes.Int32}}, &md1)
	schema2 := arrow.NewSchema([]arrow.Field{{Name: "b", Type: arrow.PrimitiveTypes.Float64}}, &md2)
	schema3 := arrow.NewSchema([]arrow.Field{{Name: "c", Type: arrow.BinaryTypes.String}}, &md3)
	tables := Tables{schema1, schema2, schema3}

	got, ok := tables.Get("TABLE-NAME")
	if !ok {
		t.Errorf("expected table to be found, but it was not")
	}
	if !reflect.DeepEqual(schema3, got) {
		t.Errorf("expected schema to be %v, but got %v", schema3, got)
	}
}
