// Package httptransport adapts HTTP requests to product use cases.
package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/product"
)

type handler struct{ service *product.Service }

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	in, err := readCreateRequest(w, r)
	if err != nil {
		fail(w, err)
		return
	}
	result, err := h.service.Create(product.CreateInput{ProductCode: in.ProductCode, Name: in.Name, Description: in.Description})
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
