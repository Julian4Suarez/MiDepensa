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

func TestPantryRepository_UpdateItem_PersistsSelectedVariants(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	pantryRepository := persistence.NewPostgresPantryRepository(pool)
	productRepository := persistence.NewPostgresProductRepository(pool)

	catalog, err := productRepository.List(ctx)
	require.NoError(t, err)
	var product entities.Product
	for _, candidate := range catalog {
		if len(candidate.Variants) >= 2 {
			product = candidate
			break
		}
	}
	require.NotEqual(t, uuid.Nil, product.ID)

	pantry, err := entities.NewPantry("Variants " + uuid.NewString()[:8])
	require.NoError(t, err)
	require.NoError(t, pantryRepository.Create(ctx, pantry, []entities.PantryItem{
		entities.NewPantryItem(pantry.ID, product),
	}))
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM pantries WHERE id = $1`, pantry.ID) })

	status := entities.StatusInCart
	selected := []uuid.UUID{product.Variants[0].ID, product.Variants[1].ID}
	updated, err := pantryRepository.UpdateItem(ctx, pantry.ID, product.ID, entities.ItemPatch{
		Status: &status, SelectedVariantIDs: &selected,
	})

	require.NoError(t, err)
	assert.Equal(t, entities.StatusInCart, updated.Status)
	assert.ElementsMatch(t, selected, updated.SelectedVariantIDs)
	assert.Len(t, updated.Product.Variants, len(product.Variants))

	empty := []uuid.UUID{}
	updated, err = pantryRepository.UpdateItem(ctx, pantry.ID, product.ID, entities.ItemPatch{
		SelectedVariantIDs: &empty,
	})
	require.NoError(t, err)
	assert.Empty(t, updated.SelectedVariantIDs)
	assert.Equal(t, entities.StatusInCart, updated.Status)
}

func TestCatalog_UsesConsolidatedProductsAndVariants(t *testing.T) {
	ctx := context.Background()
	productRepository := persistence.NewPostgresProductRepository(newPool(t))

	catalog, err := productRepository.List(ctx)
	require.NoError(t, err)

	byCode := make(map[string]entities.Product, len(catalog))
	for _, product := range catalog {
		byCode[product.Code] = product
	}

	for _, removed := range []string{
		"canned_food", "canned_peas", "canned_tuna", "canned_sardines",
		"frozen_vegetables", "margarine", "egg_whites", "napkins", "tissues",
		"toilet_paper",
	} {
		_, exists := byCode[removed]
		assert.False(t, exists, "%s should not remain an independent product", removed)
	}

	expectedVariants := map[string][]string{
		"chicken":              {"Whole chicken", "Chicken breast", "Chicken thighs", "Chicken wings"},
		"fish":                 {"White fish", "Salmon", "Hake", "Tuna", "Sardines", "Prawns"},
		"leafy_greens":         {"Mixed greens", "Lettuce", "Spinach"},
		"onion":                {"Yellow onion", "Red onion", "Sweet onion", "Leek"},
		"potato":               {"White potato", "Red potato", "Sweet potato"},
		"bell_pepper":          {"Green pepper", "Red pepper", "Yellow pepper"},
		"tomato":               {"Cherry tomatoes", "Roma tomatoes", "Beefsteak tomatoes", "Vine tomatoes"},
		"beans":                {"White beans", "Lentils", "Chickpeas", "Black beans", "Kidney beans"},
		"butter":               {"Salted butter", "Unsalted butter", "Ghee"},
		"cereal":               {"Corn flakes", "Oats", "Muesli", "Granola"},
		"spices":               {"Black pepper", "Paprika", "Oregano", "Cumin", "Cinnamon", "Curry powder"},
		"aluminium_foil":       {"Aluminium foil", "Cling film", "Baking paper"},
		"laundry_detergent":    {"Laundry detergent", "Fabric softener", "Stain remover", "Bleach"},
		"multipurpose_cleaner": {"Multipurpose cleaner", "Disinfectant", "Toilet cleaner", "Floor cleaner", "Glass cleaner", "Degreaser"},
		"dish_soap":            {"Dish soap", "Dishwasher tablets", "Dishwasher salt"},
		"sponges":              {"Sponges", "Dishcloths", "Rubber gloves"},
		"paper_towels":         {"Kitchen roll", "Napkins", "Tissues", "Toilet paper"},
		"hand_soap":            {"Hand soap", "Shower gel", "Shampoo", "Conditioner"},
	}
	for code, expected := range expectedVariants {
		product, exists := byCode[code]
		require.True(t, exists, "%s should be in the catalog", code)
		names := make([]string, 0, len(product.Variants))
		for _, variant := range product.Variants {
			names = append(names, variant.Name)
		}
		assert.Equal(t, expected, names, "unexpected variants for %s", code)
	}

	assert.Equal(t, "Beef", byCode["red_meat"].Name)
	assert.Equal(t, "Cold cuts", byCode["cooked_ham"].Name)
	assert.Equal(t, "Cilantro", byCode["coriander"].Name)
	assert.Equal(t, "Paper products", byCode["paper_towels"].Name)
	assert.Equal(t, entities.TypeEssential, byCode["avocado"].DefaultType)
	assert.Equal(t, entities.TypeEssential, byCode["orange"].DefaultType)
	assert.Equal(t, entities.TypeEssential, byCode["lemon"].DefaultType)
	assert.Equal(t, entities.CategoryDairyEggs, byCode["cream"].DefaultCategory)
	assert.Equal(t, entities.StatusPending, byCode["honey"].DefaultStatus)
	assert.Equal(t, entities.StatusPending, byCode["mustard"].DefaultStatus)
	assert.Equal(t, entities.StatusPending, byCode["grapes"].DefaultStatus)
	assert.Equal(t, entities.StatusPending, byCode["cucumber"].DefaultStatus)
	assert.Equal(t, entities.StatusPending, byCode["spices"].DefaultStatus)
	assert.Equal(t, entities.StatusArchived, byCode["cherries"].DefaultStatus)
	assert.Equal(t, entities.StatusArchived, byCode["raspberries"].DefaultStatus)
	assert.Equal(t, entities.StatusArchived, byCode["green_beans"].DefaultStatus)
	assert.Equal(t, entities.StatusArchived, byCode["bacon"].DefaultStatus)
	assert.Equal(t, entities.StatusArchived, byCode["turkey"].DefaultStatus)
	assert.Equal(t, entities.StatusArchived, byCode["tofu"].DefaultStatus)
	assert.Equal(t, entities.StatusArchived, byCode["mouthwash"].DefaultStatus)
	assert.Equal(t, entities.StatusArchived, byCode["breadcrumbs"].DefaultStatus)
	assert.Equal(t, entities.StatusArchived, byCode["tomato_sauce"].DefaultStatus)
}
