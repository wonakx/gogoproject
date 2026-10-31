package sale

import (
	"example.com/gogoproject/internal/platform/fault"
	"math"
	"sync"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	sales map[string]Sale
}

func NewMemoryRepository() *MemoryRepository { return &MemoryRepository{sales: make(map[string]Sale)} }
func (r *MemoryRepository) Create(s Sale) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sales[s.ID] = s
	return nil
}
func (r *MemoryRepository) Get(id string) (Sale, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sales[id]
	if !ok {
		return Sale{}, fault.ErrNotFound
	}
	return s, nil
}
func (r *MemoryRepository) Reserve(id, buyerID string, quantity int64) (Sale, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.sales[id]
	if !ok {
		return Sale{}, fault.ErrNotFound
	}
	if quantity <= 0 || buyerID == "" || buyerID == s.SellerID || quantity > math.MaxInt64/s.UnitPrice {
		return Sale{}, fault.ErrInvalid
	}
	if s.Stock < quantity {
		return Sale{}, fault.ErrConflict
	}
	s.Stock -= quantity
	r.sales[id] = s
	return s, nil
}
func (r *MemoryRepository) Release(id string, quantity int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.sales[id]
	if !ok {
		return fault.ErrNotFound
	}
	if quantity <= 0 || s.Stock > math.MaxInt64-quantity {
		return fault.ErrInvalid
	}
	s.Stock += quantity
	r.sales[id] = s
	return nil
}
