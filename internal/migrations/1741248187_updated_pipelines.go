package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2153234234")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"listRule": "@request.auth.id != \"\" &&\n@request.auth.roles.permissions.action ?= \"allow\" &&\n(\n  // roles or groups\n  @request.auth.roles.permissions.type ?= \"read\" &&\n  @request.auth.roles.permissions.object ?= \"_connectors\"\n  ||\n  @request.auth.groups.permission.type = \"read\" &&\n  @request.auth.groups.permission.object = source.id\n  ||\n  @request.auth.groups.permission.type = \"read\" &&\n  @request.auth.groups.permission.object = destination.id\n) && deleted = false"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2153234234")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"listRule": "@request.auth.id != \"\" &&\n@request.auth.roles.permissions.action ?= \"allow\" &&\n(\n  // roles or groups\n  @request.auth.roles.permissions.type ?= \"read\" &&\n  @request.auth.roles.permissions.object ?= \"_connectors\"\n  ||\n  @request.auth.groups.permission.type = \"read\" &&\n  @request.auth.groups.permission.object = source.id\n  ||\n  @request.auth.groups.permission.type = \"read\" &&\n  @request.auth.groups.permission.object = destination.id\n)"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
