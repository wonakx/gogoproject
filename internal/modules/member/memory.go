package member

import (
	"sync"

	"example.com/gogoproject/internal/platform/fault"
)

type MemoryRepository struct {
	mu      sync.RWMutex
	members map[string]Member
	emails  map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{members: make(map[string]Member), emails: make(map[string]string)}
}
func (r *MemoryRepository) Create(m Member) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.emails[m.Email]; ok {
		return fault.ErrConflict
	}
	r.members[m.ID] = m
	r.emails[m.Email] = m.ID
	return nil
}
func (r *MemoryRepository) Get(id string) (Member, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.members[id]
	if !ok {
		return Member{}, fault.ErrNotFound
	}
	return m, nil
}

func (r *MemoryRepository) Update(m Member) (Member, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.members[m.ID]; !ok {
		r.members[m.ID] = m
		return r.members[m.ID], nil
	}

	return Member{}, fault.ErrNotFound
}

func (r *MemoryRepository) Delete(id string) (Member, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.members[id]; !ok {
		delete(r.members, id)
		return r.members[id], nil
	}

	return Member{}, fault.ErrNotFound
}
