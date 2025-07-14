package services

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
)

func CreateUpload(upload *database.Upload, state *app.State) error {
	result := state.Database.Create(upload)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchUploadById(id int, state *app.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	result := preloadQuery(state, preload).First(upload, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchUploadByChecksum(checksum string, state *app.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	query := preloadQuery(state, preload).Where("checksum = ?", checksum)
	result := query.First(upload)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func UpdateUpload(upload *database.Upload, state *app.State) error {
	result := state.Database.Save(upload)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteUpload(upload *database.Upload, state *app.State) error {
	result := state.Database.Delete(upload)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
