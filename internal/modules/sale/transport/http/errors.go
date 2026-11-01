package httptransport

import (
	"errors"
	"log/slog"

	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/platform/fault"
)

func createSaleError(err error) api.CreateSaleResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.CreateSale400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrNotFound):
		return api.CreateSale404JSONResponse{Error: "not found"}
	case errors.Is(err, fault.ErrConflict):
		return api.CreateSale409JSONResponse{Error: "sale is unavailable or quantity is insufficient"}
	default:
		slog.Error("sale request failed", "error", err)
		return api.CreateSale500JSONResponse{Error: "internal server error"}
	}
}

func getSaleError(err error) api.GetSaleResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.GetSale400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrNotFound):
		return api.GetSale404JSONResponse{Error: "not found"}
	default:
		slog.Error("sale request failed", "error", err)
		return api.GetSale500JSONResponse{Error: "internal server error"}
	}
}
