// Package entities holds the business objects of the pantry domain.
package entities

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"midepensa/internal/domain/valueobjects"
)

// MaxPantryNameLength bounds the display name of a pantry.
const MaxPantryNameLength = 60

// ErrPantryNameEmpty is returned when a pantry is created without a name.
var ErrPantryNameEmpty = errors.New("pantry: name is required")

// Pantry is a named collection of stock levels reachable at /pantries/{slug}.
type Pantry struct {
	ID        uuid.UUID
	Slug      valueobjects.Slug
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewPantry builds a pantry from a free-text name, deriving its slug.
func NewPantry(name string) (*Pantry, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, ErrPantryNameEmpty
	}
	if len(trimmed) > MaxPantryNameLength {
		trimmed = trimmed[:MaxPantryNameLength]
	}

	slug, err := valueobjects.NewSlug(trimmed)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Pantry{
		ID:        uuid.New(),
		Slug:      slug,
		Name:      trimmed,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
