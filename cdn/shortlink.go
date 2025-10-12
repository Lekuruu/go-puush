package cdn

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

func UploadShortlink(ctx *app.Context) {
	identifier := ctx.Vars["identifier"]
	poolIdentifier := ctx.Vars["pool"]

	upload, err := ResolveShortLink(identifier, poolIdentifier, ctx.State)
	if err != nil {
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	if upload.Pool.Type == database.PoolTypePrivate && upload.Pool.Identifier != poolIdentifier {
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	// Pass the request to the normal upload handler
	ctx.Vars["filename"] = upload.Filename
	ctx.Vars["pool"] = upload.Pool.Identifier
	Upload(ctx)
}

func ThumbnailShortlink(ctx *app.Context) {
	identifier := ctx.Vars["identifier"]
	poolIdentifier := ctx.Vars["pool"]

	upload, err := ResolveShortLink(identifier, poolIdentifier, ctx.State)
	if err != nil {
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	if upload.Pool.Type == database.PoolTypePrivate && upload.Pool.Identifier != poolIdentifier {
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	// Pass the request to the normal thumbnail handler
	ctx.Vars["filename"] = upload.Filename
	ctx.Vars["pool"] = upload.Pool.Identifier
	Thumbnail(ctx)
}

// ResolveShortLink tries to resolve the given identifier as a shortlink first,
// then as a filename in the given pool if that fails.
func ResolveShortLink(identifier string, poolIdentifier string, state *app.State) (*database.Upload, error) {
	link, err := services.FetchShortLinkByIdentifier(identifier, state, "Upload", "Upload.Pool")
	if err == nil {
		return link.Upload, nil
	}

	pool, err := services.FetchPoolByIdentifier(poolIdentifier, state)
	if err != nil {
		return nil, err
	}

	upload, err := services.FetchUploadByFilenameAndPool(identifier, pool.Id, state, "Pool")
	if err == nil {
		return upload, nil
	}

	return nil, err
}
