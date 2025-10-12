package api

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
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

	user, err := UserAuthenticationFromKey(request.Key, ctx.State, "DefaultPool", "DefaultPool.Uploads")
	if err != nil {
		WritePuushError(ctx, AuthenticationFailure)
		return
	}

	if request.FileSize <= 0 {
		WritePuushError(ctx, RequestError)
		return
	}

	fileData := make([]byte, request.FileSize)
	_, err = request.File.Read(fileData)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	if !request.CompareChecksum(fileData) {
		WritePuushError(ctx, ChecksumError)
		return
	}

	err = services.UpdateUserLatestActivity(user.Id, ctx.State)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	// Check if another upload with the same checksum already exists
	existingUpload, err := services.FetchUploadByChecksum(request.FileChecksum, ctx.State, "Pool")
	if err == nil {
		response := &UploadResponse{
			UploadUrl:        ctx.State.Config.Cdn.Url + existingUpload.UrlEncoded(),
			UpdatedDiskUsage: user.DiskUsage,
		}
		WritePuushResponse(ctx, response)
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

	upload := &database.Upload{
		UserId:    user.Id,
		PoolId:    user.DefaultPoolId,
		Filename:  request.FileName,
		Filesize:  request.FileSize,
		Checksum:  request.FileChecksum,
		MimeType:  mimetype.Detect(fileData).String(),
		CreatedAt: time.Now(),
		Pool:      user.DefaultPool,
		User:      user,
	}

	err = services.CreateUpload(upload, ctx.State)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	err = ctx.State.Storage.SaveUpload(upload.Key(), fileData)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	if upload.IsImage() {
		// Try to generate a thumbnail & do nothing if it fails
		services.CreateThumbnail(upload.Key(), fileData, ctx.State)
	}

	user.DefaultPool.UploadCount = len(user.DefaultPool.Uploads) + 1
	err = services.UpdatePool(user.DefaultPool, ctx.State)
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

	shortlink, err := services.CreateShortLink(upload.Id, nil, ctx.State)
	if err != nil {
		WritePuushError(ctx, ServerError)
		return
	}

	uploadResponse := &UploadResponse{
		UploadUrl:        ctx.State.Config.Cdn.Url + shortlink.UrlEncoded(),
		UpdatedDiskUsage: user.DiskUsage,
	}
	WritePuushResponse(ctx, uploadResponse)
}

type UploadRequest struct {
	Key          string
	FileChecksum string
	FileName     string
	FileSize     int64
	File         multipart.File
}

func (request *UploadRequest) CompareChecksum(data []byte) bool {
	checksum := md5.New()
	checksum.Write(data)
	checksumBytes := checksum.Sum(nil)
	checksumHex := strings.ToLower(hex.EncodeToString(checksumBytes))
	return checksumHex == request.FileChecksum
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

	fileChecksum := app.GetMultipartFormValue(request, "c")
	if fileChecksum == "" {
		return nil, errors.New("missing file checksum")
	}

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
