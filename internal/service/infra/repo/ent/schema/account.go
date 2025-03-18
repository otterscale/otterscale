package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Annotations of the Account.
func (Account) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "account"},
	}
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.Text("id"),
		field.Text("account_id"),
		field.Text("provider_id"),
		field.Text("access_token").Optional(),
		field.Text("refresh_token").Optional(),
		field.Text("id_token").Optional(),
		field.Time("access_token_expires_at").Optional(),
		field.Time("refresh_token_expires_at").Optional(),
		field.Text("scope").Optional(),
		field.Text("password").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the Account.
func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("accounts").
			Unique().
			Required(),
	}
}
