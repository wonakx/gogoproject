// Package httptransport adapts HTTP requests to product use cases.
package httptransport

import (
	"example.com/gogoproject/internal/modules/product"
	"net/http"
)

// handler holds the application service; it contains no business policy.
type handler struct{ service *product.Service }

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	in, err := readCreateRequest(w, r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.Create(in.Name, in.Description)
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
