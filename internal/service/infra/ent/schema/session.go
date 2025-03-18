package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Annotations of the Session.
func (Session) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "session"},
	}
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.Text("id"),
		field.Time("expires_at"),
		field.Text("token").Unique(),
		field.Time("created_at"),
		field.Time("updated_at"),
		field.Text("ip_address").Optional(),
		field.Text("user_agent").Optional(),
		field.Text("impersonated_by").Optional(),
		field.Text("active_organization_id").Optional(),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("sessions").
			Unique().
			Required(),
	}
}
