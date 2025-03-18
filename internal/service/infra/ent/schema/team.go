package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Team holds the schema definition for the Team entity.
type Team struct {
	ent.Schema
}

// Annotations of the Team.
func (Team) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "team"},
	}
}

// Fields of the Team.
func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.Text("id"),
		field.Text("name"),
		field.Time("created_at"),
		field.Time("updated_at").Optional(),
	}
}

// Edges of the Team.
func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organization", Organization.Type).
			Ref("teams").
			Unique().
			Required(),
	}
}
