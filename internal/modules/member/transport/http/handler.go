// Package httptransport adapts HTTP requests to member use cases.
package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/member"
)

// handler holds the application service; it contains no business policy.
type handler struct{ service *member.Service }

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	in, err := readCreateRequest(w, r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.Register(in.Name, in.Email)
	if err != nil {
		fail(w, err)
		return
	}
	respond(w, http.StatusCreated, result)
}

func (h handler) get(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.Get(r.PathValue("id"))
	if err != nil {
		fail(w, err)
		return
	}
	respond(w, http.StatusOK, result)
}

func (h handler) update(w http.ResponseWriter, r *http.Request) {
	in, err := readUpdateRequest(w, r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.Update(member.UpdateInput{ID: r.PathValue("id"), Name: in.Name, Email: in.Email})
	if err != nil {
		fail(w, err)
		return
	}
	respond(w, http.StatusOK, result)
}

func (h handler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.PathValue("id")); err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
