package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/services"
)

// /api/del handles deletion of uploads.
// Once deleted, it returns a history response.
func PuushDelete(ctx *app.Context) {
	request, err := NewDeleteRequest(ctx.Request)
	if err != nil {
		WritePuushError(ctx, RequestError)
		return
	}

	user, err := UserAuthenticationFromKey(request.Key, ctx.State)
	if err != nil {
		WritePuushError(ctx, AuthenticationFailure)
		return
	}

	upload, err := services.FetchUploadById(request.UploadId, ctx.State, "Pool", "User")
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	if upload.UserId != user.Id {
		WritePuushError(ctx, ForbiddenError)
		return
	}

	// Remove thumbnail if it exists, do nothing on error
	ctx.State.Storage.RemoveThumbnail(upload.Key())

	// Remove the upload from storage
	err = ctx.State.Storage.RemoveUpload(upload.Key())
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	// Remove the upload from the database
	err = services.DeleteUpload(upload, ctx.State)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	// Update disk usage for user
	err = services.UpdateUserDiskUsage(user.Id, -upload.Filesize, ctx.State)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	recentUploads, err := services.FetchRecentUploadsByUser(user, ctx.State, 5, "Pool")
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	history := &HistoryResponse{
		CdnUrl:  ctx.State.Config.Cdn.Url,
		Uploads: recentUploads,
		User:    user,
	}
	WritePuushResponse(ctx, history)
}

type DeleteRequest struct {
	Key      string
	UploadId int
}

func NewDeleteRequest(request *http.Request) (*DeleteRequest, error) {
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

	junk := request.FormValue("z")
	if junk != "poop" {
		return nil, errors.New("invalid request parameter 'z'")
	}

	return &DeleteRequest{
		Key:      key,
		UploadId: uploadId,
	}, nil
}
