package sale

import (
	"time"

	"example.com/gogoproject/internal/platform/value"
)

// Sale maps to dealer.dt_sale.
type Sale struct {
	SaleNo       int64      `json:"sale_no" db:"sale_no"`
	SellMemberNo int64      `json:"sell_member_no" db:"sell_member_no"`
	Name         *string    `json:"name" db:"name"`
	Description  *string    `json:"description" db:"description"`
	ProductNo    int64      `json:"product_no" db:"product_no"`
	Quantity     int32      `json:"quantity" db:"quantity"`
	Price        Price      `json:"price" db:"price"`
	Status       string     `json:"status" db:"status"`
	StartAt      time.Time  `json:"start_at" db:"start_at"`
	EndAt        *time.Time `json:"end_at" db:"end_at"`
	CreateAt     time.Time  `json:"create_at" db:"create_at"`
	UpdateAt     *time.Time `json:"update_at" db:"update_at"`
}

func clone(v Sale) Sale {
	v.Name = value.Copy(v.Name)
	v.Description = value.Copy(v.Description)
	v.EndAt = value.Copy(v.EndAt)
	v.UpdateAt = value.Copy(v.UpdateAt)
	return v
}
