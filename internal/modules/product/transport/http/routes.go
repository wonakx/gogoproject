package httptransport

import (
	"example.com/gogoproject/internal/modules/product"
	"net/http"
)

// RegisterRoutes only connects HTTP routes to module handlers.
func RegisterRoutes(mux *http.ServeMux, service *product.Service) {
	h := handler{service: service}
	mux.HandleFunc("POST /products", h.create)
	mux.HandleFunc("GET /products/{id}", h.get)
}
