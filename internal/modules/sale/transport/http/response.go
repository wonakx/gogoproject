package httptransport

import (
	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/modules/sale"
)

// toResponse is the explicit boundary between the domain model and the public API.
func toResponse(m sale.Sale) api.SaleResponse {
	return api.SaleResponse{SaleNo: m.SaleNo, SellMemberNo: m.SellMemberNo, ProductNo: m.ProductNo, Name: m.Name, Description: m.Description, Quantity: m.Quantity, Price: m.Price.String(), Status: api.SaleStatus(m.Status), StartAt: m.StartAt, EndAt: m.EndAt, CreateAt: m.CreateAt, UpdateAt: m.UpdateAt}
}
