package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Verification holds the schema definition for the Verification entity.
type Verification struct {
	ent.Schema
}

// Annotations of the Verification.
func (Verification) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "verification"},
	}
}

// Fields of the Verification.
func (Verification) Fields() []ent.Field {
	return []ent.Field{
		field.Text("id"),
		field.Text("identifier"),
		field.Text("value"),
		field.Text("expires_at"),
		field.Text("created_at").Optional(),
		field.Text("updated_at").Optional(),
	}
}

// Edges of the Verification.
func (Verification) Edges() []ent.Edge {
	return nil
}
