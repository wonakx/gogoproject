package member

import (
	"errors"
	"strings"
	"testing"

	"example.com/gogoproject/internal/platform/fault"
)

func TestSchemaAndUniqueness(t *testing.T) {
	service := NewService(NewMemoryRepository())
	email := "same@example.com"
	first, err := service.Register(CreateInput{ID: "first", Email: &email})
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Register(CreateInput{ID: "second", Email: &email})
	if err != nil {
		t.Fatal("email is not unique in schema:", err)
	}
	if first.MemberNo != 1 || second.MemberNo != 2 || first.CreateAt.IsZero() || first.UpdateAt != nil {
		t.Fatal("invalid generated fields")
	}
	if _, err := service.Register(CreateInput{ID: "first"}); !errors.Is(err, fault.ErrConflict) {
		t.Fatal(err)
	}
	*first.Email = "outside mutation"
	saved, _ := service.Get(first.MemberNo)
	if *saved.Email != email {
		t.Fatal("pointer escaped repository")
	}
	updated, err := service.Update(UpdateInput{MemberNo: first.MemberNo, CreateInput: CreateInput{ID: "renamed"}})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Email != nil || updated.UpdateAt == nil || !updated.CreateAt.Equal(first.CreateAt) {
		t.Fatalf("wrong nullable update: %+v", updated)
	}
	if _, err := service.Update(UpdateInput{MemberNo: first.MemberNo, CreateInput: CreateInput{ID: "second"}}); !errors.Is(err, fault.ErrConflict) {
		t.Fatal(err)
	}
	if _, err := service.Register(CreateInput{ID: "first"}); err != nil {
		t.Fatal("old id not released", err)
	}
	if err := service.Delete(first.MemberNo); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Register(CreateInput{ID: "renamed"}); err != nil {
		t.Fatal("deleted id not released", err)
	}
}
func TestValidationWithoutHTTP(t *testing.T) {
	s := NewService(NewMemoryRepository())
	longEmail := strings.Repeat("a", 201)
	for _, in := range []CreateInput{{ID: ""}, {ID: strings.Repeat("가", 101)}, {ID: "valid", Email: &longEmail}} {
		if _, err := s.Register(in); !errors.Is(err, fault.ErrInvalid) {
			t.Errorf("got %v", err)
		}
	}
	m, err := s.Register(CreateInput{ID: strings.Repeat("가", 100)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(UpdateInput{MemberNo: m.MemberNo, CreateInput: CreateInput{ID: ""}}); !errors.Is(err, fault.ErrInvalid) {
		t.Fatal(err)
	}
	saved, _ := s.Get(m.MemberNo)
	if saved.ID != m.ID {
		t.Fatal("invalid update mutated member")
	}
}
