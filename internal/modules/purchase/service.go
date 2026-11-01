package purchase

import (
"context"
"time"

"example.com/gogoproject/internal/platform/fault"
)

const StatusCompleted = "COMPLETED"

type Members interface {
	RequireMember(context.Context, int64) error
}
type Sales interface {
	Reserve(ctx context.Context, saleNo, buyerNo int64, quantity int32) (sellerNo int64, err error)
}
type Repository interface {
	// CreateWithHistory persists both records atomically or leaves neither on error.
	CreateWithHistory(context.Context, Purchase) (Purchase, error)
	Get(context.Context, int64) (Purchase, error)
	History(context.Context, int64) ([]History, error)
}
type CreateInput struct {
	SaleNo, BuyMemberNo int64
	Quantity            int32
}
type Transactor interface {
	Within(context.Context, func(context.Context) error) error
}

type Service struct {
	tx      Transactor
	repo    Repository
	members Members
	sales   Sales
}

func NewService(repo Repository, members Members, sales Sales, tx Transactor) *Service {
	return &Service{repo: repo, members: members, sales: sales, tx: tx}
}
func (s *Service) Create(ctx context.Context, in CreateInput) (Purchase, error) {
	if in.SaleNo <= 0 || in.BuyMemberNo <= 0 || in.Quantity <= 0 {
		return Purchase{}, fault.ErrInvalid
	}
	var result Purchase
	err := s.tx.Within(ctx, func(ctx context.Context) error {
		if err := s.members.RequireMember(ctx, in.BuyMemberNo); err != nil {
			return err
		}
		seller, err := s.sales.Reserve(ctx, in.SaleNo, in.BuyMemberNo, in.Quantity)
		if err != nil {
			return err
		}
		now := time.Now().UTC().Truncate(time.Second)
		result, err = s.repo.CreateWithHistory(ctx, Purchase{SaleNo: in.SaleNo, SellMemberNo: seller, BuyMemberNo: in.BuyMemberNo, Quantity: in.Quantity, PurchaseStatus: StatusCompleted, PurchaseAt: now, CreateAt: now})
		return err
	})
	if err != nil {
		return Purchase{}, err
	}
	return result, nil
}
func (s *Service) Get(ctx context.Context, no int64) (Purchase, error) {
	if no <= 0 {
		return Purchase{}, fault.ErrInvalid
	}
	return s.repo.Get(ctx, no)
}
func (s *Service) History(ctx context.Context, no int64) ([]History, error) {
	if no <= 0 {
		return nil, fault.ErrInvalid
	}
	return s.repo.History(ctx, no)
}
