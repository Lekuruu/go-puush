package api

import (
	"mime/multipart"
	"net/http"
)

func GetMultipartFormValue(request *http.Request, key string) string {
	values, ok := request.MultipartForm.Value[key]
	if !ok || len(values) == 0 {
		return ""
	}

	return values[0]
}

func GetMultipartFormFile(request *http.Request, key string) *multipart.FileHeader {
	files, ok := request.MultipartForm.File[key]
	if !ok || len(files) == 0 {
		return nil
	}

	return files[0]
}
