package repo

import (
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	"github.com/openhdc/openhdc/internal/service/infra/ent"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewEntClient() (*ent.Client, func(), error) {
	db, err := sql.Open("pgx", `databaseUrl`)
	if err != nil {
		return nil, nil, err
	}
	drv := entsql.OpenDB(dialect.Postgres, db)
	c := ent.NewClient(ent.Driver(drv))
	return c, func() { c.Close() }, nil
}
