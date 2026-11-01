package httptransport

import (
	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/modules/purchase"
)

// toResponse is the explicit boundary between the domain model and the public API.
func toResponse(m purchase.Purchase) api.PurchaseResponse {
	return api.PurchaseResponse{PurchaseNo: m.PurchaseNo, SaleNo: m.SaleNo, SellMemberNo: m.SellMemberNo, BuyMemberNo: m.BuyMemberNo, Quantity: m.Quantity, PurchaseStatus: api.PurchaseResponsePurchaseStatus(m.PurchaseStatus), PurchaseAt: m.PurchaseAt, CreateAt: m.CreateAt}
}

func historyResponse(m purchase.History) api.PurchaseHistoryResponse {
	return api.PurchaseHistoryResponse{PurchaseHistoryNo: m.PurchaseHistoryNo, PurchaseNo: m.PurchaseNo, SaleNo: m.SaleNo, SellMemberNo: m.SellMemberNo, BuyMemberNo: m.BuyMemberNo, Quantity: m.Quantity, PurchaseStatus: api.PurchaseHistoryResponsePurchaseStatus(m.PurchaseStatus), PurchaseAt: m.PurchaseAt, CreateAt: m.CreateAt}
}
