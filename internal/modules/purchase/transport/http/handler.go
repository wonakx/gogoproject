// Package httptransport exposes the purchase module over HTTP.
package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/purchase"
	"example.com/gogoproject/internal/platform/httpjson"
)

type createRequest struct {
	BuyerID  string `json:"buyer_id"`
	SaleID   string `json:"sale_id"`
	Quantity int64  `json:"quantity"`
}

// RegisterRoutes registers resource paths on the supplied API router.
// The application owns the API version prefix.
func RegisterRoutes(mux *http.ServeMux, service *purchase.Service) {
	mux.HandleFunc("POST /purchases", httpjson.Create(func(in createRequest) (any, error) {
		return service.Create(in.BuyerID, in.SaleID, in.Quantity)
	}))
	mux.HandleFunc("GET /purchases/{id}", httpjson.Get(service.Get))
}
