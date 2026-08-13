package api

import (
	"errors"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/Lekuruu/go-puush/internal/server"
	"github.com/Lekuruu/go-puush/internal/services"
)

// /api/del handles deletion of uploads.
// Once deleted, it returns a history response.
func PuushDelete(ctx *server.Context) {
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

	upload, err := services.FetchUploadById(request.UploadId, ctx.State, "Pool")
	if err != nil {
		ctx.Logger.Error("Failed to fetch upload for deletion", "user_id", user.Id, "upload_id", request.UploadId, "error", err)
		WritePuushError(ctx, ServerError)
		return
	}

	if upload.UserId != user.Id {
		WritePuushError(ctx, ForbiddenError)
		return
	}

	if err := ctx.State.Storage.RemoveThumbnail(upload.Key()); err != nil && !errors.Is(err, fs.ErrNotExist) {
		ctx.Logger.Warn("Failed to remove upload thumbnail", "user_id", user.Id, "upload_id", upload.Id, "error", err)
	}

	// Remove the upload from storage
	err = ctx.State.Storage.RemoveUpload(upload.Key())
	if err != nil {
		ctx.Logger.Error("Failed to remove upload from storage", "user_id", user.Id, "upload_id", upload.Id, "error", err)
		WritePuushError(ctx, ServerError)
		return
	}

	// Remove the upload from the database
	err = services.DeleteUpload(upload, ctx.State)
	if err != nil {
		ctx.Logger.Error("Failed to delete upload record", "user_id", user.Id, "upload_id", upload.Id, "error", err)
		WritePuushError(ctx, ServerError)
		return
	}

	// Update disk usage for user
	err = services.UpdateUserDiskUsage(user.Id, -upload.Filesize, ctx.State)
	if err != nil {
		ctx.Logger.Error("Failed to update disk usage after upload deletion", "user_id", user.Id, "upload_id", upload.Id, "error", err)
		WritePuushError(ctx, ServerError)
		return
	}

	// Update pool upload count
	err = services.UpdatePoolUploadCount(upload.Pool.Id, ctx.State)
	if err != nil {
		ctx.Logger.Error("Failed to update pool count after upload deletion", "user_id", user.Id, "upload_id", upload.Id, "pool_id", upload.Pool.Id, "error", err)
		WritePuushError(ctx, ServerError)
		return
	}

	recentUploads, err := services.FetchRecentUploadsByUser(user, ctx.State, 5, "Pool")
	if err != nil {
		ctx.Logger.Error("Failed to fetch upload history after deletion", "user_id", user.Id, "upload_id", upload.Id, "error", err)
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

	return &DeleteRequest{
		Key:      key,
		UploadId: uploadId,
	}, nil
}
