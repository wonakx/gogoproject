// Package httptransport exposes the product module over HTTP.
package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/product"
	"example.com/gogoproject/internal/platform/httpjson"
)

type createRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RegisterRoutes registers resource paths on the supplied API router.
// The application owns the API version prefix.
func RegisterRoutes(mux *http.ServeMux, service *product.Service) {
	mux.HandleFunc("POST /products", httpjson.Create(func(in createRequest) (any, error) {
		return service.Create(in.Name, in.Description)
	}))
	mux.HandleFunc("GET /products/{id}", httpjson.Get(service.Get))
}
