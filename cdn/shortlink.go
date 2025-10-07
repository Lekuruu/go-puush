package cdn

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/services"
)

func UploadShortlink(ctx *app.Context) {
	identifier := ctx.Vars["identifier"]

	link, err := services.FetchShortLinkByIdentifier(identifier, ctx.State, "Upload", "Upload.Pool")
	if err != nil {
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	// Pass the request to the normal upload handler
	ctx.Vars["filename"] = link.Upload.Filename
	ctx.Vars["pool"] = link.Upload.Pool.Identifier
	Upload(ctx)
}

func ThumbnailShortlink(ctx *app.Context) {
	identifier := ctx.Vars["identifier"]

	link, err := services.FetchShortLinkByIdentifier(identifier, ctx.State, "Upload", "Upload.Pool")
	if err != nil {
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	// Pass the request to the normal thumbnail handler
	ctx.Vars["filename"] = link.Upload.Filename
	ctx.Vars["pool"] = link.Upload.Pool.Identifier
	Thumbnail(ctx)
}
