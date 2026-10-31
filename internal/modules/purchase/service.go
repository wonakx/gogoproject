// Package purchase owns orders and orchestrates stock reservation through ports.
package purchase

import (
	"errors"
	"example.com/gogoproject/internal/platform/fault"
	"example.com/gogoproject/internal/platform/identity"
	"fmt"
)

type Purchase struct {
	ID         string `json:"id"`
	BuyerID    string `json:"buyer_id"`
	SaleID     string `json:"sale_id"`
	ProductID  string `json:"product_id"`
	SellerID   string `json:"seller_id"`
	Quantity   int64  `json:"quantity"`
	UnitPrice  int64  `json:"unit_price"`
	TotalPrice int64  `json:"total_price"`
}
type Reservation struct {
	ProductID, SellerID string
	UnitPrice           int64
}
type Members interface{ RequireMember(string) error }
type Sales interface {
	Reserve(saleID, buyerID string, quantity int64) (Reservation, error)
	Release(saleID string, quantity int64) error
}
type Repository interface {
	// Create must leave no order persisted if it returns an error.
	Create(Purchase) error
	Get(string) (Purchase, error)
}
type Service struct {
	repo    Repository
	members Members
	sales   Sales
}

func NewService(repo Repository, members Members, sales Sales) *Service {
	return &Service{repo: repo, members: members, sales: sales}
}
func (s *Service) Create(buyerID, saleID string, quantity int64) (Purchase, error) {
	if buyerID == "" || saleID == "" || quantity <= 0 {
		return Purchase{}, fault.ErrInvalid
	}
	if err := s.members.RequireMember(buyerID); err != nil {
		return Purchase{}, err
	}
	reserved, err := s.sales.Reserve(saleID, buyerID, quantity)
	if err != nil {
		return Purchase{}, err
	}
	p := Purchase{ID: identity.New(), BuyerID: buyerID, SaleID: saleID, ProductID: reserved.ProductID, SellerID: reserved.SellerID, Quantity: quantity, UnitPrice: reserved.UnitPrice, TotalPrice: reserved.UnitPrice * quantity}
	if err := s.repo.Create(p); err != nil {
		if releaseErr := s.sales.Release(saleID, quantity); releaseErr != nil {
			return Purchase{}, errors.Join(err, fmt.Errorf("stock compensation failed: %w", releaseErr))
		}
		return Purchase{}, err
	}
	return p, nil
}
func (s *Service) Get(id string) (Purchase, error) { return s.repo.Get(id) }
