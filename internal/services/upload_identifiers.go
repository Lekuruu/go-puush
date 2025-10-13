package services

import (
	"errors"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
)

const minimumLinkIdentifierLength = 5
const maximumLinkIdentifierLength = 16

func CreateUploadIdentifier(uploadId int, state *app.State) (string, error) {
	identifier, err := GenerateUploadIdentifier(state)
	if err != nil {
		return "", err
	}

	// Update the upload with the shortlink identifier
	result := state.Database.Model(&database.Upload{}).
		Where("id = ?", uploadId).
		Update("identifier", identifier)

	if result.Error != nil {
		return "", result.Error
	}

	return identifier, nil
}

func FetchUploadByIdentifier(identifier string, state *app.State, preload ...string) (*database.Upload, error) {
	upload := &database.Upload{}
	query := preloadQuery(state, preload).Where("identifier = ?", identifier)
	result := query.First(upload)

	if result.Error != nil {
		return nil, result.Error
	}

	return upload, nil
}

func FetchManyUploadsByIdentifiers(identifiers []string, state *app.State, preload ...string) ([]*database.Upload, error) {
	var uploads []*database.Upload
	query := preloadQuery(state, preload).Where("identifier IN ?", identifiers)
	result := query.Find(&uploads)

	if result.Error != nil {
		return nil, result.Error
	}

	return uploads, nil
}

func UploadIdentifierExists(identifier string, state *app.State) (bool, error) {
	var count int64
	result := state.Database.Model(&database.Upload{}).Where("identifier = ?", identifier).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func GenerateUploadIdentifier(state *app.State) (string, error) {
	for i := minimumLinkIdentifierLength; i <= maximumLinkIdentifierLength; i++ {
		identifier := app.RandomString(i)
		exists, err := UploadIdentifierExists(identifier, state)
		if err != nil {
			return "", err
		}
		if !exists {
			return identifier, nil
		}
	}
	return "", errors.New("could not generate unique identifier")
}
