package product

import (
"time"
)

// Product maps to dealer.dt_product.
type Product struct {
	ProductNo   int64      `db:"product_no"`
	ProductCode string     `db:"product_code"`
	Name        *string    `db:"name"`
	Description *string    `db:"description"`
	CreateAt    time.Time  `db:"create_at"`
	UpdateAt    *time.Time `db:"update_at"`
}
