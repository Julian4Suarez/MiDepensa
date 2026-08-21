package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrInvalidItemPatch is returned when a patch carries an unknown enum value.
var ErrInvalidItemPatch = errors.New("pantry item: patch contains an invalid value")

// PantryItem is the per-pantry state of a catalog product.
type PantryItem struct {
	PantryID  uuid.UUID
	Product   Product
	Status    StockStatus
	View      PantryView
	Category  Category
	UpdatedAt time.Time
}

// NewPantryItem creates the initial item for a freshly created pantry, taking
// its category and view from the catalog defaults.
func NewPantryItem(pantryID uuid.UUID, product Product) PantryItem {
	return PantryItem{
		PantryID:  pantryID,
		Product:   product,
		Status:    StatusOK,
		View:      product.DefaultView,
		Category:  product.DefaultCategory,
		UpdatedAt: time.Now().UTC(),
	}
}

// ItemPatch is a partial update of a pantry item; nil fields are left untouched.
type ItemPatch struct {
	Status   *StockStatus
	View     *PantryView
	Category *Category
}

// IsEmpty reports whether the patch would change nothing.
func (p ItemPatch) IsEmpty() bool {
	return p.Status == nil && p.View == nil && p.Category == nil
}

// Validate rejects patches carrying values outside the known enums.
func (p ItemPatch) Validate() error {
	if p.Status != nil && !p.Status.IsValid() {
		return ErrInvalidItemPatch
	}
	if p.View != nil && !p.View.IsValid() {
		return ErrInvalidItemPatch
	}
	if p.Category != nil && !p.Category.IsValid() {
		return ErrInvalidItemPatch
	}
	return nil
}

// Apply mutates the item with the non-nil fields of the patch.
func (i *PantryItem) Apply(patch ItemPatch) {
	if patch.Status != nil {
		i.Status = *patch.Status
	}
	if patch.View != nil {
		i.View = *patch.View
	}
	if patch.Category != nil {
		i.Category = *patch.Category
	}
	i.UpdatedAt = time.Now().UTC()
}
