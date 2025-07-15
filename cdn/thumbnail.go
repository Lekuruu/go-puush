package cdn

import (
	"fmt"
	"net/http"
	"strconv"

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

	if !upload.IsImage {
		WriteResponse(415, "", ctx)
		return
	}

	data, err := ctx.State.Storage.ReadThumbnail(upload.Key())
	if err == nil {
		thumbnailFilename := fmt.Sprintf("thumbnail_%s", upload.Filename)
		ctx.Response.Header().Set("Content-Type", http.DetectContentType(data))
		ctx.Response.Header().Set("Content-Length", strconv.Itoa(len(data)))
		ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", thumbnailFilename))
		ctx.Response.WriteHeader(200)
		_, err = ctx.Response.Write(data)
		if err != nil {
			WriteResponse(500, "An error occurred while writing the response.", ctx)
		}
		return
	}

	// If thumbnail not found, try to generate it
	uploadData, err := ctx.State.Storage.ReadUpload(upload.Key())
	if err != nil {
		// TODO: Queue for deletion
		WriteResponse(404, "That puush could not be found.", ctx)
		return
	}

	thumbnail, err := services.CreateThumbnail(upload.Key(), uploadData, ctx.State)
	if err != nil {
		WriteResponse(500, "An error occurred while generating the thumbnail.", ctx)
		return
	}

	thumbnailFilename := fmt.Sprintf("thumbnail_%s", upload.Filename)
	ctx.Response.Header().Set("Content-Type", http.DetectContentType(thumbnail))
	ctx.Response.Header().Set("Content-Length", strconv.Itoa(len(thumbnail)))
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", thumbnailFilename))
	ctx.Response.WriteHeader(200)
	_, err = ctx.Response.Write(thumbnail)
	if err != nil {
		WriteResponse(500, "An error occurred while writing the response.", ctx)
		return
	}
}
