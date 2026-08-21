package services

import "errors"

// Errors returned by the application services. The HTTP adapter maps each one
// to a status code in helpers.RegisterDomainError.
var (
	// ErrPantryNotFound means no pantry exists for the requested slug.
	ErrPantryNotFound = errors.New("pantry not found")
	// ErrItemNotFound means the product is not part of the pantry.
	ErrItemNotFound = errors.New("pantry item not found")
	// ErrSlugAlreadyExists means another pantry already owns that slug.
	ErrSlugAlreadyExists = errors.New("pantry slug already exists")
	// ErrInvalidPantryName means the name has no URL-safe characters.
	ErrInvalidPantryName = errors.New("pantry name is invalid")
	// ErrEmptyPatch means the update request would change nothing.
	ErrEmptyPatch = errors.New("update must change at least one field")
	// ErrInvalidPatch means the update carries an unknown enum value.
	ErrInvalidPatch = errors.New("update contains an invalid value")
	// ErrEmptyCatalog means the catalog seed migration has not run.
	ErrEmptyCatalog = errors.New("product catalog is empty")
)
