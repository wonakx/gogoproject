// Package app is the composition root: only this layer connects modules.
package app

import (
	"example.com/gogoproject/internal/modules/member"
	"example.com/gogoproject/internal/modules/product"
	"example.com/gogoproject/internal/modules/purchase"
	"example.com/gogoproject/internal/modules/sale"
)

type App struct {
	Members   *member.Service
	Products  *product.Service
	Sales     *sale.Service
	Purchases *purchase.Service
}

func New() *App {
	members := member.NewService(member.NewMemoryRepository())
	products := product.NewService(product.NewMemoryRepository())
	directory := directoryAdapter{members, products}
	sales := sale.NewService(sale.NewMemoryRepository(), directory)
	purchases := purchase.NewService(purchase.NewMemoryRepository(), directory, salesAdapter{sales})
	return &App{Members: members, Products: products, Sales: sales, Purchases: purchases}
}

type directoryAdapter struct {
	members  *member.Service
	products *product.Service
}

func (a directoryAdapter) RequireMember(no int64) error  { _, err := a.members.Get(no); return err }
func (a directoryAdapter) RequireProduct(no int64) error { _, err := a.products.Get(no); return err }

type salesAdapter struct{ sales *sale.Service }

func (a salesAdapter) Reserve(no, buyer int64, quantity int32) (int64, error) {
	return a.sales.Reserve(no, buyer, quantity)
}
func (a salesAdapter) Release(no int64, quantity int32) error { return a.sales.Release(no, quantity) }
