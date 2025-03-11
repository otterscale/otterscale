package migrations

import (
	"errors"
	"os"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		if os.Getenv("OPENHDC_SUPERUSER_EMAIL") == "" || os.Getenv("OPENHDC_SUPERUSER_PASSWORD") == "" {
			return errors.New("environment variables 'OPENHDC_SUPERUSER_EMAIL' and 'OPENHDC_SUPERUSER_PASSWORD' are required")
		}

		superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		if err != nil {
			return err
		}

		record := core.NewRecord(superusers)

		record.Set("email", os.Getenv("OPENHDC_SUPERUSER_EMAIL"))
		record.Set("password", os.Getenv("OPENHDC_SUPERUSER_PASSWORD"))

		return app.Save(record)
	}, nil)
}
