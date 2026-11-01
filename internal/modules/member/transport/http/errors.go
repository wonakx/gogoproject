package httptransport

import (
	"errors"
	"log/slog"

	"example.com/gogoproject/internal/api"
	"example.com/gogoproject/internal/platform/fault"
)

func createMemberError(err error) api.CreateMemberResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.CreateMember400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrConflict):
		return api.CreateMember409JSONResponse{Error: "duplicate member id"}
	default:
		slog.Error("member request failed", "error", err)
		return api.CreateMember500JSONResponse{Error: "internal server error"}
	}
}

func getMemberError(err error) api.GetMemberResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.GetMember400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrNotFound):
		return api.GetMember404JSONResponse{Error: "not found"}
	default:
		slog.Error("member request failed", "error", err)
		return api.GetMember500JSONResponse{Error: "internal server error"}
	}
}

func updateMemberError(err error) api.UpdateMemberResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.UpdateMember400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrNotFound):
		return api.UpdateMember404JSONResponse{Error: "not found"}
	case errors.Is(err, fault.ErrConflict):
		return api.UpdateMember409JSONResponse{Error: "duplicate member id"}
	default:
		slog.Error("member request failed", "error", err)
		return api.UpdateMember500JSONResponse{Error: "internal server error"}
	}
}

func deleteMemberError(err error) api.DeleteMemberResponseObject {
	switch {
	case errors.Is(err, fault.ErrInvalid):
		return api.DeleteMember400JSONResponse{Error: err.Error()}
	case errors.Is(err, fault.ErrNotFound):
		return api.DeleteMember404JSONResponse{Error: "not found"}
	default:
		slog.Error("member request failed", "error", err)
		return api.DeleteMember500JSONResponse{Error: "internal server error"}
	}
}
