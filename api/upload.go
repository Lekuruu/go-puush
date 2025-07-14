package api

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

// /api/up is the main endpoint for uploading files to the puush service.
func PuushUpload(ctx *Context) {
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

	if user.DefaultPool == nil {
		WritePuushError(ctx, ServerError)
		return
	}

	// TODO: Create upload
	placeholderResponse := &UploadResponse{
		UploadUrl:        "http://i.imgur.com/nFfry2P.mp4",
		UpdatedDiskUsage: user.DiskUsage + request.FileSize,
	}
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write(placeholderResponse.Serialize())
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

func (request *UploadRequest) FilenameChecksum() string {
	checksum := md5.New()
	checksum.Write([]byte(request.FileName))
	checksumBytes := checksum.Sum(nil)
	return strings.ToLower(hex.EncodeToString(checksumBytes))
}

func NewUploadRequest(request *http.Request) (*UploadRequest, error) {
	err := request.ParseMultipartForm(100 << 20) // 100 MB
	if err != nil {
		return nil, err
	}

	key := GetMultipartFormValue(request, "k")
	if key == "" {
		return nil, errors.New("missing api key")
	}

	file := GetMultipartFormFile(request, "f")
	if err != nil {
		return nil, err
	}

	fileChecksum := GetMultipartFormValue(request, "c")
	if fileChecksum == "" {
		return nil, errors.New("missing file checksum")
	}

	fileName := file.Filename
	fileStream, err := file.Open()
	if err != nil {
		return nil, err
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
