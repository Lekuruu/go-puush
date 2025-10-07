package cdn

import (
	"fmt"
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
	ctx.Response.Header().Set("Content-Security-Policy", "default-src 'none'; sandbox")
	ctx.Response.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.Response.Header().Set("X-Frame-Options", "DENY")
}

func WriteUpload(ctx *app.Context, upload *database.Upload, data []byte) {
	ctx.Response.Header().Set("Content-Type", upload.MimeType)
	ctx.Response.Header().Set("Content-Length", strconv.Itoa(int(upload.Filesize)))
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", upload.Filename))
	ctx.Response.Header().Set("Last-Modified", upload.CreatedAt.Format(http.TimeFormat))
	ctx.Response.Header().Set("Date", time.Now().Format(http.TimeFormat))
	ctx.Response.Header().Set("ETag", fmt.Sprintf(`"%s"`, upload.Checksum))
	ctx.Response.WriteHeader(200)
	ctx.Response.Write(data)
}

func WriteThumbnail(ctx *app.Context, upload *database.Upload, thumbnail []byte) {
	thumbnailFilename := fmt.Sprintf("thumbnail_%s", upload.Filename)
	ctx.Response.Header().Set("Content-Type", http.DetectContentType(thumbnail))
	ctx.Response.Header().Set("Content-Length", strconv.Itoa(len(thumbnail)))
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", thumbnailFilename))
	ctx.Response.Header().Set("Last-Modified", upload.CreatedAt.Format(http.TimeFormat))
	ctx.Response.Header().Set("Date", time.Now().Format(http.TimeFormat))
	// Thumbnails are not the original upload, so we don't use the same etag
	ctx.Response.Header().Set("ETag", fmt.Sprintf(`"t%s"`, upload.Checksum))
	ctx.Response.WriteHeader(200)
	ctx.Response.Write(thumbnail)
}
