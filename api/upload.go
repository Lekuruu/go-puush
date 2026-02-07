package api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
	"github.com/gabriel-vasile/mimetype"
)

// /api/up is the main endpoint for uploading files to the puush service.
func PuushUpload(ctx *app.Context) {
	request, err := NewUploadRequest(ctx.Request)
	if err != nil {
		WritePuushError(ctx, RequestError)
		return
	}
	defer request.File.Close()

	user, err := UserAuthenticationFromKey(request.Key, ctx.State, "DefaultPool")
	if err != nil {
		WritePuushError(ctx, AuthenticationFailure)
		return
	}

	err = services.UpdateUserLatestActivity(user.Id, ctx.State)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	// Check if another upload with the same checksum already exists
	existingUpload, err := services.FetchUploadByChecksum(request.FileChecksum, ctx.State, "Pool")
	if err == nil && existingUpload.UserId == user.Id {
		response := &UploadResponse{
			UploadUrl:        ctx.State.Config.Cdn.Url + existingUpload.UrlEncoded(),
			UpdatedDiskUsage: user.DiskUsage,
		}
		WritePuushResponse(ctx, response)
		return
	}

	if request.FileSize <= 0 {
		WritePuushError(ctx, RequestError)
		return
	}

	if user.DefaultPool == nil {
		WritePuushError(ctx, ServerError)
		return
	}

	if request.ExceedsUploadLimit(user) {
		WritePuushError(ctx, InsufficientStorageError)
		return
	}

	// Read first 512 bytes for MIME detection
	mimeBuffer := make([]byte, 512)
	n, err := io.ReadFull(request.File, mimeBuffer)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		WritePuushError(ctx, ServerError)
		return
	}
	mimeBuffer = mimeBuffer[:n]
	mimeType := mimetype.Detect(mimeBuffer).String()

	// Create identifier for the upload
	identifierLength := user.DefaultPool.UploadIdentifierLength()
	identifier, err := services.GenerateUploadIdentifier(identifierLength, ctx.State)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	upload := &database.Upload{
		UserId:     user.Id,
		PoolId:     user.DefaultPoolId,
		Filename:   request.FileName,
		Filesize:   request.FileSize,
		Checksum:   request.FileChecksum,
		MimeType:   mimeType,
		Identifier: identifier,
		CreatedAt:  time.Now(),
		Pool:       user.DefaultPool,
		User:       user,
	}

	// Stream file to disk while computing checksum at the same time
	hash := md5.New()
	combinedReader := io.MultiReader(bytes.NewReader(mimeBuffer), request.File)
	teeReader := io.TeeReader(combinedReader, hash)

	err = ctx.State.Storage.SaveUploadStream(upload.Key(), teeReader)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	// Perform post-upload actions in a separate goroutine to respond faster
	go performPostUploadActions(upload, ctx.State)

	upload.Checksum = hex.EncodeToString(hash.Sum(nil))
	err = services.CreateUpload(upload, ctx.State)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	user.DiskUsage += upload.Filesize
	err = services.UpdateUserDiskUsage(user.Id, upload.Filesize, ctx.State)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	err = services.UpdatePoolUploadCount(upload.Pool.Id, ctx.State)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	uploadResponse := &UploadResponse{
		UploadUrl:        ctx.State.Config.Cdn.Url + upload.UrlEncoded(),
		UpdatedDiskUsage: user.DiskUsage,
	}
	WritePuushResponse(ctx, uploadResponse)
}

func performPostUploadActions(upload *database.Upload, state *app.State) {
	if !upload.IsImage() && !upload.IsVideo() {
		return
	}

	data, err := state.Storage.ReadUpload(upload.Key())
	if err != nil {
		state.Logger.Logf("Failed to read upload for post-upload actions: %v", err)
		return
	}

	CreateThumbnailFromUpload(upload, data, state)
}

type UploadRequest struct {
	Key          string
	FileChecksum string
	FileName     string
	FileSize     int64
	File         multipart.File
}

func (request *UploadRequest) ExceedsUploadLimit(user *database.User) bool {
	if user.UploadLimit() < 0 {
		// No limit for unlimited accounts
		return false
	}
	return user.DiskUsage+request.FileSize > user.UploadLimit()
}

func NewUploadRequest(request *http.Request) (*UploadRequest, error) {
	err := request.ParseMultipartForm(100 << 20) // 100 MB
	if err != nil {
		return nil, err
	}

	key := app.GetMultipartFormValue(request, "k")
	if key == "" {
		return nil, errors.New("missing api key")
	}

	file := app.GetMultipartFormFile(request, "f")
	if file == nil {
		return nil, errors.New("missing file")
	}

	// This argument is optional (ShareX doesn't provide it)
	fileChecksum := app.GetMultipartFormValue(request, "c")

	fileName := file.Filename
	fileStream, err := file.Open()
	if err != nil {
		return nil, err
	}

	// Replace commas & slashes with underscores in the filename
	// to avoid issues in history response parsing
	fileName = strings.ReplaceAll(fileName, ",", "_")
	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName = strings.Trim(fileName, " ")

	// Ensure filename is not empty
	if fileName == "" {
		return nil, errors.New("filename cannot be empty")
	}

	return &UploadRequest{
		Key:          key,
		FileChecksum: fileChecksum,
		FileName:     fileName,
		FileSize:     file.Size,
		File:         fileStream,
	}, nil
}

type UploadResponse struct {
	UploadUrl        string
	UpdatedDiskUsage int64
}

func (r *UploadResponse) Serialize() []byte {
	data := []string{
		"0", // Status "Success"
		r.UploadUrl,
		strconv.Itoa(int(r.UpdatedDiskUsage)),
		strconv.Itoa(int(r.UpdatedDiskUsage)),
	}
	return []byte(strings.Join(data, ","))
}
