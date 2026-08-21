package services

import (
	"context"

	"midepensa/internal/domain/entities"
	"midepensa/internal/domain/repositories"
)

type catalogService struct {
	products repositories.ProductRepository
}

// NewCatalogService wires the catalog use case to its persistence port.
func NewCatalogService(products repositories.ProductRepository) *catalogService {
	return &catalogService{products: products}
}

// Products returns every catalog entry.
func (s *catalogService) Products(ctx context.Context) ([]entities.Product, error) {
	return s.products.List(ctx)
}
