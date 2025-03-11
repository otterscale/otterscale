package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1335902081")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"listRule": "@request.auth.id != \"\" &&\n@request.auth.roles.permissions.action ?= \"allow\" &&\n(@request.auth.roles.permissions.object ?= \"_connectors_read\" ||\n@request.auth.groups.permissions.object ?= id)\n"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1335902081")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"listRule": "@request.auth.id != \"\" &&\n@request.auth.roles.permissions.object ?= \"_connectors_read\" &&\n@request.auth.roles.permissions.action ?= \"allow\""
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
