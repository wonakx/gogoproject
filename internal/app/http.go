package app

import (
	"encoding/json"
	"net/http"

	memberhttp "example.com/gogoproject/internal/modules/member/transport/http"
	producthttp "example.com/gogoproject/internal/modules/product/transport/http"
	purchasehttp "example.com/gogoproject/internal/modules/purchase/transport/http"
	salehttp "example.com/gogoproject/internal/modules/sale/transport/http"
)

const apiPrefix = "/api/v1"

func (a *App) Handler() http.Handler {
	api := http.NewServeMux()
	memberhttp.RegisterRoutes(api, a.Members)
	producthttp.RegisterRoutes(api, a.Products)
	salehttp.RegisterRoutes(api, a.Sales)
	purchasehttp.RegisterRoutes(api, a.Purchases)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	mux.Handle(apiPrefix+"/", http.StripPrefix(apiPrefix, api))
	return mux
}
