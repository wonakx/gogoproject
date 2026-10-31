// Package httptransport exposes the sale module over HTTP.
package httptransport

import (
	"net/http"

	"example.com/gogoproject/internal/modules/sale"
	"example.com/gogoproject/internal/platform/httpjson"
)

type createRequest struct {
	SellerID  string `json:"seller_id"`
	ProductID string `json:"product_id"`
	UnitPrice int64  `json:"unit_price"`
	Stock     int64  `json:"stock"`
}

// RegisterRoutes registers resource paths on the supplied API router.
// The application owns the API version prefix.
func RegisterRoutes(mux *http.ServeMux, service *sale.Service) {
	mux.HandleFunc("POST /sales", httpjson.Create(func(in createRequest) (any, error) {
		return service.Create(in.SellerID, in.ProductID, in.UnitPrice, in.Stock)
	}))
	mux.HandleFunc("GET /sales/{id}", httpjson.Get(service.Get))
}
