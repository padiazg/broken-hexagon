package redis

import "context"

// Cache is a stand-in for a Redis client.
type Cache struct{}

// New creates an in-memory cache.
func New() *Cache { return &Cache{} }

// Set stores a value.
func (c *Cache) Set(ctx context.Context, key, value string) {
	_ = ctx
	_ = key
	_ = value
}

// InvalidateProduct drops a product cache entry.
func InvalidateProduct(ctx context.Context, id string) {
	_ = ctx
	_ = id
}
