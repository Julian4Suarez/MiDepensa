package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"midepensa/internal/domain/entities"
)

const productColumns = `p.id, p.code, p.name, p.image, p.default_category, p.default_type, p.sort_order`

type postgresProductRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresProductRepository returns the catalog adapter backed by PostgreSQL.
func NewPostgresProductRepository(pool *pgxpool.Pool) *postgresProductRepository {
	return &postgresProductRepository{pool: pool}
}

// List returns every catalog product ordered by sort order.
func (r *postgresProductRepository) List(ctx context.Context) ([]entities.Product, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+productColumns+` FROM products p ORDER BY p.sort_order`)
	if err != nil {
		return nil, fmt.Errorf("persistence: list products: %w", err)
	}
	defer rows.Close()

	products := make([]entities.Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("persistence: list products: %w", err)
	}
	return products, nil
}

func scanProduct(row pgx.Row) (entities.Product, error) {
	var product entities.Product
	if err := row.Scan(
		&product.ID,
		&product.Code,
		&product.Name,
		&product.Image,
		&product.DefaultCategory,
		&product.DefaultType,
		&product.SortOrder,
	); err != nil {
		return entities.Product{}, fmt.Errorf("persistence: scan product: %w", err)
	}
	return product, nil
}
