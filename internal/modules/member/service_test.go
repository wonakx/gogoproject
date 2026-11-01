package member

import (
"context"
"errors"
"strings"
"testing"
"time"

"example.com/gogoproject/internal/platform/fault"
"example.com/gogoproject/internal/platform/mysqlstore"
"github.com/DATA-DOG/go-sqlmock"
"github.com/go-sql-driver/mysql"
)

func TestMySQLMemberNullsAndDuplicate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(NewMySQLRepository(mysqlstore.New(db)))
	ctx := context.Background()
	mock.ExpectExec("INSERT INTO dealer.dt_member").WithArgs("first", nil, nil, nil, sqlmock.AnyArg(), nil).WillReturnResult(sqlmock.NewResult(10, 1))
	m, err := service.Register(ctx, CreateInput{ID: "first"})
	if err != nil || m.MemberNo != 10 || m.Email != nil {
		t.Fatalf("%+v %v", m, err)
	}
	mock.ExpectQuery("SELECT .* FROM dealer.dt_member WHERE member_no").WithArgs(int64(10)).WillReturnRows(sqlmock.NewRows(strings.Split(memberColumns, ", ")).AddRow(10, "first", nil, nil, nil, m.CreateAt, nil))
	saved, err := service.Get(ctx, 10)
	if err != nil || saved.Email != nil || saved.UpdateAt != nil {
		t.Fatalf("%+v %v", saved, err)
	}
	mock.ExpectExec("INSERT INTO dealer.dt_member").WillReturnError(&mysql.MySQLError{Number: 1062})
	if _, err := service.Register(ctx, CreateInput{ID: "first"}); !errors.Is(err, fault.ErrConflict) {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
func TestMySQLMemberNoopUpdateAndDeleteMissing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(NewMySQLRepository(mysqlstore.New(db)))
	now := time.Now().UTC().Truncate(time.Second)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .* FOR UPDATE").WithArgs(int64(10)).WillReturnRows(sqlmock.NewRows(strings.Split(memberColumns, ", ")).AddRow(10, "same", nil, nil, nil, now, now))
	// MySQL reports changed rows, so a valid no-op may affect zero rows.
	mock.ExpectExec("UPDATE dealer.dt_member").WithArgs("same", nil, nil, nil, sqlmock.AnyArg(), int64(10)).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	got, err := service.Update(context.Background(), UpdateInput{MemberNo: 10, CreateInput: CreateInput{ID: "same"}})
	if err != nil || !got.CreateAt.Equal(now) {
		t.Fatalf("%+v %v", got, err)
	}
	mock.ExpectExec("DELETE FROM dealer.dt_member").WithArgs(int64(99)).WillReturnResult(sqlmock.NewResult(0, 0))
	if err := service.Delete(context.Background(), 99); !errors.Is(err, fault.ErrNotFound) {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
func TestValidationWithoutHTTP(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(NewMySQLRepository(mysqlstore.New(db)))
	longEmail := strings.Repeat("a", 201)
	for _, in := range []CreateInput{{ID: ""}, {ID: strings.Repeat("가", 101)}, {ID: "valid", Email: &longEmail}} {
		if _, err := service.Register(context.Background(), in); !errors.Is(err, fault.ErrInvalid) {
			t.Errorf("got %v", err)
		}
	}
	if _, err := service.Update(context.Background(), UpdateInput{MemberNo: 1, CreateInput: CreateInput{ID: ""}}); !errors.Is(err, fault.ErrInvalid) {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("invalid input reached DB:", err)
	}
}
