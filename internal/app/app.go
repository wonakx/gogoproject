// Package app connects modules through ports and a single transaction scope.
package app

import (
"context"
"database/sql"

"example.com/gogoproject/internal/modules/member"
"example.com/gogoproject/internal/modules/product"
"example.com/gogoproject/internal/modules/purchase"
"example.com/gogoproject/internal/modules/sale"
"example.com/gogoproject/internal/platform/mysqlstore"
)

type App struct {
	Members   *member.Service
	Products  *product.Service
	Sales     *sale.Service
	Purchases *purchase.Service
}

func New(db *sql.DB) *App {
	store := mysqlstore.New(db)

	members := member.NewService(member.NewMySQLRepository(store))
	products := product.NewService(product.NewMySQLRepository(store))
	directory := directoryAdapter{members, products}

	sales := sale.NewService(sale.NewMySQLRepository(store), directory)
	purchases := purchase.NewService(purchase.NewMySQLRepository(store), directory, salesAdapter{sales}, store)

	return &App{Members: members, Products: products, Sales: sales, Purchases: purchases}
}

type directoryAdapter struct {
	members  *member.Service
	products *product.Service
}

func (a directoryAdapter) RequireMember(ctx context.Context, no int64) error {
	_, err := a.members.Get(ctx, no)
	return err
}
func (a directoryAdapter) RequireProduct(ctx context.Context, no int64) error {
	_, err := a.products.Get(ctx, no)
	return err
}

type salesAdapter struct{ sales *sale.Service }

func (a salesAdapter) Reserve(ctx context.Context, no, buyer int64, quantity int32) (int64, error) {
	return a.sales.Reserve(ctx, no, buyer, quantity)
}
