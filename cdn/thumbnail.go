package cdn

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

func Thumbnail(ctx *app.Context) {
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

	if !upload.IsImage() {
		WriteResponse(415, "", ctx)
		return
	}

	data, err := ctx.State.Storage.ReadThumbnail(upload.Key())
	if err == nil {
		// We found the thumbnail, serve it
		WriteThumbnail(ctx, upload, data)
		return
	}

	// If thumbnail not found, try to generate it
	uploadData, err := ctx.State.Storage.ReadUpload(upload.Key())
	if err != nil {
		// TODO: Original file was not found, queue for deletion
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	thumbnail, err := services.CreateThumbnail(upload.Key(), uploadData, ctx.State)
	if err != nil {
		WriteResponse(500, "An error occurred while generating the thumbnail.", ctx)
		return
	}

	WriteThumbnail(ctx, upload, thumbnail)
}
