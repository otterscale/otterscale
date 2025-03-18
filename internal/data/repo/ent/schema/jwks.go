package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// JWKS holds the schema definition for the JWKS entity.
type JWKS struct {
	ent.Schema
}

// Annotations of the JWKS.
func (JWKS) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "jwks"},
	}
}

// Fields of the JWKS.
func (JWKS) Fields() []ent.Field {
	return []ent.Field{
		field.Text("id"),
		field.Text("public_key"),
		field.Text("private_key"),
		field.Time("created_at"),
	}
}

// Edges of the JWKS.
func (JWKS) Edges() []ent.Edge {
	return nil
}
