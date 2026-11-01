package sale

import (
"context"
"time"

"example.com/gogoproject/internal/platform/fault"
"example.com/gogoproject/internal/platform/mysqlstore"
)

type MySQLRepository struct{ store *mysqlstore.Store }

func NewMySQLRepository(store *mysqlstore.Store) *MySQLRepository {
	return &MySQLRepository{store: store}
}

var _ Repository = (*MySQLRepository)(nil)

const saleColumns = "sale_no, sell_member_no, name, description, product_no, quantity, price, status, start_at, end_at, create_at, update_at"

type scanner interface{ Scan(...any) error }

func scanSale(row scanner) (Sale, error) {
	var s Sale
	var decimal string
	err := row.Scan(&s.SaleNo, &s.SellMemberNo, &s.Name, &s.Description, &s.ProductNo, &s.Quantity, &decimal, &s.Status, &s.StartAt, &s.EndAt, &s.CreateAt, &s.UpdateAt)
	if err != nil {
		return Sale{}, mysqlstore.Error(err)
	}
	s.Price, err = ParsePrice(decimal)
	return s, err
}
func (r *MySQLRepository) Create(ctx context.Context, s Sale) (Sale, error) {
	result, err := r.store.Executor(ctx).ExecContext(ctx, "INSERT INTO dealer.dt_sale (sell_member_no, name, description, product_no, quantity, price, status, start_at, end_at, create_at, update_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", s.SellMemberNo, s.Name, s.Description, s.ProductNo, s.Quantity, s.Price.String(), s.Status, s.StartAt, s.EndAt, s.CreateAt, s.UpdateAt)
	if err != nil {
		return Sale{}, mysqlstore.Error(err)
	}
	s.SaleNo, err = result.LastInsertId()
	return s, err
}
func (r *MySQLRepository) Get(ctx context.Context, no int64) (Sale, error) {
	return scanSale(r.store.Executor(ctx).QueryRowContext(ctx, "SELECT "+saleColumns+" FROM dealer.dt_sale WHERE sale_no = ?", no))
}
func (r *MySQLRepository) Reserve(ctx context.Context, no, buyer int64, quantity int32, now time.Time) (int64, error) {
	var seller int64
	err := r.store.Within(ctx, func(ctx context.Context) error {
		s, err := scanSale(r.store.Executor(ctx).QueryRowContext(ctx, "SELECT "+saleColumns+" FROM dealer.dt_sale WHERE sale_no = ? FOR UPDATE", no))
		if err != nil {
			return err
		}
		if buyer <= 0 || quantity <= 0 || buyer == s.SellMemberNo {
			return fault.ErrInvalid
		}
		if s.Status != StatusActive || s.Quantity < quantity || now.Before(s.StartAt) || (s.EndAt != nil && !now.Before(*s.EndAt)) {
			return fault.ErrConflict
		}
		_, err = r.store.Executor(ctx).ExecContext(ctx, "UPDATE dealer.dt_sale SET quantity = quantity - ?, update_at = ? WHERE sale_no = ?", quantity, now.UTC().Truncate(time.Second), no)
		if err != nil {
			return mysqlstore.Error(err)
		}
		seller = s.SellMemberNo
		return nil
	})
	if err != nil {
		return 0, err
	}
	return seller, nil
}
