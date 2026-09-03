package ports

import (
	"context"

	"github.com/padiazg/broken-hexagon/internal/core/domain"
)

// ProductRepository is the outbound port for product persistence.
type ProductRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Save(ctx context.Context, p *domain.Product) error
}
