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
	r.mu.Lock()
	defer r.mu.Unlock()
	previous, ok := r.members[m.ID]
	if !ok {
		return Member{}, fault.ErrNotFound
	}
	if owner, exists := r.emails[m.Email]; exists && owner != m.ID {
		return Member{}, fault.ErrConflict
	}
	delete(r.emails, previous.Email)
	r.members[m.ID] = m
	r.emails[m.Email] = m.ID
	return m, nil
}

func (r *MemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.members[id]
	if !ok {
		return fault.ErrNotFound
	}
	delete(r.members, id)
	delete(r.emails, m.Email)
	return nil
}
