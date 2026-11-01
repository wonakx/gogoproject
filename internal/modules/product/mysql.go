package product

import (
"context"

"example.com/gogoproject/internal/platform/mysqlstore"
)

type MySQLRepository struct{ store *mysqlstore.Store }

func NewMySQLRepository(store *mysqlstore.Store) *MySQLRepository {
	return &MySQLRepository{store: store}
}

var _ Repository = (*MySQLRepository)(nil)

func (r *MySQLRepository) Create(ctx context.Context, p Product) (Product, error) {
	result, err := r.store.Executor(ctx).ExecContext(ctx, "INSERT INTO dealer.dt_product (product_code, name, description, create_at, update_at) VALUES (?, ?, ?, ?, ?)", p.ProductCode, p.Name, p.Description, p.CreateAt, p.UpdateAt)
	if err != nil {
		return Product{}, mysqlstore.Error(err)
	}
	p.ProductNo, err = result.LastInsertId()
	return p, err
}
func (r *MySQLRepository) Get(ctx context.Context, no int64) (Product, error) {
	var p Product
	err := r.store.Executor(ctx).QueryRowContext(ctx, "SELECT product_no, product_code, name, description, create_at, update_at FROM dealer.dt_product WHERE product_no = ?", no).Scan(&p.ProductNo, &p.ProductCode, &p.Name, &p.Description, &p.CreateAt, &p.UpdateAt)
	return p, mysqlstore.Error(err)
}
