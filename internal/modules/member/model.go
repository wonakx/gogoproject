package member

import (
	"time"

	"example.com/gogoproject/internal/platform/value"
)

// Member maps to dealer.dt_member.
type Member struct {
	MemberNo    int64      `json:"member_no" db:"member_no"`
	ID          string     `json:"id" db:"id"`
	Email       *string    `json:"email" db:"email"`
	Address     *string    `json:"address" db:"address"`
	PhoneNumber *string    `json:"phone_number" db:"phone_number"`
	CreateAt    time.Time  `json:"create_at" db:"create_at"`
	UpdateAt    *time.Time `json:"update_at" db:"update_at"`
}

func clone(v Member) Member {
	v.Email = value.Copy(v.Email)
	v.Address = value.Copy(v.Address)
	v.PhoneNumber = value.Copy(v.PhoneNumber)
	v.UpdateAt = value.Copy(v.UpdateAt)
	return v
}
