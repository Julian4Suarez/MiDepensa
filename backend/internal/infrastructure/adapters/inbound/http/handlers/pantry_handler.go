package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"midepensa/internal/application/usecases"
	"midepensa/internal/domain/entities"
	"midepensa/internal/domain/valueobjects"
	"midepensa/internal/infrastructure/adapters/inbound/http/helpers"
)

// PantryHandler exposes the pantry use cases over HTTP.
type PantryHandler struct {
	service usecases.PantryService
}

// NewPantryHandler wires the handler to the pantry use cases.
func NewPantryHandler(service usecases.PantryService) *PantryHandler {
	return &PantryHandler{service: service}
}

// Create handles POST /v1/pantries.
func (h *PantryHandler) Create(c *gin.Context) {
	var request createPantryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.RespondError(c, http.StatusBadRequest, "invalid_body", "name is required (max 60 characters)")
		return
	}

	pantry, err := h.service.Create(c.Request.Context(), request.Name)
	if err != nil {
		helpers.Respond(c, err)
		return
	}
	c.JSON(http.StatusCreated, toPantryResponse(pantry))
}

// Get handles GET /v1/pantries/:slug.
func (h *PantryHandler) Get(c *gin.Context) {
	slug, ok := parseSlug(c)
	if !ok {
		return
	}

	pantry, items, err := h.service.Get(c.Request.Context(), slug)
	if err != nil {
		helpers.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, pantryDetailResponse{
		pantryResponse: toPantryResponse(pantry),
		Items:          toItemResponses(items),
	})
}

// UpdateItem handles PATCH /v1/pantries/:slug/items/:productId.
func (h *PantryHandler) UpdateItem(c *gin.Context) {
	slug, ok := parseSlug(c)
	if !ok {
		return
	}

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		helpers.RespondError(c, http.StatusBadRequest, "invalid_product_id", "product id must be a uuid")
		return
	}

	var request updateItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.RespondError(c, http.StatusBadRequest, "invalid_body", "request body is malformed")
		return
	}

	item, err := h.service.UpdateItem(c.Request.Context(), slug, productID, entities.ItemPatch{
		Status:   request.Status,
		Type:     request.Type,
		Category: request.Category,
	})
	if err != nil {
		helpers.Respond(c, err)
		return
	}
	c.JSON(http.StatusOK, toItemResponse(*item))
}

// parseSlug validates the :slug path parameter, responding with 400 when it is
// not a well-formed slug so malformed input never reaches the database.
func parseSlug(c *gin.Context) (valueobjects.Slug, bool) {
	slug, err := valueobjects.ParseSlug(c.Param("slug"))
	if err != nil {
		helpers.RespondError(c, http.StatusBadRequest, "invalid_slug", "slug is not valid")
		return "", false
	}
	return slug, true
}
