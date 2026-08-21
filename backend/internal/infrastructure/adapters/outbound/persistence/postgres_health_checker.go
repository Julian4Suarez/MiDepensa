package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresHealthChecker struct {
	pool *pgxpool.Pool
}

// NewPostgresHealthChecker returns a readiness probe for the database.
func NewPostgresHealthChecker(pool *pgxpool.Pool) *postgresHealthChecker {
	return &postgresHealthChecker{pool: pool}
}

// Name identifies the dependency in the readiness payload.
func (c *postgresHealthChecker) Name() string { return "database" }

// Check pings the pool using the caller's deadline.
func (c *postgresHealthChecker) Check(ctx context.Context) error {
	return c.pool.Ping(ctx)
}
