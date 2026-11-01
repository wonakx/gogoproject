package product

import (
	"sync"

	"example.com/gogoproject/internal/platform/fault"
)

type MemoryRepository struct {
	mu       sync.RWMutex
	next     int64
	products map[int64]Product
	codes    map[string]int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{products: make(map[int64]Product), codes: make(map[string]int64)}
}
func (r *MemoryRepository) Create(p Product) (Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.codes[p.ProductCode]; ok {
		return Product{}, fault.ErrConflict
	}
	r.next++
	p.ProductNo = r.next
	r.products[p.ProductNo] = clone(p)
	r.codes[p.ProductCode] = p.ProductNo
	return clone(p), nil
}
func (r *MemoryRepository) Get(no int64) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.products[no]
	if !ok {
		return Product{}, fault.ErrNotFound
	}
	return clone(p), nil
}
