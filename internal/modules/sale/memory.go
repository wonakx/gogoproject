package sale

import (
	"math"
	"sync"
	"time"

	"example.com/gogoproject/internal/platform/fault"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	next  int64
	sales map[int64]Sale
}

func NewMemoryRepository() *MemoryRepository { return &MemoryRepository{sales: make(map[int64]Sale)} }
func (r *MemoryRepository) Create(s Sale) (Sale, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	s.SaleNo = r.next
	r.sales[s.SaleNo] = clone(s)
	return clone(s), nil
}
func (r *MemoryRepository) Get(no int64) (Sale, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sales[no]
	if !ok {
		return Sale{}, fault.ErrNotFound
	}
	return clone(s), nil
}
func (r *MemoryRepository) Reserve(no, buyer int64, quantity int32, now time.Time) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.sales[no]
	if !ok {
		return 0, fault.ErrNotFound
	}
	if buyer <= 0 || quantity <= 0 || buyer == s.SellMemberNo {
		return 0, fault.ErrInvalid
	}
	if s.Status != StatusActive || s.Quantity < quantity || now.Before(s.StartAt) || (s.EndAt != nil && !now.Before(*s.EndAt)) {
		return 0, fault.ErrConflict
	}
	s.Quantity -= quantity
	updated := now.UTC().Truncate(time.Second)
	s.UpdateAt = &updated
	r.sales[no] = s
	return s.SellMemberNo, nil
}
func (r *MemoryRepository) Release(no int64, quantity int32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.sales[no]
	if !ok {
		return fault.ErrNotFound
	}
	if quantity <= 0 || s.Quantity > math.MaxInt32-quantity {
		return fault.ErrInvalid
	}
	s.Quantity += quantity
	now := time.Now().UTC().Truncate(time.Second)
	s.UpdateAt = &now
	r.sales[no] = s
	return nil
}
