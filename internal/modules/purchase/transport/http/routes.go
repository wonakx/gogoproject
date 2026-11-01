package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/purchase"
)

func RegisterRoutes(mux *http.ServeMux, service *purchase.Service) {
	h := handler{service: service}
	mux.HandleFunc("POST /purchases", h.create)
	mux.HandleFunc("GET /purchases/{purchase_no}", h.get)
	mux.HandleFunc("GET /purchases/{purchase_no}/histories", h.history)
}
