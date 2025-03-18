package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user"},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Text("id"),
		field.Text("name"),
		field.Text("email").Unique(),
		field.Bool("email_verified"),
		field.Text("image").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
		field.Text("role").Optional(),
		field.Bool("banned").Optional(),
		field.Text("ban_reason").Optional(),
		field.Time("ban_expires").Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sessions", Session.Type).StorageKey(edge.Column("user_id")),
		edge.To("accounts", Account.Type).StorageKey(edge.Column("user_id")),
		edge.To("members", Member.Type).StorageKey(edge.Column("user_id")),
		edge.To("invitations", Invitation.Type).StorageKey(edge.Column("inviter_id")),
	}
}
