package repo

import (
	gosql "database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"

	"github.com/openhdc/openhdc/internal/data/repo/ent"
	"github.com/openhdc/openhdc/internal/env"
)

type Repo = ent.Client

type Data struct {
	db *ent.Client
}

func NewConfig() (*pgx.ConnConfig, error) {
	cfg, err := pgx.ParseConfig(env.OPENHDC_CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	cfg.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	return cfg, nil
}

func New(cfg *pgx.ConnConfig) (*Repo, error) {
	db, err := gosql.Open("pgx", stdlib.RegisterConnConfig(cfg))
	if err != nil {
		return nil, err
	}
	drv := sql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv)), nil
}
