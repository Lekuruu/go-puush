package api

import (
	"net/http"

	"github.com/Lekuruu/go-puush/internal/database"
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

	// TODO: Retrieve user uploads
	history := &HistoryResponse{Uploads: []*database.Upload{}}
	ctx.Response.WriteHeader(http.StatusOK)
	_, err = ctx.Response.Write(history.Serialize())
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}
}

type HistoryResponse struct {
	Uploads []*database.Upload
}

func (r *HistoryResponse) Serialize() []byte {
	return []byte("0\n")
}
