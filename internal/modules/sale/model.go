package sale

import (
"time"
)

// Sale maps to dealer.dt_sale.
type Sale struct {
	SaleNo       int64      `db:"sale_no"`
	SellMemberNo int64      `db:"sell_member_no"`
	Name         *string    `db:"name"`
	Description  *string    `db:"description"`
	ProductNo    int64      `db:"product_no"`
	Quantity     int32      `db:"quantity"`
	Price        Price      `db:"price"`
	Status       string     `db:"status"`
	StartAt      time.Time  `db:"start_at"`
	EndAt        *time.Time `db:"end_at"`
	CreateAt     time.Time  `db:"create_at"`
	UpdateAt     *time.Time `db:"update_at"`
}
