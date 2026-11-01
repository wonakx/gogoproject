package sale

import (
	"time"
	"unicode/utf8"

	"example.com/gogoproject/internal/platform/fault"
)

const StatusActive = "ACTIVE"
const StatusClosed = "CLOSED"

type Directory interface {
	RequireProduct(int64) error
	RequireMember(int64) error
}
type Repository interface {
	Create(Sale) (Sale, error)
	Get(int64) (Sale, error)
	// Reserve checks buyer and availability and decrements quantity atomically.
	Reserve(saleNo, buyerNo int64, quantity int32, now time.Time) (sellerNo int64, err error)
	Release(saleNo int64, quantity int32) error
}
type CreateInput struct {
	Name, Description       *string
	ProductNo, SellMemberNo int64
	Quantity                int32
	Price                   Price
	Status                  string
	StartAt                 time.Time
	EndAt                   *time.Time
}
type Service struct {
	repo      Repository
	directory Directory
}

func NewService(repo Repository, directory Directory) *Service {
	return &Service{repo: repo, directory: directory}
}
func (s *Service) Create(in CreateInput) (Sale, error) {
	if in.SellMemberNo <= 0 || in.ProductNo <= 0 || in.Quantity <= 0 || in.Price <= 0 || in.Price > MaxPrice || (in.Status != StatusActive && in.Status != StatusClosed) || in.StartAt.IsZero() || (in.Name != nil && utf8.RuneCountInString(*in.Name) > 200) || (in.EndAt != nil && !in.EndAt.After(in.StartAt)) {
		return Sale{}, fault.ErrInvalid
	}
	if err := s.directory.RequireProduct(in.ProductNo); err != nil {
		return Sale{}, err
	}
	if err := s.directory.RequireMember(in.SellMemberNo); err != nil {
		return Sale{}, err
	}
	return s.repo.Create(Sale{SellMemberNo: in.SellMemberNo, Name: in.Name, Description: in.Description, ProductNo: in.ProductNo, Quantity: in.Quantity, Price: in.Price, Status: in.Status, StartAt: in.StartAt, EndAt: in.EndAt, CreateAt: time.Now().UTC().Truncate(time.Second)})
}
func (s *Service) Get(no int64) (Sale, error) {
	if no <= 0 {
		return Sale{}, fault.ErrInvalid
	}
	return s.repo.Get(no)
}
func (s *Service) Reserve(no, buyer int64, quantity int32) (int64, error) {
	if no <= 0 || buyer <= 0 || quantity <= 0 {
		return 0, fault.ErrInvalid
	}
	return s.repo.Reserve(no, buyer, quantity, time.Now().UTC())
}
func (s *Service) Release(no int64, quantity int32) error { return s.repo.Release(no, quantity) }
