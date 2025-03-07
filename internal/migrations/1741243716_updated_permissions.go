package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3709660955")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"hidden": false,
			"id": "select1204587666",
			"maxSelect": 1,
			"name": "action",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"allow",
				"deny"
			]
		}`)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"hidden": false,
			"id": "select2363381545",
			"maxSelect": 1,
			"name": "type",
			"presentable": true,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"read",
				"write"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3709660955")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"hidden": false,
			"id": "select1204587666",
			"maxSelect": 1,
			"name": "action",
			"presentable": true,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"allow",
				"deny"
			]
		}`)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"hidden": false,
			"id": "select2363381545",
			"maxSelect": 1,
			"name": "type",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"read",
				"write"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
