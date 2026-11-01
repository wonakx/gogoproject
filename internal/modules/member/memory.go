package member

import (
	"sync"

	"example.com/gogoproject/internal/platform/fault"
)

type MemoryRepository struct {
	mu      sync.RWMutex
	next    int64
	members map[int64]Member
	ids     map[string]int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{members: make(map[int64]Member), ids: make(map[string]int64)}
}
func (r *MemoryRepository) Create(m Member) (Member, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.ids[m.ID]; ok {
		return Member{}, fault.ErrConflict
	}
	r.next++
	m.MemberNo = r.next
	r.members[m.MemberNo] = clone(m)
	r.ids[m.ID] = m.MemberNo
	return clone(m), nil
}
func (r *MemoryRepository) Get(no int64) (Member, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.members[no]
	if !ok {
		return Member{}, fault.ErrNotFound
	}
	return clone(m), nil
}
func (r *MemoryRepository) Update(m Member) (Member, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	old, ok := r.members[m.MemberNo]
	if !ok {
		return Member{}, fault.ErrNotFound
	}
	if owner, exists := r.ids[m.ID]; exists && owner != m.MemberNo {
		return Member{}, fault.ErrConflict
	}
	delete(r.ids, old.ID)
	m.CreateAt = old.CreateAt
	r.members[m.MemberNo] = clone(m)
	r.ids[m.ID] = m.MemberNo
	return clone(m), nil
}
func (r *MemoryRepository) Delete(no int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.members[no]
	if !ok {
		return fault.ErrNotFound
	}
	delete(r.members, no)
	delete(r.ids, m.ID)
	return nil
}
