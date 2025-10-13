package services

import (
	"errors"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
)

const minimumPoolIdentifierLength = 6
const maximumPoolIdentifierLength = 16

func CreatePool(pool *database.Pool, state *app.State) error {
	result := state.Database.Create(pool)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchPoolById(id int, state *app.State, preload ...string) (*database.Pool, error) {
	pool := &database.Pool{}
	result := preloadQuery(state, preload).First(pool, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return pool, nil
}

func FetchPoolByIdentifier(identifier string, state *app.State, preload ...string) (*database.Pool, error) {
	pool := &database.Pool{}
	query := preloadQuery(state, preload).Where("identifier = ?", identifier)
	result := query.First(pool)

	if result.Error != nil {
		return nil, result.Error
	}

	return pool, nil
}

func FetchPoolByUserAndName(userId int, name string, state *app.State, preload ...string) (*database.Pool, error) {
	pool := &database.Pool{}
	query := preloadQuery(state, preload).Where("user_id = ? AND name = ?", userId, name)
	result := query.First(pool)

	if result.Error != nil {
		return nil, result.Error
	}

	return pool, nil
}

func PoolExists(identifier string, state *app.State) (bool, error) {
	var count int64
	result := state.Database.Model(&database.Pool{}).Where("identifier = ?", identifier).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func UpdatePool(pool *database.Pool, state *app.State) error {
	result := state.Database.Save(pool)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GeneratePoolIdentifier(state *app.State) (string, error) {
	for i := minimumPoolIdentifierLength; i <= maximumPoolIdentifierLength; i++ {
		identifier := app.RandomString(i)
		exists, err := PoolExists(identifier, state)
		if err != nil {
			return "", err
		}
		if !exists {
			return identifier, nil
		}
	}
	return "", errors.New("could not generate unique identifier")
}
