package services

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/prplecake/go-thumbnail"
)

func CreateThumbnail(key string, data []byte, state *app.State) ([]byte, error) {
	generator := thumbnail.NewGenerator(thumbnail.Generator{})
	generator.Scaler = "CatmullRom"
	generator.Height = 100
	generator.Width = 100

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
