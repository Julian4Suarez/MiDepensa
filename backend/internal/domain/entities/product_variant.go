package entities

import "github.com/google/uuid"

// ProductVariant is a concrete choice under a general catalog product.
type ProductVariant struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Code      string
	Name      string
	Image     string
	SortOrder int
}
