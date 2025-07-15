package services

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
)

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
