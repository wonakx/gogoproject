package httptransport

import (
	"errors"
	"log/slog"

	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/platform/fault"
)

func createProductError(err error) api.CreateProductResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.CreateProduct400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrConflict):
		return api.CreateProduct409JSONResponse{Error: "duplicate product code"}
	default:
		slog.Error("product request failed", "error", err)
		return api.CreateProduct500JSONResponse{Error: "internal server error"}
	}
}

func getProductError(err error) api.GetProductResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.GetProduct400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrNotFound):
		return api.GetProduct404JSONResponse{Error: "not found"}
	default:
		slog.Error("product request failed", "error", err)
		return api.GetProduct500JSONResponse{Error: "internal server error"}
	}
}
