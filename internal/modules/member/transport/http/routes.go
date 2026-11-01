package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/member"
)

func RegisterRoutes(mux *http.ServeMux, service *member.Service) {
	h := handler{service: service}
	mux.HandleFunc("POST /members", h.create)
	mux.HandleFunc("GET /members/{member_no}", h.get)
	mux.HandleFunc("PUT /members/{member_no}", h.update)
	mux.HandleFunc("DELETE /members/{member_no}", h.delete)
}
