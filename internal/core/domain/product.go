package domain

import (
	"context"
	"time"

	"github.com/padiazg/broken-hexagon/internal/adapters/secondary/redis"
)

// Product is the core business entity.
type Product struct {
	ID        string
	Name      string
	Price     float64
	Active    bool
	CreatedAt time.Time
}

// NewProduct creates a Product with basic validation.
func NewProduct(id, name string, price float64) *Product {
	return &Product{ID: id, Name: name, Price: price, Active: true, CreatedAt: time.Now()}
}

// Deactivate marks the product as inactive.
// BROKEN: domain logic reaches directly into a secondary adapter.
func (p *Product) Deactivate(ctx context.Context) {
	p.Active = false
	redis.InvalidateProduct(ctx, p.ID)
}
