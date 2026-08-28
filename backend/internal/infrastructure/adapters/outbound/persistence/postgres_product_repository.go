package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"midepensa/internal/domain/entities"
)

const productColumns = `p.id, p.code, p.name, p.image, p.default_category, p.default_type, p.default_status, p.sort_order`

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
	rows.Close()
	if err := attachProductVariants(ctx, r.pool, products); err != nil {
		return nil, err
	}
	return products, nil
}

func attachProductVariants(ctx context.Context, pool *pgxpool.Pool, products []entities.Product) error {
	if len(products) == 0 {
		return nil
	}
	ids := make([]uuid.UUID, len(products))
	positions := make(map[uuid.UUID]int, len(products))
	for i := range products {
		ids[i] = products[i].ID
		positions[products[i].ID] = i
		products[i].Variants = make([]entities.ProductVariant, 0)
	}
	rows, err := pool.Query(ctx,
		`SELECT id, product_id, code, name, image, sort_order
		 FROM product_variants WHERE product_id = ANY($1) ORDER BY product_id, sort_order`, ids)
	if err != nil {
		return fmt.Errorf("persistence: list product variants: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var variant entities.ProductVariant
		if err := rows.Scan(&variant.ID, &variant.ProductID, &variant.Code, &variant.Name, &variant.Image, &variant.SortOrder); err != nil {
			return fmt.Errorf("persistence: scan product variant: %w", err)
		}
		if position, ok := positions[variant.ProductID]; ok {
			products[position].Variants = append(products[position].Variants, variant)
		}
	}
	return rows.Err()
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
		&product.DefaultStatus,
		&product.SortOrder,
	); err != nil {
		return entities.Product{}, fmt.Errorf("persistence: scan product: %w", err)
	}
	return product, nil
}
