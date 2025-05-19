package repo

import (
	"context"
	gosql "database/sql"
	"os"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/openhdc/otterscale/internal/data/repo/ent"
	"github.com/openhdc/otterscale/internal/env"
)

// Repo is an alias for the ent Client
type Repo = ent.Client

// NewConfig creates a new PostgreSQL connection configuration
func NewConfig() (*pgx.ConnConfig, error) {
	cfg, err := pgx.ParseConfig(os.Getenv(env.OPENHDC_CONNECTION_STRING))
	if err != nil {
		return nil, err
	}
	cfg.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	return cfg, nil
}

// New creates a new repository client with the given configuration
func New(cfg *pgx.ConnConfig) (*Repo, error) {
	// Open database connection
	db, err := gosql.Open("pgx", stdlib.RegisterConnConfig(cfg))
	if err != nil {
		return nil, err
	}

	// Create Ent client
	drv := sql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	// Initialize schema with timeout
	const schemaTimeout = time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), schemaTimeout)
	defer cancel()

	if err := client.Schema.Create(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
