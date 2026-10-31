// Package httptransport adapts HTTP requests to sale use cases.
package httptransport

import (
	"example.com/gogoproject/internal/modules/sale"
	"net/http"
)

// handler holds the application service; it contains no business policy.
type handler struct{ service *sale.Service }

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	in, err := readCreateRequest(w, r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.Create(in.SellerID, in.ProductID, in.UnitPrice, in.Stock)
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
