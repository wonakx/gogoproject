// Package sale owns listings, prices and available stock.
package sale

import (
	"example.com/gogoproject/internal/platform/fault"
	"example.com/gogoproject/internal/platform/identity"
	"fmt"
)

type Sale struct {
	ID        string `json:"id"`
	SellerID  string `json:"seller_id"`
	ProductID string `json:"product_id"`
	UnitPrice int64  `json:"unit_price"`
	Stock     int64  `json:"stock"`
}

// Directory is a consumer-owned port. The composition root supplies its adapter.
type Directory interface {
	RequireMember(string) error
	RequireProduct(string) error
}
type Repository interface {
	Create(Sale) error
	Get(string) (Sale, error)
	// Reserve must atomically validate buyer/quantity and decrement stock.
	Reserve(id, buyerID string, quantity int64) (Sale, error)
	Release(id string, quantity int64) error
}
type Service struct {
	repo      Repository
	directory Directory
}

func NewService(repo Repository, directory Directory) *Service {
	return &Service{repo: repo, directory: directory}
}
func (s *Service) Create(sellerID, productID string, price, stock int64) (Sale, error) {
	if sellerID == "" || productID == "" || price <= 0 || stock <= 0 {
		return Sale{}, fmt.Errorf("%w: seller, product, positive price and stock required", fault.ErrInvalid)
	}
	if err := s.directory.RequireMember(sellerID); err != nil {
		return Sale{}, err
	}
	if err := s.directory.RequireProduct(productID); err != nil {
		return Sale{}, err
	}
	sale := Sale{ID: identity.New(), SellerID: sellerID, ProductID: productID, UnitPrice: price, Stock: stock}
	if err := s.repo.Create(sale); err != nil {
		return Sale{}, err
	}
	return sale, nil
}
func (s *Service) Get(id string) (Sale, error) { return s.repo.Get(id) }
func (s *Service) Reserve(id, buyerID string, quantity int64) (Sale, error) {
	if quantity <= 0 {
		return Sale{}, fault.ErrInvalid
	}
	return s.repo.Reserve(id, buyerID, quantity)
}
func (s *Service) Release(id string, quantity int64) error {
	if quantity <= 0 {
		return fault.ErrInvalid
	}
	return s.repo.Release(id, quantity)
}
