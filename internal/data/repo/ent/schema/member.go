package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Member holds the schema definition for the Member entity.
type Member struct {
	ent.Schema
}

// Annotations of the Member.
func (Member) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "member"},
	}
}

// Fields of the Member.
func (Member) Fields() []ent.Field {
	return []ent.Field{
		field.Text("id"),
		field.Text("role"),
		field.Text("team_id").Optional(),
		field.Text("created_at"),
	}
}

// Edges of the Member.
func (Member) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("members").
			Unique().
			Required(),

		edge.From("organization", Organization.Type).
			Ref("members").
			Unique().
			Required(),
	}
}
