package purchase

import (
	"sync"

	"example.com/gogoproject/internal/platform/fault"
)

type MemoryRepository struct {
	mu                sync.RWMutex
	next, nextHistory int64
	purchases         map[int64]Purchase
	histories         map[int64][]History
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{purchases: make(map[int64]Purchase), histories: make(map[int64][]History)}
}
func (r *MemoryRepository) CreateWithHistory(p Purchase) (Purchase, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	p.PurchaseNo = r.next
	r.nextHistory++
	no := p.PurchaseNo
	h := History{Quantity: p.Quantity, PurchaseHistoryNo: r.nextHistory, PurchaseNo: &no, SaleNo: p.SaleNo, SellMemberNo: p.SellMemberNo, BuyMemberNo: p.BuyMemberNo, PurchaseStatus: p.PurchaseStatus, PurchaseAt: p.PurchaseAt, CreateAt: p.CreateAt}
	r.purchases[no] = p
	r.histories[no] = []History{h}
	return p, nil
}
func (r *MemoryRepository) Get(no int64) (Purchase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.purchases[no]
	if !ok {
		return Purchase{}, fault.ErrNotFound
	}
	return p, nil
}
func (r *MemoryRepository) History(no int64) ([]History, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.purchases[no]; !ok {
		return nil, fault.ErrNotFound
	}
	result := make([]History, len(r.histories[no]))
	for i, h := range r.histories[no] {
		result[i] = cloneHistory(h)
	}
	return result, nil
}
