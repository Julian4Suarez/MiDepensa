package entities

import "github.com/google/uuid"

// Product is a catalog entry shared by every pantry. The catalog is seeded by
// a migration and is read-only at runtime.
type Product struct {
	ID   uuid.UUID
	Code string
	Name string
	// Image is the file name of the icon under frontend assets, e.g. "tomato.svg".
	Image           string
	DefaultCategory Category
	DefaultType     ProductType
	DefaultStatus   StockStatus
	SortOrder       int
}
