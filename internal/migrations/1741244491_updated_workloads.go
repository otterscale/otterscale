package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3441398300")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"createRule": "@request.auth.id != \"\" &&\n@request.auth.roles.permissions.action ?= \"allow\" &&\n(\n  // roles or groups\n  @request.auth.roles.permissions.type ?= \"write\" &&\n  @request.auth.roles.permissions.object ?= \"_connectors\"\n  ||\n  @request.auth.groups.permission.type = \"write\" &&\n  @request.auth.groups.permission.object = connector.id\n)",
			"listRule": "@request.auth.id != \"\" &&\n@request.auth.roles.permissions.action ?= \"allow\" &&\n(\n  // roles or groups\n  @request.auth.roles.permissions.type ?= \"read\" &&\n  @request.auth.roles.permissions.object ?= \"_connectors\"\n  ||\n  @request.auth.groups.permission.type = \"read\" &&\n  @request.auth.groups.permission.object = connector.id\n)",
			"updateRule": "@request.auth.id != \"\" &&\n@request.auth.roles.permissions.action ?= \"allow\" &&\n(\n  // roles or groups\n  @request.auth.roles.permissions.type ?= \"write\" &&\n  @request.auth.roles.permissions.object ?= \"_connectors\"\n  ||\n  @request.auth.groups.permission.type = \"write\" &&\n  @request.auth.groups.permission.object = connector.id\n)"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3441398300")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"createRule": "@request.auth.id != ''",
			"listRule": "@request.auth.id = user.id",
			"updateRule": null
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
