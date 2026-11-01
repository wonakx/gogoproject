package product

import (
	"strings"
	"time"
	"unicode/utf8"

	"example.com/gogoproject/internal/platform/fault"
)

type Repository interface {
	Create(Product) (Product, error)
	Get(int64) (Product, error)
}
type CreateInput struct {
	ProductCode       string
	Name, Description *string
}
type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) Create(in CreateInput) (Product, error) {
	if strings.TrimSpace(in.ProductCode) == "" || utf8.RuneCountInString(in.ProductCode) > 100 || (in.Name != nil && utf8.RuneCountInString(*in.Name) > 200) {
		return Product{}, fault.ErrInvalid
	}
	return s.repo.Create(Product{ProductCode: in.ProductCode, Name: in.Name, Description: in.Description, CreateAt: time.Now().UTC().Truncate(time.Second)})
}
func (s *Service) Get(no int64) (Product, error) {
	if no <= 0 {
		return Product{}, fault.ErrInvalid
	}
	return s.repo.Get(no)
}
