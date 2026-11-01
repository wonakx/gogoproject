package httptransport

import (
	"errors"
	"log/slog"

	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/platform/fault"
)

func createPurchaseError(err error) api.CreatePurchaseResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.CreatePurchase400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrNotFound):
		return api.CreatePurchase404JSONResponse{Error: "not found"}
	case errors.Is(err, fault.ErrConflict):
		return api.CreatePurchase409JSONResponse{Error: "sale is unavailable or quantity is insufficient"}
	default:
		slog.Error("purchase request failed", "error", err)
		return api.CreatePurchase500JSONResponse{Error: "internal server error"}
	}
}

func getPurchaseError(err error) api.GetPurchaseResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.GetPurchase400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrNotFound):
		return api.GetPurchase404JSONResponse{Error: "not found"}
	default:
		slog.Error("purchase request failed", "error", err)
		return api.GetPurchase500JSONResponse{Error: "internal server error"}
	}
}

func listPurchaseHistoriesError(err error) api.ListPurchaseHistoriesResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.ListPurchaseHistories400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrNotFound):
		return api.ListPurchaseHistories404JSONResponse{Error: "not found"}
	default:
		slog.Error("purchase request failed", "error", err)
		return api.ListPurchaseHistories500JSONResponse{Error: "internal server error"}
	}
}
