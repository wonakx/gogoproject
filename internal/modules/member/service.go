// Package member owns members identified by member_no and unique login id.
package member

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"example.com/gogoproject/internal/platform/fault"
)

type Repository interface {
	Create(Member) (Member, error)
	Update(Member) (Member, error)
	Get(int64) (Member, error)
	Delete(int64) error
}

type CreateInput struct {
	ID          string
	Email       *string
	Address     *string
	PhoneNumber *string
}
type UpdateInput struct {
	MemberNo int64
	CreateInput
}
type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func validate(in CreateInput) error {
	if strings.TrimSpace(in.ID) == "" || utf8.RuneCountInString(in.ID) > 100 {
		return fmt.Errorf("%w: id must contain 1-100 characters", fault.ErrInvalid)
	}
	for _, field := range []struct {
		value *string
		max   int
	}{{in.Email, 200}, {in.Address, 300}, {in.PhoneNumber, 100}} {
		if field.value != nil && utf8.RuneCountInString(*field.value) > field.max {
			return fmt.Errorf("%w: member field exceeds column length", fault.ErrInvalid)
		}
	}
	return nil
}
func (s *Service) Register(in CreateInput) (Member, error) {
	if err := validate(in); err != nil {
		return Member{}, err
	}
	return s.repo.Create(Member{ID: in.ID, Email: in.Email, Address: in.Address, PhoneNumber: in.PhoneNumber, CreateAt: time.Now().UTC().Truncate(time.Second)})
}
func (s *Service) Get(no int64) (Member, error) {
	if no <= 0 {
		return Member{}, fault.ErrInvalid
	}
	return s.repo.Get(no)
}
func (s *Service) Update(in UpdateInput) (Member, error) {
	if in.MemberNo <= 0 {
		return Member{}, fault.ErrInvalid
	}
	if err := validate(in.CreateInput); err != nil {
		return Member{}, err
	}
	now := time.Now().UTC().Truncate(time.Second)
	return s.repo.Update(Member{MemberNo: in.MemberNo, ID: in.ID, Email: in.Email, Address: in.Address, PhoneNumber: in.PhoneNumber, UpdateAt: &now})
}
func (s *Service) Delete(no int64) error {
	if no <= 0 {
		return fault.ErrInvalid
	}
	return s.repo.Delete(no)
}
