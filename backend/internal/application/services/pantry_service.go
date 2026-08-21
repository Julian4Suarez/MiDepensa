// Package services implements the inbound ports declared in usecases.
package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"midepensa/internal/domain/entities"
	"midepensa/internal/domain/repositories"
	"midepensa/internal/domain/valueobjects"
)

type pantryService struct {
	pantries repositories.PantryRepository
	products repositories.ProductRepository
}

// NewPantryService wires the pantry use cases to their persistence ports.
func NewPantryService(
	pantries repositories.PantryRepository,
	products repositories.ProductRepository,
) *pantryService {
	return &pantryService{pantries: pantries, products: products}
}

// Create builds a pantry and one item per catalog product.
func (s *pantryService) Create(ctx context.Context, name string) (*entities.Pantry, error) {
	pantry, err := entities.NewPantry(name)
	if err != nil {
		return nil, ErrInvalidPantryName
	}

	catalog, err := s.products.List(ctx)
	if err != nil {
		return nil, err
	}
	if len(catalog) == 0 {
		return nil, ErrEmptyCatalog
	}

	items := make([]entities.PantryItem, 0, len(catalog))
	for _, product := range catalog {
		items = append(items, entities.NewPantryItem(pantry.ID, product))
	}

	if err := s.pantries.Create(ctx, pantry, items); err != nil {
		if errors.Is(err, repositories.ErrSlugTaken) {
			return nil, ErrSlugAlreadyExists
		}
		return nil, err
	}
	return pantry, nil
}

// Get returns a pantry with all of its items.
func (s *pantryService) Get(
	ctx context.Context,
	slug valueobjects.Slug,
) (*entities.Pantry, []entities.PantryItem, error) {
	pantry, err := s.pantries.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, nil, ErrPantryNotFound
		}
		return nil, nil, err
	}

	items, err := s.pantries.ListItems(ctx, pantry.ID)
	if err != nil {
		return nil, nil, err
	}
	return pantry, items, nil
}

// UpdateItem changes the status, view or category of a single item.
func (s *pantryService) UpdateItem(
	ctx context.Context,
	slug valueobjects.Slug,
	productID uuid.UUID,
	patch entities.ItemPatch,
) (*entities.PantryItem, error) {
	if patch.IsEmpty() {
		return nil, ErrEmptyPatch
	}
	if err := patch.Validate(); err != nil {
		return nil, ErrInvalidPatch
	}

	pantry, err := s.pantries.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrPantryNotFound
		}
		return nil, err
	}

	item, err := s.pantries.UpdateItem(ctx, pantry.ID, productID, patch)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return item, nil
}
