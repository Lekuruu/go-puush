package services

import (
	"net/http"
	"strings"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/prplecake/go-thumbnail"
)

func IsImage(data []byte) bool {
	contentType := http.DetectContentType(data[:512])
	return strings.HasPrefix(contentType, "image/")
}

func CreateThumbnail(key string, data []byte, state *app.State) error {
	generator := thumbnail.NewGenerator(thumbnail.Generator{})
	generator.Scaler = "CatmullRom"
	generator.Width = 300
	generator.Height = 300

	image, err := generator.NewImageFromByteArray(data)
	if err != nil {
		return err
	}

	thumbnailData, err := generator.CreateThumbnail(image)
	if err != nil {
		return err
	}

	return state.Storage.SaveThumbnail(key, thumbnailData)
}

func DeleteThumbnail(key string, state *app.State) error {
	return state.Storage.RemoveThumbnail(key)
}
