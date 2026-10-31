package httptransport

import (
	"example.com/gogoproject/internal/modules/purchase"
	"net/http"
)

// RegisterRoutes only connects HTTP routes to module handlers.
func RegisterRoutes(mux *http.ServeMux, service *purchase.Service) {
	h := handler{service: service}
	mux.HandleFunc("POST /purchases", h.create)
	mux.HandleFunc("GET /purchases/{id}", h.get)
}
