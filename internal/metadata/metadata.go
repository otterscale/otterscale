package metadata

import "github.com/go-openapi/inflect"

var (
	KeySchemaTableName     = inflect.Parameterize("table name")
	KeyFieldDataType       = inflect.Parameterize("data type")
	KeyFieldPrimaryKeyName = inflect.Parameterize("primary key name")
	KeyFieldUniqueName     = inflect.Parameterize("unique name")
)
