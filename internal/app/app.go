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

func (a directoryAdapter) RequireMember(id string) error  { _, err := a.members.Get(id); return err }
func (a directoryAdapter) RequireProduct(id string) error { _, err := a.products.Get(id); return err }

type salesAdapter struct{ sales *sale.Service }

func (a salesAdapter) Reserve(id, buyerID string, quantity int64) (purchase.Reservation, error) {
	s, err := a.sales.Reserve(id, buyerID, quantity)
	if err != nil {
		return purchase.Reservation{}, err
	}
	return purchase.Reservation{ProductID: s.ProductID, SellerID: s.SellerID, UnitPrice: s.UnitPrice}, nil
}
func (a salesAdapter) Release(id string, quantity int64) error { return a.sales.Release(id, quantity) }
