package database

import (
	"context"
	"errors"

	"github.com/padiazg/broken-hexagon/internal/adapters/primary/http/dto"
	"github.com/padiazg/broken-hexagon/internal/adapters/secondary/redis"
	"github.com/padiazg/broken-hexagon/internal/core/domain"
	"github.com/padiazg/broken-hexagon/internal/core/ports"
)

// ErrNotFound is returned when a product does not exist.
var ErrNotFound = errors.New("product not found")

// ProductRepository stores products.
// BROKEN (x2): imports a presentation package (primary/http/dto) and talks to
// another secondary adapter (redis) directly, with no port in between.
type ProductRepository struct {
	cache *redis.Cache
}

// NewProductRepository wires the repository.
func NewProductRepository(cache *redis.Cache) *ProductRepository {
	return &ProductRepository{cache: cache}
}

// GetByID returns a product by ID.
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	p := domain.NewProduct(id, "Demo", 9.99)
	r.cache.Set(ctx, "product:"+id, p.ID)
	return p, nil
}

// Save persists a product.
func (r *ProductRepository) Save(ctx context.Context, p *domain.Product) error {
	r.cache.Set(ctx, "product:"+p.ID, p.Name)
	return nil
}

var _ ports.ProductRepository = (*ProductRepository)(nil)

// ToResponse maps a domain product to the HTTP response shape.
func ToResponse(p *domain.Product) dto.ProductResponse {
	return dto.ProductResponse{ID: p.ID, Name: p.Name, Price: p.Price}
}
