package api

import (
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

// /api/hist returns the 5 most recent uploads of the authenticated user.
func PuushHistory(ctx *Context) {
	user, err := UserAuthenticationFromContext(ctx)
	if err != nil {
		WritePuushError(ctx, RequestError)
		return
	}

	if user == nil {
		WritePuushError(ctx, RequestError)
		return
	}

	recentUploads, err := services.FetchRecentUploadsByUser(user, ctx.State, 5)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	history := &HistoryResponse{Uploads: recentUploads}
	WritePuushResponse(ctx, history)
}

type HistoryResponse struct {
	Uploads []*database.Upload
}

func (r *HistoryResponse) Serialize() []byte {
	return []byte("0\n")
}
