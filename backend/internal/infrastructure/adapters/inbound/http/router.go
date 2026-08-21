// Package http wires the inbound HTTP adapter: middleware, routes and the
// mapping from domain errors to status codes.
package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"midepensa/internal/application/ports"
	"midepensa/internal/application/services"
	"midepensa/internal/application/usecases"
	"midepensa/internal/config"
	"midepensa/internal/infrastructure/adapters/inbound/http/handlers"
	"midepensa/internal/infrastructure/adapters/inbound/http/helpers"
	"midepensa/internal/infrastructure/adapters/inbound/http/middleware"
)

// maxRequestBody is generous for the largest payload the API accepts.
const maxRequestBody = 8 << 10 // 8 KiB

// SetupRouter builds the fully wired gin engine.
func SetupRouter(
	cfg config.Config,
	logger *slog.Logger,
	pantryService usecases.PantryService,
	catalogService usecases.CatalogService,
	checkers ...ports.HealthChecker,
) *gin.Engine {
	registerDomainErrors()

	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.Recovery(logger),
		middleware.Logger(logger),
		middleware.SecurityHeaders(),
		middleware.CORS(cfg.CORS),
		middleware.BodyLimit(maxRequestBody),
	)

	health := handlers.NewHealthHandler(checkers...)
	router.GET("/healthz", health.Healthz)
	router.GET("/readyz", health.Readyz)

	catalog := handlers.NewCatalogHandler(catalogService)
	pantries := handlers.NewPantryHandler(pantryService)

	v1 := router.Group("/v1")
	{
		v1.GET("/catalog", catalog.Get)
		v1.POST("/pantries", pantries.Create)
		v1.GET("/pantries/:slug", pantries.Get)
		v1.PATCH("/pantries/:slug/items/:productId", pantries.UpdateItem)
	}

	return router
}

// registerDomainErrors is the single place where business failures become
// status codes, so handlers never hard-code them.
func registerDomainErrors() {
	helpers.RegisterDomainError(services.ErrPantryNotFound,
		http.StatusNotFound, "pantry_not_found", "pantry not found")
	helpers.RegisterDomainError(services.ErrItemNotFound,
		http.StatusNotFound, "item_not_found", "product is not part of this pantry")
	helpers.RegisterDomainError(services.ErrSlugAlreadyExists,
		http.StatusConflict, "slug_already_exists", "a pantry with that name already exists")
	helpers.RegisterDomainError(services.ErrInvalidPantryName,
		http.StatusBadRequest, "invalid_name", "name must contain at least one letter or digit")
	helpers.RegisterDomainError(services.ErrEmptyPatch,
		http.StatusBadRequest, "empty_update", "provide at least one field to update")
	helpers.RegisterDomainError(services.ErrInvalidPatch,
		http.StatusBadRequest, "invalid_update", "status, type or category has an unknown value")
	helpers.RegisterDomainError(services.ErrEmptyCatalog,
		http.StatusServiceUnavailable, "catalog_unavailable", "product catalog is not seeded")
}
