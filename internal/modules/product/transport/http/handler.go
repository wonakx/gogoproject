// Package httptransport implements the module's generated OpenAPI operations.
package httptransport

import (
	"context"

	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/modules/product"
	"example.com/gogoproject/internal/platform/fault"
)

type Handler struct{ service *product.Service }

func NewHandler(service *product.Service) *Handler { return &Handler{service: service} }

func (h *Handler) CreateProduct(ctx context.Context, request api.CreateProductRequestObject) (api.CreateProductResponseObject, error) {
	if request.Body == nil {
		return createProductError(fault.ErrInvalid), nil
	}
	in := request.Body
	result, err := h.service.Create(ctx, product.CreateInput{ProductCode: in.ProductCode, Name: in.Name, Description: in.Description})
	if err != nil {
		return createProductError(err), nil
	}
	return api.CreateProduct201JSONResponse(toResponse(result)), nil
}

func (h *Handler) GetProduct(ctx context.Context, request api.GetProductRequestObject) (api.GetProductResponseObject, error) {
	result, err := h.service.Get(ctx, request.ProductNo)
	if err != nil {
		return getProductError(err), nil
	}
	return api.GetProduct200JSONResponse(toResponse(result)), nil
}
