//go:build integration

// Package integration exercises the PostgreSQL adapters against a real
// database. Start it with `make services-start` before running these tests.
package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"midepensa/internal/config"
	"midepensa/internal/infrastructure/adapters/outbound/persistence"
	"midepensa/internal/infrastructure/migrations"
)

func databaseConfig() config.Database {
	return config.Database{
		Host:     envOr("DB_HOST", "localhost"),
		Port:     envOr("DB_PORT", "5433"),
		User:     envOr("DB_USER", "midepensa_user"),
		Password: envOr("DB_PASSWORD", "midepensa_password"),
		Name:     envOr("DB_NAME", "midepensa_db"),
	}
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// newPool migrates the schema and returns a pool bound to the test lifetime.
func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg := databaseConfig()
	require.NoError(t, migrations.Run(cfg.DSN()), "run `make services-start` first")

	pool, err := persistence.NewPool(context.Background(), cfg)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	return pool
}
