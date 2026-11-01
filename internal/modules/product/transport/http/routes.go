package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/product"
)

func RegisterRoutes(mux *http.ServeMux, service *product.Service) {
	h := handler{service: service}
	mux.HandleFunc("POST /products", h.create)
	mux.HandleFunc("GET /products/{product_no}", h.get)
}
