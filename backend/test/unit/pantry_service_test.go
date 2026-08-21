package unit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"midepensa/internal/application/services"
	"midepensa/internal/domain/entities"
	"midepensa/internal/domain/repositories"
	"midepensa/internal/domain/valueobjects"
)

// mockPantryRepository is a hand-written stub: each method delegates to a
// function field the test sets, so no mocking framework is needed.
type mockPantryRepository struct {
	CreateFn     func(ctx context.Context, p *entities.Pantry, items []entities.PantryItem) error
	GetBySlugFn  func(ctx context.Context, slug valueobjects.Slug) (*entities.Pantry, error)
	ListItemsFn  func(ctx context.Context, pantryID uuid.UUID) ([]entities.PantryItem, error)
	UpdateItemFn func(ctx context.Context, pantryID, productID uuid.UUID, patch entities.ItemPatch) (*entities.PantryItem, error)
}

func (m *mockPantryRepository) Create(ctx context.Context, p *entities.Pantry, items []entities.PantryItem) error {
	return m.CreateFn(ctx, p, items)
}

func (m *mockPantryRepository) GetBySlug(ctx context.Context, slug valueobjects.Slug) (*entities.Pantry, error) {
	return m.GetBySlugFn(ctx, slug)
}

func (m *mockPantryRepository) ListItems(ctx context.Context, pantryID uuid.UUID) ([]entities.PantryItem, error) {
	return m.ListItemsFn(ctx, pantryID)
}

func (m *mockPantryRepository) UpdateItem(ctx context.Context, pantryID, productID uuid.UUID, patch entities.ItemPatch) (*entities.PantryItem, error) {
	return m.UpdateItemFn(ctx, pantryID, productID, patch)
}

type mockProductRepository struct {
	ListFn func(ctx context.Context) ([]entities.Product, error)
}

func (m *mockProductRepository) List(ctx context.Context) ([]entities.Product, error) {
	return m.ListFn(ctx)
}

func catalogOfTwo() []entities.Product {
	return []entities.Product{
		{ID: uuid.New(), Code: "tomato", DefaultCategory: entities.CategoryFruitVeg, DefaultView: entities.ViewPrimary},
		{ID: uuid.New(), Code: "rice", DefaultCategory: entities.CategoryDryCanned, DefaultView: entities.ViewSecondary},
	}
}

func TestCreate_WithValidName_SeedsOneItemPerCatalogProduct(t *testing.T) {
	var seeded []entities.PantryItem
	pantries := &mockPantryRepository{
		CreateFn: func(_ context.Context, _ *entities.Pantry, items []entities.PantryItem) error {
			seeded = items
			return nil
		},
	}
	products := &mockProductRepository{
		ListFn: func(context.Context) ([]entities.Product, error) { return catalogOfTwo(), nil },
	}
	service := services.NewPantryService(pantries, products)

	pantry, err := service.Create(context.Background(), "Familia Suárez")

	require.NoError(t, err)
	assert.Equal(t, "familia-suarez", pantry.Slug.String())
	require.Len(t, seeded, 2)
	assert.Equal(t, entities.StatusOK, seeded[0].Status)
	assert.Equal(t, entities.ViewPrimary, seeded[0].View)
}

func TestCreate_WithTakenSlug_ReturnsConflictError(t *testing.T) {
	pantries := &mockPantryRepository{
		CreateFn: func(context.Context, *entities.Pantry, []entities.PantryItem) error {
			return repositories.ErrSlugTaken
		},
	}
	products := &mockProductRepository{
		ListFn: func(context.Context) ([]entities.Product, error) { return catalogOfTwo(), nil },
	}
	service := services.NewPantryService(pantries, products)

	_, err := service.Create(context.Background(), "Familia")

	assert.ErrorIs(t, err, services.ErrSlugAlreadyExists)
}

func TestCreate_WithUnslugifiableName_ReturnsInvalidNameError(t *testing.T) {
	service := services.NewPantryService(&mockPantryRepository{}, &mockProductRepository{})

	_, err := service.Create(context.Background(), "///")

	assert.ErrorIs(t, err, services.ErrInvalidPantryName)
}

func TestCreate_WithEmptyCatalog_ReturnsCatalogError(t *testing.T) {
	products := &mockProductRepository{
		ListFn: func(context.Context) ([]entities.Product, error) { return nil, nil },
	}
	service := services.NewPantryService(&mockPantryRepository{}, products)

	_, err := service.Create(context.Background(), "Familia")

	assert.ErrorIs(t, err, services.ErrEmptyCatalog)
}

func TestGet_WithUnknownSlug_ReturnsNotFoundError(t *testing.T) {
	pantries := &mockPantryRepository{
		GetBySlugFn: func(context.Context, valueobjects.Slug) (*entities.Pantry, error) {
			return nil, repositories.ErrNotFound
		},
	}
	service := services.NewPantryService(pantries, &mockProductRepository{})

	_, _, err := service.Get(context.Background(), "missing")

	assert.ErrorIs(t, err, services.ErrPantryNotFound)
}

func TestUpdateItem_WithEmptyPatch_ReturnsErrorBeforeTouchingTheRepository(t *testing.T) {
	service := services.NewPantryService(&mockPantryRepository{}, &mockProductRepository{})

	_, err := service.UpdateItem(context.Background(), "familia", uuid.New(), entities.ItemPatch{})

	assert.ErrorIs(t, err, services.ErrEmptyPatch)
}

func TestUpdateItem_WithUnknownStatus_ReturnsInvalidPatchError(t *testing.T) {
	service := services.NewPantryService(&mockPantryRepository{}, &mockProductRepository{})
	unknown := entities.StockStatus("SOLD_OUT")

	_, err := service.UpdateItem(context.Background(), "familia", uuid.New(), entities.ItemPatch{Status: &unknown})

	assert.ErrorIs(t, err, services.ErrInvalidPatch)
}

func TestUpdateItem_WithUnknownProduct_ReturnsItemNotFoundError(t *testing.T) {
	pantries := &mockPantryRepository{
		GetBySlugFn: func(context.Context, valueobjects.Slug) (*entities.Pantry, error) {
			return &entities.Pantry{ID: uuid.New()}, nil
		},
		UpdateItemFn: func(context.Context, uuid.UUID, uuid.UUID, entities.ItemPatch) (*entities.PantryItem, error) {
			return nil, repositories.ErrNotFound
		},
	}
	service := services.NewPantryService(pantries, &mockProductRepository{})
	status := entities.StatusOut

	_, err := service.UpdateItem(context.Background(), "familia", uuid.New(), entities.ItemPatch{Status: &status})

	assert.ErrorIs(t, err, services.ErrItemNotFound)
}
