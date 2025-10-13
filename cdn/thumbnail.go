package cdn

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

func Thumbnail(ctx *app.Context) {
	poolIdentifier := ctx.Vars["pool"]
	poolPassword := ctx.Vars["password"]
	identifier := ctx.Vars["identifier"]

	upload, err := services.FetchUploadByIdentifier(identifier, ctx.State)
	if err != nil {
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	if upload.Pool.Type == database.PoolTypePrivate && upload.Pool.Identifier != poolIdentifier {
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	if upload.Pool.Type == database.PoolTypePasswordProtected && poolPassword != upload.Pool.PasswordHash() {
		WriteResponse(403, "Incorrect password for this puush.", ctx)
		return
	}

	data, err := ctx.State.Storage.ReadThumbnail(upload.Key())
	if err != nil {
		WriteResponse(404, "That puush does not have a thumbnail.", ctx)
		return
	}

	WriteThumbnail(ctx, upload, data)
}
