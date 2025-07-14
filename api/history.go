package api

import (
	"strconv"
	"strings"
	"time"

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

	recentUploads, err := services.FetchRecentUploadsByUser(user, ctx.State, 5, "Pool")
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	history := &HistoryResponse{
		CdnUrl:  ctx.State.Config.Cdn.Url,
		Uploads: recentUploads,
		User:    user,
	}
	WritePuushResponse(ctx, history)
}

type HistoryResponse struct {
	CdnUrl  string
	Uploads []*database.Upload
	User    *database.User
}

func (r *HistoryResponse) Serialize() []byte {
	var data = []string{strconv.Itoa(len(r.Uploads))}

	for _, upload := range r.Uploads {
		var historyItem = []string{
			strconv.Itoa(upload.Id),
			upload.CreatedAt.Format(time.RFC850),
			r.CdnUrl + upload.Url(),
			upload.Filename,
			strconv.Itoa(upload.Views),
		}
		data = append(data, strings.Join(historyItem, ","))
	}

	return []byte(strings.Join(data, "\n"))
}
