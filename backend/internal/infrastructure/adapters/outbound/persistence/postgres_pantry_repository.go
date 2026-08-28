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
	if err := r.attachVariantData(ctx, pantryID, items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *postgresPantryRepository) attachVariantData(
	ctx context.Context,
	pantryID uuid.UUID,
	items []entities.PantryItem,
) error {
	products := make([]entities.Product, len(items))
	positions := make(map[uuid.UUID]int, len(items))
	for i := range items {
		products[i] = items[i].Product
		positions[items[i].Product.ID] = i
		items[i].SelectedVariantIDs = make([]uuid.UUID, 0)
	}
	if err := attachProductVariants(ctx, r.pool, products); err != nil {
		return err
	}
	for i := range items {
		items[i].Product = products[i]
	}
	rows, err := r.pool.Query(ctx,
		`SELECT product_id, variant_id FROM pantry_item_variants
		 WHERE pantry_id = $1 ORDER BY product_id, variant_id`, pantryID)
	if err != nil {
		return fmt.Errorf("persistence: list selected variants: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var productID, variantID uuid.UUID
		if err := rows.Scan(&productID, &variantID); err != nil {
			return fmt.Errorf("persistence: scan selected variant: %w", err)
		}
		if position, ok := positions[productID]; ok {
			items[position].SelectedVariantIDs = append(items[position].SelectedVariantIDs, variantID)
		}
	}
	return rows.Err()
}

// UpdateItem applies the non-nil patch fields and returns the resulting item.
func (r *postgresPantryRepository) UpdateItem(
	ctx context.Context,
	pantryID uuid.UUID,
	productID uuid.UUID,
	patch entities.ItemPatch,
) (*entities.PantryItem, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("persistence: begin item update: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if patch.SelectedVariantIDs != nil {
		unique := make(map[uuid.UUID]struct{}, len(*patch.SelectedVariantIDs))
		for _, id := range *patch.SelectedVariantIDs {
			unique[id] = struct{}{}
		}
		if len(unique) != len(*patch.SelectedVariantIDs) {
			return nil, repositories.ErrInvalidVariant
		}
		if _, err := tx.Exec(ctx,
			`DELETE FROM pantry_item_variants WHERE pantry_id = $1 AND product_id = $2`,
			pantryID, productID,
		); err != nil {
			return nil, fmt.Errorf("persistence: clear selected variants: %w", err)
		}
		if len(*patch.SelectedVariantIDs) > 0 {
			result, err := tx.Exec(ctx,
				`INSERT INTO pantry_item_variants (pantry_id, product_id, variant_id)
				 SELECT $1, $2, variant.id
				 FROM unnest($3::uuid[]) AS requested(id)
				 JOIN product_variants variant ON variant.id = requested.id AND variant.product_id = $2`,
				pantryID, productID, *patch.SelectedVariantIDs,
			)
			if err != nil {
				return nil, fmt.Errorf("persistence: select product variants: %w", err)
			}
			if result.RowsAffected() != int64(len(*patch.SelectedVariantIDs)) {
				return nil, repositories.ErrInvalidVariant
			}
		}
	}

	row := tx.QueryRow(ctx,
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
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("persistence: commit item update: %w", err)
	}
	items := []entities.PantryItem{item}
	if err := r.attachVariantData(ctx, pantryID, items); err != nil {
		return nil, err
	}
	item = items[0]
	return &item, nil
}

// ResetActiveItems marks every non-archived pantry item as pending in one query.
func (r *postgresPantryRepository) ResetActiveItems(
	ctx context.Context,
	pantryID uuid.UUID,
) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE pantry_items
		 SET status = $2, updated_at = now()
		 WHERE pantry_id = $1 AND status <> $3`,
		pantryID,
		entities.StatusPending,
		entities.StatusArchived,
	)
	if err != nil {
		return fmt.Errorf("persistence: reset active pantry items: %w", err)
	}
	return nil
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
