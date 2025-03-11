package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2153234234")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_1335902081",
			"hidden": false,
			"id": "relation1602912115",
			"maxSelect": 999,
			"minSelect": 0,
			"name": "source",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_1335902081",
			"hidden": false,
			"id": "relation1053179562",
			"maxSelect": 999,
			"minSelect": 0,
			"name": "destination",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2153234234")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation1602912115")

		// remove field
		collection.Fields.RemoveById("relation1053179562")

		return app.Save(collection)
	})
}
