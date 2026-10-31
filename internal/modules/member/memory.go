package member

import (
	"example.com/gogoproject/internal/platform/fault"
	"sync"
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
