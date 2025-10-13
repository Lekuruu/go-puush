package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"

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

	upload, err := services.FetchUploadByIdentifier(identifier, ctx.State, "Pool")
	if err != nil {
		renderRaw(200, "image/png", defaultThumbnailData, ctx)
		return
	}
	if upload == nil {
		renderRaw(200, "image/png", defaultThumbnailData, ctx)
		return
	}

	if upload.UserId != user.Id {
		// User does not own this upload
		renderRaw(200, "image/png", defaultThumbnailData, ctx)
		return
	}

	image, err := ctx.State.Storage.ReadThumbnail(upload.Key())
	if err != nil {
		renderRaw(200, "image/png", defaultThumbnailData, ctx)
		return
	}

	ctx.Response.Header().Set("Content-Type", http.DetectContentType(image))
	ctx.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(image)))
	ctx.Response.Header().Set("Last-Modified", upload.CreatedAt.Format(http.TimeFormat))
	ctx.Response.Header().Set("Date", time.Now().Format(http.TimeFormat))
	ctx.Response.Header().Set("ETag", fmt.Sprintf(`"t%s"`, upload.Checksum))
	ctx.Response.WriteHeader(200)
	ctx.Response.Write(image)
}

func init() {
	// Load default thumbnail data
	defaultThumbnailData, _ = os.ReadFile(DefaultThumbnailPath)
}
