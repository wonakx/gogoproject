// Package member owns member registration and member identity.
package member

import (
	"fmt"
	"net/mail"
	"strings"

	"example.com/gogoproject/internal/platform/fault"
	"example.com/gogoproject/internal/platform/identity"
)

type Member struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Repository interface {
	Create(Member) error
	Update(Member) (Member, error)
	Get(string) (Member, error)
	Delete(string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(name, email string) (Member, error) {
	name, email = strings.TrimSpace(name), strings.ToLower(strings.TrimSpace(email))
	address, err := mail.ParseAddress(email)

	if name == "" || err != nil || address.Address != email {
		return Member{}, fmt.Errorf("%w: name and valid email required", fault.ErrInvalid)
	}

	m := Member{ID: identity.New(), Name: name, Email: email}
	if err := s.repo.Create(m); err != nil {
		return Member{}, err
	}
	return m, nil
}

func (s *Service) Get(id string) (Member, error) {
	return s.repo.Get(id)
}

// UpdateInput is independent of HTTP and represents the member update use case.
type UpdateInput struct {
	ID    string
	Name  string
	Email string
}

func (s *Service) Update(in UpdateInput) (Member, error) {
	m := Member{ID: in.ID, Name: in.Name, Email: in.Email}
	m.Name, m.Email = strings.TrimSpace(m.Name), strings.ToLower(strings.TrimSpace(m.Email))
	address, err := mail.ParseAddress(m.Email)
	if m.ID == "" || m.Name == "" || err != nil || address.Address != m.Email {
		return Member{}, fmt.Errorf("%w: id, name and valid email required", fault.ErrInvalid)
	}
	return s.repo.Update(m)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
