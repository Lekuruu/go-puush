package cdn

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
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
		// TODO: Original file was not found, queue for deletion
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	if !upload.IsImage() && !upload.IsVideo() {
		// Avoid xss attacks by sandboxing html files
		WriteXssHeaders(ctx)
	}

	WriteUpload(ctx, upload, data)
}
