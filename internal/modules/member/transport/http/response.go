package httptransport

import (
	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/modules/member"
)

// toResponse is the explicit boundary between the domain model and the public API.
func toResponse(m member.Member) api.MemberResponse {
	return api.MemberResponse{MemberNo: m.MemberNo, ID: m.ID, Email: m.Email, Address: m.Address, PhoneNumber: m.PhoneNumber, CreateAt: m.CreateAt, UpdateAt: m.UpdateAt}
}
