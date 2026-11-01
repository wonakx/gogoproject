package httptransport

import (
	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/modules/product"
)

// toResponse is the explicit boundary between the domain model and the public API.
func toResponse(m product.Product) api.ProductResponse {
	return api.ProductResponse{ProductNo: m.ProductNo, ProductCode: m.ProductCode, Name: m.Name, Description: m.Description, CreateAt: m.CreateAt, UpdateAt: m.UpdateAt}
}
