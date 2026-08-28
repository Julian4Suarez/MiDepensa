//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"midepensa/internal/domain/entities"
	"midepensa/internal/domain/repositories"
	"midepensa/internal/infrastructure/adapters/outbound/persistence"
)

func TestPantryRepository_CreateAndRead_RoundTripsPantryWithItems(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)

	pantryRepository := persistence.NewPostgresPantryRepository(pool)
	productRepository := persistence.NewPostgresProductRepository(pool)

	catalog, err := productRepository.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, catalog, "the seed migration must populate the catalog")

	pantry, err := entities.NewPantry("Integration " + uuid.NewString()[:8])
	require.NoError(t, err)

	items := make([]entities.PantryItem, 0, len(catalog))
	for _, product := range catalog {
		items = append(items, entities.NewPantryItem(pantry.ID, product))
	}
	require.NoError(t, pantryRepository.Create(ctx, pantry, items))
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM pantries WHERE id = $1`, pantry.ID)
	})

	stored, err := pantryRepository.GetBySlug(ctx, pantry.Slug)
	require.NoError(t, err)
	assert.Equal(t, pantry.ID, stored.ID)
	assert.Equal(t, pantry.Name, stored.Name)

	storedItems, err := pantryRepository.ListItems(ctx, pantry.ID)
	require.NoError(t, err)
	assert.Len(t, storedItems, len(catalog))
	assert.Equal(t, catalog[0].DefaultStatus, storedItems[0].Status)
	assert.NotEmpty(t, storedItems[0].Product.Code)
}

func TestPantryRepository_Create_WithDuplicateSlug_ReturnsErrSlugTaken(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repository := persistence.NewPostgresPantryRepository(pool)

	pantry, err := entities.NewPantry("Duplicate " + uuid.NewString()[:8])
	require.NoError(t, err)
	require.NoError(t, repository.Create(ctx, pantry, nil))
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM pantries WHERE id = $1`, pantry.ID)
	})

	duplicate := *pantry
	duplicate.ID = uuid.New()

	err = repository.Create(ctx, &duplicate, nil)

	assert.ErrorIs(t, err, repositories.ErrSlugTaken)
}

func TestPantryRepository_UpdateItem_AppliesOnlyProvidedFields(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)

	pantryRepository := persistence.NewPostgresPantryRepository(pool)
	productRepository := persistence.NewPostgresProductRepository(pool)

	catalog, err := productRepository.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, catalog)

	pantry, err := entities.NewPantry("Patch " + uuid.NewString()[:8])
	require.NoError(t, err)

	product := catalog[0]
	original := entities.NewPantryItem(pantry.ID, product)
	require.NoError(t, pantryRepository.Create(ctx, pantry, []entities.PantryItem{original}))
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM pantries WHERE id = $1`, pantry.ID)
	})

	status := entities.StatusInCart
	updated, err := pantryRepository.UpdateItem(ctx, pantry.ID, product.ID, entities.ItemPatch{Status: &status})

	require.NoError(t, err)
	assert.Equal(t, entities.StatusInCart, updated.Status)
	assert.Equal(t, original.Type, updated.Type, "type must be untouched")
	assert.Equal(t, original.Category, updated.Category, "category must be untouched")
	assert.Equal(t, product.Code, updated.Product.Code)
}

func TestPantryRepository_UpdateItem_WithUnknownProduct_ReturnsErrNotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repository := persistence.NewPostgresPantryRepository(pool)

	status := entities.StatusPending

	_, err := repository.UpdateItem(ctx, uuid.New(), uuid.New(), entities.ItemPatch{Status: &status})

	assert.ErrorIs(t, err, repositories.ErrNotFound)
}
