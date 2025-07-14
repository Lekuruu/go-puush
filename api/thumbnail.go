package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/services"
)

// /api/thumb returns a thumbnail for the given upload ID.
// If the thumbnail is not available, it tries to generate one.
func PuushThumbnail(ctx *app.Context) {
	request, err := NewThumbnailRequest(ctx.Request)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := UserAuthenticationFromKey(request.Key, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusUnauthorized)
		return
	}

	upload, err := services.FetchUploadById(request.UploadId, ctx.State, "Pool", "User")
	if err != nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	if upload.UserId != user.Id {
		ctx.Response.WriteHeader(http.StatusForbidden)
		return
	}

	if !upload.IsImage {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	thumbnail, err := ctx.State.Storage.ReadThumbnail(upload.Key())
	if err == nil {
		ctx.Response.Header().Set("Content-Type", http.DetectContentType(thumbnail))
		ctx.Response.WriteHeader(http.StatusOK)
		ctx.Response.Write(thumbnail)
		return
	}

	uploadData, err := ctx.State.Storage.ReadUpload(upload.Key())
	if err != nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	// Try to generate a thumbnail if it doesn't exist
	thumbnail, err = services.CreateThumbnail(upload.Key(), uploadData, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Response.Header().Set("Content-Type", http.DetectContentType(thumbnail))
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write(thumbnail)
}

type ThumbnailRequest struct {
	Key      string
	UploadId int
}

func NewThumbnailRequest(request *http.Request) (*ThumbnailRequest, error) {
	err := request.ParseForm()
	if err != nil {
		return nil, err
	}

	key := request.FormValue("k")
	if key == "" {
		return nil, errors.New("missing api key")
	}

	uploadIdStr := request.FormValue("i")
	if uploadIdStr == "" {
		return nil, errors.New("missing upload ID")
	}

	uploadId, err := strconv.Atoi(uploadIdStr)
	if err != nil {
		return nil, errors.New("invalid upload ID")
	}

	return &ThumbnailRequest{
		Key:      key,
		UploadId: uploadId,
	}, nil
}
