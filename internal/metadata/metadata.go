package metadata

import "github.com/go-openapi/inflect"

var (
	KeySchemaTableName   = inflect.Parameterize("table name")
	KeyFieldIsPrimaryKey = inflect.Parameterize("is primary key")
	KeyFieldIsUnique     = inflect.Parameterize("is unique")
)
