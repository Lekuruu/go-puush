package services

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
)

func CreateUser(user *database.User, state *app.State) error {
	result := state.Database.Create(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FetchUserById(id int, state *app.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	result := preloadQuery(state, preload).First(user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByName(name string, state *app.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	query := preloadQuery(state, preload).Where("name = ?", name)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func FetchUserByEmail(email string, state *app.State, preload ...string) (*database.User, error) {
	user := &database.User{}
	query := preloadQuery(state, preload).Where("email = ?", email)
	result := query.First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
