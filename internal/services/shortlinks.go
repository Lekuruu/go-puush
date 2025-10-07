package services

import (
	"errors"
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
)

var minimumIdentifierLength = 3
var maximumIdentifierLength = 16

func CreateShortLink(uploadId int, expiresAt *time.Time, state *app.State) (*database.ShortLink, error) {
	identifier, err := GenerateShortLinkIdentifier(state)
	if err != nil {
		return nil, err
	}

	shortlink := &database.ShortLink{
		Identifier: identifier,
		UploadId:   uploadId,
		ExpiresAt:  expiresAt,
	}

	result := state.Database.Create(shortlink)
	if result.Error != nil {
		return nil, result.Error
	}

	return shortlink, nil
}

func FetchShortLinkByIdentifier(identifier string, state *app.State, preload ...string) (*database.ShortLink, error) {
	shortlink := &database.ShortLink{}
	query := preloadQuery(state, preload).Where("identifier = ?", identifier)
	result := query.First(shortlink)

	if result.Error != nil {
		return nil, result.Error
	}

	return shortlink, nil
}

func ShortLinkExists(identifier string, state *app.State) (bool, error) {
	var count int64
	result := state.Database.Model(&database.ShortLink{}).Where("identifier = ?", identifier).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func GenerateShortLinkIdentifier(state *app.State) (string, error) {
	for i := minimumIdentifierLength; i <= maximumIdentifierLength; i++ {
		identifier := app.RandomString(i)
		exists, err := ShortLinkExists(identifier, state)
		if err != nil {
			return "", err
		}
		if !exists {
			return identifier, nil
		}
	}
	return "", errors.New("could not generate unique identifier")
}
