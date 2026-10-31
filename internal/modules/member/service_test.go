package member

import (
	"errors"
	"testing"

	"example.com/gogoproject/internal/platform/fault"
)

// Business rules must also apply to callers that do not use the HTTP adapter.
func TestUpdateValidatesWithoutHTTP(t *testing.T) {
	service := NewService(NewMemoryRepository())
	original, err := service.Register("Member", "member@example.com")
	if err != nil {
		t.Fatal(err)
	}
	for _, in := range []UpdateInput{
		{ID: original.ID, Name: " ", Email: "member@example.com"},
		{ID: original.ID, Name: "Member", Email: "invalid"},
		{ID: "", Name: "Member", Email: "member@example.com"},
	} {
		if _, err := service.Update(in); !errors.Is(err, fault.ErrInvalid) {
			t.Errorf("got %v, want invalid input", err)
		}
	}
	saved, err := service.Get(original.ID)
	if err != nil || saved != original {
		t.Fatalf("invalid update changed member: %+v %v", saved, err)
	}
	updated, err := service.Update(UpdateInput{ID: original.ID, Name: " Updated ", Email: " UPDATED@example.com "})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "Updated" || updated.Email != "updated@example.com" {
		t.Fatalf("not normalized: %+v", updated)
	}
}
