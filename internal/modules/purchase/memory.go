package purchase

import (
	"example.com/gogoproject/internal/platform/fault"
	"sync"
)

type MemoryRepository struct {
	mu        sync.RWMutex
	purchases map[string]Purchase
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{purchases: make(map[string]Purchase)}
}
func (r *MemoryRepository) Create(p Purchase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.purchases[p.ID] = p
	return nil
}
func (r *MemoryRepository) Get(id string) (Purchase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.purchases[id]
	if !ok {
		return Purchase{}, fault.ErrNotFound
	}
	return p, nil
}
