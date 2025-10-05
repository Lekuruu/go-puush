package cdn

import (
	"fmt"
	"strconv"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
	"github.com/gabriel-vasile/mimetype"
)

func Upload(ctx *app.Context) {
	poolIdentifier := ctx.Vars["identifier"]
	poolPassword := ctx.Vars["password"]
	filename := ctx.Vars["filename"]

	pool, err := services.FetchPoolByIdentifier(poolIdentifier, ctx.State)
	if err != nil {
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	if pool.Type == database.PoolTypePasswordProtected && poolPassword != pool.PasswordHash() {
		WriteResponse(403, "Incorrect password for this puush.", ctx)
		return
	}

	upload, err := services.FetchUploadByFilenameAndPool(filename, pool.Id, ctx.State)
	if err != nil {
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	data, err := ctx.State.Storage.ReadUpload(upload.Key())
	if err != nil {
		// TODO: Queue for deletion
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	// Check if mime type is set, if not, detect it
	if upload.MimeType == "" {
		upload.MimeType = mimetype.Detect(data).String()
		err = services.UpdateUpload(upload, ctx.State)

		if err != nil {
			ctx.State.Logger.Error("Failed to update upload mime type", "error", err)
			WriteResponse(500, "An error occurred while updating the upload.", ctx)
			return
		}
	}

	// Avoid xss attacks by sandboxing html files
	if !upload.IsImage() && !upload.IsVideo() {
		ctx.Response.Header().Set("Content-Security-Policy", "default-src 'none'; sandbox")
		ctx.Response.Header().Set("X-Content-Type-Options", "nosniff")
		ctx.Response.Header().Set("X-Frame-Options", "DENY")
	}

	ctx.Response.Header().Set("Content-Type", upload.MimeType)
	ctx.Response.Header().Set("Content-Length", strconv.Itoa(int(upload.Filesize)))
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", upload.Filename))
	ctx.Response.WriteHeader(200)
	_, err = ctx.Response.Write(data)
	if err != nil {
		WriteResponse(500, "An error occurred while writing the response.", ctx)
		return
	}
}
