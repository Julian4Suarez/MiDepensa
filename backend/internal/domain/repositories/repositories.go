package repositories
// Package repositories declares the outbound ports for persistence.
package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"midepensa/internal/domain/entities"
	"midepensa/internal/domain/valueobjects"
)

// ErrNotFound is returned by repositories when a row does not exist.
var ErrNotFound = errors.New("repository: resource not found")

// ErrSlugTaken is returned when a pantry slug collides with an existing one.
var ErrSlugTaken = errors.New("repository: slug already taken")

// PantryRepository persists pantries and their per-product stock levels.
type PantryRepository interface {
	// Create stores the pantry together with its initial items atomically.
	// It returns ErrSlugTaken when the slug is already in use.
	Create(ctx context.Context, pantry *entities.Pantry, items []entities.PantryItem) error

	// GetBySlug returns the pantry, or ErrNotFound.
	GetBySlug(ctx context.Context, slug valueobjects.Slug) (*entities.Pantry, error)

	// ListItems returns every item of a pantry joined with its catalog product.
	ListItems(ctx context.Context, pantryID uuid.UUID) ([]entities.PantryItem, error)

	// UpdateItem applies a patch and returns the resulting item, or ErrNotFound.
	UpdateItem(
		ctx context.Context,
		pantryID uuid.UUID,
		productID uuid.UUID,
		patch entities.ItemPatch,
	) (*entities.PantryItem, error)
}

// ProductRepository reads the seeded catalog.
type ProductRepository interface {
	// List returns every catalog product ordered by sort order.
	List(ctx context.Context) ([]entities.Product, error)
}
