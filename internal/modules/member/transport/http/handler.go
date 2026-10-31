// Package httptransport exposes the member module over HTTP.
package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/member"
	"example.com/gogoproject/internal/platform/httpjson"
)

type createRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// RegisterRoutes registers resource paths on the supplied API router.
// The application owns the API version prefix.
func RegisterRoutes(mux *http.ServeMux, service *member.Service) {
	mux.HandleFunc("POST /members", httpjson.Create(func(in createRequest) (any, error) {
		return service.Register(in.Name, in.Email)
	}))
	mux.HandleFunc("GET /members/{id}", httpjson.Get(service.Get))
}
