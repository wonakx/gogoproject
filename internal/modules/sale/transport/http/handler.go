// Package httptransport adapts HTTP requests to sale use cases.
package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/sale"
)

type handler struct{ service *sale.Service }

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	in, err := readCreateRequest(w, r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.Create(sale.CreateInput{SellMemberNo: in.SellMemberNo, ProductNo: in.ProductNo, Name: in.Name, Description: in.Description, Quantity: in.Quantity, Price: in.Price, Status: in.Status, StartAt: in.StartAt, EndAt: in.EndAt})
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
