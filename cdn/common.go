package cdn

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
)

func WriteResponse(status int, message string, ctx *app.Context) {
	ctx.Response.WriteHeader(status)
	ctx.Response.Header().Set("Content-Type", "text/plain")
	ctx.Response.Write([]byte(message))
}

func WriteXssHeaders(ctx *app.Context) {
	// Allow media streaming and basic styling while maintaining XSS protection
	ctx.Response.Header().Set(
		"Content-Security-Policy",
		"default-src 'none'; media-src 'self'; img-src 'self' data:; style-src 'unsafe-inline'; sandbox",
	)
	ctx.Response.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.Response.Header().Set("X-Frame-Options", "DENY")
}

func WriteUpload(ctx *app.Context, upload *database.Upload, stream io.ReadSeekCloser) {
	ctx.Response.Header().Set("Content-Type", upload.MimeType)
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", upload.Filename))

	defer stream.Close()
	http.ServeContent(ctx.Response, ctx.Request, upload.Filename, upload.CreatedAt, stream)
}

func WriteThumbnail(ctx *app.Context, upload *database.Upload, thumbnail []byte) {
	thumbnailFilename := fmt.Sprintf("thumbnail_%s", upload.Filename)
	ctx.Response.Header().Set("Content-Type", http.DetectContentType(thumbnail))
	ctx.Response.Header().Set("Content-Length", strconv.Itoa(len(thumbnail)))
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", thumbnailFilename))
	ctx.Response.Header().Set("Last-Modified", upload.CreatedAt.Format(http.TimeFormat))
	ctx.Response.Header().Set("Date", time.Now().Format(http.TimeFormat))
	ctx.Response.Header().Set("ETag", fmt.Sprintf(`"t%s"`, upload.Checksum))
	ctx.Response.WriteHeader(200)
	ctx.Response.Write(thumbnail)
}
