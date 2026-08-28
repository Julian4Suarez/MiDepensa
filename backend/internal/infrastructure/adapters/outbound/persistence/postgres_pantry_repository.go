package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"midepensa/internal/domain/entities"
	"midepensa/internal/domain/repositories"
	"midepensa/internal/domain/valueobjects"
)

// uniqueViolation is the PostgreSQL SQLSTATE for a duplicate key.
const uniqueViolation = "23505"

type postgresPantryRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresPantryRepository returns the pantry adapter backed by PostgreSQL.
func NewPostgresPantryRepository(pool *pgxpool.Pool) *postgresPantryRepository {
	return &postgresPantryRepository{pool: pool}
}

// Create stores the pantry and its items in a single transaction.
func (r *postgresPantryRepository) Create(
	ctx context.Context,
	pantry *entities.Pantry,
	items []entities.PantryItem,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("persistence: begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx,
		`INSERT INTO pantries (id, slug, name, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		pantry.ID, pantry.Slug.String(), pantry.Name, pantry.CreatedAt, pantry.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return repositories.ErrSlugTaken
		}
		return fmt.Errorf("persistence: insert pantry: %w", err)
	}

	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"pantry_items"},
		[]string{"pantry_id", "product_id", "status", "product_type", "category", "updated_at"},
		pgx.CopyFromSlice(len(items), func(i int) ([]any, error) {
			item := items[i]
			return []any{
				item.PantryID,
				item.Product.ID,
				string(item.Status),
				string(item.Type),
				string(item.Category),
				item.UpdatedAt,
			}, nil
		}),
	)
	if err != nil {
		return fmt.Errorf("persistence: insert pantry items: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("persistence: commit: %w", err)
	}
	return nil
}

// GetBySlug returns the pantry identified by its slug.
func (r *postgresPantryRepository) GetBySlug(
	ctx context.Context,
	slug valueobjects.Slug,
) (*entities.Pantry, error) {
	var (
		pantry     entities.Pantry
		storedSlug string
	)
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, name, created_at, updated_at FROM pantries WHERE slug = $1`,
		slug.String(),
	).Scan(&pantry.ID, &storedSlug, &pantry.Name, &pantry.CreatedAt, &pantry.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repositories.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("persistence: get pantry: %w", err)
	}

	pantry.Slug = valueobjects.Slug(storedSlug)
	return &pantry, nil
}

// ListItems returns every item of a pantry joined with its catalog product.
func (r *postgresPantryRepository) ListItems(
	ctx context.Context,
	pantryID uuid.UUID,
) ([]entities.PantryItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+productColumns+`, i.status, i.product_type, i.category, i.updated_at
		 FROM pantry_items i
		 JOIN products p ON p.id = i.product_id
		 WHERE i.pantry_id = $1
		 ORDER BY p.sort_order`,
		pantryID,
	)
	if err != nil {
		return nil, fmt.Errorf("persistence: list pantry items: %w", err)
	}
	defer rows.Close()

	items := make([]entities.PantryItem, 0)
	for rows.Next() {
		item, err := scanPantryItem(rows, pantryID)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("persistence: list pantry items: %w", err)
	}
	return items, nil
}

// UpdateItem applies the non-nil patch fields and returns the resulting item.
func (r *postgresPantryRepository) UpdateItem(
	ctx context.Context,
	pantryID uuid.UUID,
	productID uuid.UUID,
	patch entities.ItemPatch,
) (*entities.PantryItem, error) {
	row := r.pool.QueryRow(ctx,
		`WITH updated AS (
		     UPDATE pantry_items
		     SET status       = COALESCE($3, status),
		         product_type = COALESCE($4, product_type),
		         category     = COALESCE($5, category),
		         updated_at   = now()
		     WHERE pantry_id = $1 AND product_id = $2
		     RETURNING product_id, status, product_type, category, updated_at
		 )
		 SELECT `+productColumns+`, u.status, u.product_type, u.category, u.updated_at
		 FROM updated u
		 JOIN products p ON p.id = u.product_id`,
		pantryID,
		productID,
		optionalString(patch.Status),
		optionalString(patch.Type),
		optionalString(patch.Category),
	)

	item, err := scanPantryItem(row, pantryID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repositories.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func scanPantryItem(row pgx.Row, pantryID uuid.UUID) (entities.PantryItem, error) {
	item := entities.PantryItem{PantryID: pantryID}
	if err := row.Scan(
		&item.Product.ID,
		&item.Product.Code,
		&item.Product.Name,
		&item.Product.Image,
		&item.Product.DefaultCategory,
		&item.Product.DefaultType,
		&item.Product.DefaultStatus,
		&item.Product.SortOrder,
		&item.Status,
		&item.Type,
		&item.Category,
		&item.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.PantryItem{}, err
		}
		return entities.PantryItem{}, fmt.Errorf("persistence: scan pantry item: %w", err)
	}
	return item, nil
}

// optionalString turns a typed enum pointer into a SQL NULL when unset, so the
// COALESCE in the update statement keeps the existing column value.
func optionalString[T ~string](value *T) *string {
	if value == nil {
		return nil
	}
	plain := string(*value)
	return &plain
}
