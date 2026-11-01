package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/sale"
)

func RegisterRoutes(mux *http.ServeMux, service *sale.Service) {
	h := handler{service: service}
	mux.HandleFunc("POST /sales", h.create)
	mux.HandleFunc("GET /sales/{sale_no}", h.get)
}
