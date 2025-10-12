package routes

import (
	"net/http"
	"os"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/services"
)

// Response when the thumbnail is not found
const DefaultThumbnailPath = "web/static/img/unknown.png"

var defaultThumbnailData []byte

func Thumbnail(ctx *app.Context) {
	user, err := GetUserSession(ctx)
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	identifier := ctx.Vars["identifier"]
	if identifier == "" {
		renderRaw(200, "image/png", defaultThumbnailData, ctx)
		return
	}

	link, err := services.FetchShortLinkByIdentifier(identifier, ctx.State, "Upload")
	if err != nil {
		renderRaw(200, "image/png", defaultThumbnailData, ctx)
		return
	}
	if link == nil || link.Upload == nil {
		renderRaw(200, "image/png", defaultThumbnailData, ctx)
		return
	}

	if link.Upload.UserId != user.Id {
		// User does not own this upload
		renderRaw(200, "image/png", defaultThumbnailData, ctx)
		return
	}

	image, err := ctx.State.Storage.ReadThumbnail(link.Upload.Key())
	if err != nil {
		renderRaw(200, "image/png", defaultThumbnailData, ctx)
		return
	}

	renderRaw(200, http.DetectContentType(image), image, ctx)
}

func init() {
	// Load default thumbnail data
	defaultThumbnailData, _ = os.ReadFile(DefaultThumbnailPath)
}
