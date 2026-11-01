package purchase

import (
	"time"
)

// Purchase maps to dealer.dt_purchase.
type Purchase struct {
	PurchaseNo     int64     `json:"purchase_no" db:"purchase_no"`
	SaleNo         int64     `json:"sale_no" db:"sale_no"`
	SellMemberNo   int64     `json:"sell_member_no" db:"sell_member_no"`
	BuyMemberNo    int64     `json:"buy_member_no" db:"buy_member_no"`
	Quantity       int32     `json:"quantity" db:"quantity"`
	PurchaseStatus string    `json:"purchase_status" db:"purchase_status"`
	PurchaseAt     time.Time `json:"purchase_at" db:"purchase_at"`
	CreateAt       time.Time `json:"create_at" db:"create_at"`
}
