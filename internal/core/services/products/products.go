package products

import (
	"context"

	"github.com/padiazg/broken-hexagon/internal/adapters/secondary/database"
	"github.com/padiazg/broken-hexagon/internal/core/domain"
)

// Service manages product use cases.
// BROKEN: depends on the concrete secondary adapter instead of the
// ProductRepository port in internal/core/ports.
type Service struct {
	repo *database.ProductRepository
}

// New wires the service.
func New(repo *database.ProductRepository) *Service {
	return &Service{repo: repo}
}

// Get returns a product by ID.
func (s *Service) Get(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}
