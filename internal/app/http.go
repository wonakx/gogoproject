package app

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"example.com/gogoproject/internal/platform/fault"
)

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /members", create(func(in struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}) (any, error) { return a.Members.Register(in.Name, in.Email) }))
	mux.HandleFunc("POST /products", create(func(in struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}) (any, error) { return a.Products.Create(in.Name, in.Description) }))
	mux.HandleFunc("POST /sales", create(func(in struct {
		SellerID  string `json:"seller_id"`
		ProductID string `json:"product_id"`
		UnitPrice int64  `json:"unit_price"`
		Stock     int64  `json:"stock"`
	}) (any, error) { return a.Sales.Create(in.SellerID, in.ProductID, in.UnitPrice, in.Stock) }))
	mux.HandleFunc("POST /purchases", create(func(in struct {
		BuyerID  string `json:"buyer_id"`
		SaleID   string `json:"sale_id"`
		Quantity int64  `json:"quantity"`
	}) (any, error) { return a.Purchases.Create(in.BuyerID, in.SaleID, in.Quantity) }))
	mux.HandleFunc("GET /members/{id}", get(a.Members.Get))
	mux.HandleFunc("GET /products/{id}", get(a.Products.Get))
	mux.HandleFunc("GET /sales/{id}", get(a.Sales.Get))
	mux.HandleFunc("GET /purchases/{id}", get(a.Purchases.Get))
	return mux
}
func create[T any](fn func(T) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var in T
		if err := dec.Decode(&in); err != nil {
			fail(w, fault.ErrInvalid)
			return
		}
		if err := dec.Decode(new(any)); err != io.EOF {
			fail(w, fault.ErrInvalid)
			return
		}
		result, err := fn(in)
		if err != nil {
			fail(w, err)
			return
		}
		respond(w, http.StatusCreated, result)
	}
}
func get[T any](fn func(string) (T, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := fn(r.PathValue("id"))
		if err != nil {
			fail(w, err)
			return
		}
		respond(w, http.StatusOK, result)
	}
}
func fail(w http.ResponseWriter, err error) {
	status, message := http.StatusInternalServerError, "internal server error"
	switch {
	case errors.Is(err, fault.ErrInvalid):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, fault.ErrNotFound):
		status, message = http.StatusNotFound, "not found"
	case errors.Is(err, fault.ErrConflict):
		status, message = http.StatusConflict, "duplicate email or insufficient stock"
	default:
		slog.Error("request failed", "error", err)
	}
	respond(w, status, map[string]string{"error": message})
}
func respond(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
