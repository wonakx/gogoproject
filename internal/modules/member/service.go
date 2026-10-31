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

func (s *Service) Update(m Member) (Member, error) {

}
