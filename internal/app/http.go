package app

import (
	"context"
	"net/http"

	"example.com/gogoproject/internal/api"
	memberhttp "example.com/gogoproject/internal/modules/member/transport/http"
	producthttp "example.com/gogoproject/internal/modules/product/transport/http"
	purchasehttp "example.com/gogoproject/internal/modules/purchase/transport/http"
	salehttp "example.com/gogoproject/internal/modules/sale/transport/http"
	"example.com/gogoproject/internal/platform/httpapi"
)

// Embedding composes module operations without coupling the modules to each other.
type strictServer struct {
	*memberhttp.Handler
	// Named embedded aliases avoid the four Handler types having identical field names.
	*productOperations
	*saleOperations
	*purchaseOperations
}
type productOperations = producthttp.Handler
type saleOperations = salehttp.Handler
type purchaseOperations = purchasehttp.Handler

var _ api.StrictServerInterface = (*strictServer)(nil)

func (s *strictServer) Health(ctx context.Context, request api.HealthRequestObject) (api.HealthResponseObject, error) {
	return api.Health200JSONResponse{Status: api.Ok}, nil
}
func (a *App) Handler() http.Handler {
	server := &strictServer{
		Handler:            memberhttp.NewHandler(a.Members),
		productOperations:  producthttp.NewHandler(a.Products),
		saleOperations:     salehttp.NewHandler(a.Sales),
		purchaseOperations: purchasehttp.NewHandler(a.Purchases),
	}
	return httpapi.NewHandler(server)
}
