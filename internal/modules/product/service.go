package product

import (
"context"
"strings"
"time"
"unicode/utf8"

"example.com/gogoproject/internal/platform/fault"
)

type Repository interface {
	Create(context.Context, Product) (Product, error)
	Get(context.Context, int64) (Product, error)
}
type CreateInput struct {
	ProductCode       string
	Name, Description *string
}
type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, in CreateInput) (Product, error) {
	if strings.TrimSpace(in.ProductCode) == "" || utf8.RuneCountInString(in.ProductCode) > 100 || (in.Name != nil && utf8.RuneCountInString(*in.Name) > 200) {
		return Product{}, fault.ErrInvalid
	}
	return s.repo.Create(ctx, Product{ProductCode: in.ProductCode, Name: in.Name, Description: in.Description, CreateAt: time.Now().UTC().Truncate(time.Second)})
}

func (s *Service) Get(ctx context.Context, no int64) (Product, error) {
	if no <= 0 {
		return Product{}, fault.ErrInvalid
	}
	return s.repo.Get(ctx, no)
}
