package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Organization holds the schema definition for the Organization entity.
type Organization struct {
	ent.Schema
}

// Annotations of the Organization.
func (Organization) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "organization"},
	}
}

// Fields of the Organization.
func (Organization) Fields() []ent.Field {
	return []ent.Field{
		field.Text("id"),
		field.Text("name"),
		field.Text("slug").Unique(),
		field.Text("logo").Optional(),
		field.Text("created_at"),
		field.Text("metadata").Optional(),
	}
}

// Edges of the Organization.
func (Organization) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("teams", Team.Type).StorageKey(edge.Column("organization_id")),
		edge.To("members", Member.Type).StorageKey(edge.Column("organization_id")),
		edge.To("invitations", Invitation.Type).StorageKey(edge.Column("organization_id")),
	}
}
