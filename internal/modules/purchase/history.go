package purchase

import (
"time"
)

// History maps to dealer.dt_purchase_history and belongs to the purchase module.
type History struct {
	PurchaseHistoryNo int64     `db:"purchase_history_no"`
	PurchaseNo        *int64    `db:"purchase_no"`
	SaleNo            int64     `db:"sale_no"`
	SellMemberNo      int64     `db:"sell_member_no"`
	BuyMemberNo       int64     `db:"buy_member_no"`
	Quantity          int32     `db:"quantity"`
	PurchaseStatus    string    `db:"purchase_status"`
	PurchaseAt        time.Time `db:"purchase_at"`
	CreateAt          time.Time `db:"create_at"`
}
