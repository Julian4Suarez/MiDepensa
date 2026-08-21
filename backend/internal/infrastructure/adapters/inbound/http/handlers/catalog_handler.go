package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"midepensa/internal/application/usecases"
	"midepensa/internal/domain/entities"
	"midepensa/internal/infrastructure/adapters/inbound/http/helpers"
)

// CatalogHandler exposes the read-only product catalog.
type CatalogHandler struct {
	service usecases.CatalogService
}

// NewCatalogHandler wires the handler to the catalog use case.
func NewCatalogHandler(service usecases.CatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

// Get handles GET /v1/catalog.
func (h *CatalogHandler) Get(c *gin.Context) {
	products, err := h.service.Products(c.Request.Context())
	if err != nil {
		helpers.Respond(c, err)
		return
	}

	response := catalogResponse{
		Products:   make([]productResponse, 0, len(products)),
		Categories: enumValues(entities.Categories),
		Views:      enumValues(entities.PantryViews),
		Statuses:   enumValues(entities.StockStatuses),
	}
	for _, product := range products {
		response.Products = append(response.Products, toProductResponse(product))
	}
	c.JSON(http.StatusOK, response)
}

func enumValues[T ~string](values []T) []string {
	list := make([]string, 0, len(values))
	for _, value := range values {
		list = append(list, string(value))
	}
	return list
}
