// Package httptransport adapts HTTP requests to member use cases.
package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/member"
)

type handler struct{ service *member.Service }

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	in, err := readCreateRequest(w, r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.Register(
		member.CreateInput{
			ID: in.ID,
			Email: in.Email,
			Address: in.Address,
			PhoneNumber: in.PhoneNumber,
		},
	)
	if err != nil {
		fail(w, err)
		return
	}
	respond(w, http.StatusCreated, result)
}
func (h handler) get(w http.ResponseWriter, r *http.Request) {
	no, err := readNumber(r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.Get(no)
	if err != nil {
		fail(w, err)
		return
	}
	respond(w, http.StatusOK, result)
}
func (h handler) update(w http.ResponseWriter, r *http.Request) {
	no, err := readNumber(r)
	if err != nil {
		fail(w, err)
		return
	}
	in, err := readUpdateRequest(w, r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.Update(member.UpdateInput{MemberNo: no, CreateInput: member.CreateInput{ID: in.ID, Email: in.Email, Address: in.Address, PhoneNumber: in.PhoneNumber}})
	if err != nil {
		fail(w, err)
		return
	}
	respond(w, http.StatusOK, result)
}
func (h handler) delete(w http.ResponseWriter, r *http.Request) {
	no, err := readNumber(r)
	if err != nil {
		fail(w, err)
		return
	}
	if err := h.service.Delete(no); err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
