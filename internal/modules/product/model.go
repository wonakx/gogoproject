package product

import (
	"time"

	"example.com/gogoproject/internal/platform/value"
)

// Product maps to dealer.dt_product.
type Product struct {
	ProductNo   int64      `json:"product_no" db:"product_no"`
	ProductCode string     `json:"product_code" db:"product_code"`
	Name        *string    `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	CreateAt    time.Time  `json:"create_at" db:"create_at"`
	UpdateAt    *time.Time `json:"update_at" db:"update_at"`
}

func clone(v Product) Product {
	v.Name = value.Copy(v.Name)
	v.Description = value.Copy(v.Description)
	v.UpdateAt = value.Copy(v.UpdateAt)
	return v
}
