package member

import (
"context"

"example.com/gogoproject/internal/platform/fault"
"example.com/gogoproject/internal/platform/mysqlstore"
)

type MySQLRepository struct{ store *mysqlstore.Store }

func NewMySQLRepository(store *mysqlstore.Store) *MySQLRepository {
	return &MySQLRepository{store: store}
}

var _ Repository = (*MySQLRepository)(nil)

const memberColumns = "member_no, id, email, address, phone_number, create_at, update_at"

type scanner interface{ Scan(...any) error }

func scanMember(row scanner) (Member, error) {
	var m Member
	err := row.Scan(&m.MemberNo, &m.ID, &m.Email, &m.Address, &m.PhoneNumber, &m.CreateAt, &m.UpdateAt)
	return m, mysqlstore.Error(err)
}
func (r *MySQLRepository) Create(ctx context.Context, m Member) (Member, error) {
	result, err := r.store.Executor(ctx).ExecContext(ctx, "INSERT INTO dealer.dt_member (id, email, address, phone_number, create_at, update_at) VALUES (?, ?, ?, ?, ?, ?)", m.ID, m.Email, m.Address, m.PhoneNumber, m.CreateAt, m.UpdateAt)
	if err != nil {
		return Member{}, mysqlstore.Error(err)
	}
	m.MemberNo, err = result.LastInsertId()
	return m, err
}
func (r *MySQLRepository) Get(ctx context.Context, no int64) (Member, error) {
	return scanMember(r.store.Executor(ctx).QueryRowContext(ctx, "SELECT "+memberColumns+" FROM dealer.dt_member WHERE member_no = ?", no))
}
func (r *MySQLRepository) Update(ctx context.Context, m Member) (Member, error) {
	err := r.store.Within(ctx, func(ctx context.Context) error {
		old, err := scanMember(r.store.Executor(ctx).QueryRowContext(ctx, "SELECT "+memberColumns+" FROM dealer.dt_member WHERE member_no = ? FOR UPDATE", m.MemberNo))
		if err != nil {
			return err
		}
		_, err = r.store.Executor(ctx).ExecContext(ctx, "UPDATE dealer.dt_member SET id = ?, email = ?, address = ?, phone_number = ?, update_at = ? WHERE member_no = ?", m.ID, m.Email, m.Address, m.PhoneNumber, m.UpdateAt, m.MemberNo)
		if err != nil {
			return mysqlstore.Error(err)
		}
		m.CreateAt = old.CreateAt
		return nil
	})
	if err != nil {
		return Member{}, err
	}
	return m, nil
}
func (r *MySQLRepository) Delete(ctx context.Context, no int64) error {
	result, err := r.store.Executor(ctx).ExecContext(ctx, "DELETE FROM dealer.dt_member WHERE member_no = ?", no)
	if err != nil {
		return mysqlstore.Error(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fault.ErrNotFound
	}
	return nil
}
