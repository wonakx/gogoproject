package httptransport

import (
	"example.com/gogoproject/internal/modules/sale"
	"net/http"
)

// RegisterRoutes only connects HTTP routes to module handlers.
func RegisterRoutes(mux *http.ServeMux, service *sale.Service) {
	h := handler{service: service}
	mux.HandleFunc("POST /sales", h.create)
	mux.HandleFunc("GET /sales/{id}", h.get)
}
