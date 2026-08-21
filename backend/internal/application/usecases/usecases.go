// Package usecases declares the inbound ports consumed by the HTTP adapter.
package usecases

import (
	"context"

	"github.com/google/uuid"

	"midepensa/internal/domain/entities"
	"midepensa/internal/domain/valueobjects"
)

// PantryService is the entry point for everything a user does with a pantry.
type PantryService interface {
	// Create builds a pantry from a free-text name and seeds its items.
	Create(ctx context.Context, name string) (*entities.Pantry, error)

	// Get returns a pantry with all of its items.
	Get(ctx context.Context, slug valueobjects.Slug) (*entities.Pantry, []entities.PantryItem, error)

	// UpdateItem changes the status, view or category of a single item.
	UpdateItem(
		ctx context.Context,
		slug valueobjects.Slug,
		productID uuid.UUID,
		patch entities.ItemPatch,
	) (*entities.PantryItem, error)
}

// CatalogService exposes the read-only product catalog.
type CatalogService interface {
	// Products returns every catalog entry.
	Products(ctx context.Context) ([]entities.Product, error)
}
