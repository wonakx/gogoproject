// Package httptransport implements the module's generated OpenAPI operations.
package httptransport

import (
	"context"

	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/modules/member"
	"example.com/gogoproject/internal/platform/fault"
)

type Handler struct{ service *member.Service }

func NewHandler(service *member.Service) *Handler { return &Handler{service: service} }

func (h *Handler) CreateMember(ctx context.Context, request api.CreateMemberRequestObject) (api.CreateMemberResponseObject, error) {
	if request.Body == nil {
		return createMemberError(fault.ErrInvalid), nil
	}
	in := request.Body
	result, err := h.service.Register(ctx, member.CreateInput{ID: in.ID, Email: in.Email, Address: in.Address, PhoneNumber: in.PhoneNumber})
	if err != nil {
		return createMemberError(err), nil
	}
	return api.CreateMember201JSONResponse(toResponse(result)), nil
}

func (h *Handler) GetMember(ctx context.Context, request api.GetMemberRequestObject) (api.GetMemberResponseObject, error) {
	result, err := h.service.Get(ctx, request.MemberNo)
	if err != nil {
		return getMemberError(err), nil
	}
	return api.GetMember200JSONResponse(toResponse(result)), nil
}

func (h *Handler) UpdateMember(ctx context.Context, request api.UpdateMemberRequestObject) (api.UpdateMemberResponseObject, error) {
	if request.Body == nil {
		return updateMemberError(fault.ErrInvalid), nil
	}
	in := request.Body
	result, err := h.service.Update(ctx, member.UpdateInput{MemberNo: request.MemberNo, CreateInput: member.CreateInput{ID: in.ID, Email: in.Email, Address: in.Address, PhoneNumber: in.PhoneNumber}})
	if err != nil {
		return updateMemberError(err), nil
	}
	return api.UpdateMember200JSONResponse(toResponse(result)), nil
}

func (h *Handler) DeleteMember(ctx context.Context, request api.DeleteMemberRequestObject) (api.DeleteMemberResponseObject, error) {
	if err := h.service.Delete(ctx, request.MemberNo); err != nil {
		return deleteMemberError(err), nil
	}
	return api.DeleteMember204Response{}, nil
}
