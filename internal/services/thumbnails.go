package services

import (
	"bytes"

	ffmpeg "github.com/Lekuruu/ffmpeg-go"
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

func CreateThumbnailFromVideo(key string, data []byte, mimeType string, state *app.State) ([]byte, error) {
	// Use ffmpeg to extract a frame from the video data
	// Here we extract the frame at 1 second into the video
	frameStream := ffmpeg.Input("pipe:").
		Filter("select", ffmpeg.Args{"gte(n,1)"}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": "1", "format": "image2", "vcodec": "png"})

	inputBuf := bytes.NewBuffer(data)
	outputBuf := bytes.NewBuffer(nil)
	err := frameStream.WithInput(inputBuf).WithOutput(outputBuf).Run()
	if err != nil {
		return nil, err
	}

	return CreateThumbnail(key, outputBuf.Bytes(), state)
}

func DeleteThumbnail(key string, state *app.State) error {
	return state.Storage.RemoveThumbnail(key)
}
