package member

import (
"time"
)

// Member maps to dealer.dt_member.
type Member struct {
	MemberNo    int64      `db:"member_no"`
	ID          string     `db:"id"`
	Email       *string    `db:"email"`
	Address     *string    `db:"address"`
	PhoneNumber *string    `db:"phone_number"`
	CreateAt    time.Time  `db:"create_at"`
	UpdateAt    *time.Time `db:"update_at"`
}
