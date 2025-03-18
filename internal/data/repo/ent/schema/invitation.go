package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Invitation holds the schema definition for the Invitation entity.
type Invitation struct {
	ent.Schema
}

// Annotations of the Invitation.
func (Invitation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "invitation"},
	}
}

// Fields of the Invitation.
func (Invitation) Fields() []ent.Field {
	return []ent.Field{
		field.Text("id"),
		field.Text("email"),
		field.Text("role").Optional(),
		field.Text("team_id").Optional(),
		field.Text("status"),
		field.Text("expires_at"),
	}
}

// Edges of the Invitation.
func (Invitation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("invitations").
			Unique().
			Required(),

		edge.From("organization", Organization.Type).
			Ref("invitations").
			Unique().
			Required(),
	}
}
