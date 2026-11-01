// Package httptransport adapts HTTP requests to purchase use cases.
package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/purchase"
)

type handler struct{ service *purchase.Service }

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	in, err := readCreateRequest(w, r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.Create(purchase.CreateInput{SaleNo: in.SaleNo, BuyMemberNo: in.BuyMemberNo, Quantity: in.Quantity})
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
func (h handler) history(w http.ResponseWriter, r *http.Request) {
	no, err := readNumber(r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.History(no)
	if err != nil {
		fail(w, err)
		return
	}
	respond(w, http.StatusOK, result)
}
