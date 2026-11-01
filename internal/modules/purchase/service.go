package purchase

import (
	"errors"
	"fmt"
	"time"

	"example.com/gogoproject/internal/platform/fault"
)

const StatusCompleted = "COMPLETED"

type Members interface{ RequireMember(int64) error }
type Sales interface {
	Reserve(saleNo, buyerNo int64, quantity int32) (sellerNo int64, err error)
	Release(saleNo int64, quantity int32) error
}
type Repository interface {
	// CreateWithHistory persists both records atomically or leaves neither on error.
	CreateWithHistory(Purchase) (Purchase, error)
	Get(int64) (Purchase, error)
	History(int64) ([]History, error)
}
type CreateInput struct {
	SaleNo, BuyMemberNo int64
	Quantity            int32
}
type Service struct {
	repo    Repository
	members Members
	sales   Sales
}

func NewService(repo Repository, members Members, sales Sales) *Service {
	return &Service{repo: repo, members: members, sales: sales}
}
func (s *Service) Create(in CreateInput) (Purchase, error) {
	if in.SaleNo <= 0 || in.BuyMemberNo <= 0 || in.Quantity <= 0 {
		return Purchase{}, fault.ErrInvalid
	}
	if err := s.members.RequireMember(in.BuyMemberNo); err != nil {
		return Purchase{}, err
	}
	seller, err := s.sales.Reserve(in.SaleNo, in.BuyMemberNo, in.Quantity)
	if err != nil {
		return Purchase{}, err
	}
	now := time.Now().UTC().Truncate(time.Second)
	p, err := s.repo.CreateWithHistory(Purchase{SaleNo: in.SaleNo, SellMemberNo: seller, Quantity: in.Quantity, BuyMemberNo: in.BuyMemberNo, PurchaseStatus: StatusCompleted, PurchaseAt: now, CreateAt: now})
	if err != nil {
		if releaseErr := s.sales.Release(in.SaleNo, in.Quantity); releaseErr != nil {
			return Purchase{}, errors.Join(err, fmt.Errorf("stock compensation failed: %w", releaseErr))
		}
		return Purchase{}, err
	}
	return p, nil
}
func (s *Service) Get(no int64) (Purchase, error) {
	if no <= 0 {
		return Purchase{}, fault.ErrInvalid
	}
	return s.repo.Get(no)
}
func (s *Service) History(no int64) ([]History, error) {
	if no <= 0 {
		return nil, fault.ErrInvalid
	}
	return s.repo.History(no)
}
