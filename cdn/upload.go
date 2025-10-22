package cdn

import (
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

// This will limit the view count increase to once per minute per IP
var uploadViewCooldowns = app.NewCooldownManager(time.Minute)

func Upload(ctx *app.Context) {
	poolIdentifier := ctx.Vars["pool"]
	poolPassword := ctx.Vars["password"]
	identifier := ctx.Vars["identifier"]

	upload, err := services.FetchUploadByIdentifier(identifier, ctx.State, "Pool")
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

	stream, err := ctx.State.Storage.ReadUploadStream(upload.Key())
	if err != nil {
		// TODO: Original file was not found, queue for deletion
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}
	defer stream.Close()

	// Try to increase views, if cooldown is not active
	if uploadViewCooldowns.Allow(ctx.IP()) {
		upload.Views += 1
		services.UpdateUpload(upload, ctx.State)
	}

	if !upload.IsImage() && !upload.IsVideo() {
		// Avoid xss attacks by sandboxing html files
		WriteXssHeaders(ctx)
	}

	WriteUpload(ctx, upload, stream)
}
