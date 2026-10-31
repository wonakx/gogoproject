// Package product owns the product catalog, independent of sale prices and stock.
package product

import (
	"example.com/gogoproject/internal/platform/fault"
	"example.com/gogoproject/internal/platform/identity"
	"fmt"
	"strings"
)

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type Repository interface {
	Create(Product) error
	Get(string) (Product, error)
}
type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) Create(name, description string) (Product, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Product{}, fmt.Errorf("%w: name required", fault.ErrInvalid)
	}
	p := Product{ID: identity.New(), Name: name, Description: strings.TrimSpace(description)}
	if err := s.repo.Create(p); err != nil {
		return Product{}, err
	}
	return p, nil
}
func (s *Service) Get(id string) (Product, error) { return s.repo.Get(id) }
