// Package httptransport implements the module's generated OpenAPI operations.
package httptransport

import (
	"context"

	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/modules/sale"
	"example.com/gogoproject/internal/platform/fault"
)

type Handler struct{ service *sale.Service }

func NewHandler(service *sale.Service) *Handler { return &Handler{service: service} }

func (h *Handler) CreateSale(ctx context.Context, request api.CreateSaleRequestObject) (api.CreateSaleResponseObject, error) {
	if request.Body == nil {
		return createSaleError(fault.ErrInvalid), nil
	}
	in := request.Body
	price, err := sale.ParsePrice(in.Price)
	if err != nil {
		return createSaleError(err), nil
	}
	result, err := h.service.Create(ctx, sale.CreateInput{SellMemberNo: in.SellMemberNo, ProductNo: in.ProductNo, Name: in.Name, Description: in.Description, Quantity: in.Quantity, Price: price, Status: string(in.Status), StartAt: in.StartAt, EndAt: in.EndAt})
	if err != nil {
		return createSaleError(err), nil
	}
	return api.CreateSale201JSONResponse(toResponse(result)), nil
}

func (h *Handler) GetSale(ctx context.Context, request api.GetSaleRequestObject) (api.GetSaleResponseObject, error) {
	result, err := h.service.Get(ctx, request.SaleNo)
	if err != nil {
		return getSaleError(err), nil
	}
	return api.GetSale200JSONResponse(toResponse(result)), nil
}
