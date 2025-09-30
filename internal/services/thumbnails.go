package services

import (
	"net/http"
	"strings"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/prplecake/go-thumbnail"
)

func IsImage(data []byte) bool {
	if len(data) < 512 {
		return false
	}
	contentType := http.DetectContentType(data[:512])
	return strings.HasPrefix(contentType, "image/")
}

func CreateThumbnail(key string, data []byte, state *app.State) ([]byte, error) {
	generator := thumbnail.NewGenerator(thumbnail.Generator{})
	generator.Scaler = "CatmullRom"
	generator.Height = 150
	generator.Width = 150

	image, err := generator.NewImageFromByteArray(data)
	if err != nil {
		return nil, err
	}

	thumbnailData, err := generator.CreateThumbnail(image)
	if err != nil {
		return nil, err
	}

	err = state.Storage.SaveThumbnail(key, thumbnailData)
	if err != nil {
		return nil, err
	}

	return thumbnailData, nil
}

func DeleteThumbnail(key string, state *app.State) error {
	return state.Storage.RemoveThumbnail(key)
}
