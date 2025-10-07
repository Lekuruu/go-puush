package api

import (
	"strconv"
	"strings"
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

// /api/hist returns the 5 most recent uploads of the authenticated user.
func PuushHistory(ctx *app.Context) {
	user, err := UserAuthenticationFromContext(ctx)
	if err != nil {
		WritePuushError(ctx, RequestError)
		return
	}

	if user == nil {
		WritePuushError(ctx, RequestError)
		return
	}

	recentUploads, err := services.FetchRecentUploadsByUser(user, ctx.State, 5, "Pool", "Link")
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
		uploadUrl := r.CdnUrl + upload.UrlEncoded()

		if upload.Link != nil {
			// We have a short link for this upload, use it instead
			uploadUrl = r.CdnUrl + upload.Link.UrlEncoded()
		}

		var historyItem = []string{
			strconv.Itoa(upload.Id),
			upload.CreatedAt.Format(time.DateTime),
			uploadUrl,
			upload.Filename,
			strconv.Itoa(upload.Views),
		}
		data = append(data, strings.Join(historyItem, ","))
	}

	return []byte(strings.Join(data, "\n"))
}
