// Package httptransport implements the module's generated OpenAPI operations.
package httptransport

import (
	"context"

	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/modules/purchase"
	"example.com/gogoproject/internal/platform/fault"
)

type Handler struct{ service *purchase.Service }

func NewHandler(service *purchase.Service) *Handler { return &Handler{service: service} }

func (h *Handler) CreatePurchase(ctx context.Context, request api.CreatePurchaseRequestObject) (api.CreatePurchaseResponseObject, error) {
	if request.Body == nil {
		return createPurchaseError(fault.ErrInvalid), nil
	}
	in := request.Body
	result, err := h.service.Create(ctx, purchase.CreateInput{SaleNo: in.SaleNo, BuyMemberNo: in.BuyMemberNo, Quantity: in.Quantity})
	if err != nil {
		return createPurchaseError(err), nil
	}
	return api.CreatePurchase201JSONResponse(toResponse(result)), nil
}

func (h *Handler) GetPurchase(ctx context.Context, request api.GetPurchaseRequestObject) (api.GetPurchaseResponseObject, error) {
	result, err := h.service.Get(ctx, request.PurchaseNo)
	if err != nil {
		return getPurchaseError(err), nil
	}
	return api.GetPurchase200JSONResponse(toResponse(result)), nil
}

func (h *Handler) ListPurchaseHistories(ctx context.Context, request api.ListPurchaseHistoriesRequestObject) (api.ListPurchaseHistoriesResponseObject, error) {
	result, err := h.service.History(ctx, request.PurchaseNo)
	if err != nil {
		return listPurchaseHistoriesError(err), nil
	}
	out := make(api.ListPurchaseHistories200JSONResponse, len(result))
	for i, item := range result {
		out[i] = historyResponse(item)
	}
	return out, nil
}
