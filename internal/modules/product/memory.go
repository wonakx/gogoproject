package product

import (
	"example.com/gogoproject/internal/platform/fault"
	"sync"
)

type MemoryRepository struct {
	mu       sync.RWMutex
	products map[string]Product
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{products: make(map[string]Product)}
}
func (r *MemoryRepository) Create(p Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.products[p.ID] = p
	return nil
}
func (r *MemoryRepository) Get(id string) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.products[id]
	if !ok {
		return Product{}, fault.ErrNotFound
	}
	return p, nil
}
