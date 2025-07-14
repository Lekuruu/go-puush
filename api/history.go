package api

import (
	"net/http"

	"github.com/Lekuruu/go-puush/internal/database"
)

type HistoryResponse struct {
	Uploads []*database.Upload
}

func (r *HistoryResponse) Serialize() []byte {
	return []byte("0\n")
}

func PuushHistory(ctx *Context) {
	user, err := UserAuthenticationFromContext(ctx)
	if err != nil {
		WritePuushError(ctx, -1, http.StatusUnauthorized)
		return
	}

	if user == nil {
		WritePuushError(ctx, -2, http.StatusInternalServerError)
		return
	}

	// TODO: Retrieve user uploads
	history := &HistoryResponse{Uploads: []*database.Upload{}}
	ctx.Response.WriteHeader(http.StatusOK)
	_, err = ctx.Response.Write(history.Serialize())
	if err != nil {
		WritePuushError(ctx, -2, http.StatusInternalServerError)
		return
	}
}
