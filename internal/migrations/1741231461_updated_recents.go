package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1032687501")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"createRule": "@request.auth.id != ''",
			"listRule": "@request.auth.id = user.id",
			"updateRule": "@request.auth.id = user.id"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1032687501")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"createRule": "",
			"listRule": "user.id = @request.auth.id",
			"updateRule": "user.id = @request.auth.id"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
