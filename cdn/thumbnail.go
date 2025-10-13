package cdn

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

func Thumbnail(ctx *app.Context) {
	poolPassword := ctx.Vars["password"]
	poolIdentifier := ctx.Vars["pool"]
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

	data, err := ctx.State.Storage.ReadThumbnail(upload.Key())
	if err != nil {
		WriteResponse(404, "That puush does not have a thumbnail.", ctx)
		return
	}

	WriteThumbnail(ctx, upload, data)
}
