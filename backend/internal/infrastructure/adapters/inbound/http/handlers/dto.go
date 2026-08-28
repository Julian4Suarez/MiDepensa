// Package handlers contains the HTTP inbound adapter.
package handlers

import (
	"time"

	"github.com/google/uuid"

	"midepensa/internal/domain/entities"
)

// createPantryRequest is the body of POST /v1/pantries.
type createPantryRequest struct {
	Name string `json:"name" binding:"required,max=60"`
}

// updateItemRequest is the body of PATCH /v1/pantries/:slug/items/:productId.
// Every field is optional; omitted fields are left untouched.
type updateItemRequest struct {
	Status             *entities.ItemStatus  `json:"status"`
	Type               *entities.ProductType `json:"type"`
	Category           *entities.Category    `json:"category"`
	SelectedVariantIDs *[]uuid.UUID          `json:"selectedVariantIds"`
}

// pantryResponse describes a pantry without its items.
type pantryResponse struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// pantryDetailResponse is a pantry together with all of its items.
type pantryDetailResponse struct {
	pantryResponse
	Items []itemResponse `json:"items"`
}

// productResponse is a catalog entry.
type productResponse struct {
	ID       string            `json:"id"`
	Code     string            `json:"code"`
	Name     string            `json:"name"`
	Image    string            `json:"image"`
	Variants []variantResponse `json:"variants"`
}

type variantResponse struct {
	ID    string `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// itemResponse is the per-pantry state of a catalog product.
type itemResponse struct {
	Product            productResponse `json:"product"`
	Status             string          `json:"status"`
	Type               string          `json:"type"`
	Category           string          `json:"category"`
	SelectedVariantIDs []string        `json:"selectedVariantIds"`
	UpdatedAt          time.Time       `json:"updatedAt"`
}

// catalogResponse tells the frontend which products and enum values exist.
type catalogResponse struct {
	Products   []productResponse `json:"products"`
	Categories []string          `json:"categories"`
	Types      []string          `json:"types"`
	Statuses   []string          `json:"statuses"`
}

func toPantryResponse(pantry *entities.Pantry) pantryResponse {
	return pantryResponse{
		ID:        pantry.ID.String(),
		Slug:      pantry.Slug.String(),
		Name:      pantry.Name,
		CreatedAt: pantry.CreatedAt,
		UpdatedAt: pantry.UpdatedAt,
	}
}

func toProductResponse(product entities.Product) productResponse {
	variants := make([]variantResponse, 0, len(product.Variants))
	for _, variant := range product.Variants {
		variants = append(variants, variantResponse{
			ID: variant.ID.String(), Code: variant.Code, Name: variant.Name, Image: variant.Image,
		})
	}
	return productResponse{
		ID:       product.ID.String(),
		Code:     product.Code,
		Name:     product.Name,
		Image:    product.Image,
		Variants: variants,
	}
}

func toItemResponse(item entities.PantryItem) itemResponse {
	selectedVariantIDs := make([]string, 0, len(item.SelectedVariantIDs))
	for _, id := range item.SelectedVariantIDs {
		selectedVariantIDs = append(selectedVariantIDs, id.String())
	}
	return itemResponse{
		Product:            toProductResponse(item.Product),
		Status:             string(item.Status),
		Type:               string(item.Type),
		Category:           string(item.Category),
		SelectedVariantIDs: selectedVariantIDs,
		UpdatedAt:          item.UpdatedAt,
	}
}

func toItemResponses(items []entities.PantryItem) []itemResponse {
	responses := make([]itemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toItemResponse(item))
	}
	return responses
}
