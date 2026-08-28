package unit

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"midepensa/internal/domain/entities"
)

func TestNewPantry_WithValidName_DerivesSlug(t *testing.T) {
	pantry, err := entities.NewPantry("  Familia Suárez  ")

	require.NoError(t, err)
	assert.Equal(t, "Familia Suárez", pantry.Name)
	assert.Equal(t, "familia-suarez", pantry.Slug.String())
	assert.NotEqual(t, uuid.Nil, pantry.ID)
}

func TestNewPantry_WithBlankName_ReturnsError(t *testing.T) {
	_, err := entities.NewPantry("   ")

	assert.ErrorIs(t, err, entities.ErrPantryNameEmpty)
}

func TestNewPantryItem_TakesCategoryAndTypeFromCatalogDefaults(t *testing.T) {
	product := entities.Product{
		ID:              uuid.New(),
		Code:            "tomato",
		DefaultCategory: entities.CategoryFruitVeg,
		DefaultType:     entities.TypeEssential,
	}
	pantryID := uuid.New()

	item := entities.NewPantryItem(pantryID, product)

	assert.Equal(t, pantryID, item.PantryID)
	assert.Equal(t, entities.StatusOK, item.Status)
	assert.Equal(t, entities.CategoryFruitVeg, item.Category)
	assert.Equal(t, entities.TypeEssential, item.Type)
}

func TestNewPantryItem_TakesArchivedStatusFromCatalogDefault(t *testing.T) {
	product := entities.Product{
		ID:              uuid.New(),
		DefaultCategory: entities.CategoryDrinks,
		DefaultType:     entities.TypeSecondary,
		DefaultStatus:   entities.StatusArchived,
	}

	item := entities.NewPantryItem(uuid.New(), product)

	assert.Equal(t, entities.StatusArchived, item.Status)
	assert.Equal(t, entities.TypeSecondary, item.Type)
}

func TestItemPatch_Validate_RejectsUnknownEnumValues(t *testing.T) {
	unknownStatus := entities.StockStatus("SOLD_OUT")
	unknownType := entities.ProductType("FAVOURITES")
	unknownCategory := entities.Category("PETS")

	assert.ErrorIs(t, entities.ItemPatch{Status: &unknownStatus}.Validate(), entities.ErrInvalidItemPatch)
	assert.ErrorIs(t, entities.ItemPatch{Type: &unknownType}.Validate(), entities.ErrInvalidItemPatch)
	assert.ErrorIs(t, entities.ItemPatch{Category: &unknownCategory}.Validate(), entities.ErrInvalidItemPatch)
}

func TestItemPatch_Apply_OnlyChangesProvidedFields(t *testing.T) {
	item := entities.PantryItem{
		Status:   entities.StatusOK,
		Type:     entities.TypeEssential,
		Category: entities.CategoryFruitVeg,
	}
	status := entities.StatusOut

	item.Apply(entities.ItemPatch{Status: &status})

	assert.Equal(t, entities.StatusOut, item.Status)
	assert.Equal(t, entities.TypeEssential, item.Type)
	assert.Equal(t, entities.CategoryFruitVeg, item.Category)
}

func TestItemPatch_IsEmpty_ReportsWhetherAnythingWouldChange(t *testing.T) {
	status := entities.StatusLow

	assert.True(t, entities.ItemPatch{}.IsEmpty())
	assert.False(t, entities.ItemPatch{Status: &status}.IsEmpty())
}
