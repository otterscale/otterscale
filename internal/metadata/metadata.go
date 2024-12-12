package metadata

import "github.com/go-openapi/inflect"

var (
	KeySchemaTableName   = inflect.Parameterize("table name")
	KeyFieldDataType     = inflect.Parameterize("data type")
	KeyFieldIsPrimaryKey = inflect.Parameterize("is primary key")
	KeyFieldIsUnique     = inflect.Parameterize("is unique")
)
