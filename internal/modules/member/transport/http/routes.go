package httptransport

import (
	"example.com/gogoproject/internal/modules/member"
	"net/http"
)

// RegisterRoutes only connects HTTP routes to module handlers.
func RegisterRoutes(mux *http.ServeMux, service *member.Service) {
	h := handler{service: service}
	mux.HandleFunc("POST /members", h.create)
	mux.HandleFunc("GET /members/{id}", h.get)
	mux.HandleFunc("PUT /members/{id}", h.update)
	mux.HandleFunc("DELETE /members/{id}", h.delete)
}
